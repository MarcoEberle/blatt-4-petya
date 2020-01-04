package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/ob-vss-ws19/blatt-4-petya/MovieService/Service"
	MovieService "github.com/ob-vss-ws19/blatt-4-petya/MovieService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
)

const serviceName = "MovieService"

func main() {
	serv := micro.NewService(
		micro.Name(serviceName),
	)
	serv.Init()

	msrv := Service.Spawn()
	msrv.SetShowService(
		func() ShowService.ShowService {
			return ShowService.NewShowService("ShowService", serv.Client())
		})

	err := MovieService.RegisterMovieServiceHandler(serv.Server(), msrv)

	if err == nil {
		if err := serv.Run(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
