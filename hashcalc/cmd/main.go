package main

import (
	"hashkeeper/pkg/hashlog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"hashcalc/internal/grpchandlers"
	"hashcalc/pkg/grpchashcalc"

	"github.com/alexflint/go-arg"

	log "github.com/sirupsen/logrus"
)

type AppCfg struct {
	ListenAddress string `arg:"--listen-address,env:LISTEN_ADDRESS" default:":50051" help:"Address and port of gRPC server"`
}

func main() {
	var args struct {
		AppCfg
		hashlog.LogCfg
	}
	arg.MustParse(&args)

	hashlog.InitLog(args.LogCfg)

	log.Info("hashcalc is starting ...")

	log.Infof("start arguments: %+v", args)

	lis, err := net.Listen("tcp", args.ListenAddress)
	if err != nil {
		hashlog.LogErrorWithStack(err).Fatal("can`t start listen")
	}

	log.Info("started, waiting for request ...")

	// Prepare gRPC
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(grpchashcalc.RequestIDServerInterceptor))
	funcServer := &grpchandlers.HashCalcServerImpl{}

	grpchashcalc.RegisterHashCalcServer(grpcServer, funcServer)

	// Add reflection for grpcui
	reflection.Register(grpcServer)
	log.Warn("gRPC reflection started")

	// Start
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			hashlog.LogErrorWithStack(err).Fatal("can`t start server")
		}
	}()

	// Wait for exit
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM)

	<-sigTerm

	// Stop
	log.Info("hashcalc is shutting down ...")

	grpcServer.Stop()

	log.Info("hashcalc stopped")
}
