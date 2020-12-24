// Package api provides functions used to run the API service.
package api

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/api/config"
	"github.com/thingspect/atlas/internal/api/service"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/postgres"
	"google.golang.org/grpc"
)

const (
	ServiceName = "api"
)

// API holds references to the gRPC server.
type API struct {
	server *grpc.Server
}

// New builds a new Api and returns a reference to it and an error value.
func New(cfg *config.Config) (*API, error) {
	// Set up database connection.
	pg, err := postgres.New(cfg.PgURI)
	if err != nil {
		return nil, err
	}

	// Register services.
	srv := grpc.NewServer()
	api.RegisterDeviceServiceServer(srv, service.NewDevice(device.NewDAO(pg)))

	return &API{
		server: srv,
	}, nil
}

// Serve starts the listener.
func (api *API) Serve() {
	// #nosec G102 - service should listen on all interfaces
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		alog.Fatalf("Serve net.Listen: %v", err)
	}

	go func() {
		alog.Info("Listening on :50051")
		if err := api.server.Serve(lis); err != nil {
			alog.Fatalf("Serve api.server.Serve: %v", err)
		}
	}()

	// Handle graceful shutdown.
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	alog.Info("Serve received signal, exiting")
}
