package lambda_service

import "fmt"

type Function struct {
	Name string
	Ver  string
}

type Device struct {
	Sn       string
	Isp      string
	Province string
	City     string
	Funcs    []Function
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
	// get dev form ipes api
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
	cityName := []string{"city1", "city2", "city3"}

	for i := 0; i < 1000000; i++ {
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

	newDB.Shuffle()
	newDB.CreateIndex(func(v interface{}) ([]string, error) {
		pathSlice := []string{}
		d := v.(*Device)
		for _, f := range d.Funcs {
			path := fmt.Sprintf("%s/%s/%s/%s/%s", f.Name, f.Ver, d.Isp, d.Province, d.City)
			pathSlice = append(pathSlice, path)
		}

		return pathSlice, nil
	})

	setNewDBAndMap(newDB, sn2dev)
}
