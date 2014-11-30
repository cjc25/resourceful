// Package resourceful provides an implementation of resourceful routing. It
// simplifies setting up RESTful endpoints and providing handlers.
//
// Background
//
// Resourceful routes are a way to represent CRUD-manageable resources in a web
// application. They provide a conventional set of endpoints to manage
// application models. They are a standard interface for the Ruby on Rails
// framework.
//
// Endpoints
//
// The package implements a simple API, ignoring the HTML form endpoints from
// Rails-style resourceful routes.
//
//   HTTP Method  URI             Description
//   GET          /resource       Index: list all instances.
//   POST         /resource       Create: make a new instance.
//   GET          /resource/{id}  Show: get a specific instance.
//   PUT          /resource/{id}  Update: change a specific instance.
//   DELETE       /resource/{id}  Destroy: delete a specific instance.
//
// Handlers
//
// The package provides a ResourceHandlers implementation in a struct
// HandlerFuncs which responds with a 404 on all endpoints. The endpoints can
// be overridden individually.
//
// If the application implements ResourceHandlers, only a single-process
// architecture can hold complex state in memory between calls.  Since scalable
// web applications will have to shard work across multiple processes and
// implement cheap recovery when a process goes down, this is discouraged.
// Instead, maintain only basic state like a database connection and load
// per-request data separately on each request from a memcache or persistent
// store.
//
// Routers
//
// A router is created by NewResourceRouter, and uses gorilla/mux to do URI and
// HTTP method matching. Resources are added with AddResource, which takes a
// resource name and ResourceHandlers to build the routes. For endpoints that
// use a specific instance, the resource ids for a request can be extracted
// within a handler with ResourceId(resource_name, request).
package resourceful

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func resourceIdName(resourceName string) string {
	return fmt.Sprintf("_resourceful_%v_id", resourceName)
}

// Get the id of the named resource for the given request. If the resource was
// not present on the requested route, the empty string is returned.
func ResourceId(resourceName string, r *http.Request) string {
	v := mux.Vars(r)
	return v[resourceIdName(resourceName)]
}

// An http.Handler which serves resourceful routes, created with
// NewResourceRouter.
type ResourceRouter struct {
	internalRouter *mux.Router
}

// Create a new ResourceRouter.
func NewResourceRouter() *ResourceRouter {
	return &ResourceRouter{mux.NewRouter()}
}

// Add a resource to a router. The router directs RESTful resource actions to
// the handlers defined by the provided interface. This returns the modified
// router to allow chained calls.
func (router *ResourceRouter) AddResource(resourceName string, handlers ResourceHandlers) *ResourceRouter {
	resourcePath := fmt.Sprintf("/%v", resourceName)
	// Disallow empty ids.
	instancePath := fmt.Sprintf("/%v/{%v:.+}", resourceName, resourceIdName(resourceName))

	router.internalRouter.HandleFunc(resourcePath, handlers.Index).Methods("GET")
	router.internalRouter.HandleFunc(resourcePath, handlers.Create).Methods("POST")
	router.internalRouter.HandleFunc(instancePath, handlers.Show).Methods("GET")
	router.internalRouter.HandleFunc(instancePath, handlers.Update).Methods("PUT")
	router.internalRouter.HandleFunc(instancePath, handlers.Destroy).Methods("DELETE")

	return router
}

func (router *ResourceRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO(cjc25): Keep track of the router in a context (the gorilla context?),
	// it will let us know which resources are supposed to be addressed by this
	// router.
	router.internalRouter.ServeHTTP(w, r)
}
