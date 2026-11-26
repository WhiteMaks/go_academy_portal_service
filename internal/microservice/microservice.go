package microservice

import (
	"strconv"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/controller"
	"github.com/WhiteMaks/go_academy_portal_service/internal/middleware"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
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
	if err := util.RegisterBindingValidators(); err != nil {
		panic(err)
	}
	if err := model.RegisterAthleteBindingValidatorsFromGin(); err != nil {
		panic(err)
	}
	if err := model.RegisterCoachBindingValidatorsFromGin(); err != nil {
		panic(err)
	}
	if err := model.RegisterBranchBindingValidatorsFromGin(); err != nil {
		panic(err)
	}
	if err := model.RegisterHallBindingValidatorsFromGin(); err != nil {
		panic(err)
	}

	tokenMaker := util.NewJWTTokenMaker(config.TokenKey, config.TokenLifeTime)

	ms := &Microservice{
		store:      store,
		tokenMaker: tokenMaker,
		address:    "0.0.0.0:" + strconv.Itoa(config.Port),
	}

	userRepository := repository.NewUserRepository(ms.store)
	athleteRepository := repository.NewAthleteRepository(ms.store)
	coachRepository := repository.NewCoachRepository(ms.store)
	branchRepository := repository.NewBranchRepository(ms.store)
	hallRepository := repository.NewHallRepository(ms.store)

	microserviceService := service.NewMicroserviceService(userRepository)
	userService := service.NewUserService(userRepository, ms.tokenMaker)
	athleteService := service.NewAthleteService(athleteRepository)
	coachService := service.NewCoachService(coachRepository)
	branchService := service.NewBranchService(branchRepository)
	hallService := service.NewHallService(hallRepository, branchRepository)

	microserviceController := controller.NewMicroserviceController(microserviceService)
	userController := controller.NewUserController(userService)
	athleteController := controller.NewAthleteController(athleteService, userService)
	coachController := controller.NewCoachController(coachService, userService)
	branchController := controller.NewBranchController(branchService, userService)
	hallController := controller.NewHallController(hallService, userService)

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
		Use(middleware.AuthMiddleware(ms.tokenMaker, userService))

	authRoutes.GET(route.ApiUserV1, userController.GetUserV1)
	authRoutes.POST(route.ApiUserV1, userController.CreateUserV1)

	authRoutes.POST(route.ApiAthleteV1, athleteController.CreateAthleteV1)
	authRoutes.GET(route.ApiAthletesV1, athleteController.GetAthletesV1)
	authRoutes.GET(route.ApiAthleteByIDV1, athleteController.GetAthleteV1)
	authRoutes.PUT(route.ApiAthleteByIDV1, athleteController.UpdateAthleteV1)

	authRoutes.POST(route.ApiCoachV1, coachController.CreateCoachV1)
	authRoutes.GET(route.ApiCoachesV1, coachController.GetCoachesV1)
	authRoutes.GET(route.ApiCoachByIDV1, coachController.GetCoachV1)
	authRoutes.PUT(route.ApiCoachByIDV1, coachController.UpdateCoachV1)

	authRoutes.POST(route.ApiBranchV1, branchController.CreateBranchV1)
	authRoutes.GET(route.ApiBranchesV1, branchController.GetBranchesV1)
	authRoutes.GET(route.ApiBranchHallsV1, hallController.GetBranchHallsV1)
	authRoutes.GET(route.ApiBranchByIDV1, branchController.GetBranchV1)
	authRoutes.PUT(route.ApiBranchByIDV1, branchController.UpdateBranchV1)

	authRoutes.POST(route.ApiHallV1, hallController.CreateHallV1)
	authRoutes.GET(route.ApiHallByIDV1, hallController.GetHallV1)
	authRoutes.PUT(route.ApiHallByIDV1, hallController.UpdateHallV1)

	ms.Router = router

	return ms
}

func (ms *Microservice) Start() error {
	return ms.Router.Run(ms.address)
}
