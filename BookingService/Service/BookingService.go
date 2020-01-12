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

const (
	PlayerNumberOne int32 = 1
	MamboNumberFive int32 = 5
)

func Spawn() *BookingMicroService {
	return &BookingMicroService{
		bookingRepository: make(map[int32]*Booking),
		NextId:            PlayerNumberOne,
		mu:                &sync.Mutex{},
	}
}

func (b *BookingMicroService) SetShowService(shsrv func() ShowService.ShowService) {
	b.mu.Lock()
	b.ShowService = shsrv
	b.mu.Unlock()
}

func (b *BookingMicroService) SetHallService(hsrv func() HallService.HallService) {
	b.mu.Lock()
	b.HallService = hsrv
	b.mu.Unlock()
}

func (b *BookingMicroService) ResetBookings() {
	fmt.Println("-----Entered ResetBookings-----")
	b.mu.Lock()
	fmt.Println("Locked ResetBookings")

	for bkID, ele := range b.bookingRepository {
		if !ele.Confirmation.Confirmed && ele.Confirmation.time.After(ele.Confirmation.time.Add(time.Minute*time.Duration(MamboNumberFive))) {
			fmt.Printf("Booking expired: %d\n", bkID)
			fmt.Println("Freeing seats in show...")
			s := b.ShowService()
			message := &ShowService.FreeSeatMessage{
				ShowID:    ele.ShowID,
				BookingID: bkID,
			}

			_, err := s.FreeSeats(context.TODO(), message)
			if err != nil {
				fmt.Println("Freeing seats failed!")
				fmt.Println("Unlocked ResetBookings")
				b.mu.Unlock()
				fmt.Println("-----Exited ResetBookings-----")
			}
		}
	}

	fmt.Println("Unlocked ResetBookings")
	b.mu.Unlock()
	fmt.Println("-----Exited ResetBookings-----")
}

func (b *BookingMicroService) ConfirmBooking(ctx context.Context, req *BookingService.ConfirmBookingMessage, res *BookingService.ConfirmBookingResponse) error {
	fmt.Println("-----Entered ConfirmBooking-----")
	b.mu.Lock()
	fmt.Println("Locked ConfirmBooking")

	booking, ok := b.bookingRepository[req.BookingID]
	if !ok {
		fmt.Println("The booking does not exist!")
		fmt.Println("Unlocked ConfirmBooking")
		b.mu.Unlock()
		fmt.Println("-----Exited ConfirmBooking-----")
		return fmt.Errorf("the booking does not exist")
	}

	fmt.Println("Locking seats...")
	s := b.ShowService()

	message := &ShowService.LockSeatMessage{
		ShowID:    booking.ShowID,
		BookingID: req.BookingID,
	}

	bkg, err := s.LockSeats(ctx, message)
	if err != nil || !bkg.Success {
		fmt.Println("The booking was rejected.")
		fmt.Println("Unlocked ConfirmBooking")
		b.mu.Unlock()
		fmt.Println("-----Exited ConfirmBooking-----")
		return err
	}

	fmt.Println("Booking confirmed!")
	b.bookingRepository[req.BookingID].Confirmation.Confirmed = true
	b.bookingRepository[req.BookingID].Confirmation.time = time.Now()

	res.BookingID = req.BookingID

	fmt.Println("Unlocked ConfirmBooking")
	b.mu.Unlock()
	fmt.Println("-----Exited ConfirmBooking-----")
	return nil
}

func (b *BookingMicroService) CreateBooking(ctx context.Context, req *BookingService.CreateBookingMessage, res *BookingService.CreateBookingResponse) error {
	fmt.Println("-----Entered CreateBooking-----")
	b.mu.Lock()
	fmt.Println("Locked CreateBooking")

	fmt.Println("Blocking seats...")
	s := b.ShowService()

	message := &ShowService.BlockSeatMessage{
		BookingID: b.NextId,
		ShowID:    req.ShowID,
		SeatID:    req.Seats,
	}

	booking, err := s.BlockSeats(ctx, message)

	if err != nil {
		fmt.Println("The booking was rejected.")
		fmt.Println("Unlocked CreateBooking")
		b.mu.Unlock()
		fmt.Println("-----Exited CreateBooking-----")
		return err
	}

	if !booking.Success {
		fmt.Println("The booking was rejected.")
		fmt.Println("Unlocked CreateBooking")
		b.mu.Unlock()
		fmt.Println("-----Exited CreateBooking-----")
		return fmt.Errorf("the booking was rejected")
	}

	b.bookingRepository[b.NextId] = &Booking{
		UserID: req.UserID,
		ShowID: req.ShowID,
		Seats:  req.Seats,
		Confirmation: Confirmation{
			time:      time.Now(),
			Confirmed: false,
		},
	}
	res.BookingID = b.NextId

	b.NextId++
	fmt.Println("Increased NextID")

	for i, ele := range b.bookingRepository {
		fmt.Printf("Booking (ID: %d): Show: %d User: %d\n", i, ele.ShowID, ele.UserID)
	}

	fmt.Println("Unlocked CreateBooking")
	b.mu.Unlock()
	fmt.Println("-----Exited CreateBooking-----")
	return nil
}

