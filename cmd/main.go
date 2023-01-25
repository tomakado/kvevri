package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcimpl "github.com/tomakado/kvevri/internal/grpc"
	"github.com/tomakado/kvevri/internal/pb"
	"github.com/tomakado/kvevri/store"
	"google.golang.org/grpc"
)

func main() {
	var (
		store    = store.New(30 * time.Second)
		srv      = grpcimpl.NewServer(store)
		grpcOpts []grpc.ServerOption
	)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpcOpts...)
	pb.RegisterStoreServer(grpcServer, srv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	expCtx, expCancel := context.WithCancel(context.Background())
	defer expCancel()

	log.Println("starting expiration worker")
	go store.StartExpirationWorker(expCtx, 5 * time.Second)

	go func() {
		if err := grpcServer.Serve(lis); !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatal(err)
		}
	}()

	log.Println("server started")
	<-quit

	expCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	stopped := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		stopped <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		log.Fatal("server shutdown timed out")
	case <-stopped:
		log.Println("server stopped")
	}
}
