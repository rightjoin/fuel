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
    server.DefineCache("cache1", stag.NewGoCache(5*time.Second))
    server.DefineCache("cache2", stag.NewRedisCache(...))
    server.DefineCache("cache3", <implementation of stag.Cache>)
	server.AddController(&CacheController{})
}
```