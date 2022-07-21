package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dapper-labs/identity-server/mail"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) error {

	userLogin := UserLogin{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userLogin)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to decode request body"))
		return badRequestError("could not read received login information")
	}

	err = mail.Validate(userLogin.Email)
	if err != nil {
		return badRequestError("the provided email has incorrect format")
	}

	if userLogin.Password == "" {
		return badRequestError("password must not be empty")
	}

	user, err := api.userRep.GetUserWithEmail(userLogin.Email)
	if err != nil {
		return unauthorizedError("could not find your account")
	}

	auth := user.Challenge(userLogin.Password)
	if !auth {
		return unauthorizedError("wrong user name or password")
	}

	nsecs := time.Second * 60 * time.Duration(api.config.Expiration)
	token, err := newJwtToken(user, nsecs, api.config.JWT.Secret)

	err = writeJSON(w, http.StatusOK, UserLoginResponse{Token: token, Code: http.StatusOK})
	if err != nil {
		logrus.Error(errors.Wrap(err, "could not write response"))
	}
	return err
}
