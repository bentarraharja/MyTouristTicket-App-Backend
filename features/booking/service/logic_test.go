package service

import (
	"errors"
	"my-tourist-ticket/features/booking"
	"my-tourist-ticket/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserRoleById(t *testing.T) {
	mockRepo := new(mocks.BookingData)
	service := New(mockRepo)

	t.Run("Success Get User Role by ID", func(t *testing.T) {
		expectedRole := "user"
		mockRepo.On("GetUserRoleById", 1).Return(expectedRole, nil).Once()

		role, err := service.GetUserRoleById(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedRole, role)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Get User Role by ID", func(t *testing.T) {
		expectedErr := errors.New("get user role error")
		mockRepo.On("GetUserRoleById", 2).Return("", expectedErr).Once()

		role, err := service.GetUserRoleById(2)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Empty(t, role)

		mockRepo.AssertExpectations(t)
	})
}

func TestCreateBooking(t *testing.T) {
	mockRepo := new(mocks.BookingData)
	service := New(mockRepo)

	inputBooking := booking.Core{
		Bank:        "Bank Name",
		PhoneNumber: "123456789",
		Greeting:    "Hello",
		FullName:    "John Doe",
		Email:       "john@example.com",
		BookingDate: "2024-02-17",
	}

	t.Run("Success Create Booking", func(t *testing.T) {
		mockRepo.On("InsertBooking", mock.AnythingOfType("int"), inputBooking).Return(&inputBooking, nil).Once()

		result, err := service.CreateBooking(1, inputBooking)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, inputBooking, *result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Create Booking", func(t *testing.T) {
		expectedErr := errors.New("insert error")
		mockRepo.On("InsertBooking", mock.AnythingOfType("int"), inputBooking).Return(nil, expectedErr).Once()

		result, err := service.CreateBooking(1, inputBooking)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, expectedErr.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation for Empty Fields", func(t *testing.T) {
		tests := []struct {
			name     string
			booking  booking.Core
			expected error
		}{
			{
				name:     "Empty Bank",
				booking:  booking.Core{PhoneNumber: "123456789", Greeting: "Hello", FullName: "John Doe", Email: "john@example.com", BookingDate: "2024-02-17"},
				expected: errors.New("bank is required"),
			},
			{
				name:     "Empty Phone Number",
				booking:  booking.Core{Bank: "Bank Name", Greeting: "Hello", FullName: "John Doe", Email: "john@example.com", BookingDate: "2024-02-17"},
				expected: errors.New("phone number is required"),
			},
			{
				name:     "Empty Greeting",
				booking:  booking.Core{Bank: "Bank Name", PhoneNumber: "123456789", FullName: "John Doe", Email: "john@example.com", BookingDate: "2024-02-17"},
				expected: errors.New("greeting is required"),
			},
			{
				name:     "Empty Full Name",
				booking:  booking.Core{Bank: "Bank Name", PhoneNumber: "123456789", Greeting: "Hello", Email: "john@example.com", BookingDate: "2024-02-17"},
				expected: errors.New("full name number is required"),
			},
			{
				name:     "Empty Email",
				booking:  booking.Core{Bank: "Bank Name", PhoneNumber: "123456789", Greeting: "Hello", FullName: "John Doe", BookingDate: "2024-02-17"},
				expected: errors.New("email is required"),
			},
			{
				name:     "Empty Booking Date",
				booking:  booking.Core{Bank: "Bank Name", PhoneNumber: "123456789", Greeting: "Hello", FullName: "John Doe", Email: "john@example.com"},
				expected: errors.New("booking date is required"),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				_, err := service.CreateBooking(1, test.booking)
				assert.Equal(t, test.expected, err)
			})
		}
	})
}

