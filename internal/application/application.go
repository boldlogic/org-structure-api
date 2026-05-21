package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/boldlogic/org-structure-api/internal/config"
	"github.com/boldlogic/org-structure-api/internal/repository"
	"github.com/boldlogic/org-structure-api/internal/service"
	server "github.com/boldlogic/org-structure-api/internal/transport/http"
	"github.com/boldlogic/org-structure-api/pkg/database"
	"github.com/boldlogic/packages/commonconfig"
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
	db, err := database.ConnectDatabase(a.cfg.Database.GetDSN())
	if err != nil {
		return err
	}
	a.repo = repository.NewRepo(db, a.Logger)
	svc := service.NewService(a.repo)
	handler := server.NewHandler(svc, a.Logger)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	a.srv = httpserver.NewServer(mux, a.cfg.HTTP)

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
		_ = a.srv.Shutdown(context.Background())
	}

	a.wg.Wait()
	close(a.errChan)
	errWg.Wait()

	return appErr
}
