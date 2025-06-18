package service

import (
	"errors"
	"testing"
	"time"

	"go-fhir-demo/internal/domain"
	"go-fhir-demo/internal/domain/mocks"
	"go-fhir-demo/pkg/utils"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// PatientServiceTestSuite defines the test suite
type PatientServiceTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	mockRepo *mocks.MockPatientRepository
	service  PatientServiceInterface
}

// SetupTest initializes the test suite before each test
func (suite *PatientServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockPatientRepository(suite.ctrl)
	suite.service = NewPatientService(suite.mockRepo)
}

// TearDownTest cleans up after each test
func (suite *PatientServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

// TestCreatePatient_Success tests successful patient creation
func (suite *PatientServiceTestSuite) TestCreatePatient_Success() {
	// Arrange
	fhirPatient := &fhir.Patient{
		Active: utils.CreateBoolPtr(true),
		Name: []fhir.HumanName{
			{
				Family: utils.CreateStringPtr("Doe"),
				Given:  []string{"John"},
			},
		},
		Gender:    utils.GenderPtr("male"),
		BirthDate: utils.CreateStringPtr("1980-01-01"),
	}

	suite.mockRepo.EXPECT().
		Create(gomock.Any()).
		Return(nil).
		Times(1)

	// Act
	patient, err := suite.service.CreatePatient(fhirPatient)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), "Doe", patient.Family)
	assert.Equal(suite.T(), "John", patient.Given)
	assert.Equal(suite.T(), "male", patient.Gender)
}

// TestCreatePatient_Error tests patient creation error
func (suite *PatientServiceTestSuite) TestCreatePatient_Error() {
	// Arrange
	fhirPatient := &fhir.Patient{
		Active: utils.CreateBoolPtr(true),
		Name: []fhir.HumanName{
			{
				Family: utils.CreateStringPtr("Doe"),
				Given:  []string{"John"},
			},
		},
	}

	suite.mockRepo.EXPECT().
		Create(gomock.Any()).
		Return(errors.New("database error")).
		Times(1)

	// Act
	patient, err := suite.service.CreatePatient(fhirPatient)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), patient)
	assert.Contains(suite.T(), err.Error(), "database error")
}

// TestGetPatient_Success tests successful patient retrieval
func (suite *PatientServiceTestSuite) TestGetPatient_Success() {
	// Arrange
	patientID := uint(1)
	expectedPatient := &domain.Patient{
		ID:     patientID,
		Family: "Doe",
		Given:  "John",
		Gender: "male",
	}

	suite.mockRepo.EXPECT().
		GetByID(patientID).
		Return(expectedPatient, nil).
		Times(1)

	// Act
	patient, err := suite.service.GetPatient(patientID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), patientID, patient.ID)
	assert.Equal(suite.T(), "Doe", patient.Family)
}

// TestGetPatient_NotFound tests patient not found scenario
func (suite *PatientServiceTestSuite) TestGetPatient_NotFound() {
	// Arrange
	patientID := uint(999)
	suite.mockRepo.EXPECT().
		GetByID(patientID).
		Return(nil, errors.New("patient not found")).
		Times(1)

	// Act
	patient, err := suite.service.GetPatient(patientID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), patient)
}

// TestGetPatients_Success tests successful patients retrieval with pagination
func (suite *PatientServiceTestSuite) TestGetPatients_Success() {
	// Arrange
	limit, offset := 10, 0
	expectedPatients := []*domain.Patient{
		{ID: 1, Family: "Doe", Given: "John"},
		{ID: 2, Family: "Smith", Given: "Jane"},
	}
	expectedCount := int64(2)

	suite.mockRepo.EXPECT().
		GetAll(limit, offset).
		Return(expectedPatients, nil).
		Times(1)

	suite.mockRepo.EXPECT().
		Count().
		Return(expectedCount, nil).
		Times(1)

	// Act
	patients, count, err := suite.service.GetPatients(limit, offset)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patients)
	assert.Equal(suite.T(), len(expectedPatients), len(patients))
	assert.Equal(suite.T(), expectedCount, count)
}

// TestGetPatients_Error tests patients retrieval error
func (suite *PatientServiceTestSuite) TestGetPatients_Error() {
	// Arrange
	limit, offset := 10, 0
	suite.mockRepo.EXPECT().
		GetAll(limit, offset).
		Return(nil, errors.New("database error")).
		Times(1)

	// Act
	patients, count, err := suite.service.GetPatients(limit, offset)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), patients)
	assert.Equal(suite.T(), int64(0), count)
}

