package Service

import (
	"context"
	"fmt"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
	"sync"
)

const (
	PlayerNumberOne int32 = 1
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
		NextID:         PlayerNumberOne,
		mu:             &sync.Mutex{},
	}
}

func (h *HallMicroService) CreateHall(_ context.Context, req *HallService.CreateHallMessage, res *HallService.CreateHallResponse) error {
	fmt.Println("-----Entered CreateHall-----")
	h.mu.Lock()
	fmt.Println("Locked CreateHall")
	h.HallRepository[h.NextID] = &Hall{
		hallName:    req.HallName,
		rows:        req.Rows,
		seatsPerRow: req.SeatsPerRow,
	}
	fmt.Println("Created Hall")
	res.HallID = h.NextID
	fmt.Println("Increased NextID")
	h.NextID++
	h.mu.Unlock()
	fmt.Println("Unlocked CreateHall")

	fmt.Printf("Created hall: %d %s\n", res.HallID, req.HallName)
	fmt.Println("-----Exited CreateHall-----")
	return nil
}

func (h *HallMicroService) DeleteHall(ctx context.Context, req *HallService.DeleteHallMessage, res *HallService.DeleteHallResponse) error {
	fmt.Println("-----Entered DeleteHall-----")
	h.mu.Lock()
	fmt.Println("Locked DeleteHall")
	res.Success = false
	_, hall := h.HallRepository[req.HallID]
	if hall {
		fmt.Println("Found Hall")
		delete(h.HallRepository, req.HallID)

		fmt.Println("Delete halls shows.....")
		s := h.ShowService()

		message := &ShowService.KillShowsHallMessage{
			HallID: req.HallID,
		}

		_, err := s.KillShowsHall(ctx, message)
		if err != nil {
			fmt.Println("Error while deleting halls shows!")
			res.Success = false
			h.mu.Unlock()
			fmt.Println("Unlocked DeleteHall")
			fmt.Println("-----Exited DeleteHall-----")
			return err
		}
		fmt.Println("Deleted hall!")
		res.Success = true
		h.mu.Unlock()
		fmt.Println("Unlocked DeleteHall")
		fmt.Println("-----Exited DeleteHall-----")
		return nil
	}

	h.mu.Unlock()
	fmt.Println("Unlocked DeleteHall")
	fmt.Println("-----Exited DeleteHall-----")
	return fmt.Errorf("the hall could not be found")
}

func (h *HallMicroService) GetHall(_ context.Context, req *HallService.GetHallMessage, res *HallService.GetHallResponse) error {
	fmt.Println("-----Entered GetHall-----")
	fmt.Println("Locked GetHall")
	h.mu.Lock()

	_, hall := h.HallRepository[req.HallID]
	if hall {
		fmt.Println("Found Hall")
		found := h.HallRepository[req.HallID]
		res.HallID = req.HallID
		res.HallName = found.hallName
		res.SeatsPerRow = found.seatsPerRow
		res.Rows = found.rows

		h.mu.Unlock()
		fmt.Println("Unlocked GetHall")
		fmt.Println("-----Exited GetHall-----")
		return nil
	}

	res.HallID = 0
	h.mu.Unlock()
	fmt.Println("Unlocked GetHall")
	fmt.Println("-----Exited GetHall-----")
	return fmt.Errorf("the hall could not be found")
}

func (h *HallMicroService) VerifySeat(_ context.Context, req *HallService.VerifySeatMessage, res *HallService.VerifySeatResponse) error {
	fmt.Println("-----Entered VerifySeat-----")
	h.mu.Lock()
	fmt.Println("Locked VerifySeat")
	res.Success = false

	_, hall := h.HallRepository[req.HallID]
	if hall {
		fmt.Println("Found hall")
		found := h.HallRepository[req.HallID]

		for _, ele := range req.SeatID {
			if ele > found.rows*found.seatsPerRow {
				fmt.Println("The seats are not existing!")
				h.mu.Unlock()
				fmt.Println("Unlocked VerifySeat")
				fmt.Println("-----Exited VerifySeat-----")
				return fmt.Errorf("the seats are not existing")
			}
		}
	} else {
		h.mu.Unlock()
		fmt.Println("Unlocked VerifySeat")
		fmt.Println("-----Exited VerifySeat-----")
		return fmt.Errorf("the hall could not be found")
	}

	res.Success = true
	h.mu.Unlock()
	fmt.Println("Unlocked VerifySeat")
	fmt.Println("-----Exited VerifySeat-----")
	return nil
}

func (h *HallMicroService) SetBookingService(sh func() ShowService.ShowService) {
	h.mu.Lock()
	h.ShowService = sh
	h.mu.Unlock()
}

func (h *HallMicroService) SetShowService(ssrv func() ShowService.ShowService) {
	h.mu.Lock()
	h.ShowService = ssrv
	h.mu.Unlock()
}
