package main

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestNewOptions(t *testing.T) {
	o := newOptions()
	if o == nil {
		t.Fatalf("want *options, got nil")
	}
}

func TestOptionsValidate(t *testing.T) {
	o := newOptions()
	if err := o.validate(); err == nil {
		t.Fatalf("want error, got nil")
	}

	o.dbc = "x"
	if err := o.validate(); err == nil {
		t.Fatalf("want error, got nil")
	}

	o.addr = "x"
	if err := o.validate(); err != nil {
		t.Fatalf("want nil, got %s", err)
	}
}
