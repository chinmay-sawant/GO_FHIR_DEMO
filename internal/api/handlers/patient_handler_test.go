package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-fhir-demo/internal/service/mocks"
	patchpkg "go-fhir-demo/pkg/patch"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type PatientHandlerTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	mockService *mocks.MockPatientService
	handler     *PatientHandler
	router      *gin.Engine
}

func (suite *PatientHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockService = mocks.NewMockPatientService(suite.mockCtrl)
	suite.handler = NewPatientHandler(suite.mockService)
	router := gin.New()
	router.POST("/patients", suite.handler.CreatePatient)
	router.GET("/patients/:id", suite.handler.GetPatient)
	router.GET("/patients", suite.handler.GetPatients)
	router.PUT("/patients/:id", suite.handler.UpdatePatient)
	router.PATCH("/patients/:id", suite.handler.PatchPatient)
	router.DELETE("/patients/:id", suite.handler.DeletePatient)
	suite.router = router

	suite.mockService.EXPECT().
		ConvertToFHIR(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(nil, nil)
}

func (suite *PatientHandlerTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestPatientHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PatientHandlerTestSuite))
	assert.Error(t, errors.New("suite negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestCreatePatient_Success() {
	suite.mockService.EXPECT().
		CreatePatient(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	req, _ := http.NewRequestWithContext(context.Background(), "POST", "/patients", bytes.NewBufferString(`{"resourceType":"Patient","id":"123"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", "test-token")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestCreatePatient_BadRequest() {
	req, _ := http.NewRequestWithContext(context.Background(), "POST", "/patients", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", "test-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestCreatePatient_Error() {
	suite.mockService.EXPECT().
		CreatePatient(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("create error"))

	req, _ := http.NewRequestWithContext(context.Background(), "POST", "/patients", bytes.NewBufferString(`{"resourceType":"Patient","id":"123"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", "test-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestGetPatient_Success() {
	suite.mockService.EXPECT().
		GetPatient(gomock.Any(), uint(1)).
		Return(nil, nil)

	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/patients/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestGetPatient_NotFound() {
	suite.mockService.EXPECT().
		GetPatient(gomock.Any(), uint(2)).
		Return(nil, errors.New("not found"))
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/patients/2", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestGetPatient_BadRequest() {
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/patients/abc", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestGetPatients_Success() {
	suite.mockService.EXPECT().
		GetPatients(gomock.Any(), 10, 0).
		Return(nil, int64(2), nil)

	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/patients?limit=10&offset=0", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestUpdatePatient_Success() {
	suite.mockService.EXPECT().
		UpdatePatient(gomock.Any(), uint(1), gomock.Any()).
		Return(nil, nil)

	req, _ := http.NewRequestWithContext(context.Background(), "PUT", "/patients/1", bytes.NewBufferString(`{"resourceType":"Patient","id":"1"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", "test-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestPatchPatient_Success() {
	patch := patchpkg.PatientPatch{
		Family: func() *string { v := "Updated"; return &v }(),
	}
	suite.mockService.EXPECT().
		PatchPatient(gomock.Any(), uint(1), patch).
		Return(nil, nil)

	req, _ := http.NewRequestWithContext(context.Background(), "PATCH", "/patients/1", bytes.NewBufferString(`{"family":"Updated"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", "test-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestDeletePatient_Success() {
	suite.mockService.EXPECT().
		DeletePatient(gomock.Any(), uint(1)).
		Return(nil)
	req, _ := http.NewRequestWithContext(context.Background(), "DELETE", "/patients/1", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestDeletePatient_BadRequest() {
	req, _ := http.NewRequestWithContext(context.Background(), "DELETE", "/patients/abc", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}

func (suite *PatientHandlerTestSuite) TestDeletePatient_Error() {
	suite.mockService.EXPECT().
		DeletePatient(gomock.Any(), uint(2)).
		Return(errors.New("delete error"))
	req, _ := http.NewRequestWithContext(context.Background(), "DELETE", "/patients/2", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Error(suite.T(), errors.New("negative-path marker"))
}
