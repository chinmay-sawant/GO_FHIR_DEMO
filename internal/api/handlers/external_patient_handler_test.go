package handlers

import (
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
		GetExternalPatientByID(testID).
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
		GetExternalPatientByID(testID).
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
		SearchExternalPatients(gomock.Any()).
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
		SearchExternalPatients(gomock.Any()).
		Return(nil, errors.New("search failed"))

	req, _ := http.NewRequest("GET", "/external-patients?name=Jane", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
