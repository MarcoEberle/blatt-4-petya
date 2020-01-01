package UserService

import (
	"context"
	"fmt"
	UserService "github.com/ob-vss-ws19/blatt-4-petya/UserService/messages"
	"sync"
)

type User struct {
	userName string
}

type UserMicroService struct {
	userRepository map[int32]*User
	ResService     func() resproto.BookingService
	mu             *sync.RWMutex
	NextUserID     int32
}

func Spawn() *UserMicroService {
	return &UserMicroService{
		userRepository: make(map[int32]*User),
		ResService:     nil,
		mu:             &sync.RWMutex{},
		NextUserID:     1,
	}
}

func (usrv UserMicroService) CreateUser(context context.Context, req *UserService.CreateUserMessage, res *UserService.CreateUserResponse) error {
	if req.UserName != "" {
		usrv.mu.Lock()
		usrv.userRepository[usrv.NextUserID] = &User{userName: req.UserName}
		res.UserID = usrv.NextUserID
		usrv.NextUserID++
		defer usrv.mu.Unlock()
		return nil
	}

	return fmt.Errorf("The user could not be created.")
}

func (usrv UserMicroService) DeleteUser(context context.Context, req *UserService.DeleteUserMessage, res *UserService.DeleteUserResponse) error {
	res.Success = false
	_, storedUser := usrv.userRepository[req.UserID]

	if !storedUser {
		return fmt.Errorf("The user could not be deleted.")
	}

	if HasBookings(req.UserID) {
		usrv.mu.Lock()
		delete(usrv.userRepository, req.UserID)
		res.Success = true
		usrv.mu.Unlock()
	}

	return nil
}

func (usrv UserMicroService) GetUser(context context.Context, req *UserService.GetUserMessage, res *UserService.GetUserResponse) error {
	_, storedUser := usrv.userRepository[req.UserID]

	if storedUser {
		res.UserID = req.UserID
		res.UserName = usrv.userRepository[req.UserID].userName
	}

	return fmt.Errorf("The user could not be deleted.")
}

func HasBookings(userId int32) bool {

}
