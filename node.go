package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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
	if rn == nil {
		return nil, errors.New("rawNode must not be nil")
	}

	if rn.db == nil {
		return nil, errors.New("rawNode db must not be nil")
	}

	if rn.so == nil {
		rn.so = log.New(os.Stdout, "", log.LstdFlags)
	}

	if rn.se == nil {
		rn.se = log.New(os.Stdout, "", log.LstdFlags)
	}

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

	loc, err := locFromVals(vals)
	if err != nil {
		http.Error(w, txt422, 422)
		return
	}

	go func() {
		if err := n.db.AddLocations(loc); err != nil {
			n.se.Println(err)
		}
	}()

	ls, err := n.db.AdminAreasByPoint(loc.Point)
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

func locFromVals(v url.Values) (*location, error) {
	l := fmt.Sprintf("%s", v.Get("label"))

	a, err := strconv.ParseFloat(v.Get("acc"), 64)
	if err != nil {
		return nil, err
	}

	lat, err := strconv.ParseFloat(v.Get("lat"), 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(v.Get("lon"), 64)
	if err != nil {
		return nil, err
	}

	return newLocation(l, a, lat, lon), nil
}
