package main

import (
	"fmt"
	"github.com/Amir122002/hotel/internal/configs"
	"github.com/Amir122002/hotel/internal/database"
	"github.com/Amir122002/hotel/internal/handlers"
	"github.com/Amir122002/hotel/internal/repositories"
	"github.com/Amir122002/hotel/internal/router"
	"github.com/Amir122002/hotel/internal/services"
	logger2 "github.com/Amir122002/hotel/logger"
	"log"
	"net/http"
)

func main() {
	err := Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	logger, err := logger2.InitLogger()
	if err != nil {
		return nil
	}

	config, err := configs.InitConfig()
	if err != nil {
		return err
	}
	address := config.Server.Host + config.Server.Port

	newDB, err := database.Db(config)
	repository := repositories.NewRepository(newDB, logger)
	service := services.NewService(repository, logger)
	handler := handlers.NewHandler(service, logger)
	routers := router.NewRouter(handler)

	srv := http.Server{
		Addr:    address,
		Handler: routers,
	}
	fmt.Println("Start")
	err = srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
