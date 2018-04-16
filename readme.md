# FUEL UP
REST & MVC framework

## Inspriration
- Popular WebServers (like Apache,IIS) for hierarchical configuration model
- Popular MVC frameworks for service controller based routing

## Design Goals
- Simplicity
- High Developer Productivity
- Easy and High Configurability
- Low Learning Curve
- Simple Versioning
- High Performance
- Preference for JSON (over XML)

## Features
- Controllers to define endpoints (modular, structured and organized codebase)
- [Powerful configuration model](./docs/configuration.md)
- [Routing](./docs/routing.md)
- [Parameters, Query String and Context](./docs/params.md)
- Versioning
- [Caching](./docs/caching.md)
- [Stubbing](./docs/stub.md)
- [Middleware support](./docs/middleware.md)
- [Logging](./docs/logging.md)
- ~~[MVC](./docs/mvc.md)~~
- ~~Proxying~~
- ~~CRUD~~

## Hello World

Lets see how we can quickly write a Hello World Api

First, create a controller. It should compose of fuel.Controller

```go
type HelloWorldController struct {
	fuel.Controller
}
```

Now add a field to it of type fuel.GET (this is equivalent to http get). Also implement a method that returns a string. Note that field and method have same spellings, expect that method is public & field is not.

```go
type HelloWorldController struct {
	fuel.Controller
	sayHello fuel.GET
}

func (s *HelloWorldController) SayHello() string {
	return "Hello World"
}

func main() {
	server := fuel.NewServer()
	server.AddController(&HelloWorldController{})
	server.Run()
}
```
Now open your browser and hit http://localhost:8080/hello-world/say-hello

**Note:** FUEL is fully compatbile with the standard http handler semantics. Lets say you don't want to use any magic. Just simple unadulterated http request and responses. Its time to say Hola!

```go
type HelloWorldController struct {
	fuel.Controller
	sayHello fuel.GET
	sayHola fuel.GET
}

func (s *HelloWorldController) SayHello() string {
	return "Hello World"
}

// Note: the func signature
func (s *HelloWorldController) SayHola(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola")
}

func main() {
	server := fuel.NewServer()
	server.AddController(&HelloWorldController{})
	server.Run()
}
```
You can test this by visiting: http://localhost:8080/hello-world/say-hola




#### TODO
- more test cases (WIP)
- only cache 200 - OK values
- map should be string->interface
- aide helpers
- allowed functions
- mvc
  - minification
  - templating
    - css
    - view partials
  - session handling
  - auth
  - dev vs production env
  - hot code reload
  - support
    - bower
  - grunt / task runner
  - csrf
  - cors
- re-arch cache to use middleware instead
- server events | begin_reqeust and end_request
- os signals
- slash at end? url support
- allow mux plug and play (setRouter())

