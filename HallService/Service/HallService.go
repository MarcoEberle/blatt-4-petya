package Service

import (
	"context"
	"fmt"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
	"sync"
	"time"
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
	fmt.Println("-----Entered CreateHall-----")
	hsrv.mu.Lock()
	fmt.Println("Locked CreateHall")
	hsrv.HallRepository[hsrv.NextID] = &Hall{
		hallName:    req.HallName,
		rows:        req.Rows,
		seatsPerRow: req.SeatsPerRow,
	}
	fmt.Println("Created Hall")
	res.HallID = hsrv.NextID
	fmt.Println("Increased NextID")
	hsrv.NextID++
	hsrv.mu.Unlock()
	fmt.Println("Unlocked CreateHall")

	fmt.Printf("Created hall: %d %s\n", res.HallID, req.HallName)
	fmt.Println("-----Exited CreateHall-----")
	return nil
}

func (hsrv *HallMicroService) DeleteHall(ctx context.Context, req *HallService.DeleteHallMessage, res *HallService.DeleteHallResponse) error {
	fmt.Println("-----Entered DeleteHall-----")
	hsrv.mu.Lock()
	fmt.Println("Locked DeleteHall")
	res.Success = false
	_, hall := hsrv.HallRepository[req.HallID]
	if hall {
		fmt.Println("Found Hall")
		delete(hsrv.HallRepository, req.HallID)

		fmt.Println("Delete halls shows.....")
		s := hsrv.ShowService()

		message := &ShowService.KillShowsHallMessage{
			HallID: req.HallID,
		}
		const timeout = 200 * time.Second
		ctxLong, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := s.KillShowsHall(ctxLong, message)
		if err != nil {
			fmt.Println("Error while deleting halls shows!")
			res.Success = false
			hsrv.mu.Unlock()
			fmt.Println("Unlocked DeleteHall")
			fmt.Println("-----Exited DeleteHall-----")
			return err
		}
		fmt.Println("Deleted hall!")
		res.Success = true
		hsrv.mu.Unlock()
		fmt.Println("Unlocked DeleteHall")
		fmt.Println("-----Exited DeleteHall-----")
		return nil
	}

	hsrv.mu.Unlock()
	fmt.Println("Unlocked DeleteHall")
	fmt.Println("-----Exited DeleteHall-----")
	return fmt.Errorf("The hall could not be found.")
}

func (hsrv *HallMicroService) GetHall(context context.Context, req *HallService.GetHallMessage, res *HallService.GetHallResponse) error {
	fmt.Println("-----Entered GetHall-----")
	fmt.Println("Locked GetHall")
	hsrv.mu.Lock()

	_, hall := hsrv.HallRepository[req.HallID]
	if hall {
		fmt.Println("Found Hall")
		h := hsrv.HallRepository[req.HallID]
		res.HallID = req.HallID
		res.HallName = h.hallName
		res.SeatsPerRow = h.seatsPerRow
		res.Rows = h.rows

		hsrv.mu.Unlock()
		fmt.Println("Unlocked GetHall")
		fmt.Println("-----Exited GetHall-----")
		return nil
	}

	res.HallID = 0
	hsrv.mu.Unlock()
	fmt.Println("Unlocked GetHall")
	fmt.Println("-----Exited GetHall-----")
	return fmt.Errorf("The hall could not be found.")
}

func (hsrv *HallMicroService) VerifySeat(context context.Context, req *HallService.VerifySeatMessage, res *HallService.VerifySeatResponse) error {
	fmt.Println("-----Entered VerifySeat-----")
	hsrv.mu.Lock()
	fmt.Println("Locked VerifySeat")
	res.Success = false

	_, hall := hsrv.HallRepository[req.HallID]
	if hall {
		fmt.Println("Found hall")
		h := hsrv.HallRepository[req.HallID]

		for _, ele := range req.SeatID {
			if ele > h.rows*h.seatsPerRow {
				fmt.Println("The seats are not existing!")
				hsrv.mu.Unlock()
				fmt.Println("Unlocked VerifySeat")
				fmt.Println("-----Exited VerifySeat-----")
				return fmt.Errorf("The seats are not existing.")
			}
		}
	} else {
		hsrv.mu.Unlock()
		fmt.Println("Unlocked VerifySeat")
		fmt.Println("-----Exited VerifySeat-----")
		return fmt.Errorf("The hall could not be found.")
	}

	res.Success = true
	hsrv.mu.Unlock()
	fmt.Println("Unlocked VerifySeat")
	fmt.Println("-----Exited VerifySeat-----")
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
