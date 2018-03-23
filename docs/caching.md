# FUEL :: Caching

### Introuction

There is good support for caching built into FUEL. In fact FUEL supports not one but multiple cache providers. This can be really useful when you want to use one cache provider for one controller and another for a different controler (or endpoint)

To use cache, you first define the cache store at the server level. Of course, you can define mutliple cache stores but their access key or names should be unique.

The cache store must implement the following interface. (There are a few implementation so cache in rightjoin/stag project including Redis and go-cache)

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
    server.DefineCache("cache3", <your implementation of stag.Cache>)
    server.AddController(&CacheController{})
}
```

Now using this cache is straight enough:

```go
type CacheController struct {
    fuel.Controller `cache:"cache1" ttl:"1m"`
    slowCall1 fuel.GET
    slowCall2 fuel.GET `ttl:"5m"`
    slowCall3 fuel.GET `cache:"cache2" ttl:"1h"`
    slowCall4 fuel.GET `cache:"cache3" ttl:"6h"`
}

func (s *CacheController) SlowCall1() string {
    time.Sleep(1 * time.Second)
    return "Slow1"
}

func (s *CacheController) SlowCall2() string {
    time.Sleep(2 * time.Second)
    return "Slow2"
}

func (s *CacheController) SlowCall3() string {
    time.Sleep(3 * time.Second)
    return "Slow3"
}

func (s *CacheController) SlowCall4() string {
    time.Sleep(4 * time.Second)
    return "Slow4"
}
```

**Points to note**
 - Setting 'cache' and 'ttl' at controller level (fuel.Controller tag) ensures that all services/actions under this controller are cached. So even though cache tag is not set for slowCall1 (http://localhost:8080/cache/slow-call1), it still inherits it from controller and ends up getting cached for 1 minute in cache store 1.
 - slowCall2 (http://localhost:8080/cache/slow-call2) overrides 'ttl' to '5m'. Hence it gets cached for 5 minutes in cache store 1.
 - slowCall3 (http://localhost:8080/cache/slow-call3) is cached for 1 hour in cache store 2.
 - slowCall4 (http://localhost:8080/cache/slow-call4) is cached for 6 hours in cache store 3.

**How does caching atually work?**
 - FUEL caches the output of your function/handler into the given cache store. In the above examples, it would be 'string' - 'Slow4'.
 - For the cache duration, FUEL would use this cache value instead of invoking the said function/handler.
 
**Cache Index**
  - Be default FUEL uses the relative URL of endpoint for cahcing.
  - If you want to change this behavior, you can do so by upading FUEL.CacheKey function. For example, you may want to add session_id to this key to cache same URL separately for each user


