package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewNode(t *testing.T) {
	var rn *rawNode
	if _, err := newNode(rn); err == nil {
		t.Fatalf("want error, got nil")
	}

	rn = &rawNode{}
	if _, err := newNode(rn); err == nil {
		t.Fatalf("want error, got nil")
	}

	rn.db = &sDB{}
	if _, err := newNode(rn); err != nil {
		t.Fatalf("want nil, got %s", err)
	}
}

func TestNodeMux(t *testing.T) {
	n := &node{}
	n.setMux()
	if n.mux == nil {
		t.Fatalf("want mux, got nil")
	}
}

func TestNodeLogging(t *testing.T) {
	b := &bytes.Buffer{}
	rn := &rawNode{
		db: &sDB{},
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
		t.Fatalf("log is empty")
	}
}

func TestNodeOrigin(t *testing.T) {
	rn := &rawNode{
		db: &sDB{},
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
		t.Fatalf("ACAOrigin header empty")
	}
}

func TestNodeReco(t *testing.T) {
	rn := &rawNode{
		db: &sDB{},
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
