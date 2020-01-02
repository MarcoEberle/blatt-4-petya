package ShowService

import (
	"context"
	"fmt"
	BookingService "github.com/ob-vss-ws19/blatt-4-petya/BookingService/messages"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/messages"
	MovieService "github.com/ob-vss-ws19/blatt-4-petya/MovieService/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/messages"
	"sync"
)

const (
	Blocked int32 = 1
	Taken   int32 = 2
)

type Seat struct {
	status    int32
	bookingID int32
}

type Show struct {
	hallID         int32
	movieID        int32
	SeatRepository map[int32]*Seat
}

type ShowMicroService struct {
	ShowRepository map[int32]*Show
	NextID         int32
	mu             *sync.RWMutex
	HallService    func() HallService.HallService
	MovieService   func() MovieService.MovieService
	BookingService func() BookingService.BookingService
}

func Spawn() *ShowMicroService {
	return &ShowMicroService{
		ShowRepository: make(map[int32]*Show),
		NextID:         1,
		mu:             &sync.RWMutex{},
	}
}

func (shsrv ShowMicroService) SetMovieService(msrv func() MovieService.MovieService) {
	shsrv.mu.Lock()
	shsrv.MovieService = msrv
	shsrv.mu.Unlock()
}

func (shsrv ShowMicroService) SetHallService(hsrv func() HallService.HallService) {
	shsrv.mu.Lock()
	shsrv.HallService = hsrv
	shsrv.mu.Unlock()
}

func (shsrv ShowMicroService) CreateShow(context context.Context, req *ShowService.CreateShowMessage, res *ShowService.CreateShowResponse) error {
	shsrv.mu.Lock()
	m := shsrv.MovieService()
	mmes := &MovieService.GetMovieMessage{
		MovieID: req.MovieID,
	}
	_, merr := m.GetMovie(context, mmes)
	if merr != nil {
		shsrv.mu.Unlock()
		return fmt.Errorf("The movie does not exist.")
	}

	h := shsrv.HallService()
	hmes := &HallService.GetHallMessage{
		HallID: req.HallID,
	}
	_, herr := h.GetHall(context, hmes)
	if herr != nil {
		shsrv.mu.Unlock()
		return fmt.Errorf("The hall does not exist.")
	}

	shsrv.ShowRepository[shsrv.NextID] = &Show{
		hallID:         req.HallID,
		movieID:        req.MovieID,
		SeatRepository: make(map[int32]*Seat),
	}
	shsrv.mu.Unlock()
	return nil
}

func (shsrv ShowMicroService) DeleteShow(context context.Context, req *ShowService.DeleteShowMessage, res *ShowService.DeleteShowResponse) error {
	shsrv.mu.Lock()
	res.Success = false
	_, hall := shsrv.ShowRepository[req.ShowID]
	if hall {
		delete(shsrv.ShowRepository, req.ShowID)
		res.Success = true
		shsrv.mu.Unlock()
		return nil
	}

	shsrv.mu.Unlock()
	return fmt.Errorf("The show could not be deleted.")
}

func (shsrv ShowMicroService) BlockSeats(context context.Context, req *ShowService.BlockSeatMessage, res *ShowService.BlockSeatResponse) error {
	shsrv.mu.Lock()
	res.Success = false
	res.BookingID = req.BookingID

	// Show exists
	_, exists := shsrv.ShowRepository[req.ShowID]
	if !exists {
		shsrv.mu.Unlock()
		return fmt.Errorf("The show could not be found.")
	}

	// Seats exists
	hallID := shsrv.ShowRepository[req.ShowID].hallID
	h := shsrv.HallService()

	message := &HallService.VerifySeatMessage{
		HallID: hallID,
		SeatID: req.SeatID,
	}

	status, err := h.VerifySeat(context, message)
	if !status.Success || err != nil {
		return fmt.Errorf("The seats are not existing.")
	}

	// Seats available
	for _, ele := range req.SeatID {
		_, alreadyTaken := shsrv.ShowRepository[req.ShowID].SeatRepository[ele]
		if alreadyTaken {
			shsrv.mu.Unlock()
			return fmt.Errorf("The seats are not available.")
		}
	}

	// Block seats
	for _, ele := range req.SeatID {
		shsrv.ShowRepository[req.ShowID].SeatRepository[ele] = &Seat{
			status:    Blocked,
			bookingID: req.BookingID,
		}
	}
	shsrv.mu.Unlock()
	return nil
}

func (shsrv ShowMicroService) LockSeats(context context.Context, req *ShowService.LockSeatMessage, res *ShowService.LockSeatResponse) error {
	shsrv.mu.Lock()
	res.Success = false
	res.BookingID = req.BookingID

	for _, ele := range shsrv.ShowRepository[req.ShowID].SeatRepository {
		if ele.bookingID == req.BookingID && ele.status == Blocked {
			ele.status = Taken
			res.Success = true
		}
	}

	if !res.Success {
		shsrv.mu.Unlock()
		return fmt.Errorf("There are no blocked seats!")
	}
	shsrv.mu.Unlock()
	return nil
}

func (shsrv ShowMicroService) FreeSeats(context context.Context, req *ShowService.FreeSeatMessage, res *ShowService.FreeSeatResponse) error {
	shsrv.mu.Lock()
	res.Success = false

	for index, ele := range shsrv.ShowRepository[req.ShowID].SeatRepository {
		if ele.bookingID == req.BookingID {
			delete(shsrv.ShowRepository[req.ShowID].SeatRepository, index)
		}
	}

	res.Success = true
	shsrv.mu.Unlock()
	return nil
}

func (shsrv ShowMicroService) KillShowsHall(context context.Context, req *ShowService.KillShowsHallMessage, res *ShowService.KillShowsHallResponse) error {
	shsrv.mu.Lock()
	res.Success = false

	b := shsrv.BookingService()

	for index, ele := range shsrv.ShowRepository {
		if ele.hallID == req.HallID {
			message := &BookingService.KillBookingsMessage{
				ShowID: index,
			}

			b.KillBookings(context, message)
		}
	}

	res.Success = true
	shsrv.mu.Unlock()
	return nil
}

func (shsrv ShowMicroService) KillShowsMovie(context context.Context, req *ShowService.KillShowsMovieMessage, res *ShowService.KillShowsMovieResponse) error {
	shsrv.mu.Lock()
	res.Success = false

	b := shsrv.BookingService()

	for index, ele := range shsrv.ShowRepository {
		if ele.movieID == req.MovieID {
			message := &BookingService.KillBookingsMessage{
				ShowID: index,
			}

			b.KillBookings(context, message)
		}
	}

	res.Success = true
	shsrv.mu.Unlock()
	return nil
}
