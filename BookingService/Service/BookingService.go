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
	mu                *sync.RWMutex
}

func Spawn() *BookingMicroService {
	return &BookingMicroService{
		bookingRepository: make(map[int32]*Booking),
		NextId:            1,
		mu:                &sync.RWMutex{},
	}
}

func (bksrv BookingMicroService) SetShowService(shsrv func() ShowService.ShowService) {
	bksrv.mu.Lock()
	bksrv.ShowService = shsrv
	bksrv.mu.Unlock()
}

func (bksrv BookingMicroService) SetHallService(hsrv func() HallService.HallService) {
	bksrv.mu.Lock()
	bksrv.HallService = hsrv
	bksrv.mu.Unlock()
}

func (bksrv *BookingMicroService) ResetBookings() {
	bksrv.mu.Lock()
	for bkID, ele := range bksrv.bookingRepository {
		if !ele.Confirmation.Confirmed && ele.Confirmation.time.After(ele.Confirmation.time.Add(time.Minute*5)) {
			s := bksrv.ShowService()
			message := &ShowService.FreeSeatMessage{
				ShowID:    ele.ShowID,
				BookingID: bkID,
			}

			s.FreeSeats(nil, message)
		}
	}
	bksrv.mu.Unlock()
}

func (bksrv *BookingMicroService) ConfirmBooking(context context.Context, req *BookingService.ConfirmBookingMessage, res *BookingService.ConfirmBookingResponse) error {
	bksrv.ResetBookings()

	bksrv.mu.Lock()
	s := bksrv.ShowService()

	booking, ok := bksrv.bookingRepository[req.BookingID]
	if !ok {
		bksrv.mu.Unlock()
		return fmt.Errorf("The booking does not exist.")
	}

	message := &ShowService.LockSeatMessage{
		ShowID:    booking.ShowID,
		BookingID: req.BookingID,
	}

	bkg, _ := s.LockSeats(context, message)
	if !bkg.Success {
		bksrv.mu.Unlock()
		return fmt.Errorf("The booking was rejected.")
	}

	bksrv.bookingRepository[req.BookingID].Confirmation.Confirmed = true
	bksrv.bookingRepository[req.BookingID].Confirmation.time = time.Now()

	bksrv.mu.Unlock()
	return nil
}

func (bksrv *BookingMicroService) CreateBooking(context context.Context, req *BookingService.CreateBookingMessage, res *BookingService.CreateBookingResponse) error {
	bksrv.ResetBookings()

	bksrv.mu.Lock()
	s := bksrv.ShowService()

	message := &ShowService.BlockSeatMessage{
		BookingID: bksrv.NextId,
		ShowID:    req.ShowID,
	}

	booking, _ := s.BlockSeats(context, message)
	if !booking.Success {
		bksrv.mu.Unlock()
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

	bksrv.NextId++

	bksrv.mu.Unlock()

	return nil
}

func (bksrv *BookingMicroService) DeleteElement(context context.Context, bookingID int32) error {
	bksrv.mu.Lock()
	s := bksrv.ShowService()

	booking, ok := bksrv.bookingRepository[bookingID]
	if !ok {
		bksrv.mu.Unlock()
		return fmt.Errorf("The booking does not exist.")
	}

	message := &ShowService.FreeSeatMessage{
		ShowID:    booking.ShowID,
		BookingID: bookingID,
	}

	s.FreeSeats(context, message)

	delete(bksrv.bookingRepository, bookingID)

	bksrv.mu.Unlock()

	return nil
}

func (bksrv *BookingMicroService) DeleteBooking(context context.Context, req *BookingService.DeleteBookingMessage, res *BookingService.DeleteBookingResponse) error {
	return bksrv.DeleteElement(context, req.BookingID)
}

func (bksrv *BookingMicroService) GetUserBookings(context context.Context, req *BookingService.GetUserBookingsMessage, res *BookingService.GetUserBookingsResponse) error {
	bksrv.mu.Lock()
	var bookings []int32

	for index, ele := range bksrv.bookingRepository {
		if ele.UserID == req.UserID {
			bookings = append(bookings, index)
		}
	}

	res.BookingID = bookings
	res.UserID = req.UserID

	bksrv.mu.Unlock()
	return nil
}

func (bksrv *BookingMicroService) GetBooking(context context.Context, req *BookingService.GetBookingMessage, res *BookingService.GetBookingResponse) error {
	bksrv.mu.Lock()
	booking, ok := bksrv.bookingRepository[req.BookingID]
	if ok {
		res.BookingID = req.BookingID
		res.UserID = booking.UserID
		res.ShowID = booking.ShowID
		res.Seats = booking.Seats

		bksrv.mu.Unlock()
		return nil
	}

	bksrv.mu.Unlock()
	return fmt.Errorf("The booking does not exist.")
}

func (bksrv *BookingMicroService) KillBookingsShow(context context.Context, req *BookingService.KillBookingsShowMessage, res *BookingService.KillBookingsShowResponse) error {
	bksrv.mu.Lock()

	for index, ele := range bksrv.bookingRepository {
		if ele.ShowID == req.ShowID {
			bksrv.DeleteElement(context, index)
		}
	}

	bksrv.mu.Unlock()
	return nil
}

func (bksrv *BookingMicroService) KillBookingsUser(context context.Context, req *BookingService.KillBookingsUserMessage, res *BookingService.KillBookingsUserResponse) error {
	bksrv.mu.Lock()

	for index, ele := range bksrv.bookingRepository {
		if ele.UserID == req.UserID {
			bksrv.DeleteElement(context, index)
		}
	}

	bksrv.mu.Unlock()
	return nil
}
