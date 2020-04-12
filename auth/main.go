package main

import (
	"auth/handler"
	"auth/subscriber"
	"fmt"

	"basic"

	"github.com/micro-in-cn/tutorials/microservice-in-micro/part2/auth/model"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part2/basic/config"
	"github.com/micro/cli"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	auth "auth/proto/auth"
)

func main() {
	basic.Init()
	micReg := etcd.NewRegistry(registryOptions)
	// New Service
	service := micro.NewService(
		micro.Name("mu.micro.book.service.auth"),
		micro.Registry(micReg),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(micro.Action(func(c *cli.Context) error {
		// 初始化handler
		model.Init()
		// 初始化handler
		handler.Init()
		return nil
	}))

	// Register Handler
	auth.RegisterAuthHandler(service.Server(), new(handler.Auth))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("mu.micro.book.service.auth", service.Server(), new(subscriber.Auth))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
func registryOptions(ops *registry.Options) {
	etcdCfg := config.GetEtcdConfig()
	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdCfg.GetHost(), etcdCfg.GetPort())}
}
