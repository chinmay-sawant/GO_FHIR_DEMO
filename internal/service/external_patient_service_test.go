package service

import (
	"errors"
	"testing"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ExternalPatientServiceTestSuite defines the test suite
type ExternalPatientServiceTestSuite struct {
	suite.Suite
	mockClient *MockFHIRClient
	service    ExternalPatientServiceInterface
}

// MockFHIRClient is a mock implementation of the FHIR client
type MockFHIRClient struct {
	mock.Mock
}

// SetupTest initializes the test suite before each test
func (suite *ExternalPatientServiceTestSuite) SetupTest() {
	suite.mockClient = &MockFHIRClient{}
	suite.service = NewExternalPatientService(suite.mockClient)
}

// TestExternalPatientServiceTestSuite runs the test suite
func TestExternalPatientServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ExternalPatientServiceTestSuite))
}

func (m *MockFHIRClient) GetPatientByID(id string) (*fhir.Patient, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fhir.Patient), args.Error(1)
}

func (m *MockFHIRClient) SearchPatients(params map[string]string) (*fhir.Bundle, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fhir.Bundle), args.Error(1)
}

// FHIRClientInterface defines the contract for FHIR client
type FHIRClientInterface interface {
	GetPatientByID(id string) (*fhir.Patient, error)
	SearchPatients(params map[string]string) (*fhir.Bundle, error)
}

// TestGetExternalPatientByID_Success tests successful patient retrieval
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByID_Success() {
	// Arrange
	testID := "test-id-123"
	mockPatient := &fhir.Patient{
		Id: &testID,
	}
	suite.mockClient.On("GetPatientByID", testID).Return(mockPatient, nil)

	// Act
	patient, err := suite.service.GetExternalPatientByID(testID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), testID, *patient.Id)
}

// TestGetExternalPatientByID_Error tests error handling
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByID_Error() {
	// Arrange
	testID := "notfound"
	suite.mockClient.On("GetPatientByID", testID).Return(nil, errors.New("patient not found"))

	// Act
	patient, err := suite.service.GetExternalPatientByID(testID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), patient)
	assert.Contains(suite.T(), err.Error(), "patient not found")
}

// TestSearchExternalPatients_Success tests successful patient search
func (suite *ExternalPatientServiceTestSuite) TestSearchExternalPatients_Success() {
	// Arrange
	searchParams := map[string]string{"name": "John", "gender": "male"}
	mockBundle := &fhir.Bundle{
		Type: fhir.BundleTypeSearchset,
	}
	suite.mockClient.On("SearchPatients", searchParams).Return(mockBundle, nil)

	// Act
	bundle, err := suite.service.SearchExternalPatients(searchParams)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), bundle)
	assert.Equal(suite.T(), fhir.BundleTypeSearchset, bundle.Type)
}

// TestSearchExternalPatients_Error tests search error handling
func (suite *ExternalPatientServiceTestSuite) TestSearchExternalPatients_Error() {
	// Arrange
	searchParams := map[string]string{"name": "Jane"}
	suite.mockClient.On("SearchPatients", searchParams).Return(nil, errors.New("search failed"))

	// Act
	bundle, err := suite.service.SearchExternalPatients(searchParams)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), bundle)
	assert.Contains(suite.T(), err.Error(), "search failed")
}

// TestSearchExternalPatients_EmptyParams tests search with empty parameters
func (suite *ExternalPatientServiceTestSuite) TestSearchExternalPatients_EmptyParams() {
	// Arrange
	emptyParams := map[string]string{}
	mockBundle := &fhir.Bundle{
		Type: fhir.BundleTypeSearchset,
	}
	suite.mockClient.On("SearchPatients", emptyParams).Return(mockBundle, nil)

	// Act
	bundle, err := suite.service.SearchExternalPatients(emptyParams)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), bundle)
}
