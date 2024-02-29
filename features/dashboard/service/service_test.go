package service

import (
	"errors"
	"my-tourist-ticket/features/dashboard"
	"my-tourist-ticket/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDashboard(t *testing.T) {
	repo := new(mocks.DashboardData)
	srv := New(repo)
	userId := 1
	userRole := "admin"
	expectedDashboard := &dashboard.Dashboard{
		TotalCustomer:    10,
		TotalPengelola:   5,
		TotalTransaction: 20,
		RecentBooking:    []dashboard.Booking{},
		TopTours:         []dashboard.Tour{},
		TotalTour:        15,
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("GetUserRoleById", userId).Return(userRole, nil).Once()
		repo.On("GetTotalCustomer").Return(expectedDashboard.TotalCustomer, nil).Once()
		repo.On("GetTotalPengelola").Return(expectedDashboard.TotalPengelola, nil).Once()
		repo.On("GetTotalTransaction").Return(expectedDashboard.TotalTransaction, nil).Once()
		repo.On("GetRecentTransaction").Return([]dashboard.Booking{}, nil).Once()
		repo.On("GetTopTour").Return([]dashboard.Tour{}, nil).Once()
		repo.On("GetTotalTour").Return(expectedDashboard.TotalTour, nil).Once()

		dash, err := srv.Dashboard(userId)

		assert.NotNil(t, dash)
		assert.NoError(t, err)
		assert.Equal(t, expectedDashboard.TotalTour, dash.TotalTour)

		repo.AssertExpectations(t)
	})

	t.Run("User not admin", func(t *testing.T) {
		repo.On("GetUserRoleById", userId).Return("user", nil).Once()

		dash, err := srv.Dashboard(userId)

		assert.Nil(t, dash)
		assert.EqualError(t, err, "maaf anda tidak memiliki akses")

		repo.AssertExpectations(t)
	})

	t.Run("Error get user role", func(t *testing.T) {
		repo.On("GetUserRoleById", userId).Return("", errors.New("error")).Once()

		dash, err := srv.Dashboard(userId)

		assert.Nil(t, dash)
		assert.EqualError(t, err, "error")

		repo.AssertExpectations(t)
	})

	t.Run("Error get total customer", func(t *testing.T) {
		repo.On("GetUserRoleById", userId).Return(userRole, nil).Once()
		repo.On("GetTotalCustomer").Return(0, errors.New("error")).Once()

		dash, err := srv.Dashboard(userId)

		assert.Nil(t, dash)
		assert.EqualError(t, err, "error")

		repo.AssertExpectations(t)
	})

	t.Run("Error get total pengelola", func(t *testing.T) {
		repo.On("GetUserRoleById", userId).Return(userRole, nil).Once()
		repo.On("GetTotalCustomer").Return(expectedDashboard.TotalCustomer, nil).Once()
		repo.On("GetTotalPengelola").Return(0, errors.New("error")).Once()

		dash, err := srv.Dashboard(userId)

		assert.Nil(t, dash)
		assert.EqualError(t, err, "error")

		repo.AssertExpectations(t)
	})

	t.Run("Error get total transaction", func(t *testing.T) {
		repo.On("GetUserRoleById", userId).Return(userRole, nil).Once()
		repo.On("GetTotalCustomer").Return(expectedDashboard.TotalCustomer, nil).Once()
		repo.On("GetTotalPengelola").Return(expectedDashboard.TotalPengelola, nil).Once()
		repo.On("GetTotalTransaction").Return(0, errors.New("error")).Once()

		dash, err := srv.Dashboard(userId)

		assert.Nil(t, dash)
		assert.EqualError(t, err, "error")

		repo.AssertExpectations(t)
	})

	t.Run("Error get recent transaction", func(t *testing.T) {
		repo.On("GetUserRoleById", userId).Return(userRole, nil).Once()
		repo.On("GetTotalCustomer").Return(expectedDashboard.TotalCustomer, nil).Once()
		repo.On("GetTotalPengelola").Return(expectedDashboard.TotalPengelola, nil).Once()
		repo.On("GetTotalTransaction").Return(expectedDashboard.TotalTransaction, nil).Once()
		repo.On("GetRecentTransaction").Return(nil, errors.New("error")).Once()

		dash, err := srv.Dashboard(userId)

		assert.Nil(t, dash)
		assert.EqualError(t, err, "error")

		repo.AssertExpectations(t)
	})

	t.Run("Error get top tour", func(t *testing.T) {
		repo.On("GetUserRoleById", userId).Return(userRole, nil).Once()
		repo.On("GetTotalCustomer").Return(expectedDashboard.TotalCustomer, nil).Once()
		repo.On("GetTotalPengelola").Return(expectedDashboard.TotalPengelola, nil).Once()
		repo.On("GetTotalTransaction").Return(expectedDashboard.TotalTransaction, nil).Once()
		repo.On("GetRecentTransaction").Return([]dashboard.Booking{}, nil).Once()
		repo.On("GetTopTour").Return(nil, errors.New("error")).Once()

		dash, err := srv.Dashboard(userId)

		assert.Nil(t, dash)
		assert.EqualError(t, err, "error")

		repo.AssertExpectations(t)
	})

	t.Run("Error get total tour", func(t *testing.T) {
		repo.On("GetUserRoleById", userId).Return(userRole, nil).Once()
		repo.On("GetTotalCustomer").Return(expectedDashboard.TotalCustomer, nil).Once()
		repo.On("GetTotalPengelola").Return(expectedDashboard.TotalPengelola, nil).Once()
		repo.On("GetTotalTransaction").Return(expectedDashboard.TotalTransaction, nil).Once()
		repo.On("GetRecentTransaction").Return([]dashboard.Booking{}, nil).Once()
		repo.On("GetTopTour").Return([]dashboard.Tour{}, nil).Once()
		repo.On("GetTotalTour").Return(0, errors.New("error")).Once()

		dash, err := srv.Dashboard(userId)
		assert.Nil(t, dash)
		assert.EqualError(t, err, "error")

		repo.AssertExpectations(t)
	})

}
