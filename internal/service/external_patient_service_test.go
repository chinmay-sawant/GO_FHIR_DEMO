package service

import (
	"errors"
	fhirclientmocks "go-fhir-demo/pkg/fhirclient/mocks"
	"testing"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// ExternalPatientServiceTestSuite defines the test suite
type ExternalPatientServiceTestSuite struct {
	suite.Suite
	mockClient *fhirclientmocks.MockClientInterface
	service    ExternalPatientServiceInterface
}

// SetupTest initializes the test suite before each test
func (suite *ExternalPatientServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockClient = fhirclientmocks.NewMockClientInterface(ctrl)
	suite.service = NewExternalPatientService(suite.mockClient)
}

// TestExternalPatientServiceTestSuite runs the test suite
func TestExternalPatientServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ExternalPatientServiceTestSuite))
}

// TestGetExternalPatientByID_Success tests successful patient retrieval
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByID_Success() {
	// Arrange
	testID := "test-id-123"
	mockPatient := &fhir.Patient{
		Id: &testID,
	}
	suite.mockClient.EXPECT().GetPatientByID(testID).Return(mockPatient, nil)

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
	suite.mockClient.EXPECT().GetPatientByID(testID).Return(nil, errors.New("patient not found"))

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
	suite.mockClient.EXPECT().SearchPatients(searchParams).Return(mockBundle, nil)

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
	suite.mockClient.EXPECT().SearchPatients(searchParams).Return(nil, errors.New("search failed"))

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
	suite.mockClient.EXPECT().SearchPatients(emptyParams).Return(mockBundle, nil)

	// Act
	bundle, err := suite.service.SearchExternalPatients(emptyParams)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), bundle)
}

// TestCreateExternalPatient_Success tests successful creation of external patient
func (suite *ExternalPatientServiceTestSuite) TestCreateExternalPatient_Success() {
	// Arrange
	mockPatient := &fhir.Patient{Id: func() *string { s := "new-id"; return &s }()}
	suite.mockClient.EXPECT().CreatePatient(mockPatient).Return(mockPatient, nil)

	// Act
	created, err := suite.service.CreateExternalPatient(mockPatient)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), created)
	assert.Equal(suite.T(), *mockPatient.Id, *created.Id)
}

// TestCreateExternalPatient_Error tests error handling for external patient creation
func (suite *ExternalPatientServiceTestSuite) TestCreateExternalPatient_Error() {
	// Arrange
	mockPatient := &fhir.Patient{Id: func() *string { s := "fail-id"; return &s }()}
	suite.mockClient.EXPECT().CreatePatient(mockPatient).Return(nil, errors.New("create failed"))

	// Act
	created, err := suite.service.CreateExternalPatient(mockPatient)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), created)
	assert.Contains(suite.T(), err.Error(), "create failed")
}
