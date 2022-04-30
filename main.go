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
	eventRepo := repo.NewEventRepo(globalRepo)
	_ = repo.NewPlaceRepo(globalRepo)
	_ = repo.NewUserStoryRepo(globalRepo)
	_ = repo.NewRegionRepo(globalRepo)

	mailer := services.MakeMailer(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_USER"))

	router := v1.MakeRouter(
		controllers.NewUserController(userRepo, mailer),
		controllers.NewAuthController(userRepo),
		controllers.NewFileController(os.Getenv("FILE_STORAGE_PATH"), userRepo),
		controllers.NewEventController(eventRepo, userRepo))

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
