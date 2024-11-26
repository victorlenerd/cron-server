// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	models "scheduler0/pkg/models"

	mock "github.com/stretchr/testify/mock"

	utils "scheduler0/pkg/utils"
)

// JobRepo is an autogenerated mock type for the JobRepo type
type JobRepo struct {
	mock.Mock
}

// BatchGetJobsByID provides a mock function with given fields: jobIDs
func (_m *JobRepo) BatchGetJobsByID(jobIDs []uint64) ([]models.Job, *utils.GenericError) {
	ret := _m.Called(jobIDs)

	var r0 []models.Job
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func([]uint64) ([]models.Job, *utils.GenericError)); ok {
		return rf(jobIDs)
	}
	if rf, ok := ret.Get(0).(func([]uint64) []models.Job); ok {
		r0 = rf(jobIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Job)
		}
	}

	if rf, ok := ret.Get(1).(func([]uint64) *utils.GenericError); ok {
		r1 = rf(jobIDs)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// BatchGetJobsWithIDRange provides a mock function with given fields: lowerBound, upperBound
func (_m *JobRepo) BatchGetJobsWithIDRange(lowerBound uint64, upperBound uint64) ([]models.Job, *utils.GenericError) {
	ret := _m.Called(lowerBound, upperBound)

	var r0 []models.Job
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(uint64, uint64) ([]models.Job, *utils.GenericError)); ok {
		return rf(lowerBound, upperBound)
	}
	if rf, ok := ret.Get(0).(func(uint64, uint64) []models.Job); ok {
		r0 = rf(lowerBound, upperBound)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Job)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, uint64) *utils.GenericError); ok {
		r1 = rf(lowerBound, upperBound)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// BatchInsertJobs provides a mock function with given fields: jobRepos
func (_m *JobRepo) BatchInsertJobs(jobRepos []models.Job) ([]uint64, *utils.GenericError) {
	ret := _m.Called(jobRepos)

	var r0 []uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func([]models.Job) ([]uint64, *utils.GenericError)); ok {
		return rf(jobRepos)
	}
	if rf, ok := ret.Get(0).(func([]models.Job) []uint64); ok {
		r0 = rf(jobRepos)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint64)
		}
	}

	if rf, ok := ret.Get(1).(func([]models.Job) *utils.GenericError); ok {
		r1 = rf(jobRepos)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// DeleteOneByID provides a mock function with given fields: jobModel
func (_m *JobRepo) DeleteOneByID(jobModel models.Job) (uint64, *utils.GenericError) {
	ret := _m.Called(jobModel)

	var r0 uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(models.Job) (uint64, *utils.GenericError)); ok {
		return rf(jobModel)
	}
	if rf, ok := ret.Get(0).(func(models.Job) uint64); ok {
		r0 = rf(jobModel)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(models.Job) *utils.GenericError); ok {
		r1 = rf(jobModel)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// GetAllByProjectID provides a mock function with given fields: projectID, offset, limit, orderBy
func (_m *JobRepo) GetAllByProjectID(projectID uint64, offset uint64, limit uint64, orderBy string) ([]models.Job, *utils.GenericError) {
	ret := _m.Called(projectID, offset, limit, orderBy)

	var r0 []models.Job
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(uint64, uint64, uint64, string) ([]models.Job, *utils.GenericError)); ok {
		return rf(projectID, offset, limit, orderBy)
	}
	if rf, ok := ret.Get(0).(func(uint64, uint64, uint64, string) []models.Job); ok {
		r0 = rf(projectID, offset, limit, orderBy)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Job)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, uint64, uint64, string) *utils.GenericError); ok {
		r1 = rf(projectID, offset, limit, orderBy)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// GetJobsPaginated provides a mock function with given fields: projectID, offset, limit
func (_m *JobRepo) GetJobsPaginated(projectID uint64, offset uint64, limit uint64) ([]models.Job, uint64, *utils.GenericError) {
	ret := _m.Called(projectID, offset, limit)

	var r0 []models.Job
	var r1 uint64
	var r2 *utils.GenericError
	if rf, ok := ret.Get(0).(func(uint64, uint64, uint64) ([]models.Job, uint64, *utils.GenericError)); ok {
		return rf(projectID, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(uint64, uint64, uint64) []models.Job); ok {
		r0 = rf(projectID, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Job)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, uint64, uint64) uint64); ok {
		r1 = rf(projectID, offset, limit)
	} else {
		r1 = ret.Get(1).(uint64)
	}

	if rf, ok := ret.Get(2).(func(uint64, uint64, uint64) *utils.GenericError); ok {
		r2 = rf(projectID, offset, limit)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*utils.GenericError)
		}
	}

	return r0, r1, r2
}

// GetJobsTotalCount provides a mock function with given fields:
func (_m *JobRepo) GetJobsTotalCount() (uint64, *utils.GenericError) {
	ret := _m.Called()

	var r0 uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func() (uint64, *utils.GenericError)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func() *utils.GenericError); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// GetJobsTotalCountByProjectID provides a mock function with given fields: projectID
func (_m *JobRepo) GetJobsTotalCountByProjectID(projectID uint64) (uint64, *utils.GenericError) {
	ret := _m.Called(projectID)

	var r0 uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(uint64) (uint64, *utils.GenericError)); ok {
		return rf(projectID)
	}
	if rf, ok := ret.Get(0).(func(uint64) uint64); ok {
		r0 = rf(projectID)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(uint64) *utils.GenericError); ok {
		r1 = rf(projectID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// GetOneByID provides a mock function with given fields: jobModel
func (_m *JobRepo) GetOneByID(jobModel *models.Job) *utils.GenericError {
	ret := _m.Called(jobModel)

	var r0 *utils.GenericError
	if rf, ok := ret.Get(0).(func(*models.Job) *utils.GenericError); ok {
		r0 = rf(jobModel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.GenericError)
		}
	}

	return r0
}

// UpdateOneByID provides a mock function with given fields: jobModel
func (_m *JobRepo) UpdateOneByID(jobModel models.Job) (uint64, *utils.GenericError) {
	ret := _m.Called(jobModel)

	var r0 uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(models.Job) (uint64, *utils.GenericError)); ok {
		return rf(jobModel)
	}
	if rf, ok := ret.Get(0).(func(models.Job) uint64); ok {
		r0 = rf(jobModel)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(models.Job) *utils.GenericError); ok {
		r1 = rf(jobModel)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewJobRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewJobRepo creates a new instance of JobRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJobRepo(t mockConstructorTestingTNewJobRepo) *JobRepo {
	mock := &JobRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
