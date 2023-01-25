package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tomakado/kvevri/internal/config"
	grpcimpl "github.com/tomakado/kvevri/internal/grpc"
	"github.com/tomakado/kvevri/internal/pb"
	"github.com/tomakado/kvevri/store"
	"google.golang.org/grpc"
)

const asciiLogo = `
 ___  __    ___      ___ _______   ___      ___ ________  ___     
|\  \|\  \ |\  \    /  /|\  ___ \ |\  \    /  /|\   __  \|\  \    
\ \  \/  /|\ \  \  /  / | \   __/|\ \  \  /  / | \  \|\  \ \  \   
 \ \   ___  \ \  \/  / / \ \  \_|/_\ \  \/  / / \ \   _  _\ \  \  
  \ \  \\ \  \ \    / /   \ \  \_|\ \ \    / /   \ \  \\  \\ \  \ 
   \ \__\\ \__\ \__/ /     \ \_______\ \__/ /     \ \__\\ _\\ \__\
    \|__| \|__|\|__|/       \|_______|\|__|/       \|__|\|__|\|__|

`

func main() {
	fmt.Print(asciiLogo)
	log.Printf("starting kvevri with config %+v", config.Get())

	var (
		store    = store.New(config.Get().TTL)
		srv      = grpcimpl.NewServer(store)
		grpcOpts []grpc.ServerOption
	)

	lis, err := net.Listen("tcp", config.Get().ListenAddr)
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

	go store.StartExpirationWorker(expCtx, 5*time.Second)

	go func() {
		if err := grpcServer.Serve(lis); !errors.Is(err, grpc.ErrServerStopped) {
			log.Println(err)
		}
	}()

	log.Println("server started")
	<-quit

	expCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), config.Get().TTL/2)
	defer shutdownCancel()

	stopped := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		stopped <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		log.Println("server shutdown timed out")
		return
	case <-stopped:
		log.Println("server stopped")
	}
}
