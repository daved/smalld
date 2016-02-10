package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB //to share with our handlers

func safeValues(v *url.Values) bool {
	log.Printf("safe %+v", v)

	return true //for now
}

func makePoint(v *url.Values) string {
	p := fmt.Sprintf("POINT(%s %s)", v.Get("lon"), v.Get("lat"))

	return p
}

func recordlocations(v *url.Values) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	p := makePoint(v)
	lbl := fmt.Sprintf("%s", v.Get("label"))

	acc, err := strconv.ParseFloat(v.Get("acc"), 64)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("insert into locations ( label, acc, geom ) values ( $1, $2, ST_PointFromText( $3, 4326) )", lbl, acc, p)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// LocationHandler is the main entry point for smalld
// it receives the get request parses the location data from it
// and logs the values to the location table.
func LocationHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("handling url", req.URL)

	if req.Method == "GET" {
		if req.URL.RawQuery != "" {
			vals, err := url.ParseQuery(req.URL.RawQuery)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(vals)

			if safeValues(&vals) {
				p := makePoint(&vals)

				go recordlocations(&vals)

				log.Println("point:", p)

				q := "select name from adminareas where st_contains(adminareas.geom, st_geomfromtext( $1 , 4326))"
				rows, err := db.Query(q, p)
				if err != nil {
					log.Print("db error", err)
				}

				var l []string
				for rows.Next() {
					var name string
					rows.Scan(&name)
					l = append(l, name)
				}

				m := make(map[string][]string)
				m["names"] = l

				j, err := json.Marshal(m)
				if err != nil {
					log.Fatal(err)
				}

				h := w.Header()
				h.Add("Access-Control-Allow-Origin", "*")

				w.Write(j)

				return
			}
		} else {
			http.Error(w, "No Content", http.StatusNoContent)

			return
		}
	}
}

func main() {
	log.Println("smalld starting")

	dbc := os.Getenv("SMALLD_DB_CONNECTION")
	//urlBase := os.Getenv("SMALLD_URL_BASE")
	addr := os.Getenv("SMALLD_LISTEN_ADDRESS")
	//options := os.Getenv("SMALLD_OPTIONS") //override command line flags

	log.Println("SMALLD_DB_CONNECTION:", dbc)
	//log.Println("SMALLD_URL_BASE:", urlBase)
	log.Println("SMALLD_LISTEN_ADDRESS", addr)
	//log.Println("SMALLD_OPTIONS:", options)

	var err error
	db, err = sql.Open("postgres", dbc)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to database")

	http.HandleFunc("/location", LocationHandler)
	log.Println("registered LocationHandler")

	http.ListenAndServe(addr, nil)
}
