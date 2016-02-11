package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	so := log.New(os.Stdout, "", log.LstdFlags)
	so.Println("smalld starting")

	dbc := os.Getenv("SMALLD_DB_CONNECTION")
	addr := os.Getenv("SMALLD_LISTEN_ADDRESS")
	//urlBase := os.Getenv("SMALLD_URL_BASE")
	//options := os.Getenv("SMALLD_OPTIONS") //override command line flags

	so.Println("connecting to database")
	db, err := newDB(dbc)
	if err != nil {
		so.Fatalln(err)
	}

	rn := &rawNode{
		db: db,
		so: so,
	}

	n, err := NewNode(rn)
	if err != nil {
		so.Fatalln(err)
	}

	so.Printf("serving on %s", addr)
	http.ListenAndServe(addr, n)
}
