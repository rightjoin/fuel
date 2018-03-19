# FUEL :: Routing

### Introuction

There are 4 main parts of a route in FUEL. These are:
 
 - prefix
 - version
 - root
 - route


#### url = prefix + version + root + url

Lets look again at our Hello World example again:

``` go
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

In the example above, the default value for:
 - 'root' is being inferred automatically from the Controller name (HelloWorldController).
 - 'route' is being inferred automatically by the Field name (sayHello and sayHola respectively)
 - 'version' and 'prefix' are empty

Hence the two urls that you get are:
 - http://localhost:8421/hello-world/say-hello
 - http://localhost:8421/hello-world/say-hola

Lets introduce some values for prefix, root and route:

``` go
type HelloWorldController struct {
	fuel.Controller `prefix:"on-the-moon" root:"flying-around"`
	sayHello fuel.GET `version:"1.1"`
	sayHola fuel.GET `route:"whisper-hola"`
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
The new URLs are:
- http://localhost:8421/on-the-moon/v1.1/flying-around/say-hello
- http://localhost:8421/on-the-moon/flying-around/whisper-hola

Note:
 1. 