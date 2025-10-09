package microservice

import (
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/controller"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
	"github.com/WhiteMaks/go_academy_portal_service/internal/route"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
	"strconv"
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
	userService := service.NewUserService(userRepository, ms.tokenMaker)
	userController := controller.NewUserController(userService)

	router := gin.Default()

	router.POST(route.ApiUserV1, userController.CreateUserV1)
	router.POST(route.ApiUserTokenV1, userController.GenerateUserTokenV1)

	ms.Router = router

	return ms
}

func (ms *Microservice) Start() error {
	return ms.Router.Run(ms.address)
}
