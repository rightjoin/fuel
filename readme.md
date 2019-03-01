# FUEL UP
REST framework

## Inspriration
- Popular WebServers (like Apache, IIS) for hierarchical configuration model
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
- Modularized to define endpoints (structured and organized codebase)
- Hierarchical configuration
- Routing
- Parameters, Query String and Context
- Versioning
- Caching
- Mocks & Stubs
- Middleware
- ~~MVC~~
- ~~Proxying~~
- ~~CRUD~~

## Hello World

Lets see how we can quickly write a Hello World Api

First, create a service. It should compose of fuel.Service

```go
type HelloWorldService struct {
	fuel.Service
}
```

Now add a field to it of type fuel.GET (this is equivalent to http get). Also implement a method that returns a string. Note that field and method have same spellings, expect that method is public & field is not.

```go
type HelloWorldService struct {
	fuel.Service
	sayHello fuel.GET
}

func (s *HelloWorldService) SayHello() string {
	return "Hello World"
}

func main() {
	server := fuel.NewServer()
	server.AddService(&HelloWorldService{})
	server.Run()
}
```
Now open your browser and hit http://localhost:8080/hello-world/say-hello

**Note:** FUEL is fully compatbile with the standard http handler semantics. Lets say you don't want to use any magic. Just simple unadulterated http request and responses. Its time to say Hola!

```go
type HelloWorldService struct {
	fuel.Service
	sayHello fuel.GET
	sayHola fuel.GET
}

func (s *HelloWorldService) SayHello() string {
	return "Hello World"
}

// Note: the func signature
func (s *HelloWorldService) SayHola(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola")
}

func main() {
	server := fuel.NewServer()
	server.AddService(&HelloWorldService{})
	server.Run()
}
```
You can test this by visiting: http://localhost:8080/hello-world/say-hola

## Hierarchical configuration

FUEL has a powerful hierarchical configuration mode that allows you to configure all kinds of things like: middleware, cache, url structures and more. It works at 5 levels:

1. Server configuration (done programmatically through code)
2. Service configuration (set declaratively using tags)
3. Service configuration (done programmatically through code)
4. Field or Endpoint configuration (set declaratively using tags)
5. Field or Endpoint configuratoin (done programmatically through code)


**Note:**
2 overrides 1; and 3 overrides 2, and so on.

---

### How it works : an example

Lets look at this with an example:

```go
type HelloWorldService struct {
	fuel.Service
	sayHello fuel.GET
}

func (s *HelloWorldService) SayHello() string {
	return "Hello World"
}

type HolaService struct {
	fuel.Service
	sayHola fuel.GET
}

func (s *HolaService) SayHola() string {
	return "Hola"
}


func main() {
    server := fuel.NewServer()
    server.AddService(&HelloWorldService{})
    server.AddService(&HolaService{})
    server.Run()
}
```

This code gives us two endpoints:
- http://localhost:8080/hello-world/say-hello
- http://localhost:8080/hola/say-hola

To give a version to these fields, we could do it directly at the server level:

```go
func main() {
    server := fuel.NewServer()

    // NOTE:
    // This is inherited by all apis
    server.Version = "1"

    server.AddService(&HelloWorldService{})
    server.AddService(&HolaService{})
    server.Run()
}
```

Running this will give you following two endpoints. Note that all APIs are now versioned 'v1'

- http://localhost:8080/v1/hello-world/say-hello
- http://localhost:8080/v1/hola/say-hola


Now lets say we want to have all APIs in Hola service to be at version 2. This could be accomplished in two ways:

#### OPTION A

```go
func main() {
    server := fuel.NewServer()

    // NOTE:
    // This is inherited by all apis
    server.Version = "1"

    server.AddService(&HelloWorldService{})

    // We can set version to 2. This will be now used by all APIs witing HolaService
    // and it will override the server value of 1.
    hola := &HolaService{}
    hola.Version = "2"
    server.AddService(hola)

    server.Run()
}
```

Option A gives you a programmatic way to override configurations

#### OPTION B

```go
type HolaService struct {
	fuel.Service `version:"2"`
	sayHola fuel.GET
}
```

Option B gives you a declarative way to override base configurations

Both these options will give you following endpoints

- http://localhost:8080/v1/hello-world/say-hello
- http://localhost:8080/v2/hola/say-hola  (Note: v2)

Now lets say you want to have multiple different endpoints within HolaService. You could override version 2 using tags at field level.


```go
type HolaService struct {
    fuel.Service `version:"2"`
    sayHola fuel.GET
    shoutHola fuel.GET `version:"2.1"`  // Note: override at field level
}

func (s *HolaService) SayHola() string {
	return "Hola"
}

func (s *HolaService) ShoutHola() string {
	return "Hoooolaaaaa"
}
```
Now you get following endpoints:

