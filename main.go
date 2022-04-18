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
	"strconv"
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
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	fmt.Printf("Server Running on Port: %s:%d\n", host, port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), core.ServerInst.Router)
	if err != nil {
		panic("[ERROR]: Cant start HTTP server")
	}
}
