package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dapper-labs/identity-server/config"
	"github.com/dapper-labs/identity-server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type API struct {
	config  *config.Config
	userRep storage.UserRepository
	mux     *chi.Mux
}

const (
	MaxShutdownTimeoutInSeconds = 60
)

func (api *API) ListenAndServe() error {

	address := fmt.Sprintf("0.0.0.0:%d", api.config.ListenPort)

	server := &http.Server{
		Handler: api.mux,
	}

	shutdown := make(chan struct{})
	defer close(shutdown)

	go func() {

		timeout := api.config.Server.Timeout
		if timeout < 0 || timeout > MaxShutdownTimeoutInSeconds {
			timeout = MaxShutdownTimeoutInSeconds
		}

		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-sigchan:
			logrus.Info(fmt.Sprintf("received os signal - %s, triggering server shutdown", sig))
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		logrus.Error(errors.Wrap(err, fmt.Sprintf("cannot start server at %s", address)))
	}

	return nil
}

func NewAPI(config *config.Config) (*API, error) {
	userRep, err := storage.NewUserRepository(config)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	api := &API{
		userRep: userRep,
		config:  config,
		mux:     router,
	}

	router.Route("/v1", func(r chi.Router) {
		r.Post("/signup", handler(api.SignUp))
		r.Post("/login", handler(api.Login))
		r.Get("/users", handler(api.GetUsers))
		r.Put("/users", handler(api.UpdateUser))
	})

	return api, nil
}