- http://localhost:8080/v1/hello-world/say-hello (note: v1)
- http://localhost:8080/v2/hola/say-hola (note: v2)
- http://localhost:8080/v2.1/hola/shout-hola (note: v2.1)

---

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


## Routing

There are 4 main parts of a route in FUEL. These are:
 
 - prefix
 - version
 - root
 - route


### url = prefix + version + root + url

Lets look at our Hello World example again:

``` go
type HelloWorldService struct {
	fuel.Service
	sayHello fuel.GET
	sayHola fuel.GET
}

func (s *HelloWorldService) SayHello() string {
	return "Hello World"
}

// Note: the func signature
func (s *HelloWorldService) SayHola(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola")
}

func main() {
	server := fuel.NewServer()
	server.AddService(&HelloWorldService{})
	server.Run()
}
```

In the example above, the default value for:
 - 'root' is being inferred automatically from the Service name (HelloWorldService).
 - 'route' is being inferred automatically by the Field name (sayHello and sayHola respectively)
 - 'version' and 'prefix' are empty

Hence the two urls that you get are:
 - http://localhost:8080/hello-world/say-hello
 - http://localhost:8080/hello-world/say-hola

Lets introduce some values for prefix, root and route:

``` go
type HelloWorldService struct {
	fuel.Service `prefix:"on-the-moon" root:"flying-around"`
	sayHello fuel.GET `version:"1.1"`
	sayHola fuel.GET `route:"whisper-hola"`
}

func (s *HelloWorldService) SayHello() string {
	return "Hello World"
}

// Note: the func signature
func (s *HelloWorldService) SayHola(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola")
}

func main() {
	server := fuel.NewServer()
	server.AddService(&HelloWorldService{})
	server.Run()
}
```
The new URLs are:
- http://localhost:8080/on-the-moon/v1.1/flying-around/say-hello
- http://localhost:8080/on-the-moon/flying-around/whisper-hola

---

Notes:
 - FUEL uses gorilla mux for routing
 - You can use slashes in prefix, root and route. It doesn't have to be just a single word. So the above example prefix could be changed from 'on-the-moon' to 'solar-system/on-the-moon'
 - You don't have to worry about extra slashes. Double (or more) slashes are cleaned up internally
 - If you want to turn off automatic url inference, you can just say root:'-'. This would stop HelloWorldService to use 'hello-world' as its root value and just set it to be empty

## Parameters, Query Strings & more

FUEL offers automatic parameter conversion. Let's look at how this works:

```go
type HelloWorldService struct {
	fuel.Service
	saySomething fuel.GET `route:"say/{greeting}/{count}"`
}

func (s *HelloWorldService) SaySomething(greeting string, count int) string {

	repeat := ""
	for i := 0; i < count; i++ {
		repeat += greeting + ","
	}

    return repeat
}

func main() {
	server := fuel.NewServer()
	server.AddService(&HelloWorldService{})
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
type HelloWorldService struct {
	fuel.Service
	saySomething fuel.GET `route:"say/{greeting}/{count:[0-9]+}"`
}
```

Notes
- FUEL has automatic parameter conversion for string, int and uint types
- You do not need to play with Request directly for any of this basic stuff


### Query Strings

FUEL exposes the underlying Request object to you thorugh an object called Aide. To access Aide, you just add it as an additional parameter to your method handler. So the above example would become:

```go
// Note: just added fuel.Aide as the last parameter (if you need to access the underlying Request/Response object)
func (s *HelloWorldService) SaySomething(greeting string, count int, a fuel.Aide) string {

	repeat := ""
	for i := 0; i < count; i++ {
		repeat += greeting + ","
	}

    return repeat
}
```

[todo]


## Mocks & Stubs

FUEL makes it  simple to quickly create mock api stubs by only writing very little code. You basically specify a file on disk, and FUEL reads and serves back its contents

```go
type MockService struct {
	fuel.Service
	yetToCode fuel.GET `stub:"sub/directory/stub_file.txt"`
}

// And then run it
server := fuel.NewServer()
server.AddService(&MockService{})
server.Run()
```

Where is the stub file pick up from? FUEL tries to read it in this order:
 - If the file is specifed as absolute path then its simple.
 - In case of relative paths:
   - FUEL first scans it in executable directory
   - and and then looks it up in working directory
 - In case file is not found, you get 404

Note that when you use 'stub', you do not need to define any method implementation


## Caching

There is good support for caching built into FUEL. In fact FUEL supports not one but multiple cache providers. This can be really useful when you want to use one cache provider for one Service and another for a different controler (or endpoint)

To use cache, you first define the cache store at the server level. Each cache store should have a unique key/name.

The cache store must implement the following interface. (There are a few implementation so cache in rightjoin/stak project including Redis and go-cache)

