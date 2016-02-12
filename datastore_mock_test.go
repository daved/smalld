package main

type MockSDB struct{}

func (mdb *MockSDB) AddLocations(l *location) error {
	return nil
}

func (mdb *MockSDB) AdminAreasByPoint(point *geoPoint) ([]string, error) {
	s := []string{"Some Place"}

	return s, nil
}
