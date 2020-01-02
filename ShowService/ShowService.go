package ShowService

import (
	"context"
	"fmt"
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
}

func Spawn() *ShowMicroService {
	return &ShowMicroService{
		ShowRepository: make(map[int32]*Show),
		NextID:         1,
		mu:             &sync.RWMutex{},
	}
}

func (shsrv ShowMicroService) CreateShow(context context.Context, req *ShowService.CreateShowMessage, res *ShowService.CreateShowResponse) error {
	shsrv.mu.Lock()
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
		return nil
	}

	return fmt.Errorf("The show could not be deleted.")
}

func (shsrv ShowMicroService) BlockSeats(context context.Context, req *ShowService.BlockSeatMessage, res *ShowService.BlockSeatResponse) error {
	res.Success = false
	res.BookingID = req.BookingID

	// Show exists
	_, exists := shsrv.ShowRepository[req.ShowID]
	if !exists {
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

	return nil
}

func (shsrv ShowMicroService) LockSeats(context context.Context, req *ShowService.LockSeatMessage, res *ShowService.LockSeatResponse) error {
	res.Success = false
	res.BookingID = req.BookingID

	for _, ele := range shsrv.ShowRepository[req.ShowID].SeatRepository {
		if ele.bookingID == req.BookingID && ele.status == Blocked {
			ele.status = Taken
			res.Success = true
		}
	}

	if !res.Success {
		return fmt.Errorf("There are no blocked seats!")
	}

	return nil
}

func (shsrv ShowMicroService) FreeSeats(context context.Context, req *ShowService.FreeSeatMessage, res *ShowService.FreeSeatResponse) error {
	res.Success = false

	for index, ele := range shsrv.ShowRepository[req.ShowID].SeatRepository {
		if ele.bookingID == req.BookingID {
			delete(shsrv.ShowRepository[req.ShowID].SeatRepository, index)
		}
	}

	res.Success = true
	return nil
}
