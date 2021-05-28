package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/noyako/Audio-Decoder/storage"
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("..")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func startService() *UserService {
	initConfig()

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("user.database.host"),
		viper.GetString("user.database.port"),
		viper.GetString("user.database.user"),
		viper.GetString("user.database.password"),
		viper.GetString("user.database.dbname")))
	if err != nil {
		panic(err)
	}

	dp, err := storage.NewUserPostgres(db)
	dp.Migrate()
	if err != nil {
		panic(err)
	}

	return &UserService{
		db:   dp,
		addr: fmt.Sprintf("%s:%s", viper.GetString("user.service.addr"), viper.GetString("user.service.port")),
	}
}
