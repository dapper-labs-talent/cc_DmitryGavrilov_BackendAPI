package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dapper-labs/identity-server/config"
	"github.com/dapper-labs/identity-server/model"
	"github.com/dapper-labs/identity-server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestLoginSuite struct {
	suite.Suite
	api *API
}

func createTestUsers(repo storage.UserRepository, number int, firstname, lastname, email, password string) {
	if number <= 0 {
		return
	}

	index := 1
	for index <= number {
		newFirstname := fmt.Sprintf("%s%d", firstname, index)
		newLastname := fmt.Sprintf("%s%d", lastname, index)
		newEmail := fmt.Sprintf("%d%s", index, email)
		newPassword := fmt.Sprintf("%s%d", password, index)
		user, err := model.NewUser(newFirstname, newLastname, newEmail, newPassword)
		if err != nil {
			panic(err)
		}
		err = repo.Insert(user)
		if err != nil {
			panic(err)
		}
		index++
	}
}

func (s *TestLoginSuite) SetupSuite() {
	config := config.CreateTestConfiguration()
	api, _ := NewAPI(config)
	createTestUsers(api.userRep, 4, "test_firstname", "test_lastname", "test_email@fake.com", "pwd")
	s.api = api
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, &TestLoginSuite{})
}

func loginRequest(email string, password string, mux *chi.Mux) *http.Response {
	w := httptest.NewRecorder()
	user := UserLogin{Email: email, Password: password}
	buffer, err := convert(user)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, "/v1/login", buffer)
	if err != nil {
		panic(err)
	}
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	mux.ServeHTTP(w, req)
	return w.Result()
}

func (s *TestLoginSuite) TestLogin_ValidUser_Ok_Expected() {
	resp := loginRequest("1test_email@fake.com", "pwd1", s.api.mux)

	s.Assert().Equal(http.StatusOK, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	s.Assert().Nil(err)

	response := UserLoginResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusOK, response.Code)
	s.Assert().NotEmpty(response.Token)

	assertToken(s.Assert(), s.api, response.Token, "1test_email@fake.com")
}

func (s *TestLoginSuite) TestLogin_UnknownUser_401_Expected() {
	resp := loginRequest("0test_email@fake.com", "pwd1", s.api.mux)
	s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	s.Assert().Nil(err)

	response := HttpErrorResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusUnauthorized, response.Code)
	s.Assert().Equal("wrong user name or password", response.ErrorMessage)
}

func (s *TestLoginSuite) TestLogin_WrongPassword_401_Expected() {
	resp := loginRequest("1test_email@fake.com", "pwd2", s.api.mux)
	s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	s.Assert().Nil(err)

	response := HttpErrorResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusUnauthorized, response.Code)
	s.Assert().Equal("wrong user name or password", response.ErrorMessage)
}

func (s *TestLoginSuite) TestLogin_EmptyPassword_400_Expected() {

	resp := loginRequest("1test_email@fake.com", "", s.api.mux)
	s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	s.Assert().Nil(err)

	response := HttpErrorResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusBadRequest, response.Code)
	s.Assert().Equal("password must not be empty", response.ErrorMessage)
}

func (s *TestLoginSuite) TestLogin_EmptyEmail_400_Expected() {
	resp := loginRequest("", "pwd1", s.api.mux)

	s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	s.Assert().Nil(err)

	response := HttpErrorResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusBadRequest, response.Code)
	s.Assert().Equal("the provided email has incorrect format", response.ErrorMessage)
}

func assertToken(a *assert.Assertions, api *API, token, expectedEmail string) {
	jtoken, err := api.parseJwtToken(token)
	a.Nil(err)
	a.NotNil(token)
	a.True(jtoken.Valid)

	claims, ok := jtoken.Claims.(*jwtToken)
	a.True(ok)
	a.Equal(expectedEmail, claims.Email)
}
