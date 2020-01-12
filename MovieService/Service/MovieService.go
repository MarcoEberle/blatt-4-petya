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

const (
	PlayerNumberOne int32 = 1
)

func Spawn() *MovieMicroService {
	return &MovieMicroService{
		MovieRepository: make(map[int32]*Movie),
		NextID:          PlayerNumberOne,
		mu:              &sync.Mutex{},
	}
}

func (msrv *MovieMicroService) CreateMovie(ctx context.Context, in *MovieService.CreateMovieMessage, out *MovieService.CreateMovieResponse) error {
	fmt.Println("-----Entered CreateMovie-----")
	msrv.mu.Lock()
	fmt.Println("Locked CreateMovie")

	msrv.MovieRepository[msrv.NextID] = &Movie{title: in.Title}
	out.MovieID = msrv.NextID
	fmt.Println("Created Movie")
	msrv.NextID++
	fmt.Println("Increased NextID")

	msrv.mu.Unlock()
	fmt.Println("Unlocked CreateMovie")
	fmt.Println("-----Exited CreateMovie-----")
	return nil
}

func (msrv *MovieMicroService) DeleteMovie(ctx context.Context, in *MovieService.DeleteMovieMessage, out *MovieService.DeleteMovieResponse) error {
	fmt.Println("-----Entered DeleteMovie-----")
	msrv.mu.Lock()
	fmt.Println("Locked DeleteMovie")

	_, ok := msrv.MovieRepository[in.MovieID]

	if ok {
		fmt.Println("Found movie")

		fmt.Println("Deleting shows...")
		s := msrv.ShowService()

		message := &ShowService.KillShowsMovieMessage{
			MovieID: in.MovieID,
		}

		_, err := s.KillShowsMovie(ctx, message)
		if err != nil {
			fmt.Println("Deleting shows failed")
			msrv.mu.Unlock()
			fmt.Println("Unlocked DeleteMovie")
			fmt.Println("-----Exited DeleteMovie-----")
			return err
		}
		fmt.Println("Deleted shows")
	} else {
		fmt.Println("Movie was not found")
		return fmt.Errorf("movie was not found")
	}

	msrv.mu.Unlock()
	fmt.Println("Unlocked DeleteMovie")
	fmt.Println("-----Exited DeleteMovie-----")
	return nil
}

func (msrv *MovieMicroService) GetMovie(ctx context.Context, in *MovieService.GetMovieMessage, out *MovieService.GetMovieResponse) error {
	fmt.Println("-----Entered GetMovie-----")
	msrv.mu.Lock()
	fmt.Println("Locked GetMovie")

	m, ok := msrv.MovieRepository[in.MovieID]

	if ok {
		fmt.Println("Found movie")

		out.MovieID = in.MovieID
		out.Title = m.title

		msrv.mu.Unlock()
		fmt.Println("Unlocked GetMovie")
		fmt.Println("-----Exited GetMovie-----")
		return nil
	}

	out.MovieID = 0

	msrv.mu.Unlock()
	fmt.Println("Unlocked GetMovie")
	fmt.Println("-----Exited GetMovie-----")
	return fmt.Errorf("the movie could not be found")
}

func (msrv *MovieMicroService) SetShowService(ssrv func() ShowService.ShowService) {
	msrv.mu.Lock()
	msrv.ShowService = ssrv
	defer msrv.mu.Unlock()
}
