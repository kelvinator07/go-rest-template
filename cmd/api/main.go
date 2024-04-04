package main

import (
	"runtime"

	"github.com/kelvinator07/go-rest-template/cmd/api/server"
	"github.com/kelvinator07/go-rest-template/internal/config"
	"github.com/kelvinator07/go-rest-template/internal/constants"
	"github.com/kelvinator07/go-rest-template/pkg/logger"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
	}
	logger.Info("configuration loaded", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
}

func main() {
	numCPU := runtime.NumCPU()
	logger.InfoF("The project is running on %d CPU(s)", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig}, numCPU)

	if runtime.NumCPU() > 2 {
		runtime.GOMAXPROCS(numCPU / 2)
	}

	app, err := server.NewApp()
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer})
	}

	if err := app.Run(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer})
	}
}
