# FUEL :: Stubs

### Introuction

FUEL makes it super simple to quickly create mock api stubs by only writing very little code. You basically specify a file on disk, and FUEL reads and serves back its contents

```go
type MockService struct {
	RestService
	yetToCode  GET `stub:"samples/some.txt"`
}

// And then run it
server := aqua.NewRestServer()
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
