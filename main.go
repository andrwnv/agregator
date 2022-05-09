package main

import (
	"fmt"
	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/endpoints"
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
	regionRepo := repo.NewRegionRepo(globalRepo)
	eventRepo := repo.NewEventRepo(globalRepo)
	placeRepo := repo.NewPlaceRepo(globalRepo)
	userStoryRepo := repo.NewUserStoryRepo(globalRepo)

	mailer := services.MakeMailer(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_USER"))

	userEndpoint := endpoints.NewUserEndpoint(userRepo, mailer)
	eventEndpoint := endpoints.NewEventEndpoint(eventRepo, userEndpoint, regionRepo)
	placeEndpoint := endpoints.NewPlaceEndpoint(placeRepo, userEndpoint, regionRepo)
	authEndpoint := endpoints.NewAuthEndpoint(userEndpoint)
	userStoryEndpoint := endpoints.NewUserStoryEndpoint(userStoryRepo, userEndpoint, eventEndpoint, placeEndpoint)

	fileCtrl := controllers.NewFileController(os.Getenv("FILE_STORAGE_PATH"), userEndpoint)

	router := v1.MakeRouter(
		controllers.NewUserController(userEndpoint),
		controllers.NewEventController(eventEndpoint, fileCtrl),
		controllers.NewPlaceController(placeEndpoint, fileCtrl),
		controllers.NewAuthController(authEndpoint),
		fileCtrl,
		controllers.NewCommentController(eventEndpoint, placeEndpoint),
		controllers.NewUserStoryController(userStoryEndpoint, fileCtrl),
	)

	core.SERVER = &core.Server{
		Router:     router,
		JwtService: services.JWTAuthService(),
	}
}

type Location struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type DataSetTest struct {
	Text     string   `json:"text"`
	Location Location `json:"location"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		misc.ReportCritical("Cant load env variables")
	}

	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	//es := services.NewEsService()
	//_ = es.Create(dto.CreateAggregatorRecordDto{
	//	LocationName: "Соляной переулок, 17",
	//	Location: dto.LocationDto{
	//		Lat: 56.499446,
	//		Lon: 84.964161,
	//	},
	//	LocationType: dto.EVENT_LOCATION_TYPE,
	//})
	//es.SearchNearby(dto.LocationDto{
	//	Lat: 56.525840,
	//	Lon: 84.963095,
	//}, 1, 2)
	//
	//if err != nil {
	//	return
	//}

	fmt.Printf("Server Running on Port: %s:%d\n", host, port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), core.SERVER.Router)

	if err != nil {
		misc.ReportCritical("Cant start HTTP server")
	}
}
