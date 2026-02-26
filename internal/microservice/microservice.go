package microservice

import (
	"strconv"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/controller"
	"github.com/WhiteMaks/go_academy_portal_service/internal/middleware"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
	"github.com/WhiteMaks/go_academy_portal_service/internal/route"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Microservice struct {
	address    string
	store      database.Store
	tokenMaker util.TokenMaker
	Router     *gin.Engine
}

func NewMicroservice(config util.Microservice, store database.Store) *Microservice {
	tokenMaker := util.NewJWTTokenMaker(config.TokenKey, config.TokenLifeTime)

	ms := &Microservice{
		store:      store,
		tokenMaker: tokenMaker,
		address:    "0.0.0.0:" + strconv.Itoa(config.Port),
	}

	userRepository := repository.NewUserRepository(ms.store)

	microserviceService := service.NewMicroserviceService(userRepository)
	userService := service.NewUserService(userRepository, ms.tokenMaker)

	microserviceController := controller.NewMicroserviceController(microserviceService)
	userController := controller.NewUserController(userService)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET(route.ApiMicroserviceV1, microserviceController.GetMicroserviceV1)

	router.POST(route.ApiUserTokenV1, userController.GenerateUserTokenV1)
	router.POST(route.ApiUserAdminV1, userController.CreateAdminV1)

	authRoutes := router.Group("/").
		Use(middleware.AuthMiddleware(ms.tokenMaker))

	authRoutes.GET(route.ApiUserV1, userController.GetUserV1)
	authRoutes.POST(route.ApiUserV1, userController.CreateUserV1)

	ms.Router = router

	return ms
}

func (ms *Microservice) Start() error {
	return ms.Router.Run(ms.address)
}
