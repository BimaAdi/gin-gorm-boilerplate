package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/BimaAdi/ginGormBoilerplate/core"
	"github.com/BimaAdi/ginGormBoilerplate/migrations"
	"github.com/BimaAdi/ginGormBoilerplate/models"
	"github.com/BimaAdi/ginGormBoilerplate/routes"
	"github.com/BimaAdi/ginGormBoilerplate/schemas"
	"github.com/BimaAdi/ginGormBoilerplate/settings"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MigrateAuthTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *MigrateAuthTestSuite) SetupSuite() {

	settings.InitiateSettings("../.env")
	models.Initiate()
	migrations.MigrateUp("../.env", "file://../migrations/migrations_files/")
	router := gin.Default()
	suite.router = routes.GetRoutes(router)
}

func (suite *MigrateAuthTestSuite) SetupTest() {
	models.ClearAllData()

}

func (suite *MigrateAuthTestSuite) TestLoginSuccess() {
	// Given
	timeZoneAsiaJakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err.Error())
	}
	hashPasword, err := core.HashPassword("Fakepassword")
	if err != nil {
		panic(err.Error())
	}
	user_login := models.User{
		Email:       "test@test.com",
		Username:    "test",
		Password:    hashPasword,
		IsActive:    true,
		IsSuperuser: true,
		CreatedAt:   time.Date(2022, 10, 5, 10, 0, 0, 0, timeZoneAsiaJakarta),
	}
	models.DBConn.Create(&user_login)

	// When
	var param = url.Values{}
	param.Set("username", "test")
	param.Set("password", "Fakepassword")
	var payload = bytes.NewBufferString(param.Encode())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	suite.router.ServeHTTP(w, req)

	// Expect
	assert.Equal(suite.T(), 200, w.Code)
	jsonResponse := schemas.LoginResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")
}

func (suite *MigrateAuthTestSuite) TestLoginFailed() {
	// Given
	timeZoneAsiaJakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err.Error())
	}
	hashPasword, err := core.HashPassword("Fakepassword")
	if err != nil {
		panic(err.Error())
	}
	user_login := models.User{
		Email:       "test@test.com",
		Username:    "test",
		Password:    hashPasword,
		IsActive:    true,
		IsSuperuser: true,
		CreatedAt:   time.Date(2022, 10, 5, 10, 0, 0, 0, timeZoneAsiaJakarta),
	}
	models.DBConn.Create(&user_login)

	// When
	var param = url.Values{}
	param.Set("username", "test")
	param.Set("password", "wrong password")
	var payload = bytes.NewBufferString(param.Encode())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	suite.router.ServeHTTP(w, req)

	// Expect
	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *MigrateAuthTestSuite) TestLogoutSuccess() {
	// Given
	// create request user
	timeZoneAsiaJakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err.Error())
	}
	request_user := models.User{
		Email:       "a@test.com",
		Username:    "a",
		Password:    "Fakepassword",
		IsActive:    true,
		IsSuperuser: true,
		CreatedAt:   time.Date(2022, 10, 5, 10, 0, 0, 0, timeZoneAsiaJakarta),
	}
	models.DBConn.Create(&request_user)
	token, err := core.GenerateJWTTokenFromUser(models.DBConn, request_user)
	if err != nil {
		panic(err.Error())
	}

	// When
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("authorization", "Bearer "+token)
	suite.router.ServeHTTP(w, req)

	// Expect
	assert.Equal(suite.T(), 200, w.Code)
	jsonResponse := schemas.LogoutResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")
}

func (suite *MigrateAuthTestSuite) TestLogoutInvalidToken() {
	// Given
	token := "theinvalidtoken"

	// When
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("authorization", "Bearer "+token)
	suite.router.ServeHTTP(w, req)

	// Expect
	assert.Equal(suite.T(), 401, w.Code)
	jsonResponse := schemas.UnauthorizedResponse{}
	err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")
}

func (suite *MigrateAuthTestSuite) TearDownTest() {
	models.ClearAllData()
}

func TestMigrateAuthTestSuite(t *testing.T) {
	suite.Run(t, new(MigrateAuthTestSuite))
}
