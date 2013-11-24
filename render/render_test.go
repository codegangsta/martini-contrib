package render

import (
	"net/http/httptest"
	"reflect"
	"testing"
)

type Greeting struct {
	One string `json:"one"`
	Two string `json:"two"`
}

func Test_Render_JSON(t *testing.T) {
	res := httptest.NewRecorder()
	r := renderer{res}
	r.JSON(300, Greeting{"hello", "world"})
	expect(t, res.Code, 300)
	expect(t, res.Body.String(), `{"one":"hello","two":"world"}`)
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
