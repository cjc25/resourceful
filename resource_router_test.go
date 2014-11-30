package resourceful

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router *ResourceRouter = NewResourceRouter().AddResource("prefix", HandlerFuncs{
	IndexFunc: func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Index function")
	},
	ShowFunc: func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Show function for %v", ResourceId("prefix", r))
	},
})

func TestRouting_NonInstanceRoute(t *testing.T) {
	w := httptest.NewRecorder()
	// Index
	r, err := http.NewRequest("GET", "/prefix", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, r)
	response := w.Body.String()
	if response != "Index function" {
		t.Errorf("Expected response 'Index function' but got '%v'", response)
	}
}

func TestRouting_UndefinedRoute(t *testing.T) {
	w := httptest.NewRecorder()
	// Create
	r, err := http.NewRequest("POST", "/prefix", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, r)
	if w.Code != 404 {
		t.Errorf("Expected code 404 but got %v", w.Code)
	}
}

func TestResourceId_NotPresent(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/prefix/", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, r)
	if w.Code != 404 {
		t.Errorf("Expected 404, got %v", w.Code)
	}
}

func TestResourceId_Present(t *testing.T) {
	oldKeepContext := router.internalRouter.KeepContext
	router.internalRouter.KeepContext = true
	defer func() { router.internalRouter.KeepContext = oldKeepContext }()

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/prefix/prefix_thing", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, r)
	got := ResourceId("prefix", r)
	if got != "prefix_thing" {
		t.Errorf("Expected resource id 'prefix_thing', got '%v'", got)
	}
}

func ExampleResourceRouter() {
	router := NewResourceRouter()
	router.AddResource("prefix", HandlerFuncs{
		ShowFunc: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Show function for %v", ResourceId("prefix", r))
		},
	})
	router.AddResource("resource_2", HandlerFuncs{
		IndexFunc: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Index for resource_2")
		},
	})

	w := httptest.NewRecorder()
	// Show
	r, err := http.NewRequest("GET", "/prefix/42a", nil)
	if err != nil {
		log.Fatal(err)
	}

	router.ServeHTTP(w, r)
	fmt.Println(w.Body.String())

	w = httptest.NewRecorder()
	// Index
	r, err = http.NewRequest("GET", "/resource_2", nil)
	if err != nil {
		log.Fatal(err)
	}

	router.ServeHTTP(w, r)
	fmt.Println(w.Body.String())

	// Output:
	// Show function for 42a
	// Index for resource_2
}
