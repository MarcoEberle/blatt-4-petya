package Service

import (
	"context"
	"fmt"
	MovieService "github.com/ob-vss-ws19/blatt-4-petya/MovieService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
	"sync"
)

type Movie struct {
	title string
}

type MovieMicroService struct {
	MovieRepository map[int32]*Movie
	NextID          int32
	mu              *sync.Mutex
	ShowService     func() ShowService.ShowService
}

func Spawn() *MovieMicroService {
	return &MovieMicroService{
		MovieRepository: make(map[int32]*Movie),
		NextID:          1,
		mu:              &sync.Mutex{},
	}
}

func (msrv *MovieMicroService) CreateMovie(ctx context.Context, in *MovieService.CreateMovieMessage, out *MovieService.CreateMovieResponse) error {

	fmt.Println("ENTERED CREATEDMOVIE")
	msrv.mu.Lock()

	fmt.Printf("NextID: %d", msrv.NextID)
	fmt.Println()
	fmt.Printf("Title: %s", in.Title)
	fmt.Println()

	msrv.MovieRepository[msrv.NextID] = &Movie{title: in.Title}
	out.MovieID = msrv.NextID
	msrv.NextID++

	for i, ele := range msrv.MovieRepository {
		fmt.Printf("%d: %s", i, ele.title)
		fmt.Println()
	}

	fmt.Printf("NextID: %d", msrv.NextID)

	msrv.mu.Unlock()
	fmt.Println("EXITED CREATEDMOVIE")
	return nil
}

func (msrv *MovieMicroService) DeleteMovie(ctx context.Context, in *MovieService.DeleteMovieMessage, out *MovieService.DeleteMovieResponse) error {
	msrv.mu.Lock()

	_, ok := msrv.MovieRepository[in.MovieID]

	if ok {
		s := msrv.ShowService()

		message := &ShowService.KillShowsMovieMessage{
			MovieID: in.MovieID,
		}

		s.KillShowsMovie(ctx, message)
	}

	msrv.mu.Unlock()
	return nil
}

func (msrv *MovieMicroService) GetMovie(ctx context.Context, in *MovieService.GetMovieMessage, out *MovieService.GetMovieResponse) error {
	msrv.mu.Lock()
	m, ok := msrv.MovieRepository[in.MovieID]

	if ok {
		out.MovieID = in.MovieID
		out.Title = m.title
		msrv.mu.Unlock()
		return nil
	}

	out.MovieID = 0
	msrv.mu.Unlock()
	return fmt.Errorf("The movie could not be found.")
}

func (msrv *MovieMicroService) SetShowService(ssrv func() ShowService.ShowService) {
	msrv.mu.Lock()
	msrv.ShowService = ssrv
	msrv.mu.Unlock()
}
