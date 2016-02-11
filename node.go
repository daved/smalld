package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/codemodus/catena"
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
	c := catena.New(n.reco, n.logging, n.origin)

	n.mux = http.NewServeMux()

	n.mux.Handle("/location", c.EndFn(n.LocationHandler))

	return nil
}

func (n *node) reco(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				n.se.Printf("panic: %+v\n", err)

				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (n *node) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n.so.Printf("handling url %s\n", r.URL)

		next.ServeHTTP(w, r)
	})
}

func (n *node) origin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func (n *node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.mux.ServeHTTP(w, r)
}

// LocationHandler is the main entry point for smalld
// it receives the get request parses the location data from it
// and logs the values to the location table.
func (n *node) LocationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	txt422 := "bad or missing parameters"
	if r.URL.RawQuery == "" {
		http.Error(w, txt422, 422)
		return
	}

	vals, err := url.ParseQuery(r.URL.RawQuery)
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

	b, err := json.Marshal(m)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Write(b)
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
