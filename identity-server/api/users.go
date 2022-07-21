package api

import (
	"net/http"

	"github.com/dapper-labs/identity-server/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GetUsersResponse struct {
	Users *[]model.User `json:"users"`
}

func (api *API) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := api.userRep.GetUsers()
	if err != nil {
		logrus.Error(errors.Wrap(err, "could not get the users from an user repository"))
		return internalServerError("an error occurred while processing your request, please contact an administrator")
	}

	err = writeJSON(w, http.StatusOK, GetUsersResponse{Users: users})
	if err != nil {
		logrus.Error(errors.Wrap(err, "could not write response"))
	}
	return err
}

func (api *API) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}
