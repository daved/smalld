package sdb

import "testing"

func TestMockSDB(t *testing.T) {
	m := &MockSDB{}

	if _, err := m.Migrate(); err != nil {
		t.Fatal(err)
	}

	if _, err := m.RollBack(); err != nil {
		t.Fatal(err)
	}

	if err := m.Populate(); err != nil {
		t.Fatal(err)
	}

	l := &Location{}
	if err := m.AddLocations(l); err != nil {
		t.Fatal(err)
	}

	g := &GeoPoint{}
	if _, err := m.AdminAreasByPoint(g); err != nil {
		t.Fatal(err)
	}
}
