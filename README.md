# Gin-MVC Web Framework

Gin is a web framework base on [Gin](https://github.com/gin-gonic/gin).

**The key features of Gin are:**

- Auto Register Controller
- Custom Register Controller
- SecureServer With Graceful shutdown
- Auto Wrap data to standard json
- Auto Wrap error to standard json

## Getting started

### Prerequisites

- **[Go](https://go.dev/)**: any one of the **three latest major** [releases](https://go.dev/doc/devel/release) (we test
  it with these).

### Getting Gin-MVC

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/sdliang1013/gin-mvc"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `gin` package:

```sh
$ go get -u github.com/sdliang1013/gin-mvc
```

### Running Gin-MVC

First you need to import Gin-MVC package for using Gin-MVC, one simplest example likes the follow `example.go`:

```go
package main

import (
	"gin-mvc/src/caul"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
}

func (ctrl *Controller) Routes() []caul.Route {
	return []caul.Route{
		{Path: "/", Func: ctrl.GetData, Method: http.MethodGet},
	}
}

func (ctrl *Controller) GetData(ctx *gin.Context) (data any, err error) {
	// get param
	var id string
	id, err = caul.NotNilString(ctx, "id")
	data = map[string]string{
		"id": id,
	}
	return
}

func Run(addr string) {
	engine := gin.Default()
	// 设置router
	router := &caul.CRouter{IRouter: engine.Group("/api/v1")}
	// register middleware
	router.RegisterMiddleware(gin.BasicAuth(gin.Accounts{
		"admin": "admin",
		"guest": "guest",
	}))
	// register routers
	// Route方式
	router.RegisterRoute(caul.CRoute{Path: "/mvc1", Controller: &Controller{}})
	// 自动扫描
	router.RegisterController("/mvc2", &Controller{})
	// start engine
	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}

func main() {
	Run(":8080")
}

```

And use the Go command to run the demo:

```
# run example.go and visit 0.0.0.0:8080/ping on browser
$ go run example.go

# 通过Route方式注册
[GIN-debug] GET    /api/v1/mvc1/             --> gin-mvc/src/caul.(*CRouter).FuncWrapper.func1 (4 handlers)
# 自动注册
[GIN-debug] GET    /api/v1/mvc2/data         --> gin-mvc/src/caul.(*CRouter).RegisterController.func1.1 (4 handlers)

```

### Learn more examples

#### Quick Start

Learn and practice more examples, please read the demo/app.go and demo/controller.go

## Benchmarks

If Route registration is used, the performance is the same as gin,

If automatic registration is used, some performance is lost

## Contributing

Gin is the work of hundreds of contributors. We appreciate your help!

Please see [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.