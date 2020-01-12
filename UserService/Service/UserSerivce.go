package Service

import (
	"context"
	"fmt"
	BookingService "github.com/ob-vss-ws19/blatt-4-petya/BookingService/Service/messages"
	UserService "github.com/ob-vss-ws19/blatt-4-petya/UserService/Service/messages"
	"sync"
)

type User struct {
	userName string
}

type UserMicroService struct {
	userRepository map[int32]*User
	BookingService func() BookingService.BookingService
	mu             *sync.Mutex
	NextUserID     int32
}

const (
	PlayerNumberOne int32 = 1
)

func Spawn() *UserMicroService {
	return &UserMicroService{
		userRepository: make(map[int32]*User),
		BookingService: nil,
		mu:             &sync.Mutex{},
		NextUserID:     PlayerNumberOne,
	}
}

func (usrv *UserMicroService) CreateUser(context context.Context, req *UserService.CreateUserMessage, res *UserService.CreateUserResponse) error {
	fmt.Println("-----Entered CreateUser-----")
	if req.UserName != "" {
		usrv.mu.Lock()
		fmt.Println("Locked CreateUser")
		usrv.userRepository[usrv.NextUserID] = &User{userName: req.UserName}
		res.UserID = usrv.NextUserID
		fmt.Printf("Created user: %d %s\n", res.UserID, req.UserName)
		usrv.NextUserID++
		fmt.Println("Increased ID")
		defer usrv.mu.Unlock()
		fmt.Println("Unlocked CreateUser")

		fmt.Println("-----Exited CreateUser-----")
		return nil
	}
	fmt.Println("Username is empty!")
	fmt.Println("-----Exited CreateUser-----")
	return fmt.Errorf("The user could not be created.")
}

func (usrv *UserMicroService) DeleteUser(context context.Context, req *UserService.DeleteUserMessage, res *UserService.DeleteUserResponse) error {
	fmt.Println("-----Entered DeleteUser-----")
	usrv.mu.Lock()
	fmt.Println("Locked DeleteUser")
	res.Success = false
	_, storedUser := usrv.userRepository[req.UserID]

	if !storedUser {
		fmt.Println("User not found!")
		usrv.mu.Unlock()
		fmt.Println("Unlocked DeleteUser!")
		fmt.Println("-----Exited DeleteUser-----")
		return fmt.Errorf("The user could not be deleted.")
	}

	b := usrv.BookingService()

	fmt.Println("Delete users bookings...")
	message := &BookingService.GetUserBookingsMessage{
		UserID: req.UserID,
	}

	ele, _ := b.GetUserBookings(context, message)
	if len(ele.BookingID) != 0 {
		deleteMessage := &BookingService.KillBookingsUserMessage{
			UserID: req.UserID,
		}
		_, err := b.KillBookingsUser(context, deleteMessage)
		if err != nil {
			fmt.Println("Error while deleting users bookings!")
			usrv.mu.Unlock()
			fmt.Println("Unlocked DeleteUser")
			fmt.Println("-----Exited DeleteUser-----")
			return err
		}
	}

	delete(usrv.userRepository, req.UserID)
	fmt.Println("Deleted user")
	res.Success = true
	usrv.mu.Unlock()
	fmt.Println("Unlocked DeleteUser")
	fmt.Println("-----Exited DeleteUser-----")
	return nil
}

func (usrv *UserMicroService) GetUser(context context.Context, req *UserService.GetUserMessage, res *UserService.GetUserResponse) error {
	fmt.Println("-----Entered GetUser-----")
	_, storedUser := usrv.userRepository[req.UserID]

	if storedUser {
		fmt.Println("Found User")
		res.UserID = req.UserID
		res.UserName = usrv.userRepository[req.UserID].userName
		fmt.Println("-----Exited GetUser-----")
		return nil
	}
	fmt.Println("User not found!")
	fmt.Println("-----Exited GetUser-----")
	return fmt.Errorf("The user could not be found.")
}

func (usrv *UserMicroService) SetBookingService(bksrv func() BookingService.BookingService) {
	usrv.BookingService = bksrv
}
