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
