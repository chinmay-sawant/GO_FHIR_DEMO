package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

// PatientRepositoryTestSuite defines the test suite
type PatientRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository PatientRepositoryInterface
	postgres   *embeddedpostgres.EmbeddedPostgres
}

// SetupSuite initializes the test suite
func (suite *PatientRepositoryTestSuite) SetupSuite() {
	// Start embedded Postgres
	port := 54329 // Use a non-default port to avoid conflicts
	suite.postgres = embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Port(uint32(port)).
			Username("postgres").
			Password("postgres").
			Database("testdb"),
	)
	err := suite.postgres.Start()
	suite.Require().NoError(err)

	dsn := fmt.Sprintf("host=localhost port=%d user=postgres password=postgres dbname=testdb sslmode=disable", port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	suite.Require().NoError(err)

	err = db.AutoMigrate(&domain.Patient{})
	suite.Require().NoError(err)

	suite.db = db
	suite.repository = NewPatientRepository(db)
}

// SetupTest runs before each test
func (suite *PatientRepositoryTestSuite) SetupTest() {
	// Clean up data before each test
	suite.db.Exec("TRUNCATE TABLE patients RESTART IDENTITY CASCADE")
}

// TearDownSuite cleans up after all tests
func (suite *PatientRepositoryTestSuite) TearDownSuite() {
	if suite.postgres != nil {
		_ = suite.postgres.Stop()
	}
}

// TestPatientRepositoryTestSuite runs the test suite
func TestPatientRepositoryTestSuite(t *testing.T) {
	// Skip if running in CI without postgres support
	if os.Getenv("CI") != "" {
		t.Skip("Skipping embedded postgres test in CI")
	}
	suite.Run(t, new(PatientRepositoryTestSuite))
}

// TestCreate_Success tests successful patient creation
func (suite *PatientRepositoryTestSuite) TestCreate_Success() {
	active := true
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	patient := &domain.Patient{
		FHIRData:  []byte(`{"resourceType":"Patient","id":"p1"}`),
		Active:    &active,
		Family:    "Doe",
		Given:     "John",
		Gender:    "male",
		BirthDate: &birthDate,
	}
	err := suite.repository.Create(context.Background(), patient)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), patient.ID)
}

// TestCreate_Error tests creation with invalid data
func (suite *PatientRepositoryTestSuite) TestCreate_Error() {
	// Arrange - Create a patient with invalid data (missing required fields)
	patient := &domain.Patient{
		// Missing required fields to trigger constraint error
	}

	// Act
	err := suite.repository.Create(context.Background(), patient)

	// Assert
	// The database should return an error due to NOT NULL constraint violation
	assert.Error(suite.T(), err)
}

// TestGetByID_Success tests successful patient retrieval
func (suite *PatientRepositoryTestSuite) TestGetByID_Success() {
	active := true
	birthDate := time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC)
	patient := &domain.Patient{
		FHIRData:  []byte(`{"resourceType":"Patient","id":"p2"}`),
		Active:    &active,
		Family:    "Smith",
		Given:     "Jane",
		Gender:    "female",
		BirthDate: &birthDate,
	}
	suite.repository.Create(context.Background(), patient)
	got, err := suite.repository.GetByID(context.Background(), patient.ID)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), got)
	assert.Equal(suite.T(), patient.Family, got.Family)
	assert.Equal(suite.T(), patient.Given, got.Given)
	assert.Equal(suite.T(), patient.Gender, got.Gender)
}

// TestGetByID_NotFound tests retrieval of non-existent patient
func (suite *PatientRepositoryTestSuite) TestGetByID_NotFound() {
	got, err := suite.repository.GetByID(context.Background(), 9999)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), got)
}

