package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/types"
)

func main() {
	// Create a database instance, here we'll store everything on-disk
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("----- begin ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	db, err := genji.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// If needed, attach context, e.g. (*http.Request).Context().
	db = db.WithContext(context.Background())

	// Create a table with a strict schema.
	// Useful to have full control of the table content.
	// Notice that it is possible to define constraint on nested documents.
	err = db.Exec(`
        CREATE TABLE box (
            sn              TEXT     PRIMARY KEY,
            isp             TEXT,
			province        TEXT,
			city			TEXT,
            lambda (
                name        TEXT
            )
        )
    `)
	if err != nil {
		fmt.Println("create err:", err)
		return
	}

	// Create an index
	err = db.Exec("CREATE INDEX sn_index ON box (sn)")
	if err != nil {
		fmt.Println("index err:", err)
		return
	}
	// err = db.Exec("CREATE INDEX isp_index ON box (isp)")
	// if err != nil {
	// 	fmt.Println("index err:", err)
	// 	return
	// }
	// err = db.Exec("CREATE INDEX province_index ON box (province)")
	// if err != nil {
	// 	fmt.Println("index err:", err)
	// 	return
	// }
	// err = db.Exec("CREATE INDEX city_index ON box (city)")
	// if err != nil {
	// 	fmt.Println("index err:", err)
	// 	return
	// }
	// err = db.Exec("CREATE INDEX lambda_name_index ON box (lambda.name)")
	// if err != nil {
	// 	fmt.Println("index err:", err)
	// 	return
	// }

	err = db.Exec("CREATE INDEX filter_index ON box (lambda.name, isp, province, city)")
	if err != nil {
		fmt.Println("index err:", err)
		return
	}

	err = db.Exec("INSERT INTO box (sn, isp) VALUES (?, ?)", "sn1", "ct")
	if err != nil {
		fmt.Println("insert err:", err)
		return
	}

	// Go structures can be passed directly
	type Box struct {
		SN       string
		Isp      string
		Province string
		City     string
		Lambda   struct {
			Name string
		}
	}

	// Let's create a user
	b := Box{
		SN:  "sn2",
		Isp: "cmnet",
	}
	b.Lambda.Name = "haha"

	err = db.Exec(`INSERT INTO box VALUES ?`, &b)
	if err != nil {
		fmt.Println("insert struct err:", err)
		return
	}

	// Query some documents
	res, err := db.Query("SELECT * FROM box")
	if err != nil {
		fmt.Println("select err:", err)
		return
	}
	// always close the result when you're done with it
	defer res.Close()

	// Iterate over the results
	err = res.Iterate(func(d types.Document) error {
		// When querying an explicit list of fields, you can use the Scan function to scan them
		// in order. Note that the types don't have to match exactly the types stored in the table
		// as long as they are compatible.

		// var id int
		// var name string
		// var age int32
		// var address struct {
		// 	City    string
		// 	ZipCode string
		// }

		// err = document.Scan(d, &id, &name, &age, &address)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(id, name, age, address)

		// It is also possible to scan the results into a structure
		var b Box
		err = document.StructScan(d, &b)
		if err != nil {
			return err
		}

		fmt.Println(b)

		// Or scan into a map
		var m map[string]interface{}
		err = document.MapScan(d, &m)
		if err != nil {
			return err
		}

		fmt.Println(m)
		return nil
	})

	if err != nil {
		fmt.Println("iterate err:", err)
		return
	}

	runtime.ReadMemStats(&m)
	fmt.Println("----- some opration ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	fmt.Println("start insert")
	ispList := []string{"cmnet", "ct", "cnc"}
	FuncList := []string{"func1", "func2", "func3"}
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		b := Box{
			SN:       fmt.Sprintf("my_sn_%d", i),
			Isp:      ispList[i%len(ispList)],
			Province: "shangdong",
			City:     "weifang",
		}
		b.Lambda.Name = FuncList[i%len(FuncList)]

		// err = db.Exec(`INSERT INTO box (sn, isp, province, city, lambda) VALUES (?,?,?,?,?)`, b.SN, b.Isp, b.Province, b.City, b.Lambda)
		err = db.Exec(`INSERT INTO box VALUES ?`, &b)
		if err != nil {
			fmt.Println("insert struct err:", err)
			return
		}
	}
	fmt.Println("time used:", time.Since(start))

	runtime.ReadMemStats(&m)
	fmt.Println("----- after insert ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	// fmt.Println("start insert map")
	// boxMap := map[string]*Box{}
	// start := time.Now()
	// for i := 0; i < 1000000; i++ {
	// 	b := Box{
	// 		SN:       fmt.Sprintf("my_sn_%d", i),
	// 		Isp:      "cmnet",
	// 		Province: "shangdong",
	// 		City:     "weifang",
	// 	}

	// 	boxMap[b.SN] = &b
	// }
	// fmt.Println("time used:", time.Since(start))

	// runtime.ReadMemStats(&m)
	// fmt.Println("----- after insert ------")
	// fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	// fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	// fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	// fmt.Printf("NumGC = %v\n", m.NumGC)

	start = time.Now()
	res, err = db.Query("SELECT * FROM box where isp = 'ct'")
	if err != nil {
		fmt.Println("select err:", err)
		return
	}
	// always close the result when you're done with it
	defer res.Close()

	counter := 0

	// Iterate over the results
	err = res.Iterate(func(d types.Document) error {

		err = document.StructScan(d, &b)
		if err != nil {
			return err
		}

		if counter == 1 {
			fmt.Println("box:", b)
		}

		counter += 1

		return nil
	})

	fmt.Println("counter:", counter)
	fmt.Println("time used:", time.Since(start))
}
