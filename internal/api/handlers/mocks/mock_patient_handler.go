// Code generated by MockGen. DO NOT EDIT.
// Source: D:\Chinmay_Personal_Projects\Go_FHIR_Demo\internal\api\handlers\patient_handler.go
//
// Generated by this command:
//
//	mockgen -source=D:\Chinmay_Personal_Projects\Go_FHIR_Demo\internal\api\handlers\patient_handler.go -destination=D:\Chinmay_Personal_Projects\Go_FHIR_Demo\internal\api\handlers\mocks\mock_patient_handler.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"
)

// MockPatientHandlerInterface is a mock of PatientHandlerInterface interface.
type MockPatientHandlerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPatientHandlerInterfaceMockRecorder
	isgomock struct{}
}

// MockPatientHandlerInterfaceMockRecorder is the mock recorder for MockPatientHandlerInterface.
type MockPatientHandlerInterfaceMockRecorder struct {
	mock *MockPatientHandlerInterface
}

// NewMockPatientHandlerInterface creates a new mock instance.
func NewMockPatientHandlerInterface(ctrl *gomock.Controller) *MockPatientHandlerInterface {
	mock := &MockPatientHandlerInterface{ctrl: ctrl}
	mock.recorder = &MockPatientHandlerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPatientHandlerInterface) EXPECT() *MockPatientHandlerInterfaceMockRecorder {
	return m.recorder
}

// CreatePatient mocks base method.
func (m *MockPatientHandlerInterface) CreatePatient(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreatePatient", c)
}

// CreatePatient indicates an expected call of CreatePatient.
func (mr *MockPatientHandlerInterfaceMockRecorder) CreatePatient(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePatient", reflect.TypeOf((*MockPatientHandlerInterface)(nil).CreatePatient), c)
}

// DeletePatient mocks base method.
func (m *MockPatientHandlerInterface) DeletePatient(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeletePatient", c)
}

// DeletePatient indicates an expected call of DeletePatient.
func (mr *MockPatientHandlerInterfaceMockRecorder) DeletePatient(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePatient", reflect.TypeOf((*MockPatientHandlerInterface)(nil).DeletePatient), c)
}

// GetPatient mocks base method.
func (m *MockPatientHandlerInterface) GetPatient(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetPatient", c)
}

// GetPatient indicates an expected call of GetPatient.
func (mr *MockPatientHandlerInterfaceMockRecorder) GetPatient(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPatient", reflect.TypeOf((*MockPatientHandlerInterface)(nil).GetPatient), c)
}

// GetPatients mocks base method.
func (m *MockPatientHandlerInterface) GetPatients(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetPatients", c)
}

// GetPatients indicates an expected call of GetPatients.
func (mr *MockPatientHandlerInterfaceMockRecorder) GetPatients(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPatients", reflect.TypeOf((*MockPatientHandlerInterface)(nil).GetPatients), c)
}

// PatchPatient mocks base method.
func (m *MockPatientHandlerInterface) PatchPatient(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PatchPatient", c)
}

// PatchPatient indicates an expected call of PatchPatient.
func (mr *MockPatientHandlerInterfaceMockRecorder) PatchPatient(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchPatient", reflect.TypeOf((*MockPatientHandlerInterface)(nil).PatchPatient), c)
}

// UpdatePatient mocks base method.
func (m *MockPatientHandlerInterface) UpdatePatient(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdatePatient", c)
}

// UpdatePatient indicates an expected call of UpdatePatient.
func (mr *MockPatientHandlerInterfaceMockRecorder) UpdatePatient(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePatient", reflect.TypeOf((*MockPatientHandlerInterface)(nil).UpdatePatient), c)
}
