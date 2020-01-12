package Service

import (
	"context"
	"fmt"
	BookingService "github.com/ob-vss-ws19/blatt-4-petya/BookingService/Service/messages"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	MovieService "github.com/ob-vss-ws19/blatt-4-petya/MovieService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
	"sync"
	"time"
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
	mu             *sync.Mutex
	HallService    func() HallService.HallService
	MovieService   func() MovieService.MovieService
	BookingService func() BookingService.BookingService
}

func Spawn() *ShowMicroService {
	return &ShowMicroService{
		ShowRepository: make(map[int32]*Show),
		NextID:         1,
		mu:             &sync.Mutex{},
	}
}

func (shsrv *ShowMicroService) SetMovieService(msrv func() MovieService.MovieService) {
	shsrv.mu.Lock()
	shsrv.MovieService = msrv
	shsrv.mu.Unlock()
}

func (shsrv *ShowMicroService) SetHallService(hsrv func() HallService.HallService) {
	shsrv.mu.Lock()
	shsrv.HallService = hsrv
	shsrv.mu.Unlock()
}

func (shsrv *ShowMicroService) SetBookingService(bsrv func() BookingService.BookingService) {
	shsrv.mu.Lock()
	shsrv.BookingService = bsrv
	shsrv.mu.Unlock()
}

func (shsrv *ShowMicroService) CreateShow(ctx context.Context, req *ShowService.CreateShowMessage, res *ShowService.CreateShowResponse) error {
	shsrv.mu.Lock()

	fmt.Printf("Received: hallID %d, movieID %d ", req.HallID, req.MovieID)

	m := shsrv.MovieService()
	mmes := &MovieService.GetMovieMessage{
		MovieID: req.MovieID,
	}
	const timeout = 40 * time.Second

	tempErr := true
	for tempErr {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		m, merr := m.GetMovie(ctx, mmes)
		tempErr = merr != nil
		if !tempErr && m.MovieID <= 0 {
			shsrv.mu.Unlock()
			fmt.Println("The movie does not exist.")
			return fmt.Errorf("The movie does not exist.")
		}
	}

	h := shsrv.HallService()
	hmes := &HallService.GetHallMessage{
		HallID: req.HallID,
	}

	tempErr = true
	for tempErr {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		h, herr := h.GetHall(ctx, hmes)
		tempErr = herr != nil
		if !tempErr && h.HallID <= 0 {
			shsrv.mu.Unlock()
			fmt.Println("The hall does not exist.")
			return fmt.Errorf("The hall does not exist.")
		}
	}

	shsrv.ShowRepository[shsrv.NextID] = &Show{
		hallID:         req.HallID,
		movieID:        req.MovieID,
		SeatRepository: make(map[int32]*Seat),
	}

	res.ShowID = shsrv.NextID
	shsrv.NextID++
	shsrv.mu.Unlock()
	return nil
}

func (shsrv *ShowMicroService) DeleteShow(ctx context.Context, req *ShowService.DeleteShowMessage, res *ShowService.DeleteShowResponse) error {
	shsrv.mu.Lock()
	res.Success = false
	_, hall := shsrv.ShowRepository[req.ShowID]
	if hall {
		bksrv := shsrv.BookingService()
		mes := &BookingService.KillBookingsShowMessage{
			ShowID: req.ShowID,
		}

		bksrv.KillBookingsShow(ctx, mes)

		delete(shsrv.ShowRepository, req.ShowID)
		res.Success = true
		shsrv.mu.Unlock()
		return nil
	}

	shsrv.mu.Unlock()
	return fmt.Errorf("The show could not be deleted.")
}

func (shsrv *ShowMicroService) BlockSeats(ctx context.Context, req *ShowService.BlockSeatMessage, res *ShowService.BlockSeatResponse) error {
	shsrv.mu.Lock()
	res.Success = false
	res.BookingID = req.BookingID

	// Show exists
	_, exists := shsrv.ShowRepository[req.ShowID]
	if !exists {
		shsrv.mu.Unlock()
		fmt.Println("The show could not be found.")
		return fmt.Errorf("The show could not be found.")
	}

	// Seats exists
	hallID := shsrv.ShowRepository[req.ShowID].hallID
	h := shsrv.HallService()

	message := &HallService.VerifySeatMessage{
		HallID: hallID,
		SeatID: req.SeatID,
	}

	status, err := h.VerifySeat(ctx, message)
	if !status.Success || err != nil {
		fmt.Println("The seats are not existing.")
		return fmt.Errorf("The seats are not existing.")
	}

	// Seats available
	for _, ele := range req.SeatID {
		_, alreadyTaken := shsrv.ShowRepository[req.ShowID].SeatRepository[ele]
		if alreadyTaken {
			shsrv.mu.Unlock()
			fmt.Println("The seats are not available.")
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

	res.BookingID = req.BookingID
	res.Success = true
	shsrv.mu.Unlock()
	return nil
}

func (shsrv *ShowMicroService) LockSeats(ctx context.Context, req *ShowService.LockSeatMessage, res *ShowService.LockSeatResponse) error {
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

func (shsrv *ShowMicroService) FreeSeats(ctx context.Context, req *ShowService.FreeSeatMessage, res *ShowService.FreeSeatResponse) error {
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

func (shsrv *ShowMicroService) KillShowsHall(ctx context.Context, req *ShowService.KillShowsHallMessage, res *ShowService.KillShowsHallResponse) error {
	shsrv.mu.Lock()
	res.Success = false

	b := shsrv.BookingService()

	for index, ele := range shsrv.ShowRepository {
		if ele.hallID == req.HallID {
			message := &BookingService.KillBookingsShowMessage{
				ShowID: index,
			}

			_, err := b.KillBookingsShow(ctx, message)
			res.Success = false
			return err
		}
	}

	res.Success = true
	shsrv.mu.Unlock()
	return nil
}

func (shsrv *ShowMicroService) KillShowsMovie(ctx context.Context, req *ShowService.KillShowsMovieMessage, res *ShowService.KillShowsMovieResponse) error {
	shsrv.mu.Lock()
	res.Success = false

	b := shsrv.BookingService()

	for index, ele := range shsrv.ShowRepository {
		if ele.movieID == req.MovieID {
			message := &BookingService.KillBookingsShowMessage{
				ShowID: index,
			}

			b.KillBookingsShow(ctx, message)
		}
	}

	res.Success = true
	shsrv.mu.Unlock()
	return nil
}

func (shsrv *ShowMicroService) GetShows(ctx context.Context, req *ShowService.GetShowsMessage, res *ShowService.GetShowsResponse) error {
	shsrv.mu.Lock()
	shows := []*ShowService.Show{}

	for index, ele := range shsrv.ShowRepository {
		shows = append(shows, &ShowService.Show{
			MovieID: ele.movieID,
			HallID:  ele.hallID,
			ShowID:  index,
		})
	}

	res.Shows = shows

	shsrv.mu.Unlock()
	return nil
}

func (shsrv *ShowMicroService) GetShow(ctx context.Context, req *ShowService.GetShowMessage, res *ShowService.GetShowResponse) error {
	shsrv.mu.Lock()

	ele, ok := shsrv.ShowRepository[req.ShowID]

	if !ok {
		shsrv.mu.Unlock()
		return fmt.Errorf("The show was not found.")
	}

	res.Show = &ShowService.Show{
		MovieID: ele.movieID,
		HallID:  ele.hallID,
		ShowID:  req.ShowID,
	}

	shsrv.mu.Unlock()
	return nil
}
