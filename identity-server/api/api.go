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

func (api *API) ListenAndServe() error {
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
		r.Post("signup", handler(api.SignUp))
		r.Post("login", handler(api.Login))
		r.Get("users", handler(api.GetUsers))
		r.Put("users", handler(api.UpdateUser))
	})

	return api, nil
}
