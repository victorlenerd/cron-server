// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	models "scheduler0/pkg/models"

	mock "github.com/stretchr/testify/mock"

	utils "scheduler0/pkg/utils"
)

// Credential is an autogenerated mock type for the Credential type
type Credential struct {
	mock.Mock
}

// CreateNewCredential provides a mock function with given fields: credentialTransformer
func (_m *Credential) CreateNewCredential(credentialTransformer models.Credential) (uint64, *utils.GenericError) {
	ret := _m.Called(credentialTransformer)

	var r0 uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(models.Credential) (uint64, *utils.GenericError)); ok {
		return rf(credentialTransformer)
	}
	if rf, ok := ret.Get(0).(func(models.Credential) uint64); ok {
		r0 = rf(credentialTransformer)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(models.Credential) *utils.GenericError); ok {
		r1 = rf(credentialTransformer)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// DeleteOneCredential provides a mock function with given fields: id
func (_m *Credential) DeleteOneCredential(id uint64) (*models.Credential, error) {
	ret := _m.Called(id)

	var r0 *models.Credential
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (*models.Credential, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint64) *models.Credential); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Credential)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneCredentialByID provides a mock function with given fields: id
func (_m *Credential) FindOneCredentialByID(id uint64) (*models.Credential, error) {
	ret := _m.Called(id)

	var r0 *models.Credential
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (*models.Credential, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint64) *models.Credential); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Credential)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListCredentials provides a mock function with given fields: offset, limit, orderBy
func (_m *Credential) ListCredentials(offset uint64, limit uint64, orderBy string) (*models.PaginatedCredential, *utils.GenericError) {
	ret := _m.Called(offset, limit, orderBy)

	var r0 *models.PaginatedCredential
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(uint64, uint64, string) (*models.PaginatedCredential, *utils.GenericError)); ok {
		return rf(offset, limit, orderBy)
	}
	if rf, ok := ret.Get(0).(func(uint64, uint64, string) *models.PaginatedCredential); ok {
		r0 = rf(offset, limit, orderBy)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PaginatedCredential)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, uint64, string) *utils.GenericError); ok {
		r1 = rf(offset, limit, orderBy)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// UpdateOneCredential provides a mock function with given fields: credentialTransformer
func (_m *Credential) UpdateOneCredential(credentialTransformer models.Credential) (*models.Credential, error) {
	ret := _m.Called(credentialTransformer)

	var r0 *models.Credential
	var r1 error
	if rf, ok := ret.Get(0).(func(models.Credential) (*models.Credential, error)); ok {
		return rf(credentialTransformer)
	}
	if rf, ok := ret.Get(0).(func(models.Credential) *models.Credential); ok {
		r0 = rf(credentialTransformer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Credential)
		}
	}

	if rf, ok := ret.Get(1).(func(models.Credential) error); ok {
		r1 = rf(credentialTransformer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateServerAPIKey provides a mock function with given fields: apiKey, apiSecret
func (_m *Credential) ValidateServerAPIKey(apiKey string, apiSecret string) (bool, *utils.GenericError) {
	ret := _m.Called(apiKey, apiSecret)

	var r0 bool
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(string, string) (bool, *utils.GenericError)); ok {
		return rf(apiKey, apiSecret)
	}
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(apiKey, apiSecret)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string) *utils.GenericError); ok {
		r1 = rf(apiKey, apiSecret)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewCredential interface {
	mock.TestingT
	Cleanup(func())
}

// NewCredential creates a new instance of Credential. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCredential(t mockConstructorTestingTNewCredential) *Credential {
	mock := &Credential{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
