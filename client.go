package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	BookingService "github.com/ob-vss-ws19/blatt-4-petya/BookingService/Service/messages"
	HallService "github.com/ob-vss-ws19/blatt-4-petya/HallService/Service/messages"
	MovieService "github.com/ob-vss-ws19/blatt-4-petya/MovieService/Service/messages"
	ShowService "github.com/ob-vss-ws19/blatt-4-petya/ShowService/Service/messages"
	UserService "github.com/ob-vss-ws19/blatt-4-petya/UserService/Service/messages"
	"os"
	"time"
)

func main() {
	clientService := micro.NewService(micro.Name("Client"))
	clientService.Init()

	/////////////////////////////////////////
	// Create Movies
	/////////////////////////////////////////
	movieService := MovieService.NewMovieService("MovieService", clientService.Client())

	movies := []int32{}
	movieNames := []string{
		"The Deadmines", "Scarlet Monastery", "Uldaman",
		"The Tempel of Atal'Hakkar",
	}

	for i := 0; i < 4; i++ {
		movies = append(movies, createMovie(movieNames[i], movieService))
	}

	/////////////////////////////////////////
	// Create Halls
	/////////////////////////////////////////
	hallService := HallService.NewHallService("HallService", clientService.Client())

	halls := []int32{}
	hallNames := []string{
		"Halls of Stone", "Halls of Valor",
	}

	for i := 0; i < 2; i++ {
		halls = append(halls, createHall(hallNames[i], 10, 10, hallService))
	}

	/////////////////////////////////////////
	// Create Users
	/////////////////////////////////////////
	userService := UserService.NewUserService("UserService", clientService.Client())

	users := []int32{}
	userNames := []string{
		"Bob", "Alice", "John", "Martin",
	}

	for i := 0; i < 4; i++ {
		users = append(users, createUser(userNames[i], userService))
	}

	/////////////////////////////////////////
	// Create Show
	/////////////////////////////////////////

	showService := ShowService.NewShowService("ShowService", clientService.Client())

	shows := []int32{}

	shows = append(shows, createShow(halls[0], movies[0], showService))
	shows = append(shows, createShow(halls[0], movies[1], showService))
	shows = append(shows, createShow(halls[1], movies[2], showService))
	shows = append(shows, createShow(halls[1], movies[3], showService))

	/////////////////////////////////////////
	// Create Booking
	/////////////////////////////////////////

	bookingService := BookingService.NewBookingService("BookingService", clientService.Client())

	bookings := []int32{}

	bookings = append(bookings, createBooking(
		shows[0], users[0], []int32{10, 11, 12}, bookingService))

	createBooking(shows[0], users[1], []int32{10, 11, 12}, bookingService)

	bookings = append(bookings, createBooking(
		shows[1], users[1], []int32{13, 14, 15}, bookingService))

	bookings = append(bookings, createBooking(
		shows[2], users[2], []int32{1, 2, 3}, bookingService))

	bookings = append(bookings, createBooking(
		shows[3], users[3], []int32{4, 5, 6}, bookingService))

	// Confirm Booking

	confirmedBookings := []int32{}

	confirmedBookings = append(confirmedBookings, confirmBooking(
		bookings[0], users[0], bookingService))

	confirmedBookings = append(confirmedBookings, confirmBooking(
		bookings[1], users[1], bookingService))

	confirmedBookings = append(confirmedBookings, confirmBooking(
		bookings[2], users[2], bookingService))

	confirmedBookings = append(confirmedBookings, confirmBooking(
		bookings[3], users[3], bookingService))

	/////////////////////////////////////////
	// Delete Hall
	/////////////////////////////////////////
	deleteHall(halls[1], hallService)
}

func deleteHall(hallId int32, hallService HallService.HallService) {
	res, err := hallService.DeleteHall(context.TODO(), &HallService.DeleteHallMessage{
		HallID: hallId,
	})

	if err != nil || !res.Success {
		fmt.Println(err)
		fmt.Println("DeleteHall failed!")
	} else {
		fmt.Printf("Deleted hall %d\n", hallId)
	}
}

func confirmBooking(bookingID int32, userID int32, bookingService BookingService.BookingService) int32 {
	res, err := bookingService.ConfirmBooking(context.TODO(), &BookingService.ConfirmBookingMessage{
		UserID:    userID,
		BookingID: bookingID,
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println("Retrying...")
		os.Exit(-1)
	}

	fmt.Printf("Confirmed booking:%d: show %d for user %d\n", res.BookingID, bookingID, userID)

	return res.BookingID
}

func createBooking(showID int32, userID int32, seats []int32, bookingService BookingService.BookingService) int32 {
	res, err := bookingService.CreateBooking(context.TODO(), &BookingService.CreateBookingMessage{
		UserID: userID,
		ShowID: showID,
		Seats:  seats,
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println("Booking failed!")
		return -1
	} else {
		fmt.Printf("Created booking:%d: show %d for user %d\n", res.BookingID, showID, userID)
		return res.BookingID
	}
}

func createShow(hallID int32, movieID int32, showService ShowService.ShowService) int32 {
	res, err := showService.CreateShow(context.TODO(), &ShowService.CreateShowMessage{
		MovieID: movieID,
		HallID:  hallID,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("Created show:%d: movie %d in hall %d\n", res.ShowID, movieID, hallID)

	return res.ShowID
}

func createUser(name string, userService UserService.UserService) int32 {
	res, err := userService.CreateUser(context.TODO(), &UserService.CreateUserMessage{
		UserName: name,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("Created user: %s with ID %d\n", name, res.UserID)

	return res.UserID
}

func createHall(name string, rows int32, seatsPerRow int32, hallService HallService.HallService) int32 {
	res, err := hallService.CreateHall(context.TODO(), &HallService.CreateHallMessage{
		HallName:    name,
		Rows:        rows,
		SeatsPerRow: seatsPerRow,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("Created hall: %s with ID %d\n", name, res.HallID)

	return res.HallID
}

func createMovie(name string, movieService MovieService.MovieService) int32 {
	res, err := movieService.CreateMovie(context.TODO(), &MovieService.CreateMovieMessage{
		Title: name,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("Created movie: %s with ID %d\n", name, res.MovieID)

	return res.MovieID
}

func getMovie(movieID int32, movieService MovieService.MovieService) string {
	const timeout = 20 * time.Second

	test1Title := ""
	for test1Title == "" {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		res2, err2 := movieService.GetMovie(ctx, &MovieService.GetMovieMessage{
			MovieID: movieID,
		})

		if err2 != nil {
			fmt.Println(err2)
			fmt.Println("Retrying...")
		} else {
			test1Title = res2.Title
		}
	}

	return test1Title
}
