package database

import (
	"fmt"
	"github.com/Amir122002/hotel/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Pool *gorm.DB
}

func Db(config *models.Config) (*DB, error) {
	dbrUri := fmt.Sprintf("host=" + config.DB.Host + " port=" + config.DB.Port + " user=" + config.DB.User + " password=" + config.DB.Password + " database=" + config.DB.Database + " sslmode=" + config.DB.Sslmode)
	db, err := gorm.Open(postgres.Open(dbrUri), &gorm.Config{})
	if err != nil {

		return nil, err
	}
	return &DB{Pool: db}, nil
}
