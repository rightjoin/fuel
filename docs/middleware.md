# FUEL :: Middlware

### Introuction

Oh yes, FUEL supports middleware, making them all the more configurable and all the more flexible. Lets see how

First we need to define them at the server.

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

    server.AddController(<controller1>)
    server.AddController(<controller2>)
    server.Run()
```

Note that the middleware specification also follows the configuration model of FUEL. So if you want to invoke 5 middleware on 1 endopoint in a certain order, and 3 on another in some order you could simply do that using tags:

```go
type DemoController struct {
    fuel.Controller
    fiveMiddlewareChain fuel.GET `middleware:"m1,m2,m3,m4,m5"`
    threeMiddlewareChain fuel.GET `middleware:"m1,m2,m3"`
}
```
The middleware are chained and invoked in the same order in which you specify them.