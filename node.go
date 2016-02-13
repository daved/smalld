package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/codemodus/catena"
	"github.com/codemodus/mixmux"
	"github.com/laprice/smalld/sdb"
)

const (
	httpStatusUnprocessableEntity = 422
)

var (
	httpStatusText = map[int]string{
		422: "Unprocessable Entity",
	}
)

func HttpStatusText(code int) string {
	return httpStatusText[code]
}

type rawNode struct {
	db sdb.SmalldDB
	so *log.Logger
	se *log.Logger
}

type node struct {
	*rawNode
	mux *mixmux.Router
}

func newNode(rn *rawNode) (*node, error) {
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

	n.setMux()

	return n, nil
}

func (n *node) setMux() {
	opts := &mixmux.Options{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		MethodNotAllowed:       http.HandlerFunc(n.MethNAHandler),
		NotFound:               http.HandlerFunc(n.NotFoundHandler),
	}
	n.mux = mixmux.NewRouter(opts)

	c := catena.New(n.reco, n.logging, n.origin)

	n.mux.Get("/location", c.EndFn(n.LocationHandler))
}

func (n *node) reco(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				n.se.Printf("panic: %+v\n", err)

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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

func (n *node) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (n *node) MethNAHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// LocationHandler is the main entry point for smalld
// it receives the get request parses the location data from it
// and logs the values to the location table.
func (n *node) LocationHandler(w http.ResponseWriter, r *http.Request) {
	qv := r.URL.Query()
	if len(qv) == 0 {
		http.Error(w, HttpStatusText(httpStatusUnprocessableEntity), httpStatusUnprocessableEntity)
		return
	}

	locVals := &sdb.LocationVals{
		Label: qv.Get("label"),
		Acc:   qv.Get("acc"),
		Lat:   qv.Get("lat"),
		Lon:   qv.Get("lon"),
	}

	loc, err := sdb.NewLocationFromVals(locVals)
	if err != nil {
		http.Error(w, HttpStatusText(httpStatusUnprocessableEntity), httpStatusUnprocessableEntity)
		return
	}

	go func() {
		if err := n.db.AddLocations(loc); err != nil {
			n.se.Println(err)
		}
	}()

	ls, err := n.db.AdminAreasByPoint(loc.Point)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	m := make(map[string][]string)
	m["names"] = ls

	b, err := json.Marshal(m)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
