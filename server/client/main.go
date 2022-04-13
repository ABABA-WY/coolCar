package main

import (
	"context"
	trippb "coolcar/proto/gen/go"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
	dial, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("client1:", err)
	}
	tsClient := trippb.NewTripServiceClient(dial)
	respon, err := tsClient.GetTrip(context.Background(), &trippb.GetTripRequest{
		Id: "sb12",
	})
	if err != nil {
		log.Fatal("client2", err)
	}
	fmt.Println(respon)
}
