package main

import (
	"fmt"
	"user-srv/basic/config"
	"user-srv/handler"
	"user-srv/model"
	"user-srv/subscriber"

	"user-srv/basic"

	s "user-srv/proto/user"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry/etcd"
	log "github.com/micro/go-micro/v2/util/log"

	"github.com/micro/go-micro/v2/registry"
)

func main() {
	basic.Init()

	// 使用etcd注册
	micReg := etcd.NewRegistry(registryOptions)
	// New Service
	service := micro.NewService(
		micro.Name("mu.micro.book.service.user"),
		micro.Registry(micReg),
		micro.Version("latest"),
	)

	// Initialise service
	// 服务初始化
	service.Init(
		micro.Action(func(c *cli.Context) error {
			// 初始化模型层
			model.Init()
			// 初始化handler
			handler.Init()
			return nil
		}),
	)

	// Register Handler
	s.RegisterUserHandler(service.Server(), new(handler.User))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("mu.micro.book.service.user", service.Server(), new(subscriber.User))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registryOptions(ops *registry.Options) {
	etcdCfg := config.GetEtcdConfig()
	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdCfg.GetHost(), etcdCfg.GetPort())}
}
