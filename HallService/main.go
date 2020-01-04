package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/ob-vss-ws19/blatt-4-petya/HallService/Service"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
)

const serviceName = "HallService"

func main() {
	serv := micro.NewService(
		micro.Name(serviceName),
	)
	serv.Init()

	hsrv := Service.Spawn()
	hsrv.SetShowService(
		func() ShowService.ShowService {
			return ShowService.NewShowService("ShowService", serv.Client())
		})

	err := HallService.RegisterHallServiceHandler(serv.Server(), hsrv)

	if err == nil {
		if err := serv.Run(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
