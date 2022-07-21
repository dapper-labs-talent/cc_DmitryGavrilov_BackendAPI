package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dapper-labs/identity-server/config"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestUsersSuite struct {
	suite.Suite
	api               *API
	numberOfFakeUsers int
}

var (
	firstNameTemplate = "test_firstname"
	lastNameTemplate  = "test_lastname"
	emailTemplate     = "test_email@fake.com"
	passwordTemplate  = "pwd"
)

func (s *TestUsersSuite) SetupSuite() {
	config := config.CreateTestConfiguration()
	api, _ := NewAPI(config)
	s.numberOfFakeUsers = 5
	createTestUsers(api.userRep, s.numberOfFakeUsers, firstNameTemplate, lastNameTemplate, emailTemplate, passwordTemplate)
	s.api = api
}

func TestUsersTestSuite(t *testing.T) {
	suite.Run(t, &TestUsersSuite{})
}

func (s *TestUsersSuite) TestGetUsers_ValidXAuthToken_200_Expected() {
	resp := loginRequest("1test_email@fake.com", "pwd1", s.api.mux)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	respBody, err := ioutil.ReadAll(resp.Body)
	s.Assert().Nil(err)

	response := UserLoginResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusOK, response.Code)
	s.Assert().NotEmpty(response.Token)

	getAndAssertUsers(*s.Assert(), s.api.mux, "x-authentication-token", response.Token)
}

func (s *TestUsersSuite) TestGetUsers_ValidBearerToken_200_Expected() {
	resp := loginRequest("1test_email@fake.com", "pwd1", s.api.mux)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	respBody, err := ioutil.ReadAll(resp.Body)
	s.Assert().Nil(err)

	response := UserLoginResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusOK, response.Code)
	s.Assert().NotEmpty(response.Token)

	getAndAssertUsers(*s.Assert(), s.api.mux, "Authorization", "Bearer "+response.Token)
}

func (s *TestUsersSuite) TestGetUsers_NoToken_401_Expected() {
	loginResponse := loginRequest("1test_email@fake.com", "pwd1", s.api.mux)
	s.Assert().Equal(http.StatusOK, loginResponse.StatusCode)
	respBody, err := ioutil.ReadAll(loginResponse.Body)
	s.Assert().Nil(err)

	userLoginResponse := UserLoginResponse{}
	err = json.Unmarshal(respBody, &userLoginResponse)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusOK, userLoginResponse.Code)
	s.Assert().NotEmpty(userLoginResponse.Token)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/users", nil)

	s.api.mux.ServeHTTP(w, req)
	getUsersResponse := w.Result()
	s.Assert().Equal(http.StatusUnauthorized, getUsersResponse.StatusCode)
}

func getAndAssertUsers(guard assert.Assertions, api *chi.Mux, headerName, headerValue string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/users", nil)

	req.Header.Set(headerName, headerValue)
	api.ServeHTTP(w, req)
	resp := w.Result()

	guard.Equal(http.StatusOK, resp.StatusCode)
	respBody, err := ioutil.ReadAll(resp.Body)
	guard.Nil(err)

	getUsersResp := GetUsersResponse{}
	err = json.Unmarshal(respBody, &getUsersResp)
	guard.Nil(err)
	guard.Equal(5, len(*getUsersResp.Users))

	usedEmails := make(map[string]bool)
	for _, user := range *getUsersResp.Users {
		guard.Empty(user.PasswordHash)

		_, used := usedEmails[user.Email]
		guard.False(used)

		userSeqId := int(user.Email[0] - '0')
		expectedFirstname := fmt.Sprintf("%s%d", firstNameTemplate, userSeqId)
		guard.Equal(expectedFirstname, user.Firstname)

		expectedLastname := fmt.Sprintf("%s%d", lastNameTemplate, userSeqId)
		guard.Equal(expectedLastname, user.Lastname)

		expectedEmail := fmt.Sprintf("%d%s", userSeqId, emailTemplate)
		guard.Equal(expectedEmail, user.Email)

		usedEmails[user.Email] = true
	}
}
