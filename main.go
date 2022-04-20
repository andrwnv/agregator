package main

import (
	"fmt"
	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/services"
	v1 "github.com/andrwnv/event-aggregator/routers/v1"
	"github.com/andrwnv/event-aggregator/utils"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"
)

func init() {
	globalRepo := repo.NewPgRepo()
	userRepo := repo.NewUserRepo(globalRepo)
	_ = repo.NewEventRepo(globalRepo)
	_ = repo.NewPlaceRepo(globalRepo)
	_ = repo.NewUserStoryRepo(globalRepo)
	_ = repo.NewRegionRepo(globalRepo)

	userController := controllers.NewUserController(userRepo)
	autoController := controllers.NewAuthController(userRepo)

	controller := v1.NewController(userController, autoController)
	core.SERVER = &core.Server{
		Router:     controller.MakeRoutes(),
		JwtService: services.JWTAuthService(),
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		utils.ReportCritical("Cant load env variables")
	}

	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	fmt.Printf("Server Running on Port: %s:%d\n", host, port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), core.SERVER.Router)
	if err != nil {
		utils.ReportCritical("Cant start HTTP server")
	}
}
