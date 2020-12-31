# gin-docs

![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/didintri196/gin-docs)
![GitHub](https://img.shields.io/github/license/didintri196/gin-docs)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/didintri196/gin-docs)

Auto generate API Documentation Go using framework Gin

This package makes it easy for developers to create documentation automatically using the gin framework in the Golang programming language.

## Install

Install the package with:

```bash
go get github.com/didintri196/gin-docs
```

Import it with:

```go
import builder "github.com/didintri196/gin-docs"
```

and use `builder` as the package name inside the code.

## Example

Please check the example folder for details.

```go
package main

import (
	builder "github.com/didintri196/gin-docs"
)

func Middleware() {
	router := gin.Default()
	router.Use(gin.Recovery())

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type,Access-Control-Allow-Origin",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
	}))

	api := router.Group("/api")
	{
		Machine := new(controllers.Machine)
		api.GET("/list", Machine.GetDeviceCmdAll)
		api.GET("/ceksn", Machine.GetDeviceCmd)
		api.POST("/postdata/:sn", Machine.PostData)
	}
	router.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"messsage": "API ON"})
	})
	builder.Use(router) <--- using here
	router.Run(":9001")
}
```

Using function Use(*gin.Engine)
```go
	builder.Use(router)
```

## Run & Generate

How to run to generate documentation

```sh
$ ./project -doc=markdown
```
