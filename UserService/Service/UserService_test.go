package Service

import (
	"context"
	"fmt"
	UserService "github.com/ob-vss-ws19/blatt-4-petya/UserService/Service/messages"
	"testing"
)

const (
	User1 = "Loksey"
)

func TestCreateUser(t *testing.T) {
	service := Spawn()
	r := UserService.CreateUserResponse{}
	er := service.CreateUser(context.TODO(), &UserService.CreateUserMessage{
		UserName: User1,
	}, &r)

	if er == nil {
		if r.UserID > 0 {
			t.Log("Successfully created user.")
		}
	} else {
		fmt.Println(er)
	}
}

func TestGetUser(t *testing.T) {
	service := Spawn()
	r := UserService.CreateUserResponse{}
	err := service.CreateUser(context.TODO(), &UserService.CreateUserMessage{
		UserName: User1,
	}, &r)

	if err == nil {
		fmt.Println(err)
	}

	rr := UserService.GetUserResponse{}
	er := service.GetUser(context.TODO(), &UserService.GetUserMessage{
		UserID: r.UserID,
	}, &rr)

	if er == nil {
		if r.UserID > 0 {
			t.Log("Successfully got user.")
		}
	} else {
		fmt.Println(er)
	}
}
