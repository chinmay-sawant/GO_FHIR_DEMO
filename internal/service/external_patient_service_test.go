package service

import (
	"context"
	"errors"
	redisclientmock "go-fhir-demo/pkg/cache/mocks"
	fhirclientmocks "go-fhir-demo/pkg/fhirclient/mocks"
	"testing"
	"time"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// ExternalPatientServiceTestSuite defines the test suite
type ExternalPatientServiceTestSuite struct {
	suite.Suite
	mockClient      *fhirclientmocks.MockClientInterface
	mockRedisClient *redisclientmock.MockCacheInterface
	service         ExternalPatientServiceInterface
}

// SetupTest initializes the test suite before each test
func (suite *ExternalPatientServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockClient = fhirclientmocks.NewMockClientInterface(ctrl)
	suite.mockRedisClient = redisclientmock.NewMockCacheInterface(ctrl)

	suite.service = NewExternalPatientService(suite.mockClient, suite.mockRedisClient)
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
	suite.mockClient.EXPECT().GetPatientByID(gomock.Any(), testID).Return(mockPatient, nil)

	// Act
	patient, err := suite.service.GetExternalPatientByID(context.Background(), testID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), testID, *patient.Id)
}

// TestGetExternalPatientByID_Error tests error handling
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByID_Error() {
	// Arrange
	testID := "notfound"
	suite.mockClient.EXPECT().GetPatientByID(gomock.Any(), testID).Return(nil, errors.New("patient not found"))

	// Act
	patient, err := suite.service.GetExternalPatientByID(context.Background(), testID)

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
	suite.mockClient.EXPECT().SearchPatients(gomock.Any(), searchParams).Return(mockBundle, nil)

	// Act
	bundle, err := suite.service.SearchExternalPatients(context.Background(), searchParams)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), bundle)
	assert.Equal(suite.T(), fhir.BundleTypeSearchset, bundle.Type)
}

// TestSearchExternalPatients_Error tests search error handling
func (suite *ExternalPatientServiceTestSuite) TestSearchExternalPatients_Error() {
	// Arrange
	searchParams := map[string]string{"name": "Jane"}
	suite.mockClient.EXPECT().SearchPatients(gomock.Any(), searchParams).Return(nil, errors.New("search failed"))

	// Act
	bundle, err := suite.service.SearchExternalPatients(context.Background(), searchParams)

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
	suite.mockClient.EXPECT().SearchPatients(gomock.Any(), emptyParams).Return(mockBundle, nil)

	// Act
	bundle, err := suite.service.SearchExternalPatients(context.Background(), emptyParams)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), bundle)
}

// TestCreateExternalPatient_Success tests successful creation of external patient
func (suite *ExternalPatientServiceTestSuite) TestCreateExternalPatient_Success() {
	// Arrange
	mockPatient := &fhir.Patient{Id: func() *string { s := "new-id"; return &s }()}
	suite.mockClient.EXPECT().CreatePatient(gomock.Any(), mockPatient).Return(mockPatient, nil)

	// Act
	created, err := suite.service.CreateExternalPatient(context.Background(), mockPatient)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), created)
	assert.Equal(suite.T(), *mockPatient.Id, *created.Id)
}

// TestCreateExternalPatient_Error tests error handling for external patient creation
func (suite *ExternalPatientServiceTestSuite) TestCreateExternalPatient_Error() {
	// Arrange
	mockPatient := &fhir.Patient{Id: func() *string { s := "fail-id"; return &s }()}
	suite.mockClient.EXPECT().CreatePatient(gomock.Any(), mockPatient).Return(nil, errors.New("create failed"))

	// Act
	created, err := suite.service.CreateExternalPatient(context.Background(), mockPatient)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), created)
	assert.Contains(suite.T(), err.Error(), "create failed")
}

// TestGetExternalPatientByIDCached_CacheHit tests retrieval from cache
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByIDCached_CacheHit() {
	ctx := context.Background()
	testID := "cached-id"
	mockPatient := &fhir.Patient{Id: &testID}

	suite.mockRedisClient.EXPECT().GetPatient(ctx, testID).Return(mockPatient, nil)

	patient, err := suite.service.GetExternalPatientByIDCached(ctx, testID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), testID, *patient.Id)
}

