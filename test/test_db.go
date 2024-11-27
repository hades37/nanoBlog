package test

import (
	"nanoBlog/config"
	"nanoBlog/dao"
	"nanoBlog/utils"
)

func TestDB() {
	conf, err := config.LoadConfig("settings.yaml")
	if err != nil {
		utils.Logger.Fatal(err)
	}
	db, err := dao.InitDB(conf.DB)
	if err != nil {
		utils.Logger.Fatal(err)
	}
}
