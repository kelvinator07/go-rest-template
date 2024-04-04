package config

import (
	"github.com/kelvinator07/go-rest-template/internal/constants"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port        int    `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	Debug       bool   `mapstructure:"DEBUG"`

	DBPostgresDriver string `mapstructure:"DB_POSTGRES_DRIVER"`
	DBPostgresDsn    string `mapstructure:"DB_POSTGRES_DSN"`
	DBPostgresURL    string `mapstructure:"DB_POSTGRES_URL"`

	JWTSecret  string `mapstructure:"JWT_SECRET"`
	JWTExpired int    `mapstructure:"JWT_EXPIRED"`
	JWTIssuer  string `mapstructure:"JWT_ISSUER"`

	OTPEmail    string `mapstructure:"OTP_EMAIL"`
	OTPPassword string `mapstructure:"OTP_PASSWORD"`

	REDISHost     string `mapstructure:"REDIS_HOST"`
	REDISPassword string `mapstructure:"REDIS_PASS"`
	REDISExpired  int    `mapstructure:"REDIS_EXPIRED"`
}

func InitializeAppConfig() error {
	viper.SetConfigName(".env") // allow directly reading from .env file
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("internal/config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return constants.ErrLoadConfig
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return constants.ErrParseConfig
	}

	// check
	if AppConfig.Port == 0 || AppConfig.Environment == "" || AppConfig.JWTSecret == "" || AppConfig.JWTExpired == 0 || AppConfig.JWTIssuer == "" || AppConfig.OTPEmail == "" || AppConfig.OTPPassword == "" || AppConfig.REDISHost == "" || AppConfig.REDISPassword == "" || AppConfig.REDISExpired == 0 || AppConfig.DBPostgresDriver == "" {
		return constants.ErrEmptyVar
	}

	switch AppConfig.Environment {
	case constants.EnvironmentDevelopment:
		if AppConfig.DBPostgresDsn == "" {
			return constants.ErrEmptyVar
		}
	case constants.EnvironmentProduction:
		if AppConfig.DBPostgresURL == "" {
			return constants.ErrEmptyVar
		}
	}

	return nil
}
