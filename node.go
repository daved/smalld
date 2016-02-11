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
	se *log.Logger
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

// LocationHandler is the main entry point for smalld
// it receives the get request parses the location data from it
// and logs the values to the location table.
func (n *node) LocationHandler(w http.ResponseWriter, req *http.Request) {
	n.so.Println("handling url", req.URL)

	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	txt422 := "bad or missing parameters"
	if req.URL.RawQuery == "" {
		http.Error(w, txt422, 422)
		return
	}

	vals, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		http.Error(w, txt422, 422)
		return
	}

	l, a, p, err := labelAccPoint(vals)
	if err != nil {
		http.Error(w, txt422, 422)
		return
	}

	fmt.Println(l, a, p)
	go n.db.AddLocations(l, a, p)

	ls, err := n.db.LocationsNameByPoint(p)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	m := make(map[string][]string)
	m["names"] = ls

	j, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	h := w.Header()
	h.Add("Access-Control-Allow-Origin", "*")

	w.Write(j)
}

func labelAccPoint(v url.Values) (string, float64, string, error) {
	l := fmt.Sprintf("%s", v.Get("label"))

	a, err := strconv.ParseFloat(v.Get("acc"), 64)
	if err != nil {
		return "", 0, "", err
	}

	p := fmt.Sprintf("POINT(%s %s)", v.Get("lon"), v.Get("lat"))

	return l, a, p, nil
}
