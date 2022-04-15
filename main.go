package main

import (
	"fmt"
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/models"
	v1 "github.com/andrwnv/event-aggregator/routers/v1"
	"net/http"
)

func init() {
	core.ServerInst = &core.Server{}

	core.ServerInst.Router = v1.InitRouter()
	core.ServerInst.Database = core.InitDatabase()

	core.ServerInst.Database.AutoMigrate(&models.User{})
}

func main() {
	fmt.Println("Server Running on Port: ", 9090)
	err := http.ListenAndServe(":9090", core.ServerInst.Router)
	if err != nil {
		return
	}
}
