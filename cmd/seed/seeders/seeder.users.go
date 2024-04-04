package seeders

import (
	"github.com/kelvinator07/go-rest-template/internal/constants"
	"github.com/kelvinator07/go-rest-template/internal/datasources/records"
	"github.com/kelvinator07/go-rest-template/pkg/helpers"
	"github.com/kelvinator07/go-rest-template/pkg/logger"
	"github.com/sirupsen/logrus"
)

var pass string
var UserData []records.Users

func init() {
	var err error
	pass, err = helpers.GenerateHash("1234567890!")
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})
	}

	UserData = []records.Users{
		{
			Username: "Kelvin Geeky 7",
			Email:    "kelvin@gmail.com",
			Password: pass,
			Active:   true,
			RoleId:   1,
		},
		{
			Username: "john doe",
			Email:    "johndoe@gmail.com",
			Password: pass,
			Active:   false,
			RoleId:   2,
		},
	}
}
