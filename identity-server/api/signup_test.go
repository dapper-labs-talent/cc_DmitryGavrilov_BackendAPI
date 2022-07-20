// package api
package api

import (
	"testing"

	"github.com/dapper-labs/identity-server/config"
	"github.com/stretchr/testify/suite"
)

// import (
// 	"code_signal_rate_limiter/internal/config"
// 	"code_signal_rate_limiter/internal/ratelimitter"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/suite"
// )

type TestSuite struct {
	suite.Suite
	api *API
}

func (s *TestSuite) SetupSuite() {
	config := config.CreateTestConfiguration()
	api, _ := NewAPI(config)
	return &TestSuite{
		api: api,
	}

}

func TestTakeTestSuite(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// func (s *TakeTestSuite) SetupSuite() {
// 	s.router = gin.Default()

// 	config := config.CreateTestConfigWithEndpoints()
// 	s.rateApi, _ = NewRateLimitterAPI(config)

// 	v1 := s.router.Group("/v1")
// 	addV1Routes(v1, config, s.rateApi)
// }

// type TakeTestSuite struct {
// 	suite.Suite
// 	router  *gin.Engine
// 	rateApi *RateLimitterAPI
// }

// func TestTakeTestSuite(t *testing.T) {
// 	suite.Run(t, &TakeTestSuite{})
// }

// func (s *TakeTestSuite) SetupSuite() {
// 	s.router = gin.Default()

// 	config := config.CreateTestConfigWithEndpoints()
// 	s.rateApi, _ = NewRateLimitterAPI(config)

// 	v1 := s.router.Group("/v1")
// 	addV1Routes(v1, config, s.rateApi)
// }

// func (s *TakeTestSuite) TestTakeGetRespondsWithGetUsers200() {
// 	// Endpoint: "GET /user/:id",
// 	// Burst : 10,
// 	// Sustained : 30,
// 	burst := 10
// 	sustained := 30
// 	url := "/v1/take/get/user/123"
// 	sendTestRequests(s, burst, sustained, url)
// }

// func (s *TakeTestSuite) TestTakeGetRespondsWithPatchUsers200() {
// 	// Endpoint: "PATCH /user/:id",
// 	// Burst : 5,
// 	// Sustained : 30,
// 	burst := 5
// 	sustained := 30
// 	url := "/v1/take/patch/user/123"
// 	sendTestRequests(s, burst, sustained, url)
// }

// func (s *TakeTestSuite) TestTakeGetRespondsPostUserInfo200() {
// 	//Endpoint: "POST /userinfo",
// 	//Burst : 5,
// 	//Sustained : 60,
// 	burst := 5
// 	sustained := 60
// 	url := "/v1/take/post/userinfo"
// 	sendTestRequests(s, burst, sustained, url)
// }

// func sendTestRequests(s *TakeTestSuite, burst int, sustained int, url string) {
// 	template := `{"allow":%t,"code":200,"tokens":%d}`
// 	configuredLimit := ratelimitter.EveryMinute(float64(sustained))
// 	tokens := burst

// 	startTime := time.Now()

// 	for i := 0; i < 20; i++ {
// 		w := httptest.NewRecorder()
// 		req, _ := http.NewRequest(http.MethodGet, url, nil)
// 		s.router.ServeHTTP(w, req)

// 		s.Assert().Equal(http.StatusOK, w.Result().StatusCode)

// 		responseBody, err := ioutil.ReadAll(w.Result().Body)
// 		s.Assert().NoError(err)

// 		expected := ""

// 		now := time.Now()
// 		earned := int(configuredLimit.TokensFromTime(now.Sub(startTime)))

// 		if tokens > 0 || earned > 0 {
// 			if earned > 0 {
// 				tokens += earned
// 				startTime = now
// 			}
// 			tokens--
// 			expected = fmt.Sprintf(template, true, tokens)
// 		} else {
// 			expected = fmt.Sprintf(template, false, 0)
// 			tokens = 0
// 		}

// 		s.Assert().Equal(expected, string(responseBody))

// 		// let's give some time to earn a few tokens
// 		if i == 15 {
// 			time.Sleep(2 * time.Second)
// 		}
// 	}
// }
