# FUEL :: Caching

### Introuction

There is good support for caching built into FUEL. In fact FUEL supports not one but multiple cache providers. This can be really useful when you want to use one cache provider for one controller and another for a different controler (or endpoint)

To use cache, you first define the cache store at the server level. Of course, you can define mutliple cache stores but their access key or name should be unique.

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