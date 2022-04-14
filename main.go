package main

import (
	"fmt"
	v1 "github.com/andrwnv/event-aggregator/routers/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

var router *gin.Engine

func init() {
	router = v1.InitRouter()
}

func main() {
	fmt.Println("Server Running on Port: ", 9090)
	err := http.ListenAndServe(":9090", router)
	if err != nil {
		return
	}
}
