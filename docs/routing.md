# FUEL :: Routing

### Introuction

There are 4 main parts of a route in FUEL. These are:
 
 - prefix
 - version
 - root
 - route


### url = prefix + version + root + url

Lets look at our Hello World example again:

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


### Parameter Passing

FUEL offers automatic parameter conversion. Let's look at how this works:

```go
type HelloWorldController struct {
	fuel.Controller
	saySomething fuel.GET `route:"say/{greeting}/{count}"`
}

func (s *HelloWorldController) SaySomething(greeting string, count int) string {

    repeat := ""
    for i:= 0; i<count;i++ {
        repeat += greeting + ","
    }

	return repeat
}

func main() {
	server := fuel.NewServer()
	server.AddController(&HelloWorldController{})
	server.Run()
}
```
After you run this, you can simply hit the url:
 - http://localhost:8421/hello-world/say/namaste/5

This will produce the output:
```
Namaste,Namaste,Namaste,Namaste,Namaste,
```

Notes:
 - FUEL uses gorilla mux for routing.
 - You can use slashes in prefix, root and route. It doesn't have to be just a single word. So the above example prefix could be changed from 'on-the-moon' to 'solar-system/on-the-moon'. 
 - You don't have to worry about slashes. Double (or more) slashes are cleaned up internally.
 - If you want to turn off automatic url inference, you can just say root:'-'. This would stop HelloWorldController to use 'hello-world' as its root value and just set it to be empty.
