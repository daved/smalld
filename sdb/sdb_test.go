package sdb

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestNewLocation(t *testing.T) {
	l := NewLocation("str", 1.1, 2.2, 3.3)
	if l == nil {
		t.Fatalf("want *location, got nil")
	}

	sWant := "str"
	if l.Label.String != sWant {
		t.Fatalf("want %s, got %s", sWant, l.Label.String)
	}

	fWant := 1.1
	if l.Acc != fWant {
		t.Fatalf("want %f, got %f", fWant, l.Acc)
	}

	fWant = 2.2
	if l.Point.Lat != fWant {
		t.Fatalf("want %f, got %f", fWant, l.Point.Lat)
	}

	fWant = 3.3
	if l.Point.Lon != fWant {
		t.Fatalf("want %f, got %f", fWant, l.Point.Lon)
	}
}

func TestNewLocationFromVals(t *testing.T) {
	lv := &LocationVals{}
	l, err := NewLocationFromVals(lv)
	if err == nil {
		t.Fatal("want error, got nil")
	}

	lv.Label = "foo"
	l, err = NewLocationFromVals(lv)
	if err == nil {
		t.Fatal("want error, got nil")
	}

	lv.Acc = "5"
	l, err = NewLocationFromVals(lv)
	if err == nil {
		t.Fatal("want error, got nil")
	}

	lv.Lat = "44.09491559960329"
	l, err = NewLocationFromVals(lv)
	if err == nil {
		t.Fatal("want error, got nil")
	}

	lv.Lon = "-123.0965916720434"
	l, err = NewLocationFromVals(lv)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	if l == nil {
		t.Fatal("want *location, got nil")
	}
}

func TestNewDB(t *testing.T) {
	dd, err := sql.Open("postgres", "")
	if err != nil {
		t.Fatal(err)
	}

	db, err := New(dd)
	if err == nil {
		//t.Fatal("want error, got nil")
	}
	_ = db
}