// TestUpdatePatient_Success tests successful patient update
func (suite *PatientServiceTestSuite) TestUpdatePatient_Success() {
	// Arrange
	patientID := uint(1)
	existingPatient := &domain.Patient{
		ID:        patientID,
		Family:    "Doe",
		Given:     "John",
		CreatedAt: time.Now(),
	}

	updatedFhirPatient := &fhir.Patient{
		Active: utils.CreateBoolPtr(true),
		Name: []fhir.HumanName{
			{
				Family: utils.CreateStringPtr("Updated"),
				Given:  []string{"Jane"},
			},
		},
	}

	suite.mockRepo.EXPECT().
		GetByID(patientID).
		Return(existingPatient, nil).
		Times(1)

	suite.mockRepo.EXPECT().
		Update(gomock.Any()).
		Return(nil).
		Times(1)

	// Act
	patient, err := suite.service.UpdatePatient(patientID, updatedFhirPatient)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), patientID, patient.ID)
	assert.Equal(suite.T(), "Updated", patient.Family)
	assert.Equal(suite.T(), "Jane", patient.Given)
}

// TestUpdatePatient_NotFound tests update when patient not found
func (suite *PatientServiceTestSuite) TestUpdatePatient_NotFound() {
	// Arrange
	patientID := uint(999)
	updatedFhirPatient := &fhir.Patient{
		Name: []fhir.HumanName{
			{
				Family: utils.CreateStringPtr("Updated"),
			},
		},
	}

	suite.mockRepo.EXPECT().
		GetByID(patientID).
		Return(nil, errors.New("patient not found")).
		Times(1)

	// Act
	patient, err := suite.service.UpdatePatient(patientID, updatedFhirPatient)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), patient)
}

// TestDeletePatient_Success tests successful patient deletion
func (suite *PatientServiceTestSuite) TestDeletePatient_Success() {
	// Arrange
	patientID := uint(1)
	suite.mockRepo.EXPECT().
		Delete(patientID).
		Return(nil).
		Times(1)

	// Act
	err := suite.service.DeletePatient(patientID)

	// Assert
	assert.NoError(suite.T(), err)
}

// TestDeletePatient_Error tests patient deletion error
func (suite *PatientServiceTestSuite) TestDeletePatient_Error() {
	// Arrange
	patientID := uint(1)
	suite.mockRepo.EXPECT().
		Delete(patientID).
		Return(errors.New("delete failed")).
		Times(1)

	// Act
	err := suite.service.DeletePatient(patientID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "delete failed")
}

// TestConvertToFHIR_Success tests successful conversion to FHIR
func (suite *PatientServiceTestSuite) TestConvertToFHIR_Success() {
	// Arrange
	patient := &domain.Patient{
		ID:       1,
		FHIRData: []byte(`{"resourceType":"Patient","active":true,"name":[{"family":"Doe","given":["John"]}]}`),
	}

	// Act
	fhirPatient, err := suite.service.ConvertToFHIR(patient)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), fhirPatient)
	assert.True(suite.T(), *fhirPatient.Active)
	assert.Equal(suite.T(), "Doe", *fhirPatient.Name[0].Family)
}

// TestConvertToFHIR_InvalidJSON tests conversion with invalid JSON
func (suite *PatientServiceTestSuite) TestConvertToFHIR_InvalidJSON() {
	// Arrange
	patient := &domain.Patient{
		ID:       1,
		FHIRData: []byte(`invalid json`),
	}

	// Act
	fhirPatient, err := suite.service.ConvertToFHIR(patient)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), fhirPatient)
	assert.Contains(suite.T(), err.Error(), "failed to unmarshal FHIR data")
}

// TestConvertFromFHIR_Success tests successful conversion from FHIR
func (suite *PatientServiceTestSuite) TestConvertFromFHIR_Success() {
	// Arrange
	fhirPatient := &fhir.Patient{
		Active: utils.CreateBoolPtr(true),
		Name: []fhir.HumanName{
			{
				Family: utils.CreateStringPtr("Doe"),
				Given:  []string{"John"},
			},
		},
		Gender:    utils.GenderPtr("male"),
		BirthDate: utils.CreateStringPtr("1980-01-01"),
	}

	// Act
	patient, err := suite.service.ConvertFromFHIR(fhirPatient)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), "Doe", patient.Family)
	assert.Equal(suite.T(), "John", patient.Given)
	assert.Equal(suite.T(), "male", patient.Gender)
	assert.NotNil(suite.T(), patient.BirthDate)
}

// TestPatchPatient_Success tests successful patient patching
func (suite *PatientServiceTestSuite) TestPatchPatient_Success() {
	// Arrange
	patientID := uint(1)
	existingPatient := &domain.Patient{
		ID:       patientID,
		FHIRData: []byte(`{"resourceType":"Patient","active":true,"name":[{"family":"Doe","given":["John"]}]}`),
	}

	updates := map[string]interface{}{
		"family": "Updated",
		"active": false,
	}

	suite.mockRepo.EXPECT().
		GetByID(patientID).
		Return(existingPatient, nil).
		Times(1)

	suite.mockRepo.EXPECT().
		Update(gomock.Any()).
		Return(nil).
		Times(1)

	// Act
	patient, err := suite.service.PatchPatient(patientID, updates)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), patientID, patient.ID)
	assert.Equal(suite.T(), "Updated", patient.Family)
}

// TestPatientServiceTestSuite runs the test suite
func TestPatientServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PatientServiceTestSuite))
}
