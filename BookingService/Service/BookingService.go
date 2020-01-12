package Service

import (
	"context"
	"fmt"
	BookingService "github.com/ob-vss-ws19/blatt-4-petya/BookingService/Service/messages"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
	"sync"
	"time"
)

type Booking struct {
	UserID       int32
	ShowID       int32
	Seats        []int32
	Confirmation Confirmation
}

type Confirmation struct {
	time      time.Time
	Confirmed bool
}

type BookingMicroService struct {
	bookingRepository map[int32]*Booking
	HallService       func() HallService.HallService
	ShowService       func() ShowService.ShowService
	NextId            int32
	mu                *sync.Mutex
}

func Spawn() *BookingMicroService {
	return &BookingMicroService{
		bookingRepository: make(map[int32]*Booking),
		NextId:            1,
		mu:                &sync.Mutex{},
	}
}

func (bksrv *BookingMicroService) SetShowService(shsrv func() ShowService.ShowService) {
	bksrv.mu.Lock()
	bksrv.ShowService = shsrv
	bksrv.mu.Unlock()
}

func (bksrv *BookingMicroService) SetHallService(hsrv func() HallService.HallService) {
	bksrv.mu.Lock()
	bksrv.HallService = hsrv
	bksrv.mu.Unlock()
}

func (bksrv *BookingMicroService) ResetBookings() {
	fmt.Println("-----Entered ResetBookings-----")
	bksrv.mu.Lock()
	fmt.Println("Locked ResetBookings")

	for bkID, ele := range bksrv.bookingRepository {
		if !ele.Confirmation.Confirmed && ele.Confirmation.time.After(ele.Confirmation.time.Add(time.Minute*5)) {
			fmt.Printf("Booking expired: %d\n", bkID)
			fmt.Println("Freeing seats in show...")
			s := bksrv.ShowService()
			message := &ShowService.FreeSeatMessage{
				ShowID:    ele.ShowID,
				BookingID: bkID,
			}

			_, err := s.FreeSeats(nil, message)
			if err != nil {
				fmt.Println("Freeing seats failed!")
				fmt.Println("Unlocked ResetBookings")
				bksrv.mu.Unlock()
				fmt.Println("-----Exited ResetBookings-----")
			}
		}
	}

	fmt.Println("Unlocked ResetBookings")
	bksrv.mu.Unlock()
	fmt.Println("-----Exited ResetBookings-----")
}

func (bksrv *BookingMicroService) ConfirmBooking(ctx context.Context, req *BookingService.ConfirmBookingMessage, res *BookingService.ConfirmBookingResponse) error {
	fmt.Println("-----Entered ConfirmBooking-----")
	bksrv.mu.Lock()
	fmt.Println("Locked ConfirmBooking")

	booking, ok := bksrv.bookingRepository[req.BookingID]
	if !ok {
		fmt.Println("The booking does not exist!")
		fmt.Println("Unlocked ConfirmBooking")
		bksrv.mu.Unlock()
		fmt.Println("-----Exited ConfirmBooking-----")
		return fmt.Errorf("the booking does not exist")
	}

	fmt.Println("Locking seats...")
	s := bksrv.ShowService()

	message := &ShowService.LockSeatMessage{
		ShowID:    booking.ShowID,
		BookingID: req.BookingID,
	}

	bkg, err := s.LockSeats(ctx, message)
	if !bkg.Success {
		fmt.Errorf("The booking was rejected.")
		fmt.Println("Unlocked ConfirmBooking")
		bksrv.mu.Unlock()
		fmt.Println("-----Exited ConfirmBooking-----")
		return err
	}

	fmt.Println("Booking confirmed!")
	bksrv.bookingRepository[req.BookingID].Confirmation.Confirmed = true
	bksrv.bookingRepository[req.BookingID].Confirmation.time = time.Now()

	res.BookingID = req.BookingID

	fmt.Println("Unlocked ConfirmBooking")
	bksrv.mu.Unlock()
	fmt.Println("-----Exited ConfirmBooking-----")
	return nil
}

func (bksrv *BookingMicroService) CreateBooking(ctx context.Context, req *BookingService.CreateBookingMessage, res *BookingService.CreateBookingResponse) error {
	fmt.Println("-----Entered CreateBooking-----")
	bksrv.mu.Lock()
	fmt.Println("Locked CreateBooking")

	fmt.Println("Blocking seats...")
	s := bksrv.ShowService()

	message := &ShowService.BlockSeatMessage{
		BookingID: bksrv.NextId,
		ShowID:    req.ShowID,
		SeatID:    req.Seats,
	}

	booking, err := s.BlockSeats(ctx, message)

	if err != nil {
		fmt.Println("The booking was rejected.")
		fmt.Println("Unlocked CreateBooking")
		bksrv.mu.Unlock()
		fmt.Println("-----Exited CreateBooking-----")
		return err
	}

	if !booking.Success {
		fmt.Println("The booking was rejected.")
		fmt.Println("Unlocked CreateBooking")
		bksrv.mu.Unlock()
		fmt.Println("-----Exited CreateBooking-----")
		return fmt.Errorf("The booking was rejected.")
	}

	bksrv.bookingRepository[bksrv.NextId] = &Booking{
		UserID: req.UserID,
		ShowID: req.ShowID,
		Seats:  req.Seats,
		Confirmation: Confirmation{
			time:      time.Now(),
			Confirmed: false,
		},
	}
	res.BookingID = bksrv.NextId

	bksrv.NextId++
	fmt.Println("Increased NextID")

	fmt.Println("Unlocked CreateBooking")
	bksrv.mu.Unlock()
	fmt.Println("-----Exited CreateBooking-----")
	return nil
}

