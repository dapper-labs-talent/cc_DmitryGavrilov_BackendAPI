package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dapper-labs/identity-server/mail"
	"github.com/dapper-labs/identity-server/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UserSignUp struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UserSignUpResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

func (api *API) SignUp(w http.ResponseWriter, r *http.Request) error {
	userSignup := UserSignUp{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userSignup)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to decode request body"))
		return badRequestError("could not read received user information")
	}

	if userSignup.Password == "" {
		return badRequestError("to create a new user, the password must not be empty")
	}

	err = mail.Validate(userSignup.Email)
	if err != nil {
		return badRequestError("the provided email has incorrect format")
	}

	user, err := model.NewUser(userSignup.Firstname, userSignup.Lastname, userSignup.Email, userSignup.Password)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to create a new user structure"))
		return internalServerError("cannot create a new user, please contact administrator")
	}

	err = api.userRep.Insert(user)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to insert a new user"))
		return internalServerError("cannot persist a new user, please contact administrator")
	}

	nsecs := time.Second * 60 * time.Duration(api.config.Expiration)
	token, err := newJwtToken(user, nsecs, api.config.JWT.Secret)

	err = writeJSON(w, http.StatusOK, UserSignUpResponse{Token: token, Code: http.StatusOK})
	if err != nil {
		logrus.Error(errors.Wrap(err, "could not write response"))
	}
	return err
}
