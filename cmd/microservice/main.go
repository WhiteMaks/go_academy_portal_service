package main

import (
	"database/sql"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/microservice"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	gin.SetMode("release")

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	dbConnection, err := sql.Open("postgres", config.Database.ConnectionString)
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}

	store := database.NewStore(dbConnection)

	ms := microservice.NewMicroservice(config.Microservice, store)
	err = ms.Start()
	if err != nil {
		log.Fatal("failed to start the server:", err)
	}
}
