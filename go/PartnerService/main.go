package main

import (
	"partnerService/database"
	"partnerService/dependent"
	"partnerService/handler"
	"partnerService/proto/partner"

	"github.com/micro/plugins/v5/registry/consul"
	"go-micro.dev/v5"
	"go-micro.dev/v5/logger"
	"go-micro.dev/v5/registry"
)

var (
	service = "partner"
	version = "latest"
)

func main() {
	//	Connect database
	database := &database.Firestore{}
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
		micro.Address(":8013"),
	)

	//	Initialise flags
	service.Init()

	//	Register handler
	partnerHandler := handler.NewPartnerHandler(
		database,
		&dependent.CheckPartnerData{},
	) //	Dependency injection

	partner.RegisterPartnerHandler(service.Server(), partnerHandler)

	//	Start the service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
