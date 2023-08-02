package service

import (
	"context"
	"fmt"
	"github.com/DMXRoid/QDLEDController/v2/led"
	pb "github.com/DMXRoid/QDLEDController/v2/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
)

type ledControllerServer struct {
	pb.UnimplementedLEDControllerServer
}

func Init() {
	fmt.Println("starting service")
	l, err := net.Listen("tcp", ":8081")
	if err == nil {

		fmt.Println("starting grpc server")
		opts := []grpc.ServerOption{}
		grpcServer := grpc.NewServer(opts...)
		pb.RegisterLEDControllerServer(grpcServer, &ledControllerServer{})
		go grpcServer.Serve(l)

		fmt.Println("starting http server")

		ctx := context.Background()
		mux := runtime.NewServeMux()
		gwopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err = pb.RegisterLEDControllerHandlerFromEndpoint(ctx, mux, ":8081", gwopts)
		if err != nil {
			fmt.Println("error:::", err)
		}
		mux.HandlePath("GET", "/", func(w http.ResponseWriter, r *http.Request, p map[string]string) {
			fp := fmt.Sprintf("./www%s", r.URL.Path)
			fmt.Println("getting file:::", fp)
			http.ServeFile(w, r, fp)
		})
		mux.HandlePath("GET", "/styles.css", func(w http.ResponseWriter, r *http.Request, p map[string]string) {
			fp := fmt.Sprintf("./www%s", r.URL.Path)
			fmt.Println("getting file:::", fp)
			http.ServeFile(w, r, fp)
		})
		mux.HandlePath("GET", "/site-scripts.js", func(w http.ResponseWriter, r *http.Request, p map[string]string) {
			fp := fmt.Sprintf("./www%s", r.URL.Path)
			fmt.Println("getting file:::", fp)
			http.ServeFile(w, r, fp)
		})
		fmt.Println(http.ListenAndServe(":8080", mux))
		fmt.Println("foo")

	} else {
		fmt.Println(err)
	}
}

func (s *ledControllerServer) GetLEDs(ctx context.Context, req *pb.GetLEDsRequest) (*pb.GetLEDsResponse, error) {
	var err error
	leds := make([]*pb.LED, 0, len(led.RegisteredLEDs))
	for _, l := range led.RegisteredLEDs {
		leds = append(leds, l.LED)
	}
	resp := &pb.GetLEDsResponse{Metadata: &pb.ResponseMetadata{Code: 200}, Leds: leds}

	return resp, err
}

func (s *ledControllerServer) UpdateLEDs(ctx context.Context, req *pb.UpdateLEDsRequest) (*pb.UpdateLEDsResponse, error) {
	var err error
	resp := &pb.UpdateLEDsResponse{Metadata: &pb.ResponseMetadata{Code: 200}}
	for _, l := range req.Leds {
		err = led.Update(l)
	}

	return resp, err
}
