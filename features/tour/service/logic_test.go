package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"my-tourist-ticket/features/tour"
	"my-tourist-ticket/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserRoleById(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	userId := 1
	expectedRole := "admin"

	repo.On("GetUserRoleById", userId).Return(expectedRole, nil).Once()
	role, err := srv.GetUserRoleById(userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedRole, role)
	repo.AssertExpectations(t)
}

func TestInsert(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	userId := uint(1)
	input := tour.Core{
		TourName:    "Tour Name",
		Description: "Description",
		Address:     "Address",
		Latitude:    10.12345,
		Longitude:   20.6789,
		CityId:      1,
	}
	image := new(multipart.FileHeader)
	thumbnail := new(multipart.FileHeader)

	t.Run("Success Insert", func(t *testing.T) {
		repo.On("Insert", userId, input, image, thumbnail).Return(nil).Once()
		err := srv.Insert(userId, input, image, thumbnail)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Missing Fields", func(t *testing.T) {
		testCases := []struct {
			name        string
			input       tour.Core
			image       *multipart.FileHeader
			thumbnail   *multipart.FileHeader
			expectedErr string
		}{
			{
				name: "Tour name empty",
				input: tour.Core{
					Description: "Description",
					Address:     "Address",
					Latitude:    10.12345,
					Longitude:   20.6789,
					CityId:      1},
				image:       image,
				thumbnail:   thumbnail,
				expectedErr: "tour name is required",
			},
			{
				name: "Description empty",
				input: tour.Core{
					TourName:  "Tour Name",
					Address:   "Address",
					Latitude:  10.12345,
					Longitude: 20.6789,
					CityId:    1},
				image:       image,
				thumbnail:   thumbnail,
				expectedErr: "description is required",
			},
			{
				name: "Address empty",
				input: tour.Core{
					TourName:    "Tour Name",
					Description: "Description",
					Latitude:    10.12345,
					Longitude:   20.6789,
					CityId:      1},
				image:       image,
				thumbnail:   thumbnail,
				expectedErr: "address is required",
			},
			{
				name: "Image empty",
				input: tour.Core{
					TourName:    "Tour Name",
					Description: "Description",
					Address:     "Address",
					Latitude:    10.12345,
					Longitude:   20.6789,
					CityId:      1},
				image:       nil,
				thumbnail:   thumbnail,
				expectedErr: "image is required",
			},
			{
				name: "Thumbnail empty",
				input: tour.Core{
					TourName:    "Tour Name",
					Description: "Description",
					Address:     "Address",
					Latitude:    10.12345,
					Longitude:   20.6789,
					CityId:      1},
				image:       image,
				thumbnail:   nil,
				expectedErr: "thumbnail is required",
			},
			{
				name: "Latitude empty",
				input: tour.Core{
					TourName:    "Tour Name",
					Description: "Description",
					Address:     "Address",
					Longitude:   20.6789,
					CityId:      1},
				image:       image,
				thumbnail:   thumbnail,
				expectedErr: "latitude is required",
			},
			{
				name: "Longitude empty",
				input: tour.Core{
					TourName:    "Tour Name",
					Description: "Description",
					Address:     "Address",
					Latitude:    10.12345,
					CityId:      1},
				image:       image,
				thumbnail:   thumbnail,
				expectedErr: "longitude is required",
			},
			{
				name: "City ID empty",
				input: tour.Core{
					TourName:    "Tour Name",
					Description: "Description",
					Address:     "Address",
					Latitude:    10.12345,
					Longitude:   20.6789},
				image:       image,
				thumbnail:   thumbnail,
				expectedErr: "city id is required",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := srv.Insert(userId, tc.input, tc.image, tc.thumbnail)

				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr)
				repo.AssertNotCalled(t, "Insert")
			})
		}
	})

	t.Run("Error Insert", func(t *testing.T) {
		expectedErr := errors.New("insert error")
		repo.On("Insert", userId, input, image, thumbnail).Return(expectedErr).Once()

		err := srv.Insert(userId, input, image, thumbnail)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("error creating tour: %v", expectedErr))
		repo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	tourID := 1
	input := tour.Core{
		TourName:    "Tour Name",
		Description: "Description",
		Address:     "Address",
		Latitude:    10.12345,
		Longitude:   20.6789,
		CityId:      1,
	}
	image := new(multipart.FileHeader)
	thumbnail := new(multipart.FileHeader)
	expectedErr := errors.New("update error")
	repo.On("Update", tourID, input, image, thumbnail).Return(expectedErr).Once()

	err := srv.Update(tourID, input, image, thumbnail)

	assert.EqualError(t, err, expectedErr.Error())
	repo.AssertExpectations(t)
}

