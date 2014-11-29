# resourceful

Resourceful routing for Go.

Package resourceful provides an implementation of resourceful routing. It
simplifies setting up RESTful endpoints and providing handlers.

For full documentation see the godoc.

# Example

```Go
func NoId(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Handler with no resource id")
}

func WithId(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Handler with resource id", ResourceId("resource_name", r))
}

func main() {
  h := resourceful.HandlerFuncs{
    IndexFunc: NoId,
    CreateFunc: NoId,
    ShowFunc: WithId,
  }
  r := resourceful.NewResourceRouter(“resource_name”, h)
  http.Handle(“/”, r)

  log.Fatal(http.ListenAndServe(":8080", nil))
}
```

# Background

Resourceful routes are a way to represent CRUD-manageable resources in a web
application. They provide a conventional set of endpoints to manage application
models. They are a standard interface for the Ruby on Rails framework.

# Endpoints

The package implements a simple API, ignoring the HTML form endpoints from
Rails-style resourceful routes.

|HTTP Method|URI           |Description
|-----------|--------------|------------------------------------|
|GET        |/resource     |Index: list all instances.          |
|POST       |/resource     |Create: make a new instance.        |
|GET        |/resource/{id}|Show: get a specific instance.      |
|PUT        |/resource/{id}|Update: change a specific instance. |
|DELETE     |/resource/{id}|Destroy: delete a specific instance.|
