package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bhankey/go-utils/pkg/grpc/interceptors"
	"github.com/bhankey/go-utils/pkg/logger"
	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/app/container"
	configinternal "github.com/bhankey/pharmacy-automatization-pharmacy/internal/config"
	"github.com/bhankey/pharmacy-automatization-pharmacy/pkg/api/pharmacyproto"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	server    *grpc.Server
	listener  net.Listener
	container *container.Container
	logger    logger.Logger
}

const shutDownTimeoutSeconds = 10

func NewApp(configPath string) (*App, error) {
	config, err := configinternal.GetConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init app because of config error: %w", err)
	}

	log, err := logger.GetLogger(config.Logger.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to init logger error: %w", err)
	}

	log.Info("try to init data source resource")
	dataSources, err := newDataSource(config) // TODO remove dataSource struct
	if err != nil {
		return nil, err
	}

	dependencies := container.NewContainer(
		log,
		dataSources.db,
		dataSources.db,
	)

	grpcHandler := dependencies.GetPharmacyGRPCHandler()

	errorHandlingInterceptor := interceptors.NewErrorHandlingInterceptor(log)
	panicInterceptor := interceptors.NewPanicInterceptor(log)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			middleware.ChainUnaryServer(
				errorHandlingInterceptor.ServerInterceptor(),
				panicInterceptor.ServerInterceptor(),
			),
		),
	)

	reflection.Register(grpcServer)

	pharmacyproto.RegisterPharmacyServiceServer(grpcServer, grpcHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.Server.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return &App{
		server:    grpcServer,
		listener:  listener,
		container: dependencies,
		logger:    log,
	}, nil
}

func (a *App) Start() {
	a.logger.Info("staring server on addr: " + a.listener.Addr().String())
	go func() {
		if err := a.server.Serve(a.listener); err != nil {
			a.logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	a.logger.Info("received signal to shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeoutSeconds*time.Second)
	defer cancel()
	a.server.GracefulStop()

	<-ctx.Done()

	a.container.CloseAllConnections()

	a.logger.Info("server was shutdown")
}
