package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"go-fhir-demo/internal/service/mocks"
	"go-fhir-demo/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ExternalPatientHandlerTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	mockService *mocks.MockExternalPatientServiceInterface
	handler     ExternalPatientHandlerInterface
	router      *gin.Engine
}

func (suite *ExternalPatientHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockService = mocks.NewMockExternalPatientServiceInterface(suite.mockCtrl)
	suite.handler = NewExternalPatientHandler(suite.mockService)
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.GET("/external-patients/:id", suite.handler.GetExternalPatientByID)
	router.GET("/external-patients", suite.handler.SearchExternalPatients)
	suite.router = router
}

func (suite *ExternalPatientHandlerTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestExternalPatientHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ExternalPatientHandlerTestSuite))
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByID_Success() {
	testID := "test-id-123"
	mockPatient := &fhir.Patient{Id: utils.CreateStringPtr(testID)}
	suite.mockService.EXPECT().
		GetExternalPatientByID(gomock.Any(), testID).
		Return(mockPatient, nil)

	req, _ := http.NewRequest("GET", "/external-patients/"+testID, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var resp fhir.Patient
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testID, *resp.Id)
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByID_Error() {
	testID := "notfound"
	suite.mockService.EXPECT().
		GetExternalPatientByID(gomock.Any(), testID).
		Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/external-patients/"+testID, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByID_BadRequest() {
	req, _ := http.NewRequest("GET", "/external-patients/", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	// This will not match the route, so 404 is expected
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *ExternalPatientHandlerTestSuite) TestSearchExternalPatients_Success() {
	mockBundle := &fhir.Bundle{Type: fhir.BundleTypeSearchset}
	suite.mockService.EXPECT().
		SearchExternalPatients(gomock.Any(), gomock.Any()).
		Return(mockBundle, nil)

	req, _ := http.NewRequest("GET", "/external-patients?name=John&gender=male", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var resp fhir.Bundle
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fhir.BundleTypeSearchset, resp.Type)
}

func (suite *ExternalPatientHandlerTestSuite) TestSearchExternalPatients_Error() {
	suite.mockService.EXPECT().
		SearchExternalPatients(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("search failed"))

	req, _ := http.NewRequest("GET", "/external-patients?name=Jane", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByIDCached_Success() {
	testID := "cached-id-456"
	mockPatient := &fhir.Patient{Id: utils.CreateStringPtr(testID)}
	suite.mockService.EXPECT().
		GetExternalPatientByIDCached(gomock.Any(), testID).
		Return(mockPatient, nil)

	router := gin.New()
	router.GET("/external-patients/:id/cached", suite.handler.GetExternalPatientByIDCached)
	req, _ := http.NewRequest("GET", "/external-patients/"+testID+"/cached", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var resp fhir.Patient
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testID, *resp.Id)
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByIDCached_Error() {
	testID := "cached-error"
	suite.mockService.EXPECT().
		GetExternalPatientByIDCached(gomock.Any(), testID).
		Return(nil, errors.New("cache miss"))

	router := gin.New()
	router.GET("/external-patients/:id/cached", suite.handler.GetExternalPatientByIDCached)
	req, _ := http.NewRequest("GET", "/external-patients/"+testID+"/cached", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByIDCached_BadRequest() {
	router := gin.New()
	router.GET("/external-patients/:id/cached", suite.handler.GetExternalPatientByIDCached)
	req, _ := http.NewRequest("GET", "/external-patients//cached", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByIDDelayed_Success() {
	testID := "delayed-id-789"
	mockPatient := &fhir.Patient{Id: utils.CreateStringPtr(testID)}
	suite.mockService.EXPECT().
		GetExternalPatientByIDDelayed(gomock.Any(), testID, gomock.Any()).
		Return(mockPatient, nil)

	router := gin.New()
	router.GET("/external-patients/:id/delayed", suite.handler.GetExternalPatientByIDDelayed)
	req, _ := http.NewRequest("GET", "/external-patients/"+testID+"/delayed?timeout=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var resp fhir.Patient
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testID, *resp.Id)
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByIDDelayed_Timeout() {
	testID := "timeout-id"
	suite.mockService.EXPECT().
		GetExternalPatientByIDDelayed(gomock.Any(), testID, gomock.Any()).
		Return(nil, context.DeadlineExceeded)

	router := gin.New()
	router.GET("/external-patients/:id/delayed", suite.handler.GetExternalPatientByIDDelayed)
	req, _ := http.NewRequest("GET", "/external-patients/"+testID+"/delayed?timeout=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusRequestTimeout, w.Code)
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByIDDelayed_Error() {
	testID := "delayed-error"
	suite.mockService.EXPECT().
		GetExternalPatientByIDDelayed(gomock.Any(), testID, gomock.Any()).
		Return(nil, errors.New("delayed error"))

	router := gin.New()
	router.GET("/external-patients/:id/delayed", suite.handler.GetExternalPatientByIDDelayed)
	req, _ := http.NewRequest("GET", "/external-patients/"+testID+"/delayed?timeout=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByIDDelayed_BadRequest_NoID() {
	router := gin.New()
	router.GET("/external-patients/:id/delayed", suite.handler.GetExternalPatientByIDDelayed)
	req, _ := http.NewRequest("GET", "/external-patients//delayed", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByIDDelayed_BadRequest_InvalidTimeout() {
	testID := "delayed-badtimeout"
	router := gin.New()
	router.GET("/external-patients/:id/delayed", suite.handler.GetExternalPatientByIDDelayed)
	req, _ := http.NewRequest("GET", "/external-patients/"+testID+"/delayed?timeout=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}
