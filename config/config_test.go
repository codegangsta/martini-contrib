package config

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type confgBodyTest struct {
	path        string
	config_file string
	expected    Body
}

var confgBodyTests = []confgBodyTest{
	// Test empty file path
	{"/none", "", make(Body)},

	// Test an empty file
	{"/empty", "test_examples/empty.json", make(Body)},

	// Test a file with values
	{"/correct", "test_examples/correct.json", Body{"Name": "Platypus", "Order": "Monotremata"}},
}

func TestConfigBodyTests(t *testing.T) {
	for _, test := range confgBodyTests {
		m := martini.Classic()
		m.Use(File(test.config_file))

		m.Get(test.path, func(res Body) {
			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("\nExpected: %#v\nResult: %#v", test.expected, res)
			}
		})

		recorder := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", test.path, nil)

		m.ServeHTTP(recorder, r)
	}
}

func BenchmarkWitoutConfig(b *testing.B) {
	m := newBenchmarkMartini("")

	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/benchmark", nil)

	for n := 0; n < b.N; n++ {
		m.ServeHTTP(recorder, r)
	}
}

func BenchmarkWithConfig(b *testing.B) {
	m := newBenchmarkMartini("test_examples/correct.json")

	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/benchmark", nil)

	for n := 0; n < b.N; n++ {
		m.ServeHTTP(recorder, r)
	}
}

func newBenchmarkMartini(config_file string) *martini.ClassicMartini {
	router := martini.NewRouter()
	base := martini.New()
	base.Action(router.Handle)

	m := &martini.ClassicMartini{base, router}
	if config_file != "" {
		m.Use(File(config_file))
		m.Get("/benchmark", func(conf Body) {})
	} else {
		m.Get("/benchmark", func() {})
	}

	return m
}
