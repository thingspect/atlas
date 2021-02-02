// Package api provides functions used to run the API service.
package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/NYTimes/gziphandler"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/api/config"
	"github.com/thingspect/atlas/internal/api/interceptor"
	"github.com/thingspect/atlas/internal/api/service"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao/datapoint"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/dao/user"
	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/queue"
	"google.golang.org/grpc"

	// encoding/gzip imported for use by UseCompressor CallOption.
	_ "google.golang.org/grpc/encoding/gzip"
)

const (
	ServiceName = "api"
	GRPCHost    = "127.0.0.1"
	GRPCPort    = ":50051"
	httpPort    = ":8000"
)

var errPWTLength = errors.New("pwt key must be 32 bytes")

// API holds references to the gRPC and HTTP servers.
type API struct {
	grpcSrv    *grpc.Server
	httpSrv    *http.Server
	httpCancel context.CancelFunc
}

// New builds a new API and returns a reference to it and an error value.
func New(cfg *config.Config) (*API, error) {
	// Validate Config.
	if len(cfg.PWTKey) != 32 {
		return nil, errPWTLength
	}

	// Set up database connection.
	pg, err := postgres.New(cfg.PgURI)
	if err != nil {
		return nil, err
	}

	// Build the NSQ connection for publishing.
	nsq, err := queue.NewNSQ(cfg.NSQPubAddr, nil, "",
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		return nil, err
	}

	// Register gRPC services.
	skipAuth := map[string]struct{}{
		"/api.SessionService/Login": {},
	}
	skipValidate := map[string]struct{}{
		// Update actions validate after merge to support partial updates.
		"/api.DeviceService/UpdateDevice": {},
		"/api.OrgService/UpdateOrg":       {},
		"/api.UserService/UpdateUser":     {},
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.Log(nil),
		interceptor.Auth(skipAuth, cfg.PWTKey),
		interceptor.Validate(skipValidate),
	))
	api.RegisterDataPointServiceServer(srv, service.NewDataPoint(nsq,
		cfg.NSQPubTopic, datapoint.NewDAO(pg)))
	api.RegisterDeviceServiceServer(srv, service.NewDevice(device.NewDAO(pg)))
	api.RegisterOrgServiceServer(srv, service.NewOrg(org.NewDAO(pg)))
	api.RegisterSessionServiceServer(srv, service.NewSession(user.NewDAO(pg),
		cfg.PWTKey))
	api.RegisterUserServiceServer(srv, service.NewUser(user.NewDAO(pg)))

	// Register gRPC-Gateway handlers.
	ctx, cancel := context.WithCancel(context.Background())
	gwMux := runtime.NewServeMux(runtime.WithForwardResponseOption(statusCode))
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	// DataPoint.
	if err := api.RegisterDataPointServiceHandlerFromEndpoint(ctx, gwMux,
		GRPCHost+GRPCPort, opts); err != nil {
		cancel()
		return nil, err
	}

	// Device.
	if err := api.RegisterDeviceServiceHandlerFromEndpoint(ctx, gwMux,
		GRPCHost+GRPCPort, opts); err != nil {
		cancel()
		return nil, err
	}

	// Org.
	if err := api.RegisterOrgServiceHandlerFromEndpoint(ctx, gwMux,
		GRPCHost+GRPCPort, opts); err != nil {
		cancel()
		return nil, err
	}

	// Session.
	if err := api.RegisterSessionServiceHandlerFromEndpoint(ctx, gwMux,
		GRPCHost+GRPCPort, opts); err != nil {
		cancel()
		return nil, err
	}

	// User.
	if err := api.RegisterUserServiceHandlerFromEndpoint(ctx, gwMux,
		GRPCHost+GRPCPort, opts); err != nil {
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
			Handler: gziphandler.GzipHandler(mux),
		},
		httpCancel: cancel,
	}, nil
}

// Serve starts the listener.
func (api *API) Serve() {
	//#nosec G102 // service should listen on all interfaces
	lis, err := net.Listen("tcp", GRPCPort)
	if err != nil {
		alog.Fatalf("Serve net.Listen: %v", err)
	}

	// Serve gRPC.
	go func() {
		alog.Infof("Listening on %v", GRPCPort)
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
