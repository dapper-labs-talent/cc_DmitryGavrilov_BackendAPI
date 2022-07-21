// package api
package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dapper-labs/identity-server/config"
	"github.com/stretchr/testify/suite"
)

type TestSignUpSuite struct {
	suite.Suite
	api *API
}

func (s *TestSignUpSuite) SetupSuite() {
	config := config.CreateTestConfiguration()
	api, _ := NewAPI(config)
	s.api = api
}

func TestSignUpTestSuite(t *testing.T) {
	suite.Run(t, &TestSignUpSuite{})
}

func createSignUp(firstname, lastname, email, password string) *UserSignUp {
	return &UserSignUp{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Password:  password,
	}
}

func convert(data interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func (s *TestSignUpSuite) TestSignUp_ValidaUser_Ok_Expected() {
	w := httptest.NewRecorder()
	user := createSignUp("test_first_name", "test_last_name", "good@gmail.com", "123")
	buffer, err := convert(user)
	s.Assert().Nil(err)
	req, _ := http.NewRequest(http.MethodPost, "/v1/signup", buffer)
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	s.api.mux.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Result().StatusCode)

	respBody, err := ioutil.ReadAll(w.Result().Body)
	s.Assert().Nil(err)

	signUpResponse := UserSignUpResponse{}
	err = json.Unmarshal(respBody, &signUpResponse)

	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusOK, signUpResponse.Code)
	s.Assert().NotEmpty(signUpResponse.Token)

	assertToken(s.Assert(), s.api, signUpResponse.Token, user.Email)
}

func (s *TestSignUpSuite) TestSignUp_EmptyPassword_BadRequest_Expected() {
	w := httptest.NewRecorder()
	user := createSignUp("test_first_name", "test_last_name", "good@gmail.com", "")
	buffer, err := convert(user)
	s.Assert().Nil(err)
	req, _ := http.NewRequest(http.MethodPost, "/v1/signup", buffer)
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	s.api.mux.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusBadRequest, w.Result().StatusCode)

	respBody, err := ioutil.ReadAll(w.Result().Body)
	s.Assert().Nil(err)

	errorResponse := HttpErrorResponse{}
	err = json.Unmarshal(respBody, &errorResponse)

	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusBadRequest, errorResponse.Code)

	// due to limited time, I am going to match the exact string, however it would be good to define an error code that will
	// will be associated with an email validation error type
	s.Assert().Equal("to create a new user, the password must not be empty", errorResponse.ErrorMessage)
}

func (s *TestSignUpSuite) TestSignUp_BadEmailFormat_BadRequest_Expected() {

	w := httptest.NewRecorder()
	user := createSignUp("test_first_name", "test_last_name", "bademail", "123")
	buffer, err := convert(user)
	s.Assert().Nil(err)
	req, _ := http.NewRequest(http.MethodPost, "/v1/signup", buffer)
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	s.api.mux.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusBadRequest, w.Result().StatusCode)

	respBody, err := ioutil.ReadAll(w.Result().Body)
	s.Assert().Nil(err)

	errorResponse := HttpErrorResponse{}
	err = json.Unmarshal(respBody, &errorResponse)

	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusBadRequest, errorResponse.Code)

	// due to limited time, I am going to match the exact string, however it would be good to define an error code that will
	// will be associated with an email validation error type
	s.Assert().Equal("the provided email has incorrect format", errorResponse.ErrorMessage)
}
