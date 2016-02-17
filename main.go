package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/laprice/smalld/sdb"
	_ "github.com/lib/pq"
)

type options struct {
	dbc      string
	dbcFlag  string
	dbcEnv   string
	addr     string
	addrFlag string
	addrEnv  string
	// os.Getenv("SMALLD_URL_BASE")
	// os.Getenv("SMALLD_OPTIONS")
}

func newOptions() *options {
	return &options{
		dbcFlag:  "dbconf",
		dbcEnv:   "SMALLD_DB_CONNECTION",
		addrFlag: "port",
		addrEnv:  "SMALLD_LISTEN_ADDRESS",
	}
}

func (o *options) validate() error {
	if o.dbc == "" {
		return fmt.Errorf(
			"database configuration must be set using flag %s or env var %s",
			o.dbcFlag,
			o.dbcEnv,
		)
	}

	if o.addr == "" {
		return fmt.Errorf(
			"http listen port must be set using flag %s or env var %s",
			o.addrFlag,
			o.addrEnv,
		)
	}

	return nil
}

//go:generate go generate github.com/laprice/smalld/sdb

func main() {
	so := log.New(os.Stdout, "", log.LstdFlags)
	se := log.New(os.Stderr, "", log.LstdFlags)

	so.Println("smalld starting")

	o := newOptions()

	flag.StringVar(&o.dbc, o.dbcFlag, os.Getenv(o.dbcEnv),
		"database configuration (postgres)")
	flag.StringVar(&o.addr, o.addrFlag, os.Getenv(o.addrEnv),
		"port to listen for http requests")
	dbRB := flag.Bool("db-rollback", false,
		`Rollback all database migrations.`)
	dbPop := flag.Bool("db-populate", false,
		`Install all database fixtures.`)

	flag.Parse()

	if err := o.validate(); err != nil {
		se.Fatalln(err)
	}

	so.Println("connecting to database")

	p, err := sql.Open("postgres", o.dbc)
	if err != nil {
		se.Fatalln(err)
	}

	db, err := sdb.New(p)
	if err != nil {
		se.Fatalln(err)
	}

	ct, err := db.Migrate()
	so.Printf("migrated database: %d file(s).\n", ct)
	if err != nil {
		se.Fatalln(err)
	}

	if *dbRB {
		so.Println("rolling back database and exiting")
		ct, err := db.RollBack()
		so.Printf("rolled back database: %d file(s).\n", ct)
		if err != nil {
			se.Fatalln(err)
		}
		os.Exit(0)
	}

	if *dbPop {
		so.Println("installing database fixtures")
		if err = db.Populate(); err != nil {
			se.Fatalln(err)
		}
	}

	rn := &rawNode{
		db: db,
		so: so,
		se: se,
	}

	n, err := newNode(rn)
	if err != nil {
		se.Fatalln(err)
	}

	so.Printf("serving on %s", o.addr)
	http.ListenAndServe(o.addr, n)
}
