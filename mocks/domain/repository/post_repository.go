// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/repository/post_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pagination "github.com/mfsyahrz/image_feed_api/internal/common/pagination"
	entity "github.com/mfsyahrz/image_feed_api/internal/domain/entity"
	repository "github.com/mfsyahrz/image_feed_api/internal/domain/repository"
)

// MockPostRepository is a mock of PostRepository interface.
type MockPostRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPostRepositoryMockRecorder
}

// MockPostRepositoryMockRecorder is the mock recorder for MockPostRepository.
type MockPostRepositoryMockRecorder struct {
	mock *MockPostRepository
}

// NewMockPostRepository creates a new mock instance.
func NewMockPostRepository(ctrl *gomock.Controller) *MockPostRepository {
	mock := &MockPostRepository{ctrl: ctrl}
	mock.recorder = &MockPostRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostRepository) EXPECT() *MockPostRepositoryMockRecorder {
	return m.recorder
}

// FetchOne mocks base method.
func (m *MockPostRepository) FetchOne(ctx context.Context, input repository.GetPostInput) (*entity.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchOne", ctx, input)
	ret0, _ := ret[0].(*entity.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchOne indicates an expected call of FetchOne.
func (mr *MockPostRepositoryMockRecorder) FetchOne(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchOne", reflect.TypeOf((*MockPostRepository)(nil).FetchOne), ctx, input)
}

// FetchPaginated mocks base method.
func (m *MockPostRepository) FetchPaginated(ctx context.Context, input repository.GetPostInput) (entity.Posts, *pagination.PostCursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchPaginated", ctx, input)
	ret0, _ := ret[0].(entity.Posts)
	ret1, _ := ret[1].(*pagination.PostCursor)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FetchPaginated indicates an expected call of FetchPaginated.
func (mr *MockPostRepositoryMockRecorder) FetchPaginated(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchPaginated", reflect.TypeOf((*MockPostRepository)(nil).FetchPaginated), ctx, input)
}

// Save mocks base method.
func (m *MockPostRepository) Save(ctx context.Context, post *entity.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockPostRepositoryMockRecorder) Save(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockPostRepository)(nil).Save), ctx, post)
}
