package HallService

import (
	"context"
	"fmt"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/messages"
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
		res.Success = true
		return nil
	}

	return fmt.Errorf("The hall could not be deleted.")
}

func (hsrv HallMicroService) GetHall(context context.Context, req *HallService.GetHallMessage, res *HallService.GetHallResponse) error {
	_, hall := hsrv.HallRepository[req.HallID]
	if hall {
		h := hsrv.HallRepository[req.HallID]
		res.HallID = req.HallID
		res.HallName = h.hallName
		res.SeatsPerRow = h.seatsPerRow
		res.Rows = h.rows
		return nil
	}

	return fmt.Errorf("The hall could not be found.")
}

func (hsrv HallMicroService) VerifySeat(context context.Context, req *HallService.VerifySeatMessage, res *HallService.VerifySeatResponse) error {
	_, hall := hsrv.HallRepository[req.HallID]
	if hall {
		h := hsrv.HallRepository[req.HallID]

		res.Success = req.SeatID <= h.rows*h.seatsPerRow
		return nil
	}

	return fmt.Errorf("The hall could not be found.")
}
