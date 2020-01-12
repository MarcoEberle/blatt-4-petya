package Service

import (
	"context"
	"fmt"
	BookingService "github.com/ob-vss-ws19/blatt-4-petya/BookingService/Service/messages"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	MovieService "github.com/ob-vss-ws19/blatt-4-petya/MovieService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
	"sync"
)

const (
	PlayerNumberOne int32 = 1
	Blocked         int32 = 1
	Taken           int32 = 2
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
	DeleteMode     bool
}

func Spawn() *ShowMicroService {
	return &ShowMicroService{
		ShowRepository: make(map[int32]*Show),
		NextID:         PlayerNumberOne,
		mu:             &sync.Mutex{},
		DeleteMode:     false,
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
	fmt.Println("-----CreateShow-----")
	shsrv.mu.Lock()
	fmt.Println("Locked CreateShow")

	fmt.Println("Verifying Movie...")
	m := shsrv.MovieService()
	mmes := &MovieService.GetMovieMessage{
		MovieID: req.MovieID,
	}

	_, merr := m.GetMovie(ctx, mmes)
	if merr != nil {
		fmt.Println("The movie does not exist.")
		shsrv.mu.Unlock()
		fmt.Println("Unlocked CreateShow")
		fmt.Println("-----Exited CreateShow-----")
		return fmt.Errorf("The movie does not exist.")
	}

	fmt.Println("Verifying Hall...")
	h := shsrv.HallService()
	hmes := &HallService.GetHallMessage{
		HallID: req.HallID,
	}

	_, herr := h.GetHall(ctx, hmes)

	if herr != nil {
		shsrv.mu.Unlock()
		fmt.Println("The hall does not exist.")
		fmt.Println("Unlocked CreateShow")
		fmt.Println("-----Exited CreateShow-----")
		return fmt.Errorf("The hall does not exist.")
	}

	shsrv.ShowRepository[shsrv.NextID] = &Show{
		hallID:         req.HallID,
		movieID:        req.MovieID,
		SeatRepository: make(map[int32]*Seat),
	}
	fmt.Println("Created Show!")

	res.ShowID = shsrv.NextID
	shsrv.NextID++
	fmt.Println("Increased NextID")
	shsrv.mu.Unlock()
	fmt.Println("Unlocked CreateShow")
	fmt.Println("-----Exited CreateShow-----")
	return nil
}

func (shsrv *ShowMicroService) DeleteShow(ctx context.Context, req *ShowService.DeleteShowMessage, res *ShowService.DeleteShowResponse) error {
	fmt.Println("-----Entered DeleteShow-----")
	shsrv.mu.Lock()
	fmt.Println("Locked DeleteShow")
	res.Success = false

	_, show := shsrv.ShowRepository[req.ShowID]
	if show {
		fmt.Println("Found show!")

		fmt.Println("Deleting shows bookings...")
		bksrv := shsrv.BookingService()
		mes := &BookingService.KillBookingsShowMessage{
			ShowID: req.ShowID,
		}

		_, err := bksrv.KillBookingsShow(ctx, mes)
		if err != nil {
			fmt.Println("Error while deleting bookings!")
			shsrv.mu.Unlock()
			fmt.Println("Unlocked DeleteShow")
			fmt.Println("-----Exited DeleteShow-----")
			return err
		}

		delete(shsrv.ShowRepository, req.ShowID)
		fmt.Println("Deleted show!")

		res.Success = true
		shsrv.mu.Unlock()
		fmt.Println("Unlocked DeleteShow")
		fmt.Println("-----Exited DeleteShow-----")
		return nil
	}

	shsrv.mu.Unlock()
	fmt.Println("Unlocked DeleteShow")
	fmt.Println("-----Exited DeleteShow-----")
	return fmt.Errorf("The show could not be found.")
}

func (shsrv *ShowMicroService) BlockSeats(ctx context.Context, req *ShowService.BlockSeatMessage, res *ShowService.BlockSeatResponse) error {
	fmt.Println("-----Entered BlockSeats-----")
	shsrv.mu.Lock()
	fmt.Println("Locked BlockSeats")
	res.Success = false
	res.BookingID = req.BookingID

	_, exists := shsrv.ShowRepository[req.ShowID]
	if !exists {
		fmt.Println("The show could not be found.")
		shsrv.mu.Unlock()
		fmt.Println("Unlocked BlockSeats")
		fmt.Println("-----Exited BlockSeats-----")
		return fmt.Errorf("the show could not be found")
	}

	hallID := shsrv.ShowRepository[req.ShowID].hallID
	fmt.Println("Verifying hall...")
	h := shsrv.HallService()

	message := &HallService.VerifySeatMessage{
		HallID: hallID,
		SeatID: req.SeatID,
	}

	fmt.Println("Verifying seats existence...")
	status, err := h.VerifySeat(ctx, message)
	if !status.Success || err != nil {
		fmt.Println("The seats are not existing.")
		shsrv.mu.Unlock()
		fmt.Println("Unlocked BlockSeats")
		fmt.Println("-----Exited BlockSeats-----")
		return fmt.Errorf("the seats are not existing")
	}

	fmt.Println("Verifying seats availability...")
	for _, ele := range req.SeatID {
		_, alreadyTaken := shsrv.ShowRepository[req.ShowID].SeatRepository[ele]
		if alreadyTaken {
			fmt.Println("The seats are not available.")
			shsrv.mu.Unlock()
			fmt.Println("Unlocked BlockSeats")
			fmt.Println("-----Exited BlockSeats-----")
			return fmt.Errorf("The seats are not available.")
		}
	}

	fmt.Println("Blocking Seats...")
	for _, ele := range req.SeatID {
		shsrv.ShowRepository[req.ShowID].SeatRepository[ele] = &Seat{
			status:    Blocked,
			bookingID: req.BookingID,
		}
	}
	fmt.Println("Blocked seats!")
	res.BookingID = req.BookingID
	res.Success = true
	shsrv.mu.Unlock()
	fmt.Println("Unlocked BlockSeats")
	fmt.Println("-----Exited BlockSeats-----")
	return nil
}

func (shsrv *ShowMicroService) LockSeats(ctx context.Context, req *ShowService.LockSeatMessage, res *ShowService.LockSeatResponse) error {
	fmt.Println("-----Entered LockSeats-----")
	shsrv.mu.Lock()
	fmt.Println("Locked LockSeats")
	res.Success = false
	res.BookingID = req.BookingID

	fmt.Println("Searching for blocked seats...")
	for _, ele := range shsrv.ShowRepository[req.ShowID].SeatRepository {
		if ele.bookingID == req.BookingID && ele.status == Blocked {
			ele.status = Taken
			res.Success = true
		}
	}

	if !res.Success {
		shsrv.mu.Unlock()
		fmt.Println("No blocked seats found!")
		fmt.Println("Unlocked LockSeats")
		fmt.Println("-----Exited LockSeats-----")
		return fmt.Errorf("There are no blocked seats!")
	}
	fmt.Println("Locked seats!")

	shsrv.mu.Unlock()
	fmt.Println("Unlocked LockSeats")
	fmt.Println("-----Exited LockSeats-----")
	return nil
}

func (shsrv *ShowMicroService) FreeSeats(ctx context.Context, req *ShowService.FreeSeatMessage, res *ShowService.FreeSeatResponse) error {
	fmt.Println("-----Entered FreeSeats-----")
	if !shsrv.DeleteMode {
		shsrv.mu.Lock()
		fmt.Println("Locked FreeSeats")
	} else {
		fmt.Println("FreeSeats in DeleteMode")
	}

	res.Success = false

	fmt.Println("Freeing seats...")
	for index, ele := range shsrv.ShowRepository[req.ShowID].SeatRepository {
		if ele.bookingID == req.BookingID {
			delete(shsrv.ShowRepository[req.ShowID].SeatRepository, index)
		}
	}
	fmt.Println("Freed seats!")

	res.Success = true
	if !shsrv.DeleteMode {
		shsrv.mu.Unlock()
		fmt.Println("Unlocked FreeSeats")
	}

	fmt.Println("-----Exited FreeSeats-----")
	return nil
}

func (shsrv *ShowMicroService) KillShowsHall(ctx context.Context, req *ShowService.KillShowsHallMessage, res *ShowService.KillShowsHallResponse) error {
	fmt.Println("-----Entered KillShowsHall-----")
	shsrv.mu.Lock()
	fmt.Println("Locked KillShowsHall")
	shsrv.DeleteMode = true
	fmt.Println("Activated DeleteMode")
	res.Success = false

	fmt.Println("Delete shows bookings...")
	b := shsrv.BookingService()

	for index, ele := range shsrv.ShowRepository {
		if ele.hallID == req.HallID {
			message := &BookingService.KillBookingsShowMessage{
				ShowID: index,
			}
			_, err := b.KillBookingsShow(ctx, message)
			if err != nil {
				res.Success = false
				fmt.Println("Error while deleting bookings")
				fmt.Println("Unlocked KillShowsHall")
				shsrv.mu.Unlock()
				fmt.Println("-----Exited KillShowsHall-----")
				return err
			}
		}
	}

	res.Success = true
	shsrv.DeleteMode = false
	fmt.Println("Deactivated DeleteMode")
	fmt.Println("Unlocked KillShowsHall")
	shsrv.mu.Unlock()
	fmt.Println("-----Exited KillShowsHall-----")
	return nil
}

func (shsrv *ShowMicroService) KillShowsMovie(ctx context.Context, req *ShowService.KillShowsMovieMessage, res *ShowService.KillShowsMovieResponse) error {
	fmt.Println("-----Entered KillShowsMovie-----")
	shsrv.mu.Lock()
	fmt.Println("Locked KillShowsMovie")
	shsrv.DeleteMode = true
	fmt.Println("Activated DeleteMode")
	res.Success = false

	fmt.Println("Delete shows bookings...")
	b := shsrv.BookingService()

	for index, ele := range shsrv.ShowRepository {
		if ele.movieID == req.MovieID {
			message := &BookingService.KillBookingsShowMessage{
				ShowID: index,
			}

			_, err := b.KillBookingsShow(ctx, message)
			if err != nil {
				res.Success = false
				fmt.Println("Error while deleting bookings")
				fmt.Println("Unlocked KillShowsMovie")
				shsrv.mu.Unlock()
				fmt.Println("-----Exited KillShowsMovie-----")
				return err
			}
		}
	}

	res.Success = true
	shsrv.DeleteMode = false
	fmt.Println("Deactivated DeleteMode")
	fmt.Println("Unlocked KillShowsMovie")
	shsrv.mu.Unlock()
	fmt.Println("-----Exited KillShowsMovie-----")
	return nil
}

func (shsrv *ShowMicroService) GetShows(ctx context.Context, req *ShowService.GetShowsMessage, res *ShowService.GetShowsResponse) error {
	fmt.Println("-----Entered GetShows-----")
	shsrv.mu.Lock()
	fmt.Println("Locked GetShows")

	shows := []*ShowService.Show{}

	for index, ele := range shsrv.ShowRepository {
		shows = append(shows, &ShowService.Show{
			MovieID: ele.movieID,
			HallID:  ele.hallID,
			ShowID:  index,
		})
	}

	res.Shows = shows

	fmt.Println("Unlocked GetShows")
	shsrv.mu.Unlock()
	fmt.Println("-----Exited GetShows-----")
	return nil
}

func (shsrv *ShowMicroService) GetShow(ctx context.Context, req *ShowService.GetShowMessage, res *ShowService.GetShowResponse) error {
	fmt.Println("-----Entered GetShow-----")
	shsrv.mu.Lock()
	fmt.Println("Locked GetShow")
	ele, ok := shsrv.ShowRepository[req.ShowID]

	if !ok {
		fmt.Println("Show was not found!")
		fmt.Println("Unlocked GetShow")
		shsrv.mu.Unlock()
		fmt.Println("-----Exited GetShow-----")
		return fmt.Errorf("The show was not found.")
	}

	res.Show = &ShowService.Show{
		MovieID: ele.movieID,
		HallID:  ele.hallID,
		ShowID:  req.ShowID,
	}

	fmt.Println("Unlocked GetShow")
	shsrv.mu.Unlock()
	fmt.Println("-----Exited GetShow-----")
	return nil
}
