package main

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type location struct {
	ID       uint64
	Label    sql.NullString
	Acc      float64
	Point    *geoPoint `db:"geom"`
	Received time.Time
}

func newLocation(label string, acc, lat, lon float64) *location {
	return &location{
		Label: sql.NullString{label, true},
		Acc:   acc,
		Point: &geoPoint{lat, lon},
	}
}

type adminArea struct {
	GID        uint64
	OSMID      sql.NullString `db:"osm_id"`
	LastChange sql.NullString `db:"lastchange"`
	Code       sql.NullInt64
	FClass     sql.NullString `db:"fclass"`
	PostalCode sql.NullString `db:"postalcode"`
	Name       sql.NullString
}

type smalldDB interface {
	AddLocations(*location) error
	LocationsNameByPoint(string) ([]string, error)
}

type sDB struct {
	*sqlx.DB
}

func newDB(dbc string) (smalldDB, error) {
	d, err := sqlx.Connect("postgres", dbc)
	if err != nil {
		return nil, err
	}

	d.SetMaxIdleConns(128)

	_, err = d.Exec(`
		SET client_encoding = 'UTF8';
		SET standard_conforming_strings = on;
		SET check_function_bodies = false;
		SET client_min_messages = warning;
		SET search_path = public, pg_catalog;
		SET default_with_oids = false;
	`)
	if err != nil {
		return nil, err
	}

	return &sDB{d}, nil
}

func (sdb *sDB) AddLocations(l *location) error {
	tx, err := sdb.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`insert into locations ( label, acc, geom )
values ( $1, $2, ST_PointFromText( $3, 4326) )`, l.Label, l.Acc, l.Point)
	if err != nil {
		return err
	}

	err = tx.Commit()
	// START TEST location Scan implementation
	if err != nil {
		return err
	}

	nl := &location{}
	if err := sdb.QueryRowx("select * from locations limit 1").StructScan(nl); err != nil {
		return err
	}

	// fmt.Println(nl)

	// END TEST location Scan implementation

	return err
}

func (sdb *sDB) LocationsNameByPoint(point string) ([]string, error) {
	q := `select name from adminareas
where st_contains(adminareas.geom, st_geomfromtext( $1 , 4326))`

	rows, err := sdb.Query(q, point)
	if err != nil {
		return nil, err
	}

	var s []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		s = append(s, name)
	}

	return s, nil
}
