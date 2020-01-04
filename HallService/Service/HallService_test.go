package Service

import (
	"context"
	"fmt"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	"testing"
)

func TestCreateHall(t *testing.T) {
	service := Spawn()
	r := HallService.CreateHallResponse{}
	er := service.CreateHall(context.TODO(), &HallService.CreateHallMessage{
		HallName:    "Halls of Stone",
		Rows:        4,
		SeatsPerRow: 4,
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
	service.CreateHall(context.TODO(), &HallService.CreateHallMessage{
		HallName:    "Halls of Stone",
		Rows:        4,
		SeatsPerRow: 4,
	}, &r)

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
	service.CreateHall(context.TODO(), &HallService.CreateHallMessage{
		HallName:    "Halls of Stone",
		Rows:        4,
		SeatsPerRow: 4,
	}, &r)

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
