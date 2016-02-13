package sdb

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/laprice/smalld/sdb/mig/migdata"
	"github.com/laprice/smalld/sdb/mig/migtote"
	"github.com/rubenv/sql-migrate"
)

type LocationVals struct {
	Label string
	Acc   string
	Lat   string
	Lon   string
}

type Location struct {
	ID       uint64
	Label    sql.NullString
	Acc      float64
	Point    *GeoPoint `db:"geom"`
	Received time.Time
}

func NewLocation(label string, acc, lat, lon float64) *Location {
	return &Location{
		Label: sql.NullString{label, true},
		Acc:   acc,
		Point: &GeoPoint{lat, lon},
	}
}

func NewLocationFromVals(v *LocationVals) (*Location, error) {
	if v.Label == "" {
		return nil, errors.New("label must not be empty")
	}

	a, err := strconv.ParseFloat(v.Acc, 64)
	if err != nil {
		return nil, err
	}

	lat, err := strconv.ParseFloat(v.Lat, 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(v.Lon, 64)
	if err != nil {
		return nil, err
	}

	return NewLocation(v.Label, a, lat, lon), nil
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

type Migrator interface {
	Migrate() (int, error)
	RollBack() (int, error)
	Populate() error
}

type SmalldDB interface {
	Migrator
	AddLocations(*Location) error
	AdminAreasByPoint(*GeoPoint) ([]string, error)
}

type SDB struct {
	*sqlx.DB
	migSrc *migrate.AssetMigrationSource
}

//go:generate go-bindata -pkg migdata -o mig/migdata/mig.go mig/sql

//go:generate tote -in=mig/sqlbig -out=mig/migtote

func New(dd *sql.DB) (SmalldDB, error) {
	if dd == nil {
		return nil, errors.New("db must not be nil")
	}

	dd.SetMaxIdleConns(128)

	_, err := dd.Exec(`
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

	m := &migrate.AssetMigrationSource{
		Asset:    migdata.Asset,
		AssetDir: migdata.AssetDir,
		Dir:      "mig/sql",
	}

	d := sqlx.NewDb(dd, "postgres")

	return &SDB{d, m}, nil
}

func (sdb *SDB) Migrate() (int, error) {
	return sdb.migrates(true)
}

func (sdb *SDB) RollBack() (n int, err error) {
	return sdb.migrates(false)
}

func (sdb *SDB) migrates(up bool) (int, error) {
	dir := migrate.Up
	if !up {
		dir = migrate.Down
	}

	n, err := migrate.Exec(sdb.DB.DB, "postgres", sdb.migSrc, dir)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (sdb *SDB) Populate() error {
	_, err := sdb.Exec(migtote.Root.AddAdminareasData)

	return err
}

func (sdb *SDB) AddLocations(l *Location) error {
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
	/* // START TEST location Scan implementation
	if err != nil {
		return err
	}

	nl := &location{}
	if err := sdb.QueryRowx("select * from locations limit 1").StructScan(nl); err != nil {
		return err
	}

	fmt.Println(nl)

	// END TEST location Scan implementation */

	return err
}

func (sdb *SDB) AdminAreasByPoint(point *GeoPoint) ([]string, error) {
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
