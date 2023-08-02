package main

import (
    "context"
    "fmt"
    "log"

    "github.com/genjidb/genji"
    "github.com/genjidb/genji/document"
    "github.com/genjidb/genji/types"
)

func main() {
    // Create a database instance, here we'll store everything on-disk
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
            sn              INT     PRIMARY KEY,
            isp             TEXT,
			province        TEXT,
			city			TEXT,
        )
    `)

    // Create an index
    err = db.Exec("CREATE INDEX user_city_idx ON user (address.city, address.zipCode)")

    // Insert some data
    err = db.Exec("INSERT INTO user (id, name) VALUES (?, ?)", 10, "Foo1", 15)

    // Supported values can go from simple integers to richer data types like lists or documents
    err = db.Exec(`
    INSERT INTO user (id, name, age, address, friends)
    VALUES (
        11,
        'Foo2',
        20,
        {"city": "Lyon", "zipcode": "69001"},
        ["foo", "bar", "baz"]
    )`)

    // Go structures can be passed directly
    type User struct {
        ID              uint
        Name            string
        TheAgeOfTheUser float64 `genji:"age"`
        Address         struct {
            City    string
            ZipCode string
        }
    }

    // Let's create a user
    u := User{
        ID:              20,
        Name:            "foo",
        TheAgeOfTheUser: 40,
    }
    u.Address.City = "Lyon"
    u.Address.ZipCode = "69001"

    err = db.Exec(`INSERT INTO user VALUES ?`, &u)

    // Query some documents
    res, err := db.Query("SELECT id, name, age, address FROM user WHERE age >= ?", 18)
    // always close the result when you're done with it
    defer res.Close()

    // Iterate over the results
    err = res.Iterate(func(d types.Document) error {
        // When querying an explicit list of fields, you can use the Scan function to scan them
        // in order. Note that the types don't have to match exactly the types stored in the table
        // as long as they are compatible.
        var id int
        var name string
        var age int32
        var address struct {
            City    string
            ZipCode string
        }

        err = document.Scan(d, &id, &name, &age, &address)
        if err != nil {
            return err
        }

        fmt.Println(id, name, age, address)

        // It is also possible to scan the results into a structure
        var u User
        err = document.StructScan(d, &u)
        if err != nil {
            return err
        }

        fmt.Println(u)

        // Or scan into a map
        var m map[string]interface{}
        err = document.MapScan(d, &m)
        if err != nil {
            return err
        }

        fmt.Println(m)
        return nil
    })
}