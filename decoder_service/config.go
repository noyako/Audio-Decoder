package main

import (
	"database/sql"
	"fmt"
	"path"

	_ "github.com/lib/pq"
	"github.com/noyako/Audio-Decoder/model"
	"github.com/noyako/swolf"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initConfig() error {
	viper.AddConfigPath("..")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func startService() *DecodeService {
	initConfig()

	cMessage := "host=%s port=%s user=%s password=%s"
	connection := fmt.Sprintf(cMessage,
		viper.GetString("decoder.database.host"),
		viper.GetString("decoder.database.port"),
		viper.GetString("decoder.database.user"),
		viper.GetString("decoder.database.password"))
	connection += " dbname=%s sslmode=disable"

	dealer := swolf.Setup(swolf.Config{
		Driver:         "postgres",
		Connection:     connection,
		MasterDatabase: "master",
		MasterTable:    "masters",
		MasterData:     swolf.MasterFieldResolver("username", "database"),
	})

	db, _ := sql.Open("postgres", fmt.Sprintf(connection, "master"))

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	gormDB.AutoMigrate(&model.Master{})

	return &DecodeService{
		adb:  dealer,
		addr: fmt.Sprintf("%s:%s", viper.GetString("decoder.service.addr"), viper.GetString("decoder.service.port")),
	}
}

func getLocation(username, file string) string {
	return path.Join("files", username, file)
}