func TestCancelBooking(t *testing.T) {
	mockRepo := new(mocks.BookingData)
	service := New(mockRepo)

	inputBooking := booking.Core{
		Status: "cancelled",
	}

	t.Run("Success Cancel Booking", func(t *testing.T) {
		mockRepo.On("CancelBooking", 1, "1", inputBooking).Return(nil).Once()

		inputBooking.Status = ""
		err := service.CancelBooking(1, "1", inputBooking)
		inputBooking.Status = "cancelled"
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Cancel Booking", func(t *testing.T) {
		expectedErr := errors.New("cancel error")
		mockRepo.On("CancelBooking", mock.AnythingOfType("int"), "booking_id", inputBooking).Return(expectedErr).Once()

		err := service.CancelBooking(1, "booking_id", inputBooking)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestCreateBookingReview(t *testing.T) {
	mockRepo := new(mocks.BookingData)
	service := New(mockRepo)

	inputReview := booking.ReviewCore{
		TextReview: "Great experience!",
		StartRate:  5,
	}

	t.Run("Success Create Booking Review", func(t *testing.T) {
		mockRepo.On("InsertBookingReview", inputReview).Return(nil).Once()

		err := service.CreateBookingReview(inputReview)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Create Booking Review (Empty Text Review)", func(t *testing.T) {
		inputReview.TextReview = ""
		expectedErr := errors.New("text review is required")

		err := service.CreateBookingReview(inputReview)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Create Booking Review (star rate 0 Review)", func(t *testing.T) {
		inputReview.TextReview = "good"
		inputReview.StartRate = 0
		expectedErr := errors.New("rate is required")

		err := service.CreateBookingReview(inputReview)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Create Booking Review (Invalid Star Rate)", func(t *testing.T) {
		inputReview.TextReview = "Great experience!"
		inputReview.StartRate = 6
		expectedErr := errors.New("star rate is not valid")

		err := service.CreateBookingReview(inputReview)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error creating review", func(t *testing.T) {
		inputReview.StartRate = 5
		mockRepo.On("InsertBookingReview", inputReview).Return(errors.New("create review error"))

		err := service.CreateBookingReview(inputReview)

		assert.Error(t, err)
		assert.EqualError(t, err, "error creating review: create review error")

		mockRepo.AssertExpectations(t)
	})
}

func TestWebhoocksService(t *testing.T) {
	mockRepo := new(mocks.BookingData)
	service := New(mockRepo)

	t.Run("Invalid Order ID", func(t *testing.T) {
		invalidReqNotif := booking.Core{}
		err := service.WebhoocksService(invalidReqNotif)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid order id")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from WebhoocksData", func(t *testing.T) {
		validReqNotif := booking.Core{ID: "valid_id"}
		expectedErr := errors.New("webhoocks error")
		mockRepo.On("WebhoocksData", validReqNotif).Return(expectedErr).Once()
		err := service.WebhoocksService(validReqNotif)
		assert.Error(t, err)
		assert.EqualError(t, err, "webhoocks error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		validReqNotif := booking.Core{ID: "valid_id"}
		mockRepo.On("WebhoocksData", validReqNotif).Return(nil).Once()
		err := service.WebhoocksService(validReqNotif)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetBookingUser(t *testing.T) {
	// Inisialisasi mock
	mockRepo := new(mocks.BookingData)
	service := New(mockRepo)

	// Data yang diharapkan dari panggilan repositori
	expectedResults := []booking.Core{
		{ID: "1", UserID: 1, Status: "active"},
		{ID: "2", UserID: 1, Status: "cancelled"},
	}
	t.Run("Success", func(t *testing.T) {
		// Menyiapkan ekspektasi panggilan repositori (expected)
		mockRepo.On("SelectBookingUser", mock.AnythingOfType("int")).Return(expectedResults, nil).Once()

		// Memanggil metode yang diuji (actual)
		results, err := service.GetBookingUser(1)

		// Memastikan tidak ada kesalahan dan hasil yang diharapkan dikembalikan
		assert.NoError(t, err)
		assert.Equal(t, expectedResults, results)

		// Memastikan ekspektasi dipanggil dengan tepat
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from SelectBookingUser", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockRepo.On("SelectBookingUser", mock.AnythingOfType("int")).Return(nil, expectedErr).Once()

		results, err := service.GetBookingUser(1)

		// Memastikan kesalahan yang diharapkan dikembalikan
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		// Memastikan hasil nil dikembalikan
		assert.Nil(t, results)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetBookingUserDetail(t *testing.T) {
	mockRepo := new(mocks.BookingData)
	service := New(mockRepo)

	expectedBooking := &booking.Core{
		ID:     "1",
		UserID: 1,
		Status: "active",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("SelectBookingUserDetail", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(expectedBooking, nil).Once()

		result, err := service.GetBookingUserDetail(1, "booking_id")

		assert.NoError(t, err)
		assert.Equal(t, expectedBooking, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from SelectBookingUserDetail", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockRepo.On("SelectBookingUserDetail", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil, expectedErr).Once()

		result, err := service.GetBookingUserDetail(1, "booking_id")

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSelectAllBooking(t *testing.T) {
	// Inisialisasi mock
	mockRepo := new(mocks.BookingData)
	service := New(mockRepo)

	expectedBookings := []booking.Core{
		{ID: "1", UserID: 1, Status: "active"},
		{ID: "2", UserID: 1, Status: "cancelled"},
	}

	t.Run("Success with Default Values", func(t *testing.T) {
		expectedTotalPage := 1
		mockRepo.On("SelectAllBooking", 1, 8).Return(expectedBookings, expectedTotalPage, nil).Once()

		bookings, totalPage, err := service.SelectAllBooking(0, 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedBookings, bookings)
		assert.Equal(t, expectedTotalPage, totalPage)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		expectedTotalPage := 2
		mockRepo.On("SelectAllBooking", 2, 10).Return(expectedBookings, expectedTotalPage, nil).Once()

		bookings, totalPage, err := service.SelectAllBooking(2, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedBookings, bookings)
		assert.Equal(t, expectedTotalPage, totalPage)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from SelectAllBooking", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockRepo.On("SelectAllBooking", 1, 8).Return(nil, 0, expectedErr).Once()

		bookings, totalPage, err := service.SelectAllBooking(0, 0)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, bookings)
		assert.Zero(t, totalPage)
		mockRepo.AssertExpectations(t)
	})
}

func TestSelectAllBookingPengelola(t *testing.T) {
	mockRepo := new(mocks.BookingData)
	service := New(mockRepo)

	expectedBookings := []booking.Core{
		{ID: "1", UserID: 1, Status: "active"},
		{ID: "2", UserID: 1, Status: "cancelled"},
	}

	t.Run("Success with Default Values", func(t *testing.T) {
		expectedTotalPage := 1
		mockRepo.On("SelectAllBookingPengelola", 12, 1, 8).Return(expectedBookings, expectedTotalPage, nil).Once()

		bookings, totalPage, err := service.SelectAllBookingPengelola(12, 0, 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedBookings, bookings)
		assert.Equal(t, expectedTotalPage, totalPage)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		expectedTotalPage := 2
		mockRepo.On("SelectAllBookingPengelola", 12, 2, 10).Return(expectedBookings, expectedTotalPage, nil).Once()

		bookings, totalPage, err := service.SelectAllBookingPengelola(12, 2, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedBookings, bookings)
		assert.Equal(t, expectedTotalPage, totalPage)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from SelectAllBookingPengelola", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockRepo.On("SelectAllBookingPengelola", 12, 1, 8).Return(nil, 0, expectedErr).Once()

		bookings, totalPage, err := service.SelectAllBookingPengelola(12, 0, 0)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, bookings)
		assert.Zero(t, totalPage)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllBookingReview(t *testing.T) {
	mockRepo := new(mocks.BookingData)
	service := New(mockRepo)

	expectedReviews := []booking.ReviewCore{
		{
			ID:         123,
			BookingID:  "booking123",
			UserID:     456,
			TextReview: "Great experience",
			StartRate:  4.5},
		{
			ID:         123,
			BookingID:  "booking123",
			UserID:     456,
			TextReview: "Great experience",
			StartRate:  4.5},
	}

	t.Run("Success with Default Values", func(t *testing.T) {
		expectedAverageReview := 4.5
		expectedTotalReviews := 2

		mockRepo.On("GetAllBookingReview", 123, 2).Return(expectedReviews, nil).Once()
		mockRepo.On("GetAverageTourReview", 123).Return(expectedAverageReview, nil).Once()
		mockRepo.On("GetTotalTourReview", 123).Return(expectedTotalReviews, nil).Once()

		review, err := service.GetAllBookingReview(123, 0)

		assert.NoError(t, err)
		assert.NotNil(t, review)
		assert.Equal(t, expectedReviews, review.ReviewCore)
		assert.Equal(t, expectedAverageReview, review.AverageReview)
		assert.Equal(t, expectedTotalReviews, review.TotalReview)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		expectedAverageReview := 4.5
		expectedTotalReviews := 2

		mockRepo.On("GetAllBookingReview", 123, 5).Return(expectedReviews, nil).Once()
		mockRepo.On("GetAverageTourReview", 123).Return(expectedAverageReview, nil).Once()
		mockRepo.On("GetTotalTourReview", 123).Return(expectedTotalReviews, nil).Once()

		review, err := service.GetAllBookingReview(123, 5)

		assert.NoError(t, err)
		assert.NotNil(t, review)
		assert.Equal(t, expectedReviews, review.ReviewCore)
		assert.Equal(t, expectedAverageReview, review.AverageReview)
		assert.Equal(t, expectedTotalReviews, review.TotalReview)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from GetAllBookingReview", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockRepo.On("GetAllBookingReview", 123, 2).Return(nil, expectedErr).Once()

		review, err := service.GetAllBookingReview(123, 0)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, review)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from GetAverageTourReview", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockRepo.On("GetAllBookingReview", 123, 2).Return([]booking.ReviewCore{}, nil).Once()
		mockRepo.On("GetAverageTourReview", 123).Return(0.0, expectedErr).Once()

		review, err := service.GetAllBookingReview(123, 0)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, review)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from GetTotalTourReview", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockRepo.On("GetAllBookingReview", 123, 2).Return([]booking.ReviewCore{}, nil).Once()
		mockRepo.On("GetAverageTourReview", 123).Return(0.0, nil).Once()
		mockRepo.On("GetTotalTourReview", 123).Return(0, expectedErr).Once()

		review, err := service.GetAllBookingReview(123, 0)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, review)
		mockRepo.AssertExpectations(t)
	})
}