func (bksrv *BookingMicroService) DeleteElement(ctx context.Context, bookingID int32) error {
	fmt.Println("-----Entered DeleteElement-----")

	booking, ok := bksrv.bookingRepository[bookingID]
	if !ok {
		fmt.Println("The booking does not exist.")
		fmt.Println("-----Exited DeleteElement-----")
		return fmt.Errorf("the booking does not exist")
	}
	fmt.Println("Found booking")

	fmt.Println("Freeing seats...")
	s := bksrv.ShowService()

	message := &ShowService.FreeSeatMessage{
		ShowID:    booking.ShowID,
		BookingID: bookingID,
	}

	const timeout = 20 * time.Second
	ctxLong, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := s.FreeSeats(ctxLong, message)
	if err != nil {
		fmt.Println("Error while freeing seats")
		fmt.Println("-----Exited DeleteElement-----")
		return err
	}

	delete(bksrv.bookingRepository, bookingID)
	fmt.Println("Deleted booking")

	fmt.Println("-----Exited DeleteElement-----")
	return nil
}

func (bksrv *BookingMicroService) DeleteBooking(context context.Context, req *BookingService.DeleteBookingMessage, res *BookingService.DeleteBookingResponse) error {
	return bksrv.DeleteElement(context, req.BookingID)
}

func (bksrv *BookingMicroService) GetUserBookings(context context.Context, req *BookingService.GetUserBookingsMessage, res *BookingService.GetUserBookingsResponse) error {
	fmt.Println("-----Entered GetUserBookings-----")
	bksrv.mu.Lock()
	fmt.Println("Locked GetUserBookings")

	var bookings []int32

	for index, ele := range bksrv.bookingRepository {
		if ele.UserID == req.UserID {
			bookings = append(bookings, index)
		}
	}

	res.BookingID = bookings
	res.UserID = req.UserID

	fmt.Println("Unlocked GetUserBookings")
	bksrv.mu.Unlock()
	fmt.Println("-----Exited GetUserBookings-----")
	return nil
}

func (bksrv *BookingMicroService) GetBooking(context context.Context, req *BookingService.GetBookingMessage, res *BookingService.GetBookingResponse) error {
	fmt.Println("-----Entered GetBooking-----")
	bksrv.mu.Lock()
	fmt.Println("Locked GetBooking")

	booking, ok := bksrv.bookingRepository[req.BookingID]
	if ok {
		res.BookingID = req.BookingID
		res.UserID = booking.UserID
		res.ShowID = booking.ShowID
		res.Seats = booking.Seats

		fmt.Println("Unlocked GetBooking")
		bksrv.mu.Unlock()
		fmt.Println("-----Exited GetBooking-----")
		return nil
	}
	fmt.Println("Booking not found")

	fmt.Println("Unlocked GetBooking")
	bksrv.mu.Unlock()
	fmt.Println("-----Exited GetBooking-----")
	return fmt.Errorf("the booking does not exist")
}

func (bksrv *BookingMicroService) KillBookingsShow(ctx context.Context, req *BookingService.KillBookingsShowMessage, res *BookingService.KillBookingsShowResponse) error {
	fmt.Println("-----Entered KillBookingsShow-----")
	bksrv.mu.Lock()
	fmt.Println("Locked KillBookingsShow")

	fmt.Println("Deleting bookings...")
	for index, ele := range bksrv.bookingRepository {
		if ele.ShowID == req.ShowID {
			const timeout = 70 * time.Second
			ctxLong, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			err := bksrv.DeleteElement(ctxLong, index)
			if err != nil {
				fmt.Println("Deleting booking failed")
				fmt.Println("Unlocked KillBookingsShow")
				bksrv.mu.Unlock()
				fmt.Println("-----Exited KillBookingsShow-----")
				return err
			}
		}
	}
	fmt.Println("Deleted bookings")
	res.Success = true

	fmt.Println("Unlocked KillBookingsShow")
	bksrv.mu.Unlock()
	fmt.Println("-----Exited KillBookingsShow-----")
	return nil
}

func (bksrv *BookingMicroService) KillBookingsUser(ctx context.Context, req *BookingService.KillBookingsUserMessage, res *BookingService.KillBookingsUserResponse) error {
	fmt.Println("-----Entered KillBookingsUser-----")
	bksrv.mu.Lock()
	fmt.Println("Locked KillBookingsUser")

	for index, ele := range bksrv.bookingRepository {
		if ele.UserID == req.UserID {
			const timeout = 10 * time.Second
			ctxLong, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			err := bksrv.DeleteElement(ctxLong, index)
			if err != nil {
				fmt.Println("Deleting booking failed")
				fmt.Println("Unlocked KillBookingsShow")
				bksrv.mu.Unlock()
				fmt.Println("-----Exited KillBookingsShow-----")
				return err
			}
		}
	}

	fmt.Println("Unlocked KillBookingsUser")
	bksrv.mu.Unlock()
	fmt.Println("-----Exited KillBookingsUser-----")
	return nil
}
