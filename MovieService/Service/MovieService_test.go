package Service

import (
	"context"
	"fmt"
	MovieService "github.com/ob-vss-ws19/blatt-4-petya/MovieService/Service/messages"
	"testing"
)

func TestCreateMovie(t *testing.T) {
	service := Spawn()
	r := MovieService.CreateMovieResponse{}
	er := service.CreateMovie(context.TODO(), &MovieService.CreateMovieMessage{
		Title: "Scarlet Monastery",
	}, &r)

	if er == nil {
		if r.MovieID > 0 {
			t.Log("Successfully created movie.")
		}
	} else {
		fmt.Println(er)
	}
}

func TestGetMovie(t *testing.T) {
	service := Spawn()
	r := MovieService.CreateMovieResponse{}
	service.CreateMovie(context.TODO(), &MovieService.CreateMovieMessage{
		Title: "Scarlet Monastery",
	}, &r)

	rr := MovieService.GetMovieResponse{}
	er := service.GetMovie(context.TODO(), &MovieService.GetMovieMessage{
		MovieID: r.MovieID,
	}, &rr)

	if er == nil {
		if r.MovieID > 0 {
			t.Log("Successfully got movie.")
		}
	} else {
		fmt.Println(er)
	}
}
