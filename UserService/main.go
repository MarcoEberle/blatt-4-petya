package main

import (
	"fmt"
	"github.com/micro/go-micro"
	BookingService "github.com/ob-vss-ws19/blatt-4-petya/BookingService/Service/messages"
	"github.com/ob-vss-ws19/blatt-4-petya/UserService/Service"
	UserService "github.com/ob-vss-ws19/blatt-4-petya/UserService/Service/messages"
)

const serviceName = "UserService"

func main() {
	serv := micro.NewService(
		micro.Name(serviceName),
	)
	serv.Init()

	usrv := Service.Spawn()
	usrv.SetBookingService(
		func() BookingService.BookingService {
			return BookingService.NewBookingService("BookingService", serv.Client())
		})

	err := UserService.RegisterUserServiceHandler(serv.Server(), usrv)

	if err == nil {
		if err := serv.Run(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
