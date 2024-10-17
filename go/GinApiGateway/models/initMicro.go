package models

import (
	"github.com/micro/plugins/v5/registry/consul"
	"go-micro.dev/v5"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/registry"
)

var MicroClient client.Client

func init() {
	//	Registry consul
	consulReg := consul.NewRegistry(
		registry.Addrs("localhost:8499"),
	)

	//	Create service
	srv := micro.NewService(
		micro.Registry(consulReg),
	)

	srv.Init()

	MicroClient = srv.Client()

}
