package api

import (
	"github.com/dapper-labs/identity-server/config"
	"github.com/dapper-labs/identity-server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	config  *config.Config
	userRep storage.UserRepository
	mux     *chi.Mux
}

func NewAPI(config *config.Config) (*API, error) {
	userRep, err := storage.NewUserRepository(config)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	router.Route("/v1", func(r chi.Router) {
		// POST signup
		// GET users
		// POST login
		// PUT users
	})

	api := &API{
		userRep: userRep,
		config:  config,
		mux:     router,
	}

	return api, nil
}
