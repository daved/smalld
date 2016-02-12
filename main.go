package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type options struct {
	dbc    string
	dbcFN  string
	dbcEV  string
	addr   string
	addrFN string
	addrEV string
	// os.Getenv("SMALLD_URL_BASE")
	// os.Getenv("SMALLD_OPTIONS")
}

func newOptions() *options {
	return &options{
		dbcFN:  "dbconf",
		dbcEV:  "SMALLD_DB_CONNECTION",
		addrFN: "port",
		addrEV: "SMALLD_LISTEN_ADDRESS",
	}
}

func (o *options) validate() error {
	if o.dbc == "" {
		return fmt.Errorf(
			"database configuration must be set using flag %s or env var %s",
			o.dbcFN,
			o.dbcEV,
		)
	}

	if o.addr == "" {
		return fmt.Errorf(
			"http listen port must be set using flag %s or env var %s",
			o.addrFN,
			o.addrEV,
		)
	}

	return nil
}

func main() {
	so := log.New(os.Stdout, "", log.LstdFlags)
	se := log.New(os.Stderr, "", log.LstdFlags)

	so.Println("smalld starting")

	o := newOptions()

	flag.StringVar(&o.dbc, o.dbcFN, os.Getenv(o.dbcEV),
		"database configuration (postgres)")
	flag.StringVar(&o.addr, o.addrFN, os.Getenv(o.addrEV),
		"port to listen for http requests")
	flag.Parse()

	if err := o.validate(); err != nil {
		se.Fatalln(err)
	}

	so.Println("connecting to database")

	db, err := newDB(o.dbc)
	if err != nil {
		se.Fatalln(err)
	}

	rn := &rawNode{
		db: db,
		so: so,
		se: se,
	}

	n, err := NewNode(rn)
	if err != nil {
		se.Fatalln(err)
	}

	so.Printf("serving on %s", o.addr)
	http.ListenAndServe(o.addr, n)
}
