# FUEL :: Parameters, Query String and Context

### Parameters

FUEL offers automatic parameter conversion. Let's look at how this works:

```go
type HelloWorldController struct {
	fuel.Controller
	saySomething fuel.GET `route:"say/{greeting}/{count}"`
}

func (s *HelloWorldController) SaySomething(greeting string, count int) string {

	repeat := ""
	for i := 0; i < count; i++ {
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
 - http://localhost:8080/hello-world/say/namaste/5

This will produce the output:
```
Namaste,Namaste,Namaste,Namaste,Namaste,
```

It is important to note that we passed two different parameters - a string and an int. These are automatically inferred from the URL, and converted to the right types, and passed to the SaySomething handler.

Also note that since the underlying router is Gorilla mux, you can use regular expressions in routes. In the above example, you could curtail count to only accept numbers like:

```go
type HelloWorldController struct {
	fuel.Controller
	saySomething fuel.GET `route:"say/{greeting}/{count:[0-9]+}"`
}
```

Notes
- FUEL has automatic parameter conversion for string, int and uint types
- You do not need to play with Request directly for any of this basic stuff


---

### Query Strings

FUEL exposes the underlying Request object to you thorugh an object called Aide. To access Aide, you just add it as an additional parameter to your method handler. So the above example would become:

```go
// Note: just added fuel.Aide as the last parameter (if you need to access the underlying Request/Response object)
func (s *HelloWorldController) SaySomething(greeting string, count int, a fuel.Aide) string {

	repeat := ""
	for i := 0; i < count; i++ {
		repeat += greeting + ","
	}

    return repeat
}
```



[todo]


---

### Context (=Aide)