func TestSelectTourById(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	tourID := 123
	expectedTour := tour.Core{
		TourName:    "Tour Name",
		Description: "Description",
		Address:     "Address",
		Latitude:    10.12345,
		Longitude:   20.6789,
		CityId:      1,
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("SelectTourById", tourID).Return(expectedTour, nil).Once()
		tourData, err := srv.SelectTourById(tourID)

		assert.NoError(t, err)
		assert.Equal(t, expectedTour, tourData)

		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("select tour by id error")
		repo.On("SelectTourById", tourID).Return(tour.Core{}, expectedError).Once()
		tourData, err := srv.SelectTourById(tourID)

		assert.Error(t, err)
		assert.EqualError(t, err, "select tour by id error")
		assert.Equal(t, tour.Core{}, tourData)

		repo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	cityId := 1

	t.Run("Success Delete Tour", func(t *testing.T) {
		repo.On("Delete", cityId).Return(nil).Once()
		err := srv.Delete(cityId)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		invalidTourId := 0

		err := srv.Delete(invalidTourId)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid id")
		repo.AssertNotCalled(t, "Delete")
	})

}

func TestSelectAllTour(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	limit := 12
	expectedTourList := []tour.Core{
		{ID: uint(1), TourName: "Ancol", Description: "Ancol adalah", Address: "Jl.Ancol", Image: "image.com", Thumbnail: "thumbnail.com"},
		{ID: uint(2), TourName: "Dufan", Description: "Dufan adalah", Address: "Jl.Dufan", Image: "image.com", Thumbnail: "thumbnail.com"},
	}
	expectedTotalPage := 1

	t.Run("Success with Page and Limit Default Values", func(t *testing.T) {
		expectedPage := 1
		expectedLimit := 12

		repo.On("SelectAllTour", expectedPage, expectedLimit).Return(expectedTourList, expectedTotalPage, nil).Once()

		tours, totalPage, err := srv.SelectAllTour(0, 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedTourList, tours)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		page := 1

		repo.On("SelectAllTour", page, limit).Return(expectedTourList, expectedTotalPage, nil).Once()
		tours, totalPage, err := srv.SelectAllTour(page, limit)

		assert.NoError(t, err)

		assert.Equal(t, expectedTourList, tours)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Error with Repository Failure", func(t *testing.T) {
		page := 1

		expectedErr := errors.New("select error")
		repo.On("SelectAllTour", page, limit).Return(nil, 0, expectedErr).Once()

		tours, totalPage, err := srv.SelectAllTour(page, limit)

		assert.Error(t, err)
		assert.EqualError(t, err, "select error")

		assert.Nil(t, tours)
		assert.Zero(t, totalPage)

		repo.AssertExpectations(t)
	})
}

func TestSelectTourByPengelola(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	userId := 1
	limit := 12
	expectedTourList := []tour.Core{
		{ID: 1, TourName: "Ancol", Description: "Ancol adalah", Address: "Jl.Ancol", Image: "image.com", Thumbnail: "thumbnail.com"},
		{ID: 2, TourName: "Dufan", Description: "Dufan adalah", Address: "Jl.Dufan", Image: "image.com", Thumbnail: "thumbnail.com"},
	}
	expectedTotalPage := 1

	t.Run("Success with Page and Limit Default Values", func(t *testing.T) {
		expectedPage := 1
		expectedLimit := 12

		repo.On("SelectTourByPengelola", userId, expectedPage, expectedLimit).Return(expectedTourList, expectedTotalPage, nil).Once()

		tours, totalPage, err := srv.SelectTourByPengelola(userId, 0, 0)
		assert.NoError(t, err)

		assert.Equal(t, expectedTourList, tours)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		page := 1

		repo.On("SelectTourByPengelola", userId, page, limit).Return(expectedTourList, expectedTotalPage, nil).Once()

		tours, totalPage, err := srv.SelectTourByPengelola(userId, page, limit)

		assert.NoError(t, err)

		assert.Equal(t, expectedTourList, tours)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Error with Repository Failure", func(t *testing.T) {
		page := 1
		expectedErr := errors.New("select error")

		repo.On("SelectTourByPengelola", userId, page, limit).Return(nil, 0, expectedErr).Once()
		tours, totalPage, err := srv.SelectTourByPengelola(userId, page, limit)

		assert.Error(t, err)
		assert.EqualError(t, err, "select error")

		assert.Nil(t, tours)
		assert.Zero(t, totalPage)

		repo.AssertExpectations(t)
	})
}

func TestGetTourByCityID(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	cityId := uint(1)
	limit := 6
	expectedTourList := []tour.Core{
		{ID: 1, TourName: "Ancol", Description: "Ancol adalah", Address: "Jl.Ancol", Image: "image.com", Thumbnail: "thumbnail.com"},
		{ID: 2, TourName: "Dufan", Description: "Dufan adalah", Address: "Jl.Dufan", Image: "image.com", Thumbnail: "thumbnail.com"},
	}
	expectedTotalPage := 1

	t.Run("Success with Page and Limit Default Values", func(t *testing.T) {
		expectedPage := 1
		expectedLimit := 6

		repo.On("GetTourByCityID", cityId, expectedPage, expectedLimit).Return(expectedTourList, expectedTotalPage, nil).Once()

		tours, totalPage, err := srv.GetTourByCityID(cityId, 0, 0)

		assert.NoError(t, err)

		assert.Equal(t, expectedTourList, tours)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		page := 1

		repo.On("GetTourByCityID", cityId, page, limit).Return(expectedTourList, expectedTotalPage, nil).Once()

		tours, totalPage, err := srv.GetTourByCityID(cityId, page, limit)

		assert.NoError(t, err)
		assert.Equal(t, expectedTourList, tours)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Error with Repository Failure", func(t *testing.T) {
		page := 1

		expectedErr := errors.New("select error")

		repo.On("GetTourByCityID", cityId, page, limit).Return(nil, 0, expectedErr).Once()

		tours, totalPage, err := srv.GetTourByCityID(cityId, page, limit)

		assert.Error(t, err)
		assert.EqualError(t, err, "select error")

		assert.Nil(t, tours)
		assert.Zero(t, totalPage)

		repo.AssertExpectations(t)
	})
}

func TestInsertReportTour(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	userId := 1
	tourId := 2
	reportInput := tour.ReportCore{
		TextReport: "Tour ini butuh improvement",
	}

	t.Run("Error", func(t *testing.T) {
		expectedErr := errors.New("insert report error")
		repo.On("InsertReportTour", userId, tourId, reportInput).Return(expectedErr).Once()

		err := srv.InsertReportTour(userId, tourId, reportInput)

		assert.EqualError(t, err, "error creating report: insert report error")

		repo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		repo.On("InsertReportTour", userId, tourId, reportInput).Return(nil).Once()

		err := srv.InsertReportTour(userId, tourId, reportInput)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestSelectReportTour(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	tourID := 1
	expectedReports := []tour.ReportCore{
		{ID: 1, UserId: 1, TextReport: "Tour jelek!"},
		{ID: 2, UserId: 2, TextReport: "Bad Tour!"},
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("SelectReportTour", tourID).Return(expectedReports, nil).Once()

		reports, err := srv.SelectReportTour(tourID)

		assert.NoError(t, err)
		assert.Equal(t, expectedReports, reports)

		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedErr := errors.New("select report error")
		repo.On("SelectReportTour", tourID).Return(nil, expectedErr).Once()

		reports, err := srv.SelectReportTour(tourID)

		assert.Error(t, err)
		assert.Nil(t, reports)
		assert.EqualError(t, err, "error get report: select report error")

		repo.AssertExpectations(t)
	})
}

func TestSearchTour(t *testing.T) {
	repo := new(mocks.TourData)
	srv := NewTour(repo)

	query := "nama tour"
	expectedTours := []tour.Core{
		{ID: 1, TourName: "Dufan"},
		{ID: 2, TourName: "Ancol"},
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("SearchTour", query).Return(expectedTours, nil).Once()

		tours, err := srv.SearchTour(query)

		assert.NoError(t, err)
		assert.Equal(t, expectedTours, tours)

		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedErr := errors.New("search error")
		repo.On("SearchTour", query).Return(nil, expectedErr).Once()

		tours, err := srv.SearchTour(query)
		assert.Error(t, err)
		assert.Nil(t, tours)
		assert.EqualError(t, err, "search error")

		repo.AssertExpectations(t)
	})
}
