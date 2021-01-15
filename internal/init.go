package internal

import (
	"log"
	"worldmaomao/harddisk/internal/config"
	"worldmaomao/harddisk/internal/constant"
	"worldmaomao/harddisk/internal/dao"
	"worldmaomao/harddisk/internal/database"
	"worldmaomao/harddisk/internal/rest"
	"worldmaomao/harddisk/internal/service"

	"github.com/BurntSushi/toml"
	"github.com/sarulabs/di"
)

var (
	container di.Container
)

func loadConfig() (*config.Configuration, error) {
	var (
		config config.Configuration
	)
	_, err := toml.DecodeFile("./res/configuration.toml", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func Initial() {
	config, err := loadConfig()
	if err != nil {
		log.Println(err.Error())
		return
	}
	dbClient, err := database.NewDbClient(config.Database)
	if err != nil {
		log.Println(err.Error())
		return
	}
	diBuilder, _ := di.NewBuilder()
	diBuilder.Add(di.Def{
		Name: constant.Configuration,
		Build: func(ctn di.Container) (interface{}, error) {
			return config, nil
		},
	})

	diBuilder.Add(di.Def{
		Name: constant.DiskDao,
		Build: func(ctn di.Container) (interface{}, error) {
			return dao.NewDiskDao(dbClient), nil
		},
	})

	diBuilder.Add(di.Def{
		Name: constant.DiskItemDao,
		Build: func(ctn di.Container) (interface{}, error) {
			return dao.NewDiskItemDao(dbClient), nil
		},
	})

	diBuilder.Add(di.Def{
		Name: constant.UserDao,
		Build: func(ctn di.Container) (interface{}, error) {
			return dao.NewUserDao(dbClient), nil
		},
	})
	container = diBuilder.Build()
	//初始化用户信息
	err = service.NewUserService(container).Initial()
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 启动web server
	rest.NewServer(container).Start()
}