```go
type Cache interface {
    Set(key string, data []byte, expireIn time.Duration) error
    Get(key string) ([]byte, error)

    PrepareIndex(key string) string
    Delete(key string) error
    Close() error
}
```

Lets add the cache to the server:

```go

// Note: 
// rightjoin/stag/GoCache implements rightjoin/stag/Cache interface

func main() {
    server := fuel.NewServer()
    server.DefineCache("cache1", stag.NewGoCache(5*time.Minute))
    server.DefineCache("cache2", stag.NewRedisCache(...))
    server.DefineCache("cache3", <your implementation of stak.Cache>)
    server.AddService(&CacheService{})
}
```

Now using this cache is straight enough:

```go
type CacheService struct {
    fuel.Service `cache:"cache1" ttl:"1m"`
    slowCall1 fuel.GET
    slowCall2 fuel.GET `ttl:"5m"`
    slowCall3 fuel.GET `cache:"cache2" ttl:"1h"`
    slowCall4 fuel.GET `cache:"cache3" ttl:"6h"`
}

func (s *CacheService) SlowCall1() string {
    time.Sleep(1 * time.Second)
    return "Slow1"
}

func (s *CacheService) SlowCall2() string {
    time.Sleep(2 * time.Second)
    return "Slow2"
}

func (s *CacheService) SlowCall3() string {
    time.Sleep(3 * time.Second)
    return "Slow3"
}

func (s *CacheService) SlowCall4() string {
    time.Sleep(4 * time.Second)
    return "Slow4"
}
```

**Points to note**
 - Setting 'cache' and 'ttl' at Service level (fuel.Service tag) ensures that all services/actions under this Service are cached. So even though cache tag is not set for slowCall1 (http://localhost:8080/cache/slow-call1), it still inherits it from Service and ends up getting cached for 1 minute in cache store 1.
 - slowCall2 (http://localhost:8080/cache/slow-call2) overrides 'ttl' to '5m'. Hence it gets cached for 5 minutes in cache store 1.
 - slowCall3 (http://localhost:8080/cache/slow-call3) is cached for 1 hour in cache store 2.
 - slowCall4 (http://localhost:8080/cache/slow-call4) is cached for 6 hours in cache store 3.

**How does caching atually work?**
 - FUEL caches the output of your function/handler into the given cache store. In the above examples, it would be 'string' - 'Slow4'.
 - For the cache duration, FUEL would use this cache value instead of invoking the said function/handler.
 
**Cache Index**
  - Be default FUEL uses the relative URL of endpoint for cahcing.
  - If you want to change this behavior, you can do so by upading FUEL.CacheKey function. For example, you may want to add session_id to this key to cache same URL separately for each user


### Middleware

FUEL supports middleware, making them all the more configurable and all the more flexible. Lets see how

First we need to define them at the server

```go
    server := fuel.NewServer();
    server.Define("m1", <returns func(http.Handler) http.Handler>)
    server.Define("m2", <returns func(http.Handler) http.Handler>)
    server.Define("m3", <returns func(http.Handler) http.Handler>)
    // and so on
```

To take a concrete example, lets define a middleware to create access logs

```go
// logs every request as info
func MidAccessLog() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info(r.RequestURI, "time", fmt.Sprintf("%.3fs", time.Now().Sub(start).Seconds()))
		})
	}
}

// logs slow requests as warnings
func MidSlowLog(slowSeconds float64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			span := time.Now().Sub(start).Seconds()
			if span > slowSeconds {
				log.Warn(r.RequestURI, "time", fmt.Sprintf("%.3fs", span))
			}
		})
	}
}

```

To use it throughout the project on every single endpoint/action, we just enable it at the server level:

```go
    server := fuel.NewServer();

    // attach the middlewares to the server
    server.DefineMiddleware("access", MidAccessLog())
    server.DefineMiddleware("slow", MidSlowLog(0.5))

    // invoke the middleware for every endpoint in the given order
    server.Middleware = "access, slow"

    server.AddService(<Service1>)
    server.AddService(<Service2>)
    server.Run()
```

Note that the middleware specification also follows the configuration model of FUEL. So if you want to invoke 5 middleware on 1 endopoint in a certain order, and 3 on another in some order you could simply do that using tags:

```go
type DemoService struct {
    fuel.Service
    fiveMiddlewareChain fuel.GET `middleware:"m1,m2,m3,m4,m5"`
    threeMiddlewareChain fuel.GET `middleware:"m1,m2,m3"`
}
```
The middleware are chained and invoked in the same order in which you specify them.




#### TODO
- more test cases (WIP)
- only cache 200 - OK values
- map should be string->interface
- aide helpers
- allowed functions
- auth
- hot code reload
- re-arch cache to use middleware instead
- server events | begin_reqeust and end_request
- os signals
- slash at end? url support
- allow mux plug and play (setRouter())

