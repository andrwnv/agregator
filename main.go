package main

import (
	"fmt"
	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/services"
	"github.com/andrwnv/event-aggregator/misc"
	v1 "github.com/andrwnv/event-aggregator/routers/v1"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		misc.ReportCritical("Cant load env variables")
	}

	globalRepo := repo.NewPgRepo()
	userRepo := repo.NewUserRepo(globalRepo)
	_ = repo.NewEventRepo(globalRepo)
	_ = repo.NewPlaceRepo(globalRepo)
	_ = repo.NewUserStoryRepo(globalRepo)
	_ = repo.NewRegionRepo(globalRepo)

	userController := controllers.NewUserController(userRepo)
	autoController := controllers.NewAuthController(userRepo)
	fileController := controllers.NewFileController(os.Getenv("FILE_STORAGE_PATH"), userRepo)

	router := v1.MakeRouter(userController, autoController, fileController)

	core.SERVER = &core.Server{
		Router:     router.InitRouter(),
		JwtService: services.JWTAuthService(),
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		misc.ReportCritical("Cant load env variables")
	}

	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	fmt.Printf("Server Running on Port: %s:%d\n", host, port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), core.SERVER.Router)
	if err != nil {
		misc.ReportCritical("Cant start HTTP server")
	}
}
