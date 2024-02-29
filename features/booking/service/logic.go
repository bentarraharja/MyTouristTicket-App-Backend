package service

import (
	"errors"
	"fmt"
	"my-tourist-ticket/features/booking"
)

type bookingService struct {
	bookingData booking.BookingDataInterface
}

func New(repo booking.BookingDataInterface) booking.BookingServiceInterface {
	return &bookingService{
		bookingData: repo,
	}
}

// GetUserRoleById implements booking.BookingServiceInterface.
func (service *bookingService) GetUserRoleById(userId int) (string, error) {
	return service.bookingData.GetUserRoleById(userId)
}

// CreateBooking implements booking.BookingServiceInterface.
func (service *bookingService) CreateBooking(userIdLogin int, inputBooking booking.Core) (*booking.Core, error) {
	if inputBooking.Bank == "" {
		return nil, errors.New("bank is required")
	}
	if inputBooking.PhoneNumber == "" {
		return nil, errors.New("phone number is required")
	}
	if inputBooking.Greeting == "" {
		return nil, errors.New("greeting is required")
	}
	if inputBooking.FullName == "" {
		return nil, errors.New("full name number is required")
	}
	if inputBooking.Email == "" {
		return nil, errors.New("email is required")
	}
	if inputBooking.BookingDate == "" {
		return nil, errors.New("booking date is required")
	}

	payment, err := service.bookingData.InsertBooking(userIdLogin, inputBooking)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (service *bookingService) CancelBooking(userIdLogin int, bookingId string, bookingCore booking.Core) error {
	if bookingCore.Status == "" {
		bookingCore.Status = "cancelled"
	}

	err := service.bookingData.CancelBooking(userIdLogin, bookingId, bookingCore)
	return err
}

// CreateBookingReview implements booking.BookingServiceInterface.
func (service *bookingService) CreateBookingReview(inputReview booking.ReviewCore) error {
	if inputReview.TextReview == "" {
		return errors.New("text review is required")
	}
	if inputReview.StartRate == 0 {
		return errors.New("rate is required")
	} else if inputReview.StartRate > 5 {
		return errors.New("star rate is not valid")
	}

	err := service.bookingData.InsertBookingReview(inputReview)
	if err != nil {
		return fmt.Errorf("error creating review: %w", err)
	}
	return nil
}

// WebhoocksService implements booking.BookingServiceInterface.
func (service *bookingService) WebhoocksService(reqNotif booking.Core) error {
	if reqNotif.ID == "" {
		return errors.New("invalid order id")
	}

	err := service.bookingData.WebhoocksData(reqNotif)
	if err != nil {
		return err
	}

	return nil
}

func (service *bookingService) GetBookingUser(userIdLogin int) ([]booking.Core, error) {
	results, err := service.bookingData.SelectBookingUser(userIdLogin)
	return results, err
}

func (service *bookingService) GetBookingUserDetail(userIdLogin int, bookingId string) (*booking.Core, error) {
	result, err := service.bookingData.SelectBookingUserDetail(userIdLogin, bookingId)
	return result, err
}

// SelectAllBooking implements booking.BookingServiceInterface.
func (service *bookingService) SelectAllBooking(page int, limit int) ([]booking.Core, int, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 8
	}

	bookings, totalPage, err := service.bookingData.SelectAllBooking(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return bookings, totalPage, nil
}

// SelectAllPengelola implements booking.BookingServiceInterface.
func (service *bookingService) SelectAllBookingPengelola(pengelolaID int, page int, limit int) ([]booking.Core, int, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 8
	}

	bookings, totalPage, err := service.bookingData.SelectAllBookingPengelola(pengelolaID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return bookings, totalPage, nil
}

func (service *bookingService) GetAllBookingReview(tourId, limit int) (*booking.ReviewTourCore, error) {
	if limit == 0 {
		limit = 2
	}
	rev := &booking.ReviewTourCore{}

	reviews, err := service.bookingData.GetAllBookingReview(tourId, limit)
	if err != nil {
		return nil, err
	}
	rev.ReviewCore = reviews

	// Calculate average review and total reviews
	averageReview, err := service.bookingData.GetAverageTourReview(tourId)
	if err != nil {
		return nil, err
	}
	rev.AverageReview = averageReview

	totalReviews, err := service.bookingData.GetTotalTourReview(tourId)
	if err != nil {
		return nil, err
	}
	rev.TotalReview = int(totalReviews)

	return rev, nil
}
