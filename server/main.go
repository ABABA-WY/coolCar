package main

import (
	"context"
	trip "coolcar/TripService"
	trippb "coolcar/proto/gen/go"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)
	go startGRPCGateway()
	fmt.Println("127.0.0.1:8080")
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	trippb.RegisterTripServiceServer(s, &trip.Service{})
	log.Fatalln(s.Serve(listen))

}

//开启一个GRPC Gateway代理
func startGRPCGateway() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{},
	))
	//RegisterTripServiceHandlerFromEndpoint is same as RegisterTripServiceHandler
	//but automatically dials to "endpoint" and closes the connection when "ctx" gets done.
	//err := trippb.RegisterTripServiceHandlerFromEndpoint(
	//	//ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption
	//	c,
	//	mux,
	//	"localhost:8080",
	//	[]grpc.DialOption{grpc.WithInsecure()},
	//)
	//创建一个连接
	err := trippb.RegisterTripServiceHandlerFromEndpoint(
		c,
		mux,
		"localhost:8080",
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		fmt.Println("1")
		panic(err)
	}
	//另一个监听
	err = http.ListenAndServe("127.0.0.1:8000", mux)
	if err != nil {
		fmt.Println("2")
		panic(err)
	}
}
