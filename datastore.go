package main

import "database/sql"

type smalldDB interface {
	AddLocations(string, float64, string) error
	LocationsNameByPoint(string) ([]string, error)
}

type sDB struct {
	*sql.DB
}

func newDB(dbc string) (smalldDB, error) {
	d, err := sql.Open("postgres", dbc)
	if err != nil {
		return nil, err
	}

	d.SetMaxIdleConns(128)

	err = d.Ping()
	if err != nil {
		return nil, err
	}

	return &sDB{d}, nil
}

func (sdb *sDB) AddLocations(label string, acc float64, point string) error {
	tx, err := sdb.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`insert into locations ( label, acc, geom )
values ( $1, $2, ST_PointFromText( $3, 4326) )`, label, acc, point)
	if err != nil {
		return err
	}

	err = tx.Commit()

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
