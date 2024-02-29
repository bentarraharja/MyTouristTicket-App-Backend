package service

import (
	"errors"
	packages "my-tourist-ticket/features/package"
	"my-tourist-ticket/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPackageCreate(t *testing.T) {
	mockRepo := new(mocks.PackageData)
	service := New(mockRepo)

	benefits := []string{"benefit1", "benefit2"}
	input := packages.Core{
		JumlahTiket: 0,
	}

	t.Run("Success", func(t *testing.T) {
		expectedInput := input
		expectedInput.JumlahTiket = 1

		mockRepo.On("Insert", benefits, expectedInput).Return(nil).Once()

		err := service.Create(benefits, input)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Insert", func(t *testing.T) {
		expectedInput := input
		expectedInput.JumlahTiket = 1

		mockRepo.On("Insert", benefits, expectedInput).Return(errors.New("error inserting package")).Once()

		err := service.Create(benefits, input)

		assert.EqualError(t, err, "error inserting package")
	})
}

func TestPackageGetByTourId(t *testing.T) {
	mockRepo := new(mocks.PackageData)
	service := New(mockRepo)

	tourID := uint(123)
	expectedPackages := []packages.Core{
		{ID: 1, TourID: tourID, PackageName: "Package 1", Price: 100},
		{ID: 2, TourID: tourID, PackageName: "Package 2", Price: 150},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("SelectByTourId", tourID).Return(expectedPackages, nil).Once()

		packages, err := service.GetByTourId(tourID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPackages, packages)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error SelectByTourId", func(t *testing.T) {
		mockRepo.On("SelectByTourId", tourID).Return(nil, errors.New("error fetching packages")).Once()

		packages, err := service.GetByTourId(tourID)

		assert.Error(t, err)
		assert.Nil(t, packages)
	})
}

func TestPackageDelete(t *testing.T) {
	mockRepo := new(mocks.PackageData)
	service := New(mockRepo)

	packageID := 123

	t.Run("Success", func(t *testing.T) {
		// Mocking pemanggilan SelectAllBenefitsByPackageId dengan id paket yang valid
		mockRepo.On("SelectAllBenefitsByPackageId", packageID).Return([]packages.BenefitCore{{ID: 1}, {ID: 2}}, nil).Once()
		// Mocking pemanggilan DeleteBenefits
		mockRepo.On("DeleteBenefits", mock.Anything).Return(nil).Twice()
		// Mocking pemanggilan Delete
		mockRepo.On("Delete", packageID).Return(nil).Once()

		err := service.Delete(packageID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		err := service.Delete(0)

		assert.EqualError(t, err, "invalid id")
	})

	t.Run("Error SelectAllBenefitsByPackageId", func(t *testing.T) {
		mockRepo.On("SelectAllBenefitsByPackageId", packageID).Return(nil, errors.New("error getting benefits")).Once()

		err := service.Delete(packageID)

		assert.EqualError(t, err, "error getting benefits")
	})

	t.Run("Error DeleteBenefits", func(t *testing.T) {
		mockRepo.On("SelectAllBenefitsByPackageId", packageID).Return([]packages.BenefitCore{{ID: 1}, {ID: 2}}, nil).Once()
		mockRepo.On("DeleteBenefits", mock.Anything).Return(errors.New("error deleting benefits")).Once()

		err := service.Delete(packageID)

		assert.EqualError(t, err, "error deleting benefits")
	})

	t.Run("Error Delete", func(t *testing.T) {
		mockRepo.On("SelectAllBenefitsByPackageId", packageID).Return([]packages.BenefitCore{{ID: 1}, {ID: 2}}, nil).Once()
		mockRepo.On("DeleteBenefits", mock.Anything).Return(nil).Twice()
		mockRepo.On("Delete", packageID).Return(errors.New("error deleting package")).Once()

		err := service.Delete(packageID)

		assert.EqualError(t, err, "error deleting package")
	})
}
