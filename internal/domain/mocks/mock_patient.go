// Code generated by MockGen. DO NOT EDIT.
// Source: D:\Chinmay_Personal_Projects\Go_FHIR_Demo\internal\domain\patient.go
//
// Generated by this command:
//
//	mockgen -source=D:\Chinmay_Personal_Projects\Go_FHIR_Demo\internal\domain\patient.go -destination=D:\Chinmay_Personal_Projects\Go_FHIR_Demo\internal\domain\mocks\mock_patient.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	domain "go-fhir-demo/internal/domain"
	reflect "reflect"

	fhir "github.com/samply/golang-fhir-models/fhir-models/fhir"
	gomock "go.uber.org/mock/gomock"
)

// MockPatientRepository is a mock of PatientRepository interface.
type MockPatientRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPatientRepositoryMockRecorder
	isgomock struct{}
}

// MockPatientRepositoryMockRecorder is the mock recorder for MockPatientRepository.
type MockPatientRepositoryMockRecorder struct {
	mock *MockPatientRepository
}

// NewMockPatientRepository creates a new mock instance.
func NewMockPatientRepository(ctrl *gomock.Controller) *MockPatientRepository {
	mock := &MockPatientRepository{ctrl: ctrl}
	mock.recorder = &MockPatientRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPatientRepository) EXPECT() *MockPatientRepositoryMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockPatientRepository) Count() (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockPatientRepositoryMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockPatientRepository)(nil).Count))
}

// Create mocks base method.
func (m *MockPatientRepository) Create(patient *domain.Patient) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", patient)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockPatientRepositoryMockRecorder) Create(patient any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPatientRepository)(nil).Create), patient)
}

// Delete mocks base method.
func (m *MockPatientRepository) Delete(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPatientRepositoryMockRecorder) Delete(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPatientRepository)(nil).Delete), id)
}

// GetAll mocks base method.
func (m *MockPatientRepository) GetAll(limit, offset int) ([]*domain.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", limit, offset)
	ret0, _ := ret[0].([]*domain.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPatientRepositoryMockRecorder) GetAll(limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPatientRepository)(nil).GetAll), limit, offset)
}

// GetByID mocks base method.
func (m *MockPatientRepository) GetByID(id uint) (*domain.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(*domain.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockPatientRepositoryMockRecorder) GetByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockPatientRepository)(nil).GetByID), id)
}

// Update mocks base method.
func (m *MockPatientRepository) Update(patient *domain.Patient) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", patient)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPatientRepositoryMockRecorder) Update(patient any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPatientRepository)(nil).Update), patient)
}

// MockPatientService is a mock of PatientService interface.
type MockPatientService struct {
	ctrl     *gomock.Controller
	recorder *MockPatientServiceMockRecorder
	isgomock struct{}
}

// MockPatientServiceMockRecorder is the mock recorder for MockPatientService.
type MockPatientServiceMockRecorder struct {
	mock *MockPatientService
}

// NewMockPatientService creates a new mock instance.
func NewMockPatientService(ctrl *gomock.Controller) *MockPatientService {
	mock := &MockPatientService{ctrl: ctrl}
	mock.recorder = &MockPatientServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPatientService) EXPECT() *MockPatientServiceMockRecorder {
	return m.recorder
}

// ConvertFromFHIR mocks base method.
func (m *MockPatientService) ConvertFromFHIR(fhirPatient *fhir.Patient) (*domain.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertFromFHIR", fhirPatient)
	ret0, _ := ret[0].(*domain.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConvertFromFHIR indicates an expected call of ConvertFromFHIR.
func (mr *MockPatientServiceMockRecorder) ConvertFromFHIR(fhirPatient any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertFromFHIR", reflect.TypeOf((*MockPatientService)(nil).ConvertFromFHIR), fhirPatient)
}

// ConvertToFHIR mocks base method.
func (m *MockPatientService) ConvertToFHIR(patient *domain.Patient) (*fhir.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertToFHIR", patient)
	ret0, _ := ret[0].(*fhir.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConvertToFHIR indicates an expected call of ConvertToFHIR.
func (mr *MockPatientServiceMockRecorder) ConvertToFHIR(patient any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertToFHIR", reflect.TypeOf((*MockPatientService)(nil).ConvertToFHIR), patient)
}

// CreatePatient mocks base method.
func (m *MockPatientService) CreatePatient(fhirPatient *fhir.Patient) (*domain.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePatient", fhirPatient)
	ret0, _ := ret[0].(*domain.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePatient indicates an expected call of CreatePatient.
func (mr *MockPatientServiceMockRecorder) CreatePatient(fhirPatient any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePatient", reflect.TypeOf((*MockPatientService)(nil).CreatePatient), fhirPatient)
}

// DeletePatient mocks base method.
func (m *MockPatientService) DeletePatient(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePatient", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePatient indicates an expected call of DeletePatient.
func (mr *MockPatientServiceMockRecorder) DeletePatient(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePatient", reflect.TypeOf((*MockPatientService)(nil).DeletePatient), id)
}

// GetPatient mocks base method.
func (m *MockPatientService) GetPatient(id uint) (*domain.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPatient", id)
	ret0, _ := ret[0].(*domain.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPatient indicates an expected call of GetPatient.
func (mr *MockPatientServiceMockRecorder) GetPatient(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPatient", reflect.TypeOf((*MockPatientService)(nil).GetPatient), id)
}

// GetPatients mocks base method.
func (m *MockPatientService) GetPatients(limit, offset int) ([]*domain.Patient, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPatients", limit, offset)
	ret0, _ := ret[0].([]*domain.Patient)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPatients indicates an expected call of GetPatients.
func (mr *MockPatientServiceMockRecorder) GetPatients(limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPatients", reflect.TypeOf((*MockPatientService)(nil).GetPatients), limit, offset)
}

// PatchPatient mocks base method.
func (m *MockPatientService) PatchPatient(id uint, updates map[string]any) (*domain.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchPatient", id, updates)
	ret0, _ := ret[0].(*domain.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchPatient indicates an expected call of PatchPatient.
func (mr *MockPatientServiceMockRecorder) PatchPatient(id, updates any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchPatient", reflect.TypeOf((*MockPatientService)(nil).PatchPatient), id, updates)
}

// UpdatePatient mocks base method.
func (m *MockPatientService) UpdatePatient(id uint, fhirPatient *fhir.Patient) (*domain.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePatient", id, fhirPatient)
	ret0, _ := ret[0].(*domain.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePatient indicates an expected call of UpdatePatient.
func (mr *MockPatientServiceMockRecorder) UpdatePatient(id, fhirPatient any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePatient", reflect.TypeOf((*MockPatientService)(nil).UpdatePatient), id, fhirPatient)
}
