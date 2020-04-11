package main

import (
	"fmt"
	"net/http"
	"user-web/basic/config"
	"user-web/handler"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	log "github.com/micro/go-micro/v2/util/log"

	"github.com/micro/go-micro/v2/web"

	"user-web/basic"
)

func main() {
	basic.Init()

	// 使用etcd注册
	micReg := etcd.NewRegistry(registryOptions)
	// create new web service
	service := web.NewService(
		web.Name("mu.micro.book.web.user"),
		web.Version("latest"),
		web.Registry(micReg),
		web.Address(":9088"),
	)

	// initialise service
	if err := service.Init(
		web.Action(func(c *cli.Context) {
			handler.Init()
		}),
	); err != nil {
		log.Fatal(err)
	}

	// register html handler
	service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	service.HandleFunc("/user/call", handler.UserCall)

	// 注册登录接口
	service.HandleFunc("/user/login", handler.Login)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
func registryOptions(ops *registry.Options) {
	etcdCfg := config.GetEtcdConfig()
	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdCfg.GetHost(), etcdCfg.GetPort())}
}
