// Package api provides functions used to run the API service.
package api

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	grpcHost    = "127.0.0.1"
	grpcPort    = ":50051"
	httpPort    = ":8000"
)

// API holds references to the gRPC and HTTP servers.
type API struct {
	grpcSrv    *grpc.Server
	httpSrv    *http.Server
	httpCancel context.CancelFunc
}

// New builds a new API and returns a reference to it and an error value.
func New(cfg *config.Config) (*API, error) {
	// Set up database connection.
	pg, err := postgres.New(cfg.PgURI)
	if err != nil {
		return nil, err
	}

	// Register gRPC services.
	srv := grpc.NewServer()
	api.RegisterDeviceServiceServer(srv, service.NewDevice(device.NewDAO(pg)))

	// Register gRPC-Gateway handlers.
	ctx, cancel := context.WithCancel(context.Background())
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	// Device.
	if err := api.RegisterDeviceServiceHandlerFromEndpoint(ctx, gwMux,
		grpcHost+grpcPort, opts); err != nil {
		cancel()
		return nil, err
	}

	// OpenAPI.
	mux := http.NewServeMux()
	mux.Handle("/v1/", gwMux)
	mux.Handle("/", http.FileServer(http.Dir("web")))

	return &API{
		grpcSrv: srv,
		httpSrv: &http.Server{
			Addr:    httpPort,
			Handler: mux,
		},
		httpCancel: cancel,
	}, nil
}

// Serve starts the listener.
func (api *API) Serve() {
	// #nosec G102 - service should listen on all interfaces
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		alog.Fatalf("Serve net.Listen: %v", err)
	}

	// Serve gRPC.
	go func() {
		alog.Infof("Listening on %v", grpcPort)
		if err := api.grpcSrv.Serve(lis); err != nil {
			alog.Fatalf("Serve api.grpcSrv.Serve: %v", err)
		}
	}()

	// Serve gRPC-gateway.
	go func() {
		alog.Infof("Listening on %v", httpPort)
		if err := api.httpSrv.ListenAndServe(); err != nil {
			alog.Fatalf("Serve api.httpSrv.ListenAndServe: %v", err)
		}
	}()

	// Handle graceful shutdown.
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	alog.Info("Serve received signal, exiting")
	api.httpCancel()
}
