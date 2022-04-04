package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ab3llo/go-product-svc/pkg/config"
	"github.com/ab3llo/go-product-svc/pkg/db"
	"github.com/ab3llo/go-product-svc/pkg/product/pb"
	"github.com/ab3llo/go-product-svc/pkg/services"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed to load config", err)
	}
	db := db.Connect(cfg)

	lis, err := net.Listen("tcp", cfg.Port)

	if err != nil {
		log.Fatalln("Failed to listen to server:", err)
	}

	fmt.Println("Product svc listening on port: ", cfg.Port)

	s := services.Server{
		DbConnection: db,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterProductServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to start server", err)
	}
}
