package db

import (
	"github.com/fxckcode/BussinessBot/env"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DSN = env.ViperEnvVariable("DATABASE_URL")
var log = logrus.New()

func DBConnection() {
	var error error
	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if error != nil {
		log.Fatal("Database connection error ", error)
	} else {
		log.Info("Database connected")
	}
}