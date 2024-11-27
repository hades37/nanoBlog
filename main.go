package main

import (
	"fmt"
	"nanoBlog/config"
	"nanoBlog/dao"
	"nanoBlog/utils"
)

func main() {
	conf, err := config.LoadConfig("settings.yaml")
	if err != nil {
		utils.GetLogger().Fatal(err.Error())
	}
	err = dao.InitDB(&conf.DB)
	if err != nil {
		utils.GetLogger().Fatal(err.Error())
	} else {
		utils.GetLogger().Info("数据库连接成功")
	}
	fmt.Println(conf)
}
