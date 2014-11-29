package resourceful

import "net/http"

// ResourceHandlers define the methods to call for each action on a resource.
type ResourceHandlers interface {
	// List all instances.
	Index(http.ResponseWriter, *http.Request)
	// Create a new instance.
	Create(http.ResponseWriter, *http.Request)
	// Show a specific instance.
	Show(http.ResponseWriter, *http.Request)
	// Update a specific instance.
	Update(http.ResponseWriter, *http.Request)
	// Delete an instance.
	Destroy(http.ResponseWriter, *http.Request)
}

// A basic implementation of ResourceHandlers which responds with a not found
// error on all endpoints. Each action can be overridden independently.
type HandlerFuncs struct {
	IndexFunc   func(http.ResponseWriter, *http.Request)
	CreateFunc  func(http.ResponseWriter, *http.Request)
	ShowFunc    func(http.ResponseWriter, *http.Request)
	UpdateFunc  func(http.ResponseWriter, *http.Request)
	DestroyFunc func(http.ResponseWriter, *http.Request)
}

func (f HandlerFuncs) Index(w http.ResponseWriter, r *http.Request) {
	if f.IndexFunc == nil {
		http.NotFound(w, r)
	} else {
		f.IndexFunc(w, r)
	}
}

func (f HandlerFuncs) Create(w http.ResponseWriter, r *http.Request) {
	if f.CreateFunc == nil {
		http.NotFound(w, r)
	} else {
		f.CreateFunc(w, r)
	}
}

func (f HandlerFuncs) Show(w http.ResponseWriter, r *http.Request) {
	if f.ShowFunc == nil {
		http.NotFound(w, r)
	} else {
		f.ShowFunc(w, r)
	}
}

func (f HandlerFuncs) Update(w http.ResponseWriter, r *http.Request) {
	if f.UpdateFunc == nil {
		http.NotFound(w, r)
	} else {
		f.UpdateFunc(w, r)
	}
}

func (f HandlerFuncs) Destroy(w http.ResponseWriter, r *http.Request) {
	if f.DestroyFunc == nil {
		http.NotFound(w, r)
	} else {
		f.DestroyFunc(w, r)
	}
}
