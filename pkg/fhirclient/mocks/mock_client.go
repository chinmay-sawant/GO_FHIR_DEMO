// Code generated by MockGen. DO NOT EDIT.
// Source: D:\Chinmay_Personal_Projects\Go_FHIR_Demo\pkg\fhirclient\client.go
//
// Generated by this command:
//
//	mockgen -source=D:\Chinmay_Personal_Projects\Go_FHIR_Demo\pkg\fhirclient\client.go -destination=D:\Chinmay_Personal_Projects\Go_FHIR_Demo\pkg\fhirclient\mocks\mock_client.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	fhir "github.com/samply/golang-fhir-models/fhir-models/fhir"
	gomock "go.uber.org/mock/gomock"
)

// MockClientInterface is a mock of ClientInterface interface.
type MockClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientInterfaceMockRecorder
	isgomock struct{}
}

// MockClientInterfaceMockRecorder is the mock recorder for MockClientInterface.
type MockClientInterfaceMockRecorder struct {
	mock *MockClientInterface
}

// NewMockClientInterface creates a new mock instance.
func NewMockClientInterface(ctrl *gomock.Controller) *MockClientInterface {
	mock := &MockClientInterface{ctrl: ctrl}
	mock.recorder = &MockClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientInterface) EXPECT() *MockClientInterfaceMockRecorder {
	return m.recorder
}

// CreatePatient mocks base method.
func (m *MockClientInterface) CreatePatient(ctx context.Context, patient *fhir.Patient) (*fhir.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePatient", ctx, patient)
	ret0, _ := ret[0].(*fhir.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePatient indicates an expected call of CreatePatient.
func (mr *MockClientInterfaceMockRecorder) CreatePatient(ctx, patient any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePatient", reflect.TypeOf((*MockClientInterface)(nil).CreatePatient), ctx, patient)
}

// GetPatientByID mocks base method.
func (m *MockClientInterface) GetPatientByID(ctx context.Context, id string) (*fhir.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPatientByID", ctx, id)
	ret0, _ := ret[0].(*fhir.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPatientByID indicates an expected call of GetPatientByID.
func (mr *MockClientInterfaceMockRecorder) GetPatientByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPatientByID", reflect.TypeOf((*MockClientInterface)(nil).GetPatientByID), ctx, id)
}

// SearchPatients mocks base method.
func (m *MockClientInterface) SearchPatients(ctx context.Context, queryParams map[string]string) (*fhir.Bundle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchPatients", ctx, queryParams)
	ret0, _ := ret[0].(*fhir.Bundle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchPatients indicates an expected call of SearchPatients.
func (mr *MockClientInterfaceMockRecorder) SearchPatients(ctx, queryParams any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchPatients", reflect.TypeOf((*MockClientInterface)(nil).SearchPatients), ctx, queryParams)
}
