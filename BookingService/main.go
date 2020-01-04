package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/ob-vss-ws19/blatt-4-petya/BookingService/Service"
	BookingService "github.com/ob-vss-ws19/blatt-4-petya/BookingService/Service/messages"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
)

const serviceName = "BookingService"

func main() {
	serv := micro.NewService(
		micro.Name(serviceName),
	)
	serv.Init()

	bksrv := Service.Spawn()
	bksrv.SetHallService(
		func() HallService.HallService {
			return HallService.NewHallService("HallService", serv.Client())
		})

	bksrv.SetShowService(
		func() ShowService.ShowService {
			return ShowService.NewShowService("ShowService", serv.Client())
		})

	err := BookingService.RegisterBookingServiceHandler(serv.Server(), bksrv)

	if err == nil {
		if err := serv.Run(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
