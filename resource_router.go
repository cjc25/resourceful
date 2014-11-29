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

// Create a new ResourceRouter. The router directs RESTful resource actions to
// the handlers defined by the provided interface.
func NewResourceRouter(resourceName string, handlers ResourceHandlers) *ResourceRouter {
	resourcePath := fmt.Sprintf("/%v", resourceName)
	// Disallow empty ids.
	instancePath := fmt.Sprintf("/%v/{%v:.+}", resourceName, resourceIdName(resourceName))

	gr := mux.NewRouter()
	gr.HandleFunc(resourcePath, handlers.Index).Methods("GET")
	gr.HandleFunc(resourcePath, handlers.Create).Methods("POST")
	gr.HandleFunc(instancePath, handlers.Show).Methods("GET")
	gr.HandleFunc(instancePath, handlers.Update).Methods("PUT")
	gr.HandleFunc(instancePath, handlers.Destroy).Methods("DELETE")

	return &ResourceRouter{gr}
}

func (router *ResourceRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO(cjc25): Keep track of the router in a context (the gorilla context?),
	// it will let us know which resources are supposed to be addressed by this
	// router.

	router.internalRouter.ServeHTTP(w, r)
}
