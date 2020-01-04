package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	MovieService "github.com/ob-vss-ws19/blatt-4-petya/MovieService/Service/messages"
)

func main() {
	clientService := micro.NewService(micro.Name("Client"))
	clientService.Init()

	movieService := MovieService.NewMovieService("MovieService", clientService.Client())

	res, err := movieService.CreateMovie(context.TODO(), &MovieService.CreateMovieMessage{
		Title: "Dreamcatcher2",
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.MovieID)

	res2, err2 := movieService.GetMovie(context.TODO(), &MovieService.GetMovieMessage{
		MovieID: res.MovieID,
	})

	if err2 != nil {
		fmt.Println(err2)
		return
	}

	fmt.Println(string(res2.MovieID) + ":" + res2.Title)
}
