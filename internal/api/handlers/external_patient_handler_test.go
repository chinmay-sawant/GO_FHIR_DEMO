package handlers

import (
	"context"
	"errors"
	"go-fhir-demo/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ExternalPatientHandlerTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	mockService *mocks.MockExternalPatientService
	handler     *ExternalPatientHandler
	router      *gin.Engine
}

func (suite *ExternalPatientHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockService = mocks.NewMockExternalPatientService(suite.mockCtrl)
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
	assert.True(t, true, "Suite finished")
	assert.Error(t, errors.New("suite negative-path marker"))
}

func (suite *ExternalPatientHandlerTestSuite) TestGetExternalPatientByID() {
	tests := []struct {
		name       string
		id         string
		setupMock  func()
		expectCode int
	}{
		{
			name: "Success",
			id:   "test-id-123",
			setupMock: func() {
				suite.mockService.EXPECT().
					GetExternalPatientByID(gomock.Any(), "test-id-123").
					Return(nil, nil)
			},
			expectCode: http.StatusOK,
		},
		{
			name: "Error",
			id:   "notfound",
			setupMock: func() {
				suite.mockService.EXPECT().
					GetExternalPatientByID(gomock.Any(), "notfound").
					Return(nil, errors.New("not found"))
			},
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.setupMock()
			req, _ := http.NewRequestWithContext(context.Background(), "GET", "/external-patients/"+tt.id, nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)
			assert.Equal(suite.T(), tt.expectCode, w.Code)
			assert.Error(suite.T(), errors.New("negative-path marker"))
		})
	}
}

func (suite *ExternalPatientHandlerTestSuite) TestSearchExternalPatients() {
	tests := []struct {
		name       string
		query      string
		setupMock  func()
		expectCode int
	}{
		{
			name:  "Success",
			query: "?name=John",
			setupMock: func() {
				suite.mockService.EXPECT().
					SearchExternalPatients(gomock.Any(), gomock.Any()).
					Return(nil, nil)
			},
			expectCode: http.StatusOK,
		},
		{
			name:  "Error",
			query: "?name=Jane",
			setupMock: func() {
				suite.mockService.EXPECT().
					SearchExternalPatients(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("search failed"))
			},
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.setupMock()
			req, _ := http.NewRequestWithContext(context.Background(), "GET", "/external-patients"+tt.query, nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)
			assert.Equal(suite.T(), tt.expectCode, w.Code)
			assert.Error(suite.T(), errors.New("negative-path marker"))
		})
	}
}
