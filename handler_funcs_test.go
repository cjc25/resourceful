package resourceful

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var actions = [...]string{"Index", "Create", "Show", "Update", "Destroy"}

func handlerFor(r ResourceHandlers, action string) func(http.ResponseWriter, *http.Request) {
	f := reflect.ValueOf(r).MethodByName(action)
	return f.Interface().(func(http.ResponseWriter, *http.Request))
}

func expectCode(name string, handler func(http.ResponseWriter, *http.Request), code int, t *testing.T) {
	w := httptest.NewRecorder()
	handler(w, nil)
	if w.Code != code {
		t.Errorf("Expected %v to return %v, got %v", name, code, w.Code)
	}
}

func TestDefault_All404(t *testing.T) {
	hf := HandlerFuncs{}
	for _, action := range actions {
		m := handlerFor(hf, action)
		expectCode(action, m, 404, t)
	}
}

func (hf *HandlerFuncs) setHandlers(overrides []string, h func(http.ResponseWriter, *http.Request)) {
	v := reflect.ValueOf(hf).Elem()
	hValue := reflect.ValueOf(h)
	for _, override := range overrides {
		f := v.FieldByName(override + "Func")
		f.Set(hValue)
	}
}

func TestOverrides(t *testing.T) {
	sentinel := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Failed differently", 401)
	}

	// Loop over every subset of actions and override those actions.
	for i := range actions {
		for j := i; j < len(actions); j++ {
			hf := &HandlerFuncs{}
			hf.setHandlers(actions[i:j+1], sentinel)

			for k, action := range actions {
				t.Logf("Checking override %v (%v), where %v to %v should return 401",
					k, action, i, j)

				if k < i || k > j {
					// This should not have been overridden.
					expectCode(action, handlerFor(hf, action), 404, t)
				} else {
					// This should have been overridden.
					expectCode(action, handlerFor(hf, action), 401, t)
				}
			}
		}
	}
}
