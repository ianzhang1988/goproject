package lambda_service

import (
	"fmt"
	"strings"
	"time"

	"github.com/bluele/gcache"
)

type Function struct {
	Name string
	Ver  string
}

type Device struct {
	Sn       string
	Isp      string
	Province string
	City     string
	Funcs    []Function // necessary?
}

// 从db或者其他途径获取设备信息
type DevsGeter interface {
	GetDevs(FuncName string, version string) (map[string]*Device, error)
}

// 按照affinity筛选设备，实现可支持cache，index等功能。
type AffinityDevs interface {
	GetAffinityDevs(FuncName string, version string, aff Affinity) ([]*Device, error)
}

type TestDevs struct {
	devs []*Device
}

func (td *TestDevs) GetDevs(FuncName string, version string) (map[string]*Device, error) {
	sn2dev := map[string]*Device{}
	for _, d := range td.devs {
		for _, f := range d.Funcs {
			if f.Name != FuncName {
				continue
			}

			if version == "all" {
				sn2dev[d.Sn] = d
			} else {
				if f.Ver == version {
					sn2dev[d.Sn] = d
				}
			}
		}
	}

	return sn2dev, nil
}

func NewTestDevs() *TestDevs {
	td := TestDevs{}
	// test data
	functions := []Function{
		{Name: "func1", Ver: "v1"},
		{Name: "func1", Ver: "v2"},
		{Name: "func2", Ver: "v1"},
	}

	ispName := []string{"ct", "cnc", "cmnet"}
	provienceName := []string{"shangdong", "hebei", "guangdong", "hunan"}
	cityName := []string{"city1", "city2", "city3", "city4", "city5"}

	for i := 0; i < 100; i++ {
		dev := &Device{
			Sn:       fmt.Sprintf("sn_%d", i),
			Isp:      ispName[i%len(ispName)],
			Province: provienceName[i%len(provienceName)],
			City:     cityName[i%len(cityName)],
			Funcs:    functions,
		}
		td.devs = append(td.devs, dev)
	}

	return &td
}

func (aff Affinity) match(d *Device) bool {
	if aff.All {
		return true
	}

	flag := true

	equal_string := func(aff_value, dev_value string) {
		if aff_value == "" {
			flag = flag && true
		} else if aff_value == dev_value {
			flag = flag && true
		} else {
			flag = flag && false
		}
	}

	equal_string(aff.City, d.City)
	equal_string(aff.Isp, d.Isp)
	equal_string(aff.Province, d.Province)

	return flag
}

type CacahedAffinityDevs struct {
	Cache gcache.Cache
}

func (cd *CacahedAffinityDevs) GetAffinityDevs(FuncName string, version string, aff Affinity) ([]*Device, error) {
	devs := []*Device{}
	devs_i, err := cd.Cache.Get(FuncName + "_" + version)
	if err != nil {
		return nil, err
	}

	sn2devs := devs_i.(map[string]*Device)

	if len(aff.Sn) > 0 {
		for _, sn := range aff.Sn {
			if d, ok := sn2devs[sn]; ok {
				devs = append(devs, d)
			}
		}
	} else {
		for _, d := range sn2devs {
			if aff.match(d) {
				devs = append(devs, d)
			}
		}
	}

	return devs, nil
}

func NewCacahedAffinityDevs(geter DevsGeter, expire time.Duration) *CacahedAffinityDevs {
	return &CacahedAffinityDevs{
		Cache: gcache.New(20).
			LRU().
			LoaderExpireFunc(func(key interface{}) (interface{}, *time.Duration, error) {
				// get dev and cache by function name and version
				keystr := key.(string)
				name_ver := strings.Split(keystr, "_")
				if len(name_ver) != 2 {
					return nil, nil, fmt.Errorf("%s is not valid", keystr)
				}

				devs, err := geter.GetDevs(name_ver[0], name_ver[1])
				if err != nil {
					return nil, nil, fmt.Errorf("get new cache err: %s", err)
				}

				return devs, &expire, nil
			}).
			Build(),
	}
}

// 缓存有index的数据加速查询
type CacahedIndexedAffinityDevs struct {
	Cache gcache.Cache
}

func NewCacahedIndexedAffinityDevs(geter DevsGeter, expire time.Duration) *CacahedAffinityDevs {
	return &CacahedAffinityDevs{
		Cache: gcache.New(20).
			LRU().
			LoaderExpireFunc(func(key interface{}) (interface{}, *time.Duration, error) {
				// get dev and cache by function name and version
				keystr := key.(string)
				name_ver := strings.Split(keystr, "_")
				if len(name_ver) != 2 {
					return nil, nil, fmt.Errorf("%s is not valid", keystr)
				}

				devs, err := geter.GetDevs(name_ver[0], name_ver[1])
				if err != nil {
					return nil, nil, fmt.Errorf("get new cache err: %s", err)
				}

				db := NewMemDB(10000)
				for _, d := range devs {
					db.Insert(d)
				}

				// build index, example UpdateDB()

				return db, &expire, nil
			}).
			Build(),
	}
}

var deviceDB *MemDB
var sn2Dev map[string]*Device

func GetDeviceMemDB() *MemDB {
	return deviceDB
}

func setNewDBAndMap(new *MemDB, newMap map[string]*Device) {
	// 在已经获得旧的db的协程中继续使用
	// 上面的协程都执行完后，旧db被GC释放
	deviceDB = new
	sn2Dev = newMap
}

func GetSn2DevMap() map[string]*Device {
	return sn2Dev
}

func UpdateDB() {
	// get dev form ipes api or db
	newDB := NewMemDB(1000000)
	sn2dev := map[string]*Device{}

	// test data
	functions := []Function{
		{Name: "func1", Ver: "v1"},
		{Name: "func1", Ver: "v2"},
		{Name: "func2", Ver: "v1"},
	}

	ispName := []string{"ct", "cnc", "cmnet"}
	provienceName := []string{"shangdong", "hebei", "guangdong", "hunan"}
	cityName := []string{"city1", "city2", "city3", "city4", "city5"}

	for i := 0; i < 100; i++ {
		dev := &Device{
			Sn:       fmt.Sprintf("sn_%d", i),
			Isp:      ispName[i%len(ispName)],
			Province: provienceName[i%len(provienceName)],
			City:     cityName[i%len(cityName)],
			Funcs:    functions,
		}
		newDB.Insert(dev)
		sn2dev[dev.Sn] = dev
	}
	// test data end

	// 这里shuffle没有完全起到作用。radix会重建一定的结构，例如search isp.province
	// 那么假设有city1 city2，就会先返回所有city1的设备。所以应该在后面从数据库中取出后shuffle。
	// newDB.Shuffle()
	newDB.CreateIndex(func(v interface{}) ([]string, error) {
		pathSlice := []string{}
		funcNames := []string{}
		d := v.(*Device)

		for _, f := range d.Funcs {
			path := fmt.Sprintf("%s/%s/%s/%s/%s", f.Name, f.Ver, d.Isp, d.Province, d.City)
			pathSlice = append(pathSlice, path)
		}

		for _, f := range d.Funcs {

			funcName := f.Name
			for _, n := range funcNames {
				if n == f.Name {
					funcName = ""
					break
				}
			}

			if funcName != "" {
				funcNames = append(funcNames, funcName)
				path := fmt.Sprintf("%s/%s/%s/%s/%s", f.Name, "all", d.Isp, d.Province, d.City)
				pathSlice = append(pathSlice, path)
			}
		}

		return pathSlice, nil
	})

	setNewDBAndMap(newDB, sn2dev)
}
