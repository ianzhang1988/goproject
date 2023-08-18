package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/hashicorp/go-memdb"
)

type Box struct {
	SN       string
	Isp      string
	Province string
	City     string
}

func main() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("----- begin ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"box": &memdb.TableSchema{
				Name: "box",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "SN"},
					},
					"isp": &memdb.IndexSchema{
						Name:    "isp",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Isp"},
					},
					"province": &memdb.IndexSchema{
						Name:    "province",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Province"},
					},
					"city": &memdb.IndexSchema{
						Name:    "city",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "City"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	fmt.Println("start insert")
	start := time.Now()
	// Create a write transaction
	txn := db.Txn(true)
	ispList := []string{"cmnet", "ct", "cnc"}
	for i := 0; i < 1000000; i++ {
		b := Box{
			SN:       fmt.Sprintf("my_sn_%d", i),
			Isp:      ispList[i%len(ispList)],
			Province: "shangdong",
			City:     "weifang",
		}

		if err := txn.Insert("box", &b); err != nil {
			panic(err)
		}
	}

	// Commit the transaction
	txn.Commit()

	fmt.Println("time used:", time.Since(start))

	txn = db.Txn(false)
	// defer txn.Abort()

	// Lookup by email
	raw, err := txn.First("box", "id", "my_sn_10")
	if err != nil {
		panic(err)
	}

	// Say hi!
	fmt.Printf("Hello %v!\n", raw.(*Box))

	raw, err = txn.First("box", "isp", "ct")
	if err != nil {
		panic(err)
	}

	fmt.Printf("isp %v!\n", raw.(*Box))

	txn.Abort()

	runtime.ReadMemStats(&m)
	fmt.Println("----- after insert ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	runtime.GC()

	runtime.ReadMemStats(&m)
	fmt.Println("----- after GC ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	txn = db.Txn(true)

	txn.DeleteAll("box", "isp", "cmnet")
	txn.DeleteAll("box", "isp", "cnc")
	txn.Commit()

	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Println("----- after Delete ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	txn = db.Txn(false)

	raw, err = txn.First("box", "isp", "cnc")
	if err != nil || raw == nil {
		fmt.Println("cnc no more", err)
	} else {
		fmt.Printf("Hello %v!\n", raw.(*Box))
	}

	raw, err = txn.First("box", "isp", "ct")
	if err != nil || raw == nil {
		fmt.Println("ct no more", err)
	} else {
		fmt.Printf("Hello %v!\n", raw.(*Box))
	}

	txn.Abort()
}
