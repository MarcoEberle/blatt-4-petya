package main

import (
	"context"
	"fmt"
	micro "github.com/micro/go-micro"
	proto "github.com/ob-vss-ws19/blatt-4-petya/messages"
)

func main() {
	service := micro.NewService(micro.Name("greeter.client"))
	service.Init()

	greeter := proto.NewGreeterService("greeter", service.Client())

	rsp, err := greeter.Hello(context.TODO(), &proto.Request{Name: "John"})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rsp.Greeting)
}
