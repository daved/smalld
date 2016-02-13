package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/laprice/smalld/sdb"
)

func TestNewNode(t *testing.T) {
	var rn *rawNode
	if _, err := newNode(rn); err == nil {
		t.Fatal("want error, got nil")
	}

	rn = &rawNode{}
	if _, err := newNode(rn); err == nil {
		t.Fatal("want error, got nil")
	}

	rn.db = &sdb.MockSDB{}
	if _, err := newNode(rn); err != nil {
		t.Fatalf("want nil, got %s", err)
	}
}

func TestNodeMux(t *testing.T) {
	n := &node{}
	n.setMux()
	if n.mux == nil {
		t.Fatal("want mux set, got nil")
	}
}

func TestNodeLogging(t *testing.T) {
	b := &bytes.Buffer{}
	rn := &rawNode{
		db: &sdb.MockSDB{},
		so: log.New(b, "", 0),
	}

	n, err := newNode(rn)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	n.logging(http.NotFoundHandler()).ServeHTTP(w, r)

	l := b.Len()
	if l == 0 {
		t.Fatal("log is empty")
	}
}

func TestNodeOrigin(t *testing.T) {
	rn := &rawNode{
		db: &sdb.MockSDB{},
	}

	n, err := newNode(rn)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	n.origin(http.NotFoundHandler()).ServeHTTP(w, r)

	oh := w.Header().Get("Access-Control-Allow-Origin")
	if len(oh) == 0 {
		t.Fatal("ACAOrigin header empty")
	}
}

func TestNodeReco(t *testing.T) {
	rn := &rawNode{
		db: &sdb.MockSDB{},
		se: log.New(ioutil.Discard, "", 0),
	}

	n, err := newNode(rn)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	ph := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("but don't")
	})

	n.reco(ph).ServeHTTP(w, r)
}

func TestNotFoundHandler(t *testing.T) {
	rn := &rawNode{
		db: &sdb.MockSDB{},
		se: log.New(ioutil.Discard, "", 0),
	}

	n, err := newNode(rn)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)
	n.NotFoundHandler(w, r)

	want := 404
	if w.Code != want {
		t.Fatalf("want %d, got %d", want, w.Code)
	}
}

func TestMethNAHandler(t *testing.T) {
	rn := &rawNode{
		db: &sdb.MockSDB{},
		se: log.New(ioutil.Discard, "", 0),
	}

	n, err := newNode(rn)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)
	n.MethNAHandler(w, r)

	want := 405
	if w.Code != want {
		t.Fatalf("want %d, got %d", want, w.Code)
	}
}

func TestNodeLocationHandler(t *testing.T) {
	rn := &rawNode{
		db: &sdb.MockSDB{},
		se: log.New(ioutil.Discard, "", 0),
	}

	n, err := newNode(rn)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)
	n.LocationHandler(w, r)

	want := 422
	if w.Code != want {
		t.Fatalf("want %d, got %d", want, w.Code)
	}

	qv := r.URL.Query()
	qv.Set("lat", "44.09491559960329")
	qv.Set("lon", "-123.0965916720434")
	qv.Set("acc", "5")
	r.URL.RawQuery = qv.Encode()

	w = httptest.NewRecorder()
	n.LocationHandler(w, r)

	want = 422
	if w.Code != want {
		t.Fatalf("want %d, got %d", want, w.Code)
	}

	qv.Set("label", "foo")
	r.URL.RawQuery = qv.Encode()

	w = httptest.NewRecorder()
	n.LocationHandler(w, r)

	want = 200
	if w.Code != want {
		t.Fatalf("want %d, got %d", want, w.Code)
	}
}
