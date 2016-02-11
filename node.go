package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type rawNode struct {
	db smalldDB
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
	p := makePoint(v)
	lbl := fmt.Sprintf("%s", v.Get("label"))

	acc, err := strconv.ParseFloat(v.Get("acc"), 64)
	if err != nil {
		log.Fatal(err)
	}

	n.db.AddLocations(lbl, acc, p)
}

// LocationHandler is the main entry point for smalld
// it receives the get request parses the location data from it
// and logs the values to the location table.
func (n *node) LocationHandler(w http.ResponseWriter, req *http.Request) {
	n.so.Println("handling url", req.URL)

	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	if req.URL.RawQuery == "" {
		http.Error(w, "bad or missing parameters", 422)
		return
	}

	vals, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		// TODO: respond
		return
	}

	log.Println(vals)

	if !safeValues(&vals) {
		// TODO: respond
		return
	}

	p := makePoint(&vals)
	log.Println("point:", p)

	lbl := fmt.Sprintf("%s", vals.Get("label"))

	acc, err := strconv.ParseFloat(vals.Get("acc"), 64)
	if err != nil {
		log.Fatal(err)
	}

	go n.db.AddLocations(lbl, acc, p)

	l, err := n.db.LocationsNameByPoint(p)
	if err != nil {
		// TODO: respond
		return
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

func safeValues(v *url.Values) bool {
	log.Printf("safe %+v", v)

	return true //for now
}

func makePoint(v *url.Values) string {
	p := fmt.Sprintf("POINT(%s %s)", v.Get("lon"), v.Get("lat"))

	return p
}
