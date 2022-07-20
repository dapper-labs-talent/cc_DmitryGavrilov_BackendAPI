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

func (s *TestLoginSuite) TestLogin_ValidUser_Ok_Expected() {
	w := httptest.NewRecorder()
	user := UserLogin{Email: "1test_email@fake.com", Password: "pwd1"}
	buffer, err := convert(user)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, "/v1/login", buffer)
	if err != nil {
		panic(err)
	}
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	s.api.mux.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Result().StatusCode)

	respBody, err := ioutil.ReadAll(w.Result().Body)
	s.Assert().Nil(err)

	response := UserLoginResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusOK, response.Code)
	s.Assert().NotEmpty(response.Token)

	assertToken(s.Assert(), response.Token, user.Email, s.api.config.Secret)
}

func (s *TestLoginSuite) TestLogin_UnknownUser_401_Expected() {
	w := httptest.NewRecorder()
	user := UserLogin{Email: "0test_email@fake.com", Password: "pwd1"}
	buffer, err := convert(user)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, "/v1/login", buffer)
	if err != nil {
		panic(err)
	}
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	s.api.mux.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusUnauthorized, w.Result().StatusCode)

	respBody, err := ioutil.ReadAll(w.Result().Body)
	s.Assert().Nil(err)

	response := HttpErrorResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusUnauthorized, response.Code)
	s.Assert().Equal("wrong user name or password", response.ErrorMessage)
}

func (s *TestLoginSuite) TestLogin_WrongPassword_401_Expected() {
	w := httptest.NewRecorder()
	user := UserLogin{Email: "1test_email@fake.com", Password: "pwd2"}
	buffer, err := convert(user)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, "/v1/login", buffer)
	if err != nil {
		panic(err)
	}
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	s.api.mux.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusUnauthorized, w.Result().StatusCode)

	respBody, err := ioutil.ReadAll(w.Result().Body)
	s.Assert().Nil(err)

	response := HttpErrorResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusUnauthorized, response.Code)
	s.Assert().Equal("wrong user name or password", response.ErrorMessage)
}

func (s *TestLoginSuite) TestLogin_EmptyPassword_400_Expected() {
	w := httptest.NewRecorder()
	user := UserLogin{Email: "1test_email@fake.com", Password: ""}
	buffer, err := convert(user)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, "/v1/login", buffer)
	if err != nil {
		panic(err)
	}
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	s.api.mux.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusBadRequest, w.Result().StatusCode)

	respBody, err := ioutil.ReadAll(w.Result().Body)
	s.Assert().Nil(err)

	response := HttpErrorResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusBadRequest, response.Code)
	s.Assert().Equal("password must not be empty", response.ErrorMessage)
}

func (s *TestLoginSuite) TestLogin_EmptyEmail_400_Expected() {
	w := httptest.NewRecorder()
	user := UserLogin{Email: "", Password: "pwd1"}
	buffer, err := convert(user)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, "/v1/login", buffer)
	if err != nil {
		panic(err)
	}
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	s.api.mux.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusBadRequest, w.Result().StatusCode)

	respBody, err := ioutil.ReadAll(w.Result().Body)
	s.Assert().Nil(err)

	response := HttpErrorResponse{}
	err = json.Unmarshal(respBody, &response)
	s.Assert().Nil(err)
	s.Assert().Equal(http.StatusBadRequest, response.Code)
	s.Assert().Equal("the provided email has incorrect format", response.ErrorMessage)
}

func assertToken(a *assert.Assertions, token, expectedEmail, secret string) {
	jtoken, err := parseJwtToken(token, secret)
	a.Nil(err)
	a.NotNil(token)
	a.True(jtoken.Valid)

	claims, ok := jtoken.Claims.(*jwtToken)
	a.True(ok)
	a.Equal(expectedEmail, claims.Email)
}
