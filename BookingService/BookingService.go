package BookingService

import (
	"context"
	"github.com/gogo/protobuf/types"
	BookingService "github.com/ob-vss-ws19/blatt-4-petya/BookingService/messages"
	"sync"
)

type Booking struct {
	UserID       int32
	ShowID       int32
	Seats        []int32
	Confirmation Confirmation
}

type Confirmation struct {
	time      types.Timestamp
	Confirmed bool
}

type BookingMicroService struct {
	bookingRepository map[int32]*Booking
	HallService       func() HallService.HallService
	ShowService       func() ShowService.ShowService
	mu                *sync.RWMutex
}

func (bksrv *BookingMicroService) PrepareBooking(context context.Context, req *BookingService.PrepareBookingMessage, res *BookingService.PrepareBookingResponse) error {
	// Prüfung Sitzplätze passen

	// Prüfung Sitze frei

	// Blocke Sitze
}

func (bksrv *BookingMicroService) CreateBooking(context context.Context, req *BookingService.CreateBookingMessage, res *BookingService.CreateBookingResponse) error {
	// Prüfe Sitzplätze-Block

	// Safe Sitzplätze
}

func (bksrv *BookingMicroService) DeleteBooking(context context.Context, req *BookingService.DeleteBookingMessage, res *BookingService.DeleteBookingResponse) error {
	// Check Blocked oder Safe

	// Delete Booking
}

func (bksrv *BookingMicroService) GetUserBookings(context context.Context, req *BookingService.GetUserBookingsMessage, res *BookingService.GetUserBookingsResponse) error {
	// Check Blocked oder Safe

}

func (bksrv *BookingMicroService) GetBookings(context context.Context, req *BookingService.GetBookingMessage, res *BookingService.GetBookingResponse) error {
	// Check Blocked oder Safe

}
