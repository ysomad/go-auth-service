// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entity "github.com/ysomad/go-auth-service/internal/entity"
)

// User is an autogenerated mock type for the User type
type User struct {
	mock.Mock
}

// Archive provides a mock function with given fields: ctx, id, isArchive
func (_m *User) Archive(ctx context.Context, id int, isArchive bool) error {
	ret := _m.Called(ctx, id, isArchive)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, bool) error); ok {
		r0 = rf(ctx, id, isArchive)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, req
func (_m *User) Create(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error) {
	ret := _m.Called(ctx, req)

	var r0 *entity.User
	if rf, ok := ret.Get(0).(func(context.Context, entity.CreateUserRequest) *entity.User); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, entity.CreateUserRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *User) GetByID(ctx context.Context, id int) (*entity.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *entity.User
	if rf, ok := ret.Get(0).(func(context.Context, int) *entity.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PartialUpdate provides a mock function with given fields: ctx, id, req
func (_m *User) PartialUpdate(ctx context.Context, id int, req entity.PartialUpdateRequest) (*entity.User, error) {
	ret := _m.Called(ctx, id, req)

	var r0 *entity.User
	if rf, ok := ret.Get(0).(func(context.Context, int, entity.PartialUpdateRequest) *entity.User); ok {
		r0 = rf(ctx, id, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, entity.PartialUpdateRequest) error); ok {
		r1 = rf(ctx, id, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}