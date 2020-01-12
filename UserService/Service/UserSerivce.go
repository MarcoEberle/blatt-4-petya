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

func Spawn() *UserMicroService {
	return &UserMicroService{
		userRepository: make(map[int32]*User),
		BookingService: nil,
		mu:             &sync.Mutex{},
		NextUserID:     1,
	}
}

func (usrv *UserMicroService) CreateUser(context context.Context, req *UserService.CreateUserMessage, res *UserService.CreateUserResponse) error {
	fmt.Println("-----Entered CreateUser-----")
	if req.UserName != "" {
		usrv.mu.Lock()
		fmt.Println("Locked CreateUser")
		usrv.userRepository[usrv.NextUserID] = &User{userName: req.UserName}
		res.UserID = usrv.NextUserID
		fmt.Println("Added CreateUser")
		usrv.NextUserID++
		fmt.Println("Increased ID")
		defer usrv.mu.Unlock()
		fmt.Println("Unlocked CreateUser")

		fmt.Printf("Created user: %d %s", res.UserID, req.UserName)
		return nil
	}
	fmt.Println("Username is empty!")
	fmt.Println("-----Exited CreateUser-----")
	return fmt.Errorf("The user could not be created.")
}

func (usrv *UserMicroService) DeleteUser(context context.Context, req *UserService.DeleteUserMessage, res *UserService.DeleteUserResponse) error {
	res.Success = false
	_, storedUser := usrv.userRepository[req.UserID]

	if !storedUser {
		return fmt.Errorf("The user could not be deleted.")
	}

	b := usrv.BookingService()

	message := &BookingService.GetUserBookingsMessage{
		UserID: req.UserID,
	}

	ele, _ := b.GetUserBookings(context, message)
	if len(ele.BookingID) != 0 {
		deleteMessage := &BookingService.KillBookingsUserMessage{
			UserID: req.UserID,
		}
		b.KillBookingsUser(context, deleteMessage)
	}

	delete(usrv.userRepository, req.UserID)
	res.Success = true
	usrv.mu.Unlock()
	return nil
}

func (usrv *UserMicroService) GetUser(context context.Context, req *UserService.GetUserMessage, res *UserService.GetUserResponse) error {
	_, storedUser := usrv.userRepository[req.UserID]

	if storedUser {
		res.UserID = req.UserID
		res.UserName = usrv.userRepository[req.UserID].userName
	}

	return fmt.Errorf("The user could not be deleted.")
}

func (usrv *UserMicroService) SetBookingService(bksrv func() BookingService.BookingService) {
	usrv.BookingService = bksrv
}
