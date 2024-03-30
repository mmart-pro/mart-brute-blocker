package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	bucketstorage "github.com/mmart-pro/mart-brute-blocker/internal/bucket/storage"
	"github.com/mmart-pro/mart-brute-blocker/internal/config"
	grpcserver "github.com/mmart-pro/mart-brute-blocker/internal/grpc"
	"github.com/mmart-pro/mart-brute-blocker/internal/service/mbbservice"
	memorystorage "github.com/mmart-pro/mart-brute-blocker/internal/storage/memory"
	sqlstorage "github.com/mmart-pro/mart-brute-blocker/internal/storage/sql"
	"github.com/mmart-pro/mart-brute-blocker/pkg/logger"
)

type ConnectableStorage interface {
	mbbservice.Storage
	Connect(ctx context.Context) error
	Close() error
}

type App struct {
	cfg config.APIConfig
}

func NewApp(cfg config.APIConfig) *App {
	return &App{
		cfg: cfg,
	}
}

func (app App) Startup(ctx context.Context) error {
	// logger
	logg, err := logger.NewLogger(app.cfg.LoggerConfig.Level, app.cfg.LoggerConfig.LogFile)
	if err != nil {
		return fmt.Errorf("can't start logg: %w", err)
	}
	defer logg.Close()

	// IRL exclude secrets
	logg.Debugf("service config: %+v", app.cfg)

	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// storage
	var storage ConnectableStorage
	if app.cfg.StorageConfig.UseDB {
		storage = sqlstorage.NewStorage(
			app.cfg.StorageConfig.Host,
			app.cfg.StorageConfig.Port,
			app.cfg.StorageConfig.User,
			app.cfg.StorageConfig.Password,
			app.cfg.StorageConfig.Database)
	} else {
		storage = memorystorage.NewStorage()
	}
	if err := storage.Connect(ctx); err != nil {
		logg.Fatalf("can't connect to database %s", err.Error())
	}
	defer storage.Close()

	// events service
	ipBucketStorage := bucketstorage.NewBucketMemoryStorage()
	loginBucketStorage := bucketstorage.NewBucketMemoryStorage()
	pwdBucketStorage := bucketstorage.NewBucketMemoryStorage()
	mbbService := mbbservice.NewMBBService(logg, storage, ipBucketStorage, loginBucketStorage, pwdBucketStorage, app.cfg.ServiceConfig)

	logg.Infof("starting mbb api...")

	grpc := grpcserver.NewServer(app.cfg.GrpcConfig.GetEndpoint(), logg, mbbService)
	go func() {
		if err := grpc.Start(); err != nil {
			logg.Errorf("failed to start grpc server %s", err.Error())
			cancel()
		}
	}()

	<-ctx.Done()

	grpc.Stop()

	logg.Infof("mbb api stopped")

	return nil
}
