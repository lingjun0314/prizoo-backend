package main

import (
	"prizeService/handler"
	"prizeService/handler/dependent"
	"prizeService/proto/prize"

	"github.com/micro/plugins/v5/registry/consul"
	"go-micro.dev/v5"
	"go-micro.dev/v5/logger"
	"go-micro.dev/v5/registry"
)

var (
	service = "prize"
	version = "latest"
)

func main() {
	//	Connect database
	database := &dependent.FirestoreDataAccess{}
	err := database.Connect()
	if err != nil {
		logger.Fatal(err)
	}
	defer database.Close()

	//	Register to consul
	consulReg := consul.NewRegistry(
		registry.Addrs("localhost:8499"),
	)
	//	Create new service
	service := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consulReg),
		micro.Address(":8011"),
	)

	//	Initialise flags
	service.Init()

	//	Register handler
	checkId := handler.NewCheckId(database)
	prizeHandler := handler.NewPrize(
		database,
		checkId,
		&dependent.NameHandler{},
		&dependent.ImagePathHandler{},
	) //	Dependency injection

	prize.RegisterPrizeHandler(service.Server(), prizeHandler)

	//	Start the service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
