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
- Powerful configuration model
- MVC
- Versioning
- Caching
- Stubbing
- ~~Proxying~~
- ~~CRUD~~
- Middleware support
- Logging (using middleware)

## Hello World

Lets see how we can quickly write a Hello World Api

First, create a controller. It should compose of fuel.Controller

```
type HelloWorldController struct {
	fuel.Controller
}
```

Now add a field to it of type fuel.GET (this is equivalent to http get). Also implement a method that returns a string. Note that field and method have same spellings, expect that method is public & field is not.

```
type HelloWorldController struct {
	fuel.Controller
	sayHello fuel.GET
}

func (s *HelloWorldController) SayHello() string {
	return "Hello World"
}
```
Now open your browser and hit http://localhost:8421/hello-world/say-hello

Lets say you don't want to use any magic. Just simple unadulterated http request and responses. Its time to say Hola!

```
type HelloWorldController struct {
	fuel.Controller
	sayHello fuel.GET
	sayHola fuel.GET
}

func (s *HelloWorldController) SayHello() string {
	return "Hello World"
}

func (s *HelloWorldController) SayHola(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola")
}
```
You can test this by visiting: http://localhost:8421/hello-world/say-hola

**Note:** FUEL is fully compatbile with the standard http handler semantics.

## Documentation
#### [Routing](./docs/routing.md)



##### TODO
- test cases
- map should be string->interface
- set json encoder / decoder
- ability to use a custom mux like httprouter
- allow mux plug and play (setRouter())
- slash at end? url support