// TestGetAll_Success tests successful retrieval of all patients
func (suite *PatientRepositoryTestSuite) TestGetAll_Success() {
	active := true
	birthDate1 := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	birthDate2 := time.Date(1992, 2, 2, 0, 0, 0, 0, time.UTC)
	p1 := &domain.Patient{
		FHIRData:  []byte(`{"resourceType":"Patient","id":"p3"}`),
		Active:    &active,
		Family:    "Alpha",
		Given:     "A",
		Gender:    "male",
		BirthDate: &birthDate1,
	}
	p2 := &domain.Patient{
		FHIRData:  []byte(`{"resourceType":"Patient","id":"p4"}`),
		Active:    &active,
		Family:    "Beta",
		Given:     "B",
		Gender:    "female",
		BirthDate: &birthDate2,
	}
	suite.repository.Create(context.Background(), p1)
	suite.repository.Create(context.Background(), p2)
	list, err := suite.repository.GetAll(context.Background(), 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), list, 2)
}

// TestGetAll_WithPagination tests pagination functionality
func (suite *PatientRepositoryTestSuite) TestGetAll_WithPagination() {
	// Arrange
	for i := 0; i < 5; i++ {
		patient := &domain.Patient{
			FHIRData:  []byte(`{"resourceType":"Patient","id":"p5"}`),
			Active:    utils.CreateBoolPtr(true),
			Family:    "Patient",
			Given:     fmt.Sprintf("%d", i),
			Gender:    "male",
			BirthDate: utils.CreateTimePtr(time.Now().AddDate(-20-i, 0, 0).String()),
		}
		suite.repository.Create(context.Background(), patient)
	}

	// Act
	firstPage, err := suite.repository.GetAll(context.Background(), 2, 0)
	secondPage, err2 := suite.repository.GetAll(context.Background(), 2, 2)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), err2)
	assert.Len(suite.T(), firstPage, 2)
	assert.Len(suite.T(), secondPage, 2)
}

// TestUpdate_Success tests successful patient update
func (suite *PatientRepositoryTestSuite) TestUpdate_Success() {
	active := true
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	patient := &domain.Patient{
		FHIRData:  []byte(`{"resourceType":"Patient","id":"p5"}`),
		Active:    &active,
		Family:    "Gamma",
		Given:     "G",
		Gender:    "male",
		BirthDate: &birthDate,
	}
	suite.repository.Create(context.Background(), patient)
	patient.Family = "Delta"
	err := suite.repository.Update(context.Background(), patient)
	assert.NoError(suite.T(), err)
	got, _ := suite.repository.GetByID(context.Background(), patient.ID)
	assert.Equal(suite.T(), "Delta", got.Family)
}

// TestDelete_Success tests successful patient deletion
func (suite *PatientRepositoryTestSuite) TestDelete_Success() {
	active := true
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	patient := &domain.Patient{
		FHIRData:  []byte(`{"resourceType":"Patient","id":"p6"}`),
		Active:    &active,
		Family:    "Epsilon",
		Given:     "E",
		Gender:    "female",
		BirthDate: &birthDate,
	}
	suite.repository.Create(context.Background(), patient)
	err := suite.repository.Delete(context.Background(), patient.ID)
	assert.NoError(suite.T(), err)
	got, err := suite.repository.GetByID(context.Background(), patient.ID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), got)
}

// TestCount_Success tests successful patient count
func (suite *PatientRepositoryTestSuite) TestCount_Success() {
	active := true
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	p1 := &domain.Patient{
		FHIRData:  []byte(`{"resourceType":"Patient","id":"p7"}`),
		Active:    &active,
		Family:    "Zeta",
		Given:     "Z",
		Gender:    "male",
		BirthDate: &birthDate,
	}
	p2 := &domain.Patient{
		FHIRData:  []byte(`{"resourceType":"Patient","id":"p8"}`),
		Active:    &active,
		Family:    "Eta",
		Given:     "E",
		Gender:    "female",
		BirthDate: &birthDate,
	}
	suite.repository.Create(context.Background(), p1)
	suite.repository.Create(context.Background(), p2)
	count, err := suite.repository.Count(context.Background())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), count)
}

// TestCount_EmptyTable tests count on empty table
func (suite *PatientRepositoryTestSuite) TestCount_EmptyTable() {
	// Act
	count, err := suite.repository.Count(context.Background())

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
}
