package service

import (
	"errors"
	"mime/multipart"
	"my-tourist-ticket/features/city"
	"my-tourist-ticket/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserRoleById(t *testing.T) {
	repo := new(mocks.CityData)
	srv := NewCity(repo)

	userId := 1
	expectedRole := "admin"

	repo.On("GetUserRoleById", userId).Return(expectedRole, nil).Once()
	role, err := srv.GetUserRoleById(userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedRole, role)
	repo.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	repo := new(mocks.CityData)
	srv := NewCity(repo)

	inputData := city.Core{
		CityName:    "City Name",
		Description: "Description",
		Image:       "Image",
		Thumbnail:   "Thumbnail",
	}

	t.Run("Success Create City", func(t *testing.T) {
		image := new(multipart.FileHeader)
		thumbnail := new(multipart.FileHeader)

		repo.On("Insert", inputData, image, thumbnail).Return(nil).Once()
		err := srv.Create(inputData, image, thumbnail)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Error Create City", func(t *testing.T) {
		image := new(multipart.FileHeader)
		thumbnail := new(multipart.FileHeader)
		expectedErr := errors.New("insert error")

		repo.On("Insert", inputData, image, thumbnail).Return(expectedErr).Once()
		err := srv.Create(inputData, image, thumbnail)

		assert.Error(t, err)
		assert.EqualError(t, err, "error creating city: insert error")
		repo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	repo := new(mocks.CityData)
	srv := NewCity(repo)

	cityId := 1
	inputData := city.Core{
		CityName:    "Updated City Name",
		Description: "Updated Description",
		Image:       "Updated Image",
		Thumbnail:   "Updated Thumbnail",
	}

	t.Run("Success Update City", func(t *testing.T) {
		image := new(multipart.FileHeader)
		thumbnail := new(multipart.FileHeader)

		repo.On("Update", cityId, inputData, image, thumbnail).Return(nil).Once()
		err := srv.Update(cityId, inputData, image, thumbnail)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Error Update City", func(t *testing.T) {
		image := new(multipart.FileHeader)
		thumbnail := new(multipart.FileHeader)
		expectedErr := errors.New("update error")

		repo.On("Update", cityId, inputData, image, thumbnail).Return(expectedErr).Once()
		err := srv.Update(cityId, inputData, image, thumbnail)

		assert.Error(t, err)
		assert.EqualError(t, err, "error update city: update error")
		repo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	repo := new(mocks.CityData)
	srv := NewCity(repo)

	cityId := 1

	t.Run("Success Delete City", func(t *testing.T) {
		repo.On("Delete", cityId).Return(nil).Once()
		err := srv.Delete(cityId)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		invalidCityId := 0

		err := srv.Delete(invalidCityId)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid id")
		repo.AssertNotCalled(t, "Delete")
	})

}

func TestSelectCityById(t *testing.T) {
	repo := new(mocks.CityData)
	srv := NewCity(repo)

	cityId := 1
	expectedCity := city.Core{
		ID:          uint(cityId),
		CityName:    "City Name",
		Description: "Description",
		Image:       "Image",
		Thumbnail:   "Thumbnail",
	}

	t.Run("Success Select City By ID", func(t *testing.T) {
		repo.On("SelectCityById", cityId).Return(expectedCity, nil).Once()
		data, err := srv.SelectCityById(cityId)

		assert.NoError(t, err)
		assert.Equal(t, expectedCity, data)
		repo.AssertExpectations(t)
	})

	t.Run("Error Select City By ID", func(t *testing.T) {
		expectedErr := errors.New("select error")

		repo.On("SelectCityById", cityId).Return(city.Core{}, expectedErr).Once()
		data, err := srv.SelectCityById(cityId)

		assert.Error(t, err)
		assert.EqualError(t, err, "select error")
		assert.Equal(t, city.Core{}, data)
		repo.AssertExpectations(t)
	})
}

func TestSelectAllCity(t *testing.T) {
	repo := new(mocks.CityData)
	srv := NewCity(repo)

	limit := 10
	expectedCityList := []city.Core{
		{ID: uint(1), CityName: "Jakarata"},
		{ID: uint(2), CityName: "Bandung"},
	}
	expectedTotalPage := 2

	t.Run("Success", func(t *testing.T) {
		page := 0
		repo.On("SelectAllCity", 1, limit).Return(expectedCityList, expectedTotalPage, nil).Once()

		cities, totalPage, err := srv.SelectAllCity(page, limit)
		assert.NoError(t, err)
		assert.Equal(t, expectedCityList, cities)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Error with Repository Failure", func(t *testing.T) {
		page := 1
		expectedErr := errors.New("select error")

		repo.On("SelectAllCity", page, limit).Return(nil, 0, expectedErr).Once()
		cities, totalPage, err := srv.SelectAllCity(page, limit)

		assert.Error(t, err)
		assert.EqualError(t, err, "select error")
		assert.Nil(t, cities)
		assert.Zero(t, totalPage)

		repo.AssertExpectations(t)
	})
}
