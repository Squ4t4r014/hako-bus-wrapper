package main

import (
	"github.com/gin-gonic/gin"
	"hako-bus-wrapper/app/infrastructure"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := infrastructure.NewRouting()
	_ = r.Run()
}
