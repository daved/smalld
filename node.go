package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type rawNode struct {
	db *sql.DB
	so *log.Logger
}

type node struct {
	*rawNode
	mux *http.ServeMux
}

func NewNode(rn *rawNode) (*node, error) {
	// TODO: validate s.

	n := &node{
		rawNode: rn,
	}

	if err := n.setMux(); err != nil {
		return nil, err
	}

	return n, nil
}

func (n *node) setMux() error {
	n.mux = http.NewServeMux()

	n.mux.HandleFunc("/location", n.LocationHandler)

	return nil
}

func (n *node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.mux.ServeHTTP(w, r)
}

func (n *node) recordlocations(v *url.Values) {
	tx, err := n.db.Begin()
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
func (n *node) LocationHandler(w http.ResponseWriter, req *http.Request) {
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

				go n.recordlocations(&vals)

				log.Println("point:", p)

				q := "select name from adminareas where st_contains(adminareas.geom, st_geomfromtext( $1 , 4326))"
				rows, err := n.db.Query(q, p)
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

func safeValues(v *url.Values) bool {
	log.Printf("safe %+v", v)

	return true //for now
}

func makePoint(v *url.Values) string {
	p := fmt.Sprintf("POINT(%s %s)", v.Get("lon"), v.Get("lat"))

	return p
}
