package api

import (
	"encoding/json"
	"net/http"

	"github.com/dapper-labs/identity-server/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GetUsersResponse struct {
	Users []model.User `json:"users,omitempty"`
}

func (api *API) GetUsers(w http.ResponseWriter, r *http.Request) error {

	users, err := api.userRep.GetUsers()
	if err != nil {
		logrus.Error(errors.Wrap(err, "could not get the users from an user repository"))
		return internalServerError("an error occurred while processing your request, please contact an administrator")
	}

	response := GetUsersResponse{}
	if users != nil {
		response.Users = *users
	}

	err = writeJSON(w, http.StatusOK, response)
	if err != nil {
		logrus.Error(errors.Wrap(err, "could not write response"))
	}
	return err
}

func (api *API) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	claims, err := getClaims(r.Context())
	if err != nil {
		return unauthorizedError("user is not authorized to process this request")
	}
	updateUser := model.UpdateUser{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&updateUser)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to decode request body"))
		return badRequestError("could not read received update userinformation")
	}

	err = api.userRep.UpdateUserWithEmail(&updateUser, claims.Email)
	if err != nil {
		return internalServerError("cannot update user, please contact administrator")
	}

	user, err := api.userRep.GetUserWithEmail(claims.Email)
	if err != nil {
		return internalServerError("cannot update user, please contact administrator")
	}
	err = writeJSON(w, http.StatusOK, user)
	if err != nil {
		logrus.Error(errors.Wrap(err, "could not write response"))
	}
	return err

}
