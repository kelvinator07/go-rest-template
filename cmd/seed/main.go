package main

import (
	"github.com/kelvinator07/go-rest-template/cmd/seed/seeders"
	"github.com/kelvinator07/go-rest-template/internal/config"
	"github.com/kelvinator07/go-rest-template/internal/constants"
	"github.com/kelvinator07/go-rest-template/internal/utils"
	"github.com/kelvinator07/go-rest-template/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
	}
	logger.Info("configuration loaded", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
}

func main() {
	db, err := utils.SetupPostgresConnection()
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})
	}
	defer db.Close()

	logger.Info("seeding...", logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})

	seeder := seeders.NewSeeder(db)
	err = seeder.UserSeeder(seeders.UserData)
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})
	}
}
