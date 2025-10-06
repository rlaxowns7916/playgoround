package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"playgoround/chat/internal/config"
	"playgoround/chat/internal/database"
	"playgoround/chat/internal/database/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	cfg        *config.Config
	httpServer *echo.Echo
	rdb        *database.Database
	userRepo   *repository.UserRepository
	roomRepo   *repository.RoomRepository
	chatRepo   *repository.ChatRepository
}

func New(cfg *config.Config) (*Application, error) {
	dbCfg := &database.DatabaseConfig{
		Database: cfg.Database.Database,
		URL:      cfg.Database.URL,
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
	}
	rdb, err := database.NewDatabase(dbCfg)
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(rdb.DB())
	roomRepo := repository.NewRoomRepository(rdb.DB())
	chatRepo := repository.NewChatRepository(rdb.DB())

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Server = &http.Server{
		Addr:              net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	app := &Application{
		httpServer: e,
		rdb:        rdb,
		cfg:        cfg,
		userRepo:   userRepo,
		roomRepo:   roomRepo,
		chatRepo:   chatRepo,
	}

	return app, nil
}

func (a *Application) Start(ctx context.Context) error {
	a.httpServer.Server.BaseContext = func(net.Listener) context.Context { return ctx }
	ln, err := net.Listen("tcp", net.JoinHostPort(a.cfg.Server.Host, strconv.Itoa(a.cfg.Server.Port)))
	if err != nil {
		return fmt.Errorf("failed to bind %s:%d: %w", a.cfg.Server.Host, a.cfg.Server.Port, err)
	}
	go func() {
		if err := a.httpServer.Server.Serve(ln); err != nil && err != http.ErrServerClosed {
			a.httpServer.Logger.Errorf("http server error: %v", err)
		}
	}()
	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	var err1 error
	if a.httpServer != nil {
		err1 = a.httpServer.Shutdown(ctx)
	}
	var err2 error
	if a.rdb != nil {
		err2 = a.rdb.Close()
	}
	return errors.Join(err1, err2)
}

// Repository getters
func (a *Application) UserRepo() *repository.UserRepository {
	return a.userRepo
}

func (a *Application) RoomRepo() *repository.RoomRepository {
	return a.roomRepo
}

func (a *Application) ChatRepo() *repository.ChatRepository {
	return a.chatRepo
}
