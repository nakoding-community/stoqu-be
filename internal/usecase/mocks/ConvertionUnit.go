// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"

	dto "gitlab.com/stoqu/stoqu-be/internal/model/dto"

	mock "github.com/stretchr/testify/mock"
)

// ConvertionUnit is an autogenerated mock type for the ConvertionUnit type
type ConvertionUnit struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, payload
func (_m *ConvertionUnit) Create(ctx context.Context, payload dto.CreateConvertionUnitRequest) (dto.ConvertionUnitResponse, error) {
	ret := _m.Called(ctx, payload)

	var r0 dto.ConvertionUnitResponse
	if rf, ok := ret.Get(0).(func(context.Context, dto.CreateConvertionUnitRequest) dto.ConvertionUnitResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(dto.ConvertionUnitResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.CreateConvertionUnitRequest) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, payload
func (_m *ConvertionUnit) Delete(ctx context.Context, payload dto.ByIDRequest) (dto.ConvertionUnitResponse, error) {
	ret := _m.Called(ctx, payload)

	var r0 dto.ConvertionUnitResponse
	if rf, ok := ret.Get(0).(func(context.Context, dto.ByIDRequest) dto.ConvertionUnitResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(dto.ConvertionUnitResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.ByIDRequest) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: ctx, filterParam
func (_m *ConvertionUnit) Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.ConvertionUnitResponse, abstraction.PaginationInfo, error) {
	ret := _m.Called(ctx, filterParam)

	var r0 []dto.ConvertionUnitResponse
	if rf, ok := ret.Get(0).(func(context.Context, abstraction.Filter) []dto.ConvertionUnitResponse); ok {
		r0 = rf(ctx, filterParam)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.ConvertionUnitResponse)
		}
	}

	var r1 abstraction.PaginationInfo
	if rf, ok := ret.Get(1).(func(context.Context, abstraction.Filter) abstraction.PaginationInfo); ok {
		r1 = rf(ctx, filterParam)
	} else {
		r1 = ret.Get(1).(abstraction.PaginationInfo)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, abstraction.Filter) error); ok {
		r2 = rf(ctx, filterParam)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FindByID provides a mock function with given fields: ctx, payload
func (_m *ConvertionUnit) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ConvertionUnitResponse, error) {
	ret := _m.Called(ctx, payload)

	var r0 dto.ConvertionUnitResponse
	if rf, ok := ret.Get(0).(func(context.Context, dto.ByIDRequest) dto.ConvertionUnitResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(dto.ConvertionUnitResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.ByIDRequest) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, payload
func (_m *ConvertionUnit) Update(ctx context.Context, payload dto.UpdateConvertionUnitRequest) (dto.ConvertionUnitResponse, error) {
	ret := _m.Called(ctx, payload)

	var r0 dto.ConvertionUnitResponse
	if rf, ok := ret.Get(0).(func(context.Context, dto.UpdateConvertionUnitRequest) dto.ConvertionUnitResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(dto.ConvertionUnitResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.UpdateConvertionUnitRequest) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewConvertionUnit interface {
	mock.TestingT
	Cleanup(func())
}

// NewConvertionUnit creates a new instance of ConvertionUnit. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConvertionUnit(t mockConstructorTestingTNewConvertionUnit) *ConvertionUnit {
	mock := &ConvertionUnit{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
