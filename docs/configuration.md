# FUEL :: Configuration Model

### Hierarchical configuration

FUEL has a powerful hierarchical configuration mode that allows you to configure all kinds of things like: middleware, cache, url structures and more. It works at 5 levels:

1. Server configuration (done programmatically through code)
2. Controller configuration (set declaratively using tag)
3. Controller configuration (done programmatically through code)
4. Field or Endpoint configuration (set declaratively using tag)
5. Field or Endpoint configuratoin (done programmatically through code)


**Note:**
2 overrides 1; and 3 overrides 2, and so on.

### How it works : an example

FUEL comes with a powerful configuration model. Lets look at this with an example:

```go
type HelloWorldController struct {
	fuel.Controller
	sayHello fuel.GET
}

func (s *HelloWorldController) SayHello() string {
	return "Hello World"
}

type HolaController struct {
	fuel.Controller
	sayHola fuel.GET
}

func (s *HolaController) SayHola() string {
	return "Hola"
}


func main() {
    server := fuel.NewServer()
    server.AddController(&HelloWorldController{})
    server.AddController(&HolaController{})
    server.Run()
}
```

This code gives us two endpoints:
- http://localhost:8421/hello-world/say-hello
- http://localhost:8421/hola/say-hola

To give a version to these fields, we could do it directly at the server level:

```go
func main() {
    server := fuel.NewServer()

    // NOTE:
    // This is inherited by all apis
    server.Version = "1"

    server.AddController(&HelloWorldController{})
    server.AddController(&HolaController{})
    server.Run()
}
```

Running this will give you following two endpoints. Note that all APIs are now versioned 'v1'

- http://localhost:8421/v1/hello-world/say-hello
- http://localhost:8421/v1/hola/say-hola


Now lets say we want to have all APIs in Hola controller to be at version 2. This could be accomplished in two ways:

#### OPTION A

```go
func main() {
    server := fuel.NewServer()

    // NOTE:
    // This is inherited by all apis
    server.Version = "1"

    server.AddController(&HelloWorldController{})


    // We can set version to 2. This will be now used by all APIs witing HolaController
    // and it will override the server value of 1.
    hola := &HolaController{}
    hola.Version = "2"
    server.AddController(hola)

    server.Run()
}
```

Option A gives you a programmatic way to override configurations

#### OPTION B

```go
type HolaController struct {
	fuel.Controller `version:"2"`
	sayHola fuel.GET
}
```

Option B gives you a declarative way to override base configurations

Both these options will give you following endpoints

- http://localhost:8421/v1/hello-world/say-hello
- http://localhost:8421/v2/hola/say-hola  (Note: v2)

Now lets say you want to have multiple different endpoints within HolaController. You could override version 2 using tags at field level.


```go
type HolaController struct {
    fuel.Controller `version:"2"`
    sayHola fuel.GET
    shoutHola fuel.GET `version:"2.1"`  // Note: override at field level
}

func (s *HolaController) SayHola() string {
	return "Hola"
}

func (s *HolaController) ShoutHola() string {
	return "Hoooolaaaaa"
}
```
Now you get following endpoints:

- http://localhost:8421/v1/hello-world/say-hello (note: v1)
- http://localhost:8421/v2/hola/say-hola (note: v2)
- http://localhost:8421/v2.1/hola/shout-hola (note: v2.1)



### Configurations available

| Tag                | Usage            
| ------------------ |-----------------
| prefix, pre        | Url prefix as in: http://abc.com/[prefix]/v1/root/route
| root               | Url root as in: http://abc.com/prefix/v1/[root]/route                             
| route              | Url suffix as in in: http://abc.com/prefix/v1/root/[route]                            
| version, ver       | Url version as in: http://abc.com/prefix/v[1]/root/route                               
| cache              | The name of cache provider to use
| ttl                | Duration to cache for (e.g. 5s or 10m)
| stub               | Relative or absolute path to the file containing the mock stub
| middle, middleware | Middlewares associated with the specific endpoint (comma separated list)

