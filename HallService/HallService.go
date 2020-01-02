package HallService

import (
	"context"
	"fmt"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/messages"
	"sync"
)

type Hall struct {
	hallName    string
	rows        int32
	seatsPerRow int32
}

type HallMicroService struct {
	HallRepository map[int32]*Hall
	NextID         int32
	mu             *sync.RWMutex
	ShowService    func() ShowService.ShowService
}

func Spawn() *HallMicroService {
	return &HallMicroService{
		HallRepository: make(map[int32]*Hall),
		NextID:         1,
		mu:             &sync.RWMutex{},
	}
}

func (hsrv HallMicroService) CreateHall(context context.Context, req *HallService.CreateHallMessage, res *HallService.CreateHallResponse) error {
	hsrv.mu.Lock()
	hsrv.HallRepository[hsrv.NextID] = &Hall{
		hallName:    req.HallName,
		rows:        req.Rows,
		seatsPerRow: req.SeatsPerRow,
	}
	res.HallID = hsrv.NextID

	hsrv.NextID++
	hsrv.mu.Unlock()

	return nil
}

func (hsrv HallMicroService) DeleteHall(context context.Context, req *HallService.DeleteHallMessage, res *HallService.DeleteHallResponse) error {
	hsrv.mu.Lock()
	res.Success = false
	_, hall := hsrv.HallRepository[req.HallID]
	if hall {
		delete(hsrv.HallRepository, req.HallID)

		s := hsrv.ShowService()

		message := &ShowService.KillShowsMessage{
			HallID: req.HallID,
		}

		s.KillShows(context, message)

		res.Success = true
		hsrv.mu.Unlock()
		return nil
	}

	hsrv.mu.Unlock()
	return fmt.Errorf("The hall could not be deleted.")
}

func (hsrv HallMicroService) GetHall(context context.Context, req *HallService.GetHallMessage, res *HallService.GetHallResponse) error {
	hsrv.mu.Lock()
	_, hall := hsrv.HallRepository[req.HallID]
	if hall {
		h := hsrv.HallRepository[req.HallID]
		res.HallID = req.HallID
		res.HallName = h.hallName
		res.SeatsPerRow = h.seatsPerRow
		res.Rows = h.rows
		hsrv.mu.Unlock()
		return nil
	}

	hsrv.mu.Unlock()
	return fmt.Errorf("The hall could not be found.")
}

func (hsrv HallMicroService) VerifySeat(context context.Context, req *HallService.VerifySeatMessage, res *HallService.VerifySeatResponse) error {
	hsrv.mu.Lock()
	res.Success = false
	_, hall := hsrv.HallRepository[req.HallID]
	if hall {
		h := hsrv.HallRepository[req.HallID]

		for _, ele := range req.SeatID {
			if ele > h.rows*h.seatsPerRow {
				hsrv.mu.Unlock()
				return fmt.Errorf("The seats are not existing.")
			}
		}
	} else {
		hsrv.mu.Unlock()
		return fmt.Errorf("The hall could not be found.")
	}

	res.Success = true
	hsrv.mu.Unlock()
	return nil
}

func (usrv HallMicroService) SetBookingService(shsrv func() ShowService.ShowService) {
	usrv.mu.Lock()
	usrv.ShowService = shsrv
	usrv.mu.Unlock()
}
