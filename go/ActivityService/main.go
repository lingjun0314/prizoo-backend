package main

import (
	"activity/database"
	"activity/handler"
	"activity/models"
	"activity/proto/activity"

	"github.com/micro/plugins/v5/registry/consul"
	"go-micro.dev/v5"
	"go-micro.dev/v5/logger"
	"go-micro.dev/v5/registry"
)

var (
	service = "activity"
	version = "latest"
)

func main() {
	//	Connect database
	err := database.ConnectFirestore()
	if err != nil {
		logger.Fatal("error while connect database: ",err)
	}
	defer database.CloseFirestore()

	consulReg := consul.NewRegistry(
		registry.Addrs("localhost:8499"),
	)
	//	Create new service
	service := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consulReg),
		micro.Address(":8010"),
	)

	//	Initialise flags
	service.Init()

	models.MicroClient = service.Client()
	handler.InitService()

	//	Register handler
	activity.RegisterActivityHandler(service.Server(), new(handler.Activity))

	//	Start the service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
