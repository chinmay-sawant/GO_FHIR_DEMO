package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-fhir-demo/internal/domain"
	"go-fhir-demo/internal/domain/mocks"
	"go-fhir-demo/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type PatientHandlerTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	mockService *mocks.MockPatientService
	handler     PatientHandlerInterface
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

	// Globally mock ConvertToFHIR for any input
	suite.mockService.EXPECT().
		ConvertToFHIR(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&fhir.Patient{Id: utils.CreateStringPtr("mocked")}, nil)
}

func (suite *PatientHandlerTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestPatientHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PatientHandlerTestSuite))
}

func (suite *PatientHandlerTestSuite) TestCreatePatient_Success() {
	id := utils.CreateStringPtr("123")
	active := true
	gender := fhir.AdministrativeGenderMale
	birthDate := "1990-01-01"
	fhirPatient := &fhir.Patient{
		Id:        id,
		Active:    &active,
		Gender:    &gender,
		BirthDate: &birthDate,
		Name: []fhir.HumanName{
			{
				Family: utils.CreateStringPtr("Doe"),
				Given:  []string{"John"},
			},
		},
	}
	domainPatient := &domain.Patient{ID: 1}
	suite.mockService.EXPECT().
		CreatePatient(gomock.Any(), fhirPatient).
		Return(domainPatient, nil)

	body, _ := json.Marshal(fhirPatient)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	var resp fhir.Patient
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "mocked", *resp.Id)
}

func (suite *PatientHandlerTestSuite) TestCreatePatient_BadRequest() {
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *PatientHandlerTestSuite) TestCreatePatient_Error() {
	fhirPatient := &fhir.Patient{Id: utils.CreateStringPtr("123")}
	suite.mockService.EXPECT().
		CreatePatient(gomock.Any(), fhirPatient).
		Return(&domain.Patient{}, errors.New("create error"))

	body, _ := json.Marshal(fhirPatient)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *PatientHandlerTestSuite) TestGetPatient_Success() {
	domainPatient := &domain.Patient{ID: 1}
	suite.mockService.EXPECT().
		GetPatient(gomock.Any(), uint(1)).
		Return(domainPatient, nil)

	req, _ := http.NewRequest("GET", "/patients/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *PatientHandlerTestSuite) TestGetPatient_NotFound() {
	suite.mockService.EXPECT().
		GetPatient(gomock.Any(), uint(2)).
		Return(nil, errors.New("not found"))
	req, _ := http.NewRequest("GET", "/patients/2", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *PatientHandlerTestSuite) TestGetPatient_BadRequest() {
	req, _ := http.NewRequest("GET", "/patients/abc", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *PatientHandlerTestSuite) TestGetPatients_Success() {
	active := true

	domainPatients := []*domain.Patient{
		{ID: 1, Family: "Doe", Given: "John", Gender: "male", Active: &active},
		{ID: 2, Family: "Smith", Given: "Jane", Gender: "female", Active: &active},
	}
	suite.mockService.EXPECT().
		GetPatients(gomock.Any(), 10, 0).
		Return(domainPatients, int64(2), nil)

	req, _ := http.NewRequest("GET", "/patients?limit=10&offset=0", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *PatientHandlerTestSuite) TestUpdatePatient_Success() {
	id := utils.CreateStringPtr("1")
	active := true
	gender := fhir.AdministrativeGenderMale
	birthDate := "1990-01-01"
	fhirPatient := &fhir.Patient{
		Id:        id,
		Active:    &active,
		Gender:    &gender,
		BirthDate: &birthDate,
		Name: []fhir.HumanName{
			{
				Family: utils.CreateStringPtr("Gamma"),
				Given:  []string{"G"},
			},
		},
	}
	domainPatient := &domain.Patient{ID: 1}
	suite.mockService.EXPECT().
		UpdatePatient(gomock.Any(), uint(1), fhirPatient).
		Return(domainPatient, nil)

	body, _ := json.Marshal(fhirPatient)
	req, _ := http.NewRequest("PUT", "/patients/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *PatientHandlerTestSuite) TestPatchPatient_Success() {
	patch := map[string]interface{}{"family": "Updated"}
	domainPatient := &domain.Patient{ID: 1, Family: "Updated"}
	suite.mockService.EXPECT().
		PatchPatient(gomock.Any(), uint(1), patch).
		Return(domainPatient, nil)

	body, _ := json.Marshal(patch)
	req, _ := http.NewRequest("PATCH", "/patients/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *PatientHandlerTestSuite) TestDeletePatient_Success() {
	suite.mockService.EXPECT().
		DeletePatient(gomock.Any(), uint(1)).
		Return(nil)
	req, _ := http.NewRequest("DELETE", "/patients/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *PatientHandlerTestSuite) TestDeletePatient_BadRequest() {
	req, _ := http.NewRequest("DELETE", "/patients/abc", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *PatientHandlerTestSuite) TestDeletePatient_Error() {
	suite.mockService.EXPECT().
		DeletePatient(gomock.Any(), uint(2)).
		Return(errors.New("delete error"))
	req, _ := http.NewRequest("DELETE", "/patients/2", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
