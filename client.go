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

	test1ID := createMovie("test1", movieService)
	fmt.Println("test1ID:" + string(test1ID))
	test2ID := createMovie("test1", movieService)
	fmt.Println("test2ID:" + string(test2ID))

	res2, err2 := movieService.GetMovie(context.TODO(), &MovieService.GetMovieMessage{
		MovieID: test1ID,
	})

	if err2 != nil {
		fmt.Println(err2)
		return
	}

	fmt.Println(string(res2.MovieID) + ":" + res2.Title)
}

func createMovie(name string, movieService MovieService.MovieService) int32 {
	res, err := movieService.CreateMovie(context.TODO(), &MovieService.CreateMovieMessage{
		Title: name,
	})

	if err != nil {
		fmt.Println(err)
		return -1
	}
	fmt.Println(res.MovieID)

	return res.MovieID
}