func (b *BookingMicroService) DeleteElement(ctx context.Context, bookingID int32) error {
	fmt.Println("-----Entered DeleteElement-----")

	booking, ok := b.bookingRepository[bookingID]
	if !ok {
		fmt.Println("The booking does not exist.")
		fmt.Println("-----Exited DeleteElement-----")
		return fmt.Errorf("the booking does not exist")
	}
	fmt.Println("Found booking")

	fmt.Println("Freeing seats...")
	s := b.ShowService()

	message := &ShowService.FreeSeatMessage{
		ShowID:    booking.ShowID,
		BookingID: bookingID,
	}

	_, err := s.FreeSeats(ctx, message)
	if err != nil {
		fmt.Println("Error while freeing seats")
		fmt.Println("-----Exited DeleteElement-----")
		return err
	}

	delete(b.bookingRepository, bookingID)
	fmt.Println("Deleted booking")

	fmt.Println("-----Exited DeleteElement-----")
	return nil
}

func (b *BookingMicroService) DeleteBooking(context context.Context, req *BookingService.DeleteBookingMessage, _ *BookingService.DeleteBookingResponse) error {
	return b.DeleteElement(context, req.BookingID)
}

func (b *BookingMicroService) GetUserBookings(_ context.Context, req *BookingService.GetUserBookingsMessage, res *BookingService.GetUserBookingsResponse) error {
	fmt.Println("-----Entered GetUserBookings-----")
	b.mu.Lock()
	fmt.Println("Locked GetUserBookings")

	var bookings []int32

	for index, ele := range b.bookingRepository {
		if ele.UserID == req.UserID {
			bookings = append(bookings, index)
		}
	}

	res.BookingID = bookings
	res.UserID = req.UserID

	fmt.Println("Unlocked GetUserBookings")
	b.mu.Unlock()
	fmt.Println("-----Exited GetUserBookings-----")
	return nil
}

func (b *BookingMicroService) GetBooking(_ context.Context, req *BookingService.GetBookingMessage, res *BookingService.GetBookingResponse) error {
	fmt.Println("-----Entered GetBooking-----")
	b.mu.Lock()
	fmt.Println("Locked GetBooking")

	booking, ok := b.bookingRepository[req.BookingID]
	if ok {
		res.BookingID = req.BookingID
		res.UserID = booking.UserID
		res.ShowID = booking.ShowID
		res.Seats = booking.Seats

		fmt.Println("Unlocked GetBooking")
		b.mu.Unlock()
		fmt.Println("-----Exited GetBooking-----")
		return nil
	}
	fmt.Println("Booking not found")

	fmt.Println("Unlocked GetBooking")
	b.mu.Unlock()
	fmt.Println("-----Exited GetBooking-----")
	return fmt.Errorf("the booking does not exist")
}

func (b *BookingMicroService) KillBookingsShow(ctx context.Context, req *BookingService.KillBookingsShowMessage, res *BookingService.KillBookingsShowResponse) error {
	fmt.Println("-----Entered KillBookingsShow-----")
	b.mu.Lock()
	fmt.Println("Locked KillBookingsShow")

	fmt.Println("Deleting bookings...")
	for index, ele := range b.bookingRepository {
		if ele.ShowID == req.ShowID {
			err := b.DeleteElement(ctx, index)
			if err != nil {
				fmt.Println("Deleting booking failed")
				fmt.Println("Unlocked KillBookingsShow")
				b.mu.Unlock()
				fmt.Println("-----Exited KillBookingsShow-----")
				return err
			}
		}
	}
	fmt.Println("Deleted bookings")
	res.Success = true

	fmt.Println("Unlocked KillBookingsShow")
	b.mu.Unlock()
	fmt.Println("-----Exited KillBookingsShow-----")
	return nil
}

func (b *BookingMicroService) KillBookingsUser(ctx context.Context, req *BookingService.KillBookingsUserMessage, res *BookingService.KillBookingsUserResponse) error {
	fmt.Println("-----Entered KillBookingsUser-----")
	b.mu.Lock()
	fmt.Println("Locked KillBookingsUser")

	for index, ele := range b.bookingRepository {
		if ele.UserID == req.UserID {
			err := b.DeleteElement(ctx, index)
			if err != nil {
				fmt.Println("Deleting booking failed")
				res.Success = false
				fmt.Println("Unlocked KillBookingsShow")
				b.mu.Unlock()
				fmt.Println("-----Exited KillBookingsShow-----")
				return err
			}
		}
	}
	res.Success = true
	fmt.Println("Unlocked KillBookingsUser")
	b.mu.Unlock()
	fmt.Println("-----Exited KillBookingsUser-----")
	return nil
}
