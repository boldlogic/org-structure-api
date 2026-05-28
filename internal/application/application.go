package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/boldlogic/org-structure-api/internal/repository"
	"github.com/boldlogic/org-structure-api/internal/service"
	server "github.com/boldlogic/org-structure-api/internal/transport/http"
	"github.com/boldlogic/org-structure-api/internal/transport/middleware"
	"github.com/boldlogic/org-structure-api/pkg/config"
	"github.com/boldlogic/packages/commonconfig"
	"github.com/boldlogic/packages/dbgorm"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/packages/transport/httpserver"
	"go.uber.org/zap"
)

const (
	defaultConfigPath = "config.yaml"
	errChanBufSize    = 1
)

type Application struct {
	cfg     *config.Config
	Logger  *zap.Logger
	srv     *httpserver.Server
	repo    *repository.Repo
	closeDB func() error
	errChan chan error
	wg      sync.WaitGroup
}

func New() (*Application, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)

	conf, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}
	log := logger.New(conf.Logger)

	return &Application{
		cfg:     conf,
		Logger:  log,
		errChan: make(chan error, errChanBufSize),
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	gormDB, err := dbgorm.ConnectDatabase(ctx, &a.cfg.Database)
	if err != nil {
		return err
	}
	a.closeDB = func() error {
		return dbgorm.CloseDatabase(gormDB)
	}

	a.repo = repository.NewRepo(gormDB, a.Logger)
	svc := service.NewService(a.repo)
	handler := server.NewHandler(svc, a.Logger)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	middle := middleware.NewMiddleware(a.Logger)
	rec := middle.Recover(mux)
	wrapped := middle.WithLogging(rec)
	a.srv = httpserver.NewServer(wrapped, a.cfg.HTTP)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errChan <- fmt.Errorf("http server остановлен с ошибкой: %w", err)
		}
	}()

	a.Logger.Info("http server запущен",
		zap.String("addr", fmt.Sprintf("%s:%d", a.cfg.HTTP.ListenIp, a.cfg.HTTP.ListenPort)),
	)

	return nil
}

func (a *Application) Wait(ctx context.Context, cancel context.CancelFunc) error {
	var appErr error

	errWg := sync.WaitGroup{}
	errWg.Add(1)

	go func() {
		defer errWg.Done()
		for err := range a.errChan {
			cancel()
			appErr = err
		}
	}()

	<-ctx.Done()

	if a.srv != nil {
		shCtx, shCancel := context.WithTimeout(context.Background(), time.Duration(a.cfg.HTTP.Opts.ShutdownTimeout)*time.Second)
		defer shCancel()
		_ = a.srv.Shutdown(shCtx)
	}

	a.wg.Wait()

	if a.closeDB != nil {
		if err := a.closeDB(); err != nil {
			a.Logger.Error("не удалось остановить соединение с БД", zap.Error(err))
		}
	}

	close(a.errChan)
	errWg.Wait()

	return appErr
}
