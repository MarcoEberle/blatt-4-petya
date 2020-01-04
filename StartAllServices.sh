#!/bin/bash

echo "Starte HallService"
go run ./HallService/main.go &
echo "Starte MovieService"
go run ./MovieService/main.go &
echo "Starte ShowService"
go run ./ShowService/main.go &
echo "Starte BookingService"
go run ./BookingService/main.go &
echo "Starte UserService"
go run ./UserService/main.go &