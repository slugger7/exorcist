// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/library_path/library_path.go
//
// Generated by this command:
//
//	mockgen -source=./internal/service/library_path/library_path.go
//

// Package mock_libraryPathService is a generated GoMock package.
package mock_libraryPathService

import (
	reflect "reflect"

	model "github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	gomock "go.uber.org/mock/gomock"
)

// MockILibraryPathService is a mock of ILibraryPathService interface.
type MockILibraryPathService struct {
	ctrl     *gomock.Controller
	recorder *MockILibraryPathServiceMockRecorder
	isgomock struct{}
}

// MockILibraryPathServiceMockRecorder is the mock recorder for MockILibraryPathService.
type MockILibraryPathServiceMockRecorder struct {
	mock *MockILibraryPathService
}

// NewMockILibraryPathService creates a new mock instance.
func NewMockILibraryPathService(ctrl *gomock.Controller) *MockILibraryPathService {
	mock := &MockILibraryPathService{ctrl: ctrl}
	mock.recorder = &MockILibraryPathServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockILibraryPathService) EXPECT() *MockILibraryPathServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m_2 *MockILibraryPathService) Create(m *model.LibraryPath) (*model.LibraryPath, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Create", m)
	ret0, _ := ret[0].(*model.LibraryPath)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockILibraryPathServiceMockRecorder) Create(m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockILibraryPathService)(nil).Create), m)
}

// GetAll mocks base method.
func (m *MockILibraryPathService) GetAll() ([]model.LibraryPath, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]model.LibraryPath)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockILibraryPathServiceMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockILibraryPathService)(nil).GetAll))
}
