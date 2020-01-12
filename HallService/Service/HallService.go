package Service

import (
	"context"
	"fmt"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
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
	mu             *sync.Mutex
	ShowService    func() ShowService.ShowService
}

func Spawn() *HallMicroService {
	return &HallMicroService{
		HallRepository: make(map[int32]*Hall),
		NextID:         1,
		mu:             &sync.Mutex{},
	}
}

func (hsrv *HallMicroService) CreateHall(context context.Context, req *HallService.CreateHallMessage, res *HallService.CreateHallResponse) error {
	hsrv.mu.Lock()
	hsrv.HallRepository[hsrv.NextID] = &Hall{
		hallName:    req.HallName,
		rows:        req.Rows,
		seatsPerRow: req.SeatsPerRow,
	}
	res.HallID = hsrv.NextID

	hsrv.NextID++
	hsrv.mu.Unlock()

	fmt.Printf("Created hall: %d %s", res.HallID, req.HallName)
	return nil
}

func (hsrv *HallMicroService) DeleteHall(context context.Context, req *HallService.DeleteHallMessage, res *HallService.DeleteHallResponse) error {
	hsrv.mu.Lock()
	res.Success = false
	_, hall := hsrv.HallRepository[req.HallID]
	if hall {
		delete(hsrv.HallRepository, req.HallID)

		s := hsrv.ShowService()

		message := &ShowService.KillShowsHallMessage{
			HallID: req.HallID,
		}

		_, err := s.KillShowsHall(context, message)
		if err != nil {
			res.Success = false
			return err
		}

		res.Success = true
		hsrv.mu.Unlock()
		return nil
	}

	hsrv.mu.Unlock()
	return fmt.Errorf("The hall could not be deleted.")
}

func (hsrv *HallMicroService) GetHall(context context.Context, req *HallService.GetHallMessage, res *HallService.GetHallResponse) error {
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

	for i, ele := range hsrv.HallRepository {
		fmt.Printf("%d: %s %dx%d\n", i, ele.hallName, ele.rows, ele.seatsPerRow)
	}

	res.HallID = 0
	hsrv.mu.Unlock()
	return fmt.Errorf("The hall could not be found.")
}

func (hsrv *HallMicroService) VerifySeat(context context.Context, req *HallService.VerifySeatMessage, res *HallService.VerifySeatResponse) error {
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

func (hsrv *HallMicroService) SetBookingService(shsrv func() ShowService.ShowService) {
	hsrv.mu.Lock()
	hsrv.ShowService = shsrv
	hsrv.mu.Unlock()
}

func (hsrv *HallMicroService) SetShowService(ssrv func() ShowService.ShowService) {
	hsrv.mu.Lock()
	hsrv.ShowService = ssrv
	hsrv.mu.Unlock()
}
