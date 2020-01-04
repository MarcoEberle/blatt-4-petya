package main

import (
	"fmt"
	"github.com/micro/go-micro"
	BookingService "github.com/ob-vss-ws19/blatt-4-petya/BookingService/Service/messages"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	MovieService "github.com/ob-vss-ws19/blatt-4-petya/MovieService/Service/messages"
	"github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
)

const serviceName = "ShowService"

func main() {
	serv := micro.NewService(
		micro.Name(serviceName),
	)
	serv.Init()

	ssrv := Service.Spawn()
	ssrv.SetHallService(
		func() HallService.HallService {
			return HallService.NewHallService("HallService", serv.Client())
		})

	ssrv.SetMovieService(
		func() MovieService.MovieService {
			return MovieService.NewMovieService("MovieService", serv.Client())
		})

	ssrv.SetBookingService(
		func() BookingService.BookingService {
			return BookingService.NewBookingService("BookingService", serv.Client())
		})

	err := ShowService.RegisterShowServiceHandler(serv.Server(), ssrv)

	if err == nil {
		if err := serv.Run(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
