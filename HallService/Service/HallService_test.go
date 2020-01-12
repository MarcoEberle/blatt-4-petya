package Service

import (
	"context"
	"fmt"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	"testing"
)

const (
	MightyNumber4 int32 = 4
)

func TestCreateHall(t *testing.T) {
	service := Spawn()
	r := HallService.CreateHallResponse{}
	er := service.CreateHall(context.TODO(), &HallService.CreateHallMessage{
		HallName:    "Halls of Stone",
		Rows:        MightyNumber4,
		SeatsPerRow: MightyNumber4,
	}, &r)

	if er == nil {
		if r.HallID > 0 {
			t.Log("Successfully created hall.")
		}
	} else {
		fmt.Println(er)
	}
}

func TestGetHall(t *testing.T) {
	service := Spawn()
	r := HallService.CreateHallResponse{}
	err := service.CreateHall(context.TODO(), &HallService.CreateHallMessage{
		HallName:    "Halls of Stone",
		Rows:        MightyNumber4,
		SeatsPerRow: MightyNumber4,
	}, &r)

	if err != nil {
		fmt.Println(err)
	}

	rr := HallService.GetHallResponse{}
	er := service.GetHall(context.TODO(), &HallService.GetHallMessage{
		HallID: r.HallID,
	}, &rr)

	if er == nil {
		if r.HallID > 0 {
			t.Log("Successfully got hall.")
		}
	} else {
		fmt.Println(er)
	}
}

func TestVerifySeatl(t *testing.T) {
	service := Spawn()
	r := HallService.CreateHallResponse{}
	err := service.CreateHall(context.TODO(), &HallService.CreateHallMessage{
		HallName:    "Halls of Stone",
		Rows:        MightyNumber4,
		SeatsPerRow: MightyNumber4,
	}, &r)

	if err != nil {
		fmt.Println(err)
	}

	rr := HallService.VerifySeatResponse{}
	er := service.VerifySeat(context.TODO(), &HallService.VerifySeatMessage{
		SeatID: []int32{1, 5, 16},
	}, &rr)

	if er == nil {
		if r.HallID > 0 {
			t.Log("Successfully verified seat.")
		}
	} else {
		fmt.Println(er)
	}
}
