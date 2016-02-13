package sdb

type MockSDB struct{}

func (mdb *MockSDB) Migrate() (int, error) {
	return 0, nil
}

func (mdb *MockSDB) RollBack() (n int, err error) {
	return 0, nil
}

func (mdb *MockSDB) AddLocations(l *Location) error {
	return nil
}

func (mdb *MockSDB) AdminAreasByPoint(point *GeoPoint) ([]string, error) {
	s := []string{"Some Place"}

	return s, nil
}
