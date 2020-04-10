package main

import (
	"user-srv/handler"
	"user-srv/subscriber"

	"user-srv/basic"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	user "user-srv/proto/user"
)

func main() {
	basic.Init()
	// New Service
	service := micro.NewService(
		micro.Name("mu.micro.book.service.user"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	user.RegisterUserHandler(service.Server(), new(handler.User))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("mu.micro.book.service.user", service.Server(), new(subscriber.User))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
