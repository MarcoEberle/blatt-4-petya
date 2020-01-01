package main

// https://ewanvalentine.io/microservices-in-golang-part-1/
// Mutex auf Datastore -> bei uns wahrscheinlich der Speicher der IDs, etc.
// Wobei wir eh so ne Art Repo für die einzelnen Sachen benutzen
// Erinnerung: so ähnlich wie bei Software Engineering, nur anders
//

import (
	"context"
	"fmt"
	micro "github.com/micro/go-micro"
	proto "github.com/ob-vss-ws19/blatt-4-petya/messages"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	println("Hello World")

	service := micro.NewService(
		micro.Name("greeter"),
	)

	service.Init()

	proto.RegisterGreeterHandler(service.Server(), new(Greeter))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
