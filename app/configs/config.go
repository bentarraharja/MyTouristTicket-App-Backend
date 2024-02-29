package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

var (
	JWT_SECRET string
	CLD_URL    string
	MID_KEY    string
)

type AppConfig struct {
	DB_USERNAME    string
	DB_PASSWORD    string
	DB_HOSTNAME    string
	DB_PORT        int
	DB_NAME        string
	REDIS_ADDR     string
	REDIS_PASSWORD string
	REDIS_DB       int
}

func InitConfig() *AppConfig {
	return ReadEnv()
}

func ReadEnv() *AppConfig {
	app := AppConfig{}
	isRead := true

	if val, found := os.LookupEnv("DBUSER"); found {
		app.DB_USERNAME = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPASS"); found {
		app.DB_PASSWORD = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBHOST"); found {
		app.DB_HOSTNAME = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPORT"); found {
		cnv, _ := strconv.Atoi(val)
		app.DB_PORT = cnv
		isRead = false
	}
	if val, found := os.LookupEnv("DBNAME"); found {
		app.DB_NAME = val
		isRead = false
	}
	if val, found := os.LookupEnv("JWTSECRET"); found {
		JWT_SECRET = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLDURL"); found {
		CLD_URL = val
		isRead = false
	}
	if val, found := os.LookupEnv("MIDKEY"); found {
		MID_KEY = val
		isRead = false
	}
	if val, found := os.LookupEnv("REDIS_ADDR"); found {
		app.REDIS_ADDR = val
		isRead = false
	}
	if val, found := os.LookupEnv("REDIS_PASSWORD"); found {
		app.REDIS_PASSWORD = val
		isRead = false
	}
	if val, found := os.LookupEnv("REDIS_DB"); found {
		app.REDIS_DB, _ = strconv.Atoi(val)
		isRead = false
	}

	if isRead {
		viper.AddConfigPath(".")
		viper.SetConfigName("local")
		viper.SetConfigType("env")

		err := viper.ReadInConfig()
		if err != nil {
			log.Println("error read config : ", err.Error())
			return nil
		}

		CLD_URL = viper.GetString("CLDURL")
		JWT_SECRET = viper.GetString("JWTSECRET")
		MID_KEY = viper.GetString("MIDKEY")
		app.REDIS_ADDR = viper.GetString("REDIS_ADDR")
		app.REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")
		app.REDIS_DB, _ = strconv.Atoi(viper.Get("REDIS_DB").(string))
		app.DB_USERNAME = viper.Get("DBUSER").(string)
		app.DB_PASSWORD = viper.Get("DBPASS").(string)
		app.DB_HOSTNAME = viper.Get("DBHOST").(string)
		app.DB_PORT, _ = strconv.Atoi(viper.Get("DBPORT").(string))
		app.DB_NAME = viper.Get("DBNAME").(string)
	}

	return &app
}
