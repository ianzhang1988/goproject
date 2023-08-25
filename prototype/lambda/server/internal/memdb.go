package lambda_service

import (
	"math/rand"

	"github.com/armon/go-radix"
)

type ParsePath func(value interface{}) ([]string, error)

type MemDB struct {
	Trie *radix.Tree
	Data []interface{}
}

func NewMemDB(size uint32) *MemDB {
	return &MemDB{
		Trie: radix.New(),
		Data: make([]interface{}, 0, size),
	}
}

func (db *MemDB) Insert(value interface{}) {
	db.Data = append(db.Data, value)
}

func (db *MemDB) Shuffle() {
	rand.Shuffle(len(db.Data), func(i, j int) {
		db.Data[i], db.Data[j] = db.Data[j], db.Data[i]
	})
}

func (db *MemDB) CreateIndex(f ParsePath) error {
	for i, v := range db.Data {
		pathList, err := f(v)
		if err != nil {
			return err
		}

		for _, path := range pathList {
			idxSliceI, ok := db.Trie.Get(path)
			if !ok {
				// fmt.Println(path, "no found")
				newSlice := &[]uint32{}
				db.Trie.Insert(path, newSlice)
				// 下面的写法加的实际上不是指针
				// newSlice := []uint32{}
				// db.Trie.Insert(path, &newSlice)
				idxSliceI = newSlice
			}

			idxSlice := idxSliceI.(*([]uint32))
			//fmt.Println("create index get", *idxSlice)
			*idxSlice = append(*idxSlice, uint32(i))
			// fmt.Println("create index", *idxSlice)
		}
	}

	return nil
}

func (db *MemDB) FindPrefixIdx(prefix string) []uint32 {
	idxSlice := []uint32{}
	db.Trie.WalkPrefix(prefix, func(s string, v interface{}) bool {
		idxSlice = append(idxSlice, (v.([]uint32))...)
		return false
	})

	return idxSlice
}

func (db *MemDB) FindPrefix(prefix string) []interface{} {
	interfaceSlice := []interface{}{}
	db.Trie.WalkPrefix(prefix, func(s string, v interface{}) bool {
		idxSlice := v.(*([]uint32))
		for _, idx := range *idxSlice {
			interfaceSlice = append(interfaceSlice, db.Data[idx])
		}
		return false
	})

	return interfaceSlice
}
