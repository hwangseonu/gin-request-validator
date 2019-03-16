package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
)

type HelloResource struct {
	*gin_restful.Resource
}

func (r HelloResource) Get() string {
	return "Hello, World!"
}

func main() {
	e := gin.Default()
	api := gin_restful.NewApi(e, "/")
	api.AddResource(HelloResource{gin_restful.InitResource()}, "/hello")
	_ = e.Run(":5000")
}