// TestGetExternalPatientByIDCached_CacheMiss_SuccessfulFetch tests cache miss and successful fetch from FHIR server
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByIDCached_CacheMiss_SuccessfulFetch() {
	ctx := context.Background()
	testID := "miss-id"
	mockPatient := &fhir.Patient{Id: &testID}

	suite.mockRedisClient.EXPECT().GetPatient(ctx, testID).Return(nil, nil)
	suite.mockClient.EXPECT().GetPatientByID(ctx, testID).Return(mockPatient, nil)
	suite.mockRedisClient.EXPECT().SetPatient(ctx, testID, mockPatient, gomock.Any()).Return(nil)

	patient, err := suite.service.GetExternalPatientByIDCached(ctx, testID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), testID, *patient.Id)
}

// TestGetExternalPatientByIDCached_CacheMiss_FetchError tests cache miss and error from FHIR server
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByIDCached_CacheMiss_FetchError() {
	ctx := context.Background()
	testID := "error-id"

	suite.mockRedisClient.EXPECT().GetPatient(ctx, testID).Return(nil, nil)
	suite.mockClient.EXPECT().GetPatientByID(ctx, testID).Return(nil, errors.New("not found"))

	patient, err := suite.service.GetExternalPatientByIDCached(ctx, testID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), patient)
	assert.Contains(suite.T(), err.Error(), "not found")
}

// TestGetExternalPatientByIDCached_CacheGetError tests error when getting from cache, but fetch from FHIR server succeeds
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByIDCached_CacheGetError() {
	ctx := context.Background()
	testID := "cache-get-error"
	mockPatient := &fhir.Patient{Id: &testID}

	suite.mockRedisClient.EXPECT().GetPatient(ctx, testID).Return(nil, errors.New("redis down"))
	suite.mockClient.EXPECT().GetPatientByID(ctx, testID).Return(mockPatient, nil)
	suite.mockRedisClient.EXPECT().SetPatient(ctx, testID, mockPatient, gomock.Any()).Return(nil)

	patient, err := suite.service.GetExternalPatientByIDCached(ctx, testID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), testID, *patient.Id)
}

// TestGetExternalPatientByIDCached_CacheSetError tests error when setting cache after fetch
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByIDCached_CacheSetError() {
	ctx := context.Background()
	testID := "cache-set-error"
	mockPatient := &fhir.Patient{Id: &testID}

	suite.mockRedisClient.EXPECT().GetPatient(ctx, testID).Return(nil, nil)
	suite.mockClient.EXPECT().GetPatientByID(ctx, testID).Return(mockPatient, nil)
	suite.mockRedisClient.EXPECT().SetPatient(ctx, testID, mockPatient, gomock.Any()).Return(errors.New("set failed"))

	patient, err := suite.service.GetExternalPatientByIDCached(ctx, testID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), testID, *patient.Id)
}

// TestGetExternalPatientByIDDelayed_Success tests successful retrieval within timeout
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByIDDelayed_Success() {
	ctx := context.Background()
	testID := "delayed-id"
	mockPatient := &fhir.Patient{Id: &testID}

	suite.mockClient.EXPECT().GetPatientByID(ctx, testID).Return(mockPatient, nil)

	patient, err := suite.service.GetExternalPatientByIDDelayed(ctx, testID, 2*time.Second)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), patient)
	assert.Equal(suite.T(), testID, *patient.Id)
}

// TestGetExternalPatientByIDDelayed_Error tests error from FHIR server within timeout
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByIDDelayed_Error() {
	ctx := context.Background()
	testID := "delayed-error"

	suite.mockClient.EXPECT().GetPatientByID(ctx, testID).Return(nil, errors.New("delayed error"))

	patient, err := suite.service.GetExternalPatientByIDDelayed(ctx, testID, 2*time.Second)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), patient)
	assert.Contains(suite.T(), err.Error(), "delayed error")
}

// TestGetExternalPatientByIDDelayed_Timeout tests timeout scenario
func (suite *ExternalPatientServiceTestSuite) TestGetExternalPatientByIDDelayed_Timeout() {
	ctx := context.Background()
	testID := "timeout-id"

	// Simulate a long-running GetPatientByID by blocking until context is done
	suite.mockClient.EXPECT().GetPatientByID(ctx, testID).DoAndReturn(func(_ context.Context, _ string) (*fhir.Patient, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	})

	patient, err := suite.service.GetExternalPatientByIDDelayed(ctx, testID, 10*time.Millisecond)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), patient)
	assert.Contains(suite.T(), err.Error(), "context deadline exceeded")
}
