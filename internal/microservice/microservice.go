package microservice

import (
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/controller"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
	"github.com/WhiteMaks/go_academy_portal_service/internal/route"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/gin-gonic/gin"
)

type Microservice struct {
	store  database.Store
	Router *gin.Engine
}

func NewMicroservice(store database.Store) *Microservice {
	result := &Microservice{store: store}

	userRepository := repository.NewUserRepository(result.store)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	router := gin.Default()

	router.POST(route.ApiUserV1, userController.CreateUserV1)

	result.Router = router

	return result
}

func (ms *Microservice) Start(address string) error {
	return ms.Router.Run(address)
}
