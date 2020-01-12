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
	err := service.CreateMovie(context.TODO(), &MovieService.CreateMovieMessage{
		Title: "Scarlet Monastery",
	}, &r)

	if err == nil {
		fmt.Println(err)
	}

	rr := MovieService.GetMovieResponse{}
	er := service.GetMovie(context.TODO(), &MovieService.GetMovieMessage{
		MovieID: r.MovieID,
	}, &rr)

	if er == nil {
		if r.MovieID > 0 {
			t.Logf("Successfully got movie: %d %s", rr.MovieID, rr.Title)
		}
	} else {
		fmt.Println(er)
	}
}
