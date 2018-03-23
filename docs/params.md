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

It is important to note that we passed two different parameters - a string and an int. These are automatically inferred from the URL, and converted to the right types, and passed to the SaySomething function.

Also note that since the underlying router is Gorilla mux, you can use regular expressions in routes. In the above example, you could curtail count to only accept numbers like:

```go
type HelloWorldController struct {
	fuel.Controller
	saySomething fuel.GET `route:"say/{greeting}/{count:[0-9]+}"`
}
```

Notes
- FUEL has automatic parameter conversion for string, int and uint types
- Note that you don't need to play with Request directly for any of this basic stuff

---

### Query Strings

[todo]


---

### Context (=Aide)
