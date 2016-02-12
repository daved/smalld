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

/*func TestLocationHandlerResponseOK(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/location", nil)

	LocationHandler(res, req)
	log.Println("/location response:", res.Code)

	if res.Code != http.StatusNoContent {
		t.Fatalf("Bad Response")
	}
}

func TestLocationHandlerResponseQuery(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost:8000/location?lat=44.09491559960329&lon=-123.0965916720434&acc=5&label=foo", nil)

	LocationHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("Bad Response")
	}
	log.Println(res)
}

func Testrecordlocations(t *testing.T) {
	vals := url.Values{}
	vals.Set("acc", "5")
	vals.Set("lat", "44.09491559960329")
	vals.Set("lon", "-123.0965916720434")
	vals.Set("label", "foo")

	recordlocations(&vals)

	var res Location
	query := "select label, acc, st_y(geom) lat, st_x(geom) lon from locations where label='foo' limit 1"
	err := db.QueryRow(query).Scan(&res.label, &res.acc, &res.lat, &res.lon)
	if err != nil {
		log.Println(err)
		t.Fatalf("could not talk to database")
	}

	if res.label != "foo" {
		t.Fatalf("label inserted does not match")
	}

	if res.acc != 5 {
		t.Fatalf("acc inserted does not match")
	}

	log.Println(res)
}

type Location struct {
	label string
	lat   float64
	lon   float64
	acc   float64
}

func init() {
	log.Println("smalld testing")

	dbc := os.Getenv("SMALLD_DB_CONNECTION")
	//urlBase := os.Getenv("SMALLD_URL_BASE")
	//options := os.Getenv("SMALLD_OPTIONS") //override command line flags

	log.Println("SMALLD_DB_CONNECTION:", dbc)
	//log.Println("SMALLD_URL_BASE:", urlBase)
	//log.Println("SMALLD_OPTIONS:", options)

	var err error
	db, err = sql.Open("postgres", dbc)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected to database")
}*/
