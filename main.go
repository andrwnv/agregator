package main

import (
	"fmt"
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/models"
	"github.com/andrwnv/event-aggregator/core/services"
	v1 "github.com/andrwnv/event-aggregator/routers/v1"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func init() {
	core.ServerInst = &core.Server{}

	core.ServerInst.Router = v1.InitRouter()
	core.ServerInst.Database = core.InitDatabase()
	core.ServerInst.JwtService = services.JWTAuthService()
	core.ServerInst.Database.AutoMigrate(&models.User{})
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("[ERROR]: Cant load env variables")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	fmt.Println("Server Running on Port: ", port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), core.ServerInst.Router)
	if err != nil {
		panic("[ERROR]: Cant start HTTP server")
	}
}
