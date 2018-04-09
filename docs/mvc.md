# FUEL :: MVC

### Introuction


### Static File Server

Setting up a file server or static server to serve files from disk is fairly simple. 

Lets say you want to redirect all content from path **/assets/** to a folder, whose relative path is **./static/content/**

You can achive this by simply doing the following:

```go
type FileController struct {
    fuel.Controller
    assets fuel.GET `root:"-" folder:"./static/content/"`
}

func main() {
    server := fuel.NewServer()
    server.AddController(&FileController{})
    server.Run()
}

// Now use http://.../assets/index.html to access file ./static/content/index.html, and 
//     use http://.../assets/folderA/folderB/other.css to access the file ./static/content/folderA/folderB/other.css
```

Note that:
 - The complete route would become **/file/assets** as root name is derived from Controller name. To turn this off root was explicitly set to **-**
 - All you need to setup a server is use the **folder** tag. You set it to the directory containing the static content and that is it


### Views directory

By default, FUEL uses **views** sub-directory to pull out all the templates. Should you wish to change it, you can do so using MvcOptions:

```go
func main() {
    server := fuel.NewServer()

    // Note: Use "templats" folder in place of the default "views" folder
    server.MvcOptions.Views = "templates"

    server.AddController(&FileController{})
    server.Run()
}
```


### Layouts

To set a default layout to be used for all the views, you can simply use MvcOptions. 

```go
func main() {
    server := fuel.NewServer()

    // Note: Uses "layout.html" as the default layout
    server.MvcOptions.Layout = "layout"

    server.AddController(&FileController{})
    server.Run()
}
```

Of course, you can override it at the view level as well:

```go
type MvcController struct {
    fuel.Controller
    hello     fuel.GET
}

func (s *MvcController) Hello() fuel.View {
    return fuel.View{
        Layout: "alternate-layout",
    }
}
```

### Another Hello World - Understanding Views

```go
type MvcController struct {
	fuel.Controller
	hello fuel.GET
	hola  fuel.GET `route:"{country}"`
}

func (s *MvcController) Hello() fuel.View {
	return fuel.View{}
}

func (s *MvcController) Hola(country string) fuel.View {
	return fuel.View{
		Data: map[string]interface{}{
			"country": country,
		},
	}
}

func (c Canvas) Fuel() {
	server := fuel.NewServer()
	server.MvcOptions.Layout = "layout"
	server.AddController(&MvcController{})
	server.Run()
}
```

And we use the following templates:

```html
<!-- views/layout.html file -->
<html>
    <body>
        <h1>Hello {{ yield }}</h1>
    </body>
</html>
```

```html
<!-- views/mvc/hello.html file -->
There!
```

```html
<!-- views/mvc/hola.html file -->
{{.country}}!
```

Note that: 
 - When you hit http://localhost:8080/mvc/hello, you get "**Hello There!**"
 - When you hit http://localhost:8080/mvc/Brazil, you get "**Hello Brazil!**"
 - "mvc" directory is based on controller name (MvcController)

