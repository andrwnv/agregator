package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/services"
	"github.com/andrwnv/event-aggregator/core/usecases"
	"github.com/andrwnv/event-aggregator/misc"
	v1 "github.com/andrwnv/event-aggregator/routers/v1"
	"github.com/joho/godotenv"
)

func init() {
	// gin.SetMode(gin.ReleaseMode)

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
	likedRepo := repo.NewLikedRepo(globalRepo)

	es := services.NewEsService()
	mailer := services.MakeMailer(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_USER"))

	//es.Search("Тест событие", "place", 0, 100)

	userUsecase := usecases.NewUserUsecase(userRepo, mailer)
	eventUsecase := usecases.NewEventUsecase(eventRepo, userUsecase, regionRepo, es)
	placeUsecase := usecases.NewPlaceUsecase(placeRepo, userUsecase, regionRepo, es)
	authUsecase := usecases.NewAuthUsecase(userUsecase)
	storyUsecase := usecases.NewUserStoryUsecase(userStoryRepo, userUsecase, eventUsecase, placeUsecase)
	likeUsecase := usecases.NewLikeUsecase(likedRepo, userUsecase, eventUsecase, placeUsecase)
	searchUsecase := usecases.NewSearchUsecase(placeRepo, eventRepo, es)

	fileCtrl := controllers.NewFileController(os.Getenv("FILE_STORAGE_PATH"), userUsecase)

	router := v1.MakeRouter(
		controllers.NewUserController(userUsecase),
		controllers.NewEventController(eventUsecase, fileCtrl),
		controllers.NewPlaceController(placeUsecase, fileCtrl),
		controllers.NewAuthController(authUsecase),
		fileCtrl,
		controllers.NewCommentController(eventUsecase, placeUsecase),
		controllers.NewUserStoryController(storyUsecase, fileCtrl),
		controllers.NewLikeController(likeUsecase),
		controllers.NewSearchController(searchUsecase),
	)

	core.SERVER = &core.Server{
		Router:     router,
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
