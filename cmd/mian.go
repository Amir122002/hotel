package main

import (
	"github.com/Amir122002/hotel/internal/configs"
	"github.com/Amir122002/hotel/internal/database"
	"net/http"
)

func main() {
	config, err := configs.InitConfig()
	if err != nil {
		return
	}
	address := config.Host+config.Port

	newDB,err:=database.Db()
	//repositories:=

	srv:=http.Server{
		Addr: address,
		Handler:
	}
}
