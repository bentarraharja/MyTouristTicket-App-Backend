package handler

import "my-tourist-ticket/features/dashboard"

type DashboardResponse struct {
	TotalCustomer    int `json:"total_costumer"`
	TotalPengelola   int `json:"total_pengelola"`
	TotalTransaction int `json:"total_transaction"`
	TotalTour        int `json:"total_tour"`

	RecentBooking []BookingResponse `json:"recent_booking"`

	TopTours []TourResponse `json:"top_tours"`
}
type BookingResponse struct {
	ID string `json:"booking_id"`
	// UserID      uint             `json:"user_id"`
	// TourID      uint             `json:"tour_id"`
	// PackageID   uint             `json:"package_id"`
	// VoucherID   *uint            `json:"voucher_id"`
	// PaymentType string           `json:"payment_type"`
	GrossAmount int `json:"gross_amount"`
	// Status      string           `json:"status"`
	// VaNumber    string           `json:"va_number"`
	// Bank        string           `json:"bank"`
	// PhoneNumber string           `json:"phone_number"`
	// Greeting    string           `json:"greeting"`
	// FullName    string           `json:"full_name"`
	// Email       string           `json:"email"`
	// Quantity    int              `json:"quantity"`
	// ExpiredAt   string           `json:"payment_expired"`
	// CreatedAt   string           `json:"created_at"`
	Tour TourResponseName `json:"tour"`
	// Package     PackageResponse  `json:"package"`
}

// type PackageResponse struct {
// 	Price int `json:"price"`
// }

type TourResponseName struct {
	TourName string `json:"tour_name"`
}

type TourResponse struct {
	ID uint `json:"id"`
	// CityId      uint         `json:"city_id"`
	// UserId      uint         `json:"user_id"`
	TourName string `json:"tour_name"`
	// Description string       `json:"description"`
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
	// Address     string       `json:"address"`
	// Latitude    float64      `json:"latitude"`
	// Longitude   float64      `json:"longitude"`
	// CreatedAt   string       `json:"created_at"`
	// UpdatedAt   string       `json:"updated_at"`
	City CityResponse `json:"city"`
}

type CityResponse struct {
	CityName string `json:"city_name"`
}

func CoreToBookingResponseList(bookings []dashboard.Booking) []BookingResponse {
	var responseList []BookingResponse
	for _, booking := range bookings {
		responseList = append(responseList, BookingResponse{
			ID:          booking.ID,
			GrossAmount: booking.GrossAmount,
			Tour: TourResponseName{
				TourName: booking.Tour.TourName,
			},
			// Package: PackageResponse{
			// 	Price: booking.Package.Price,
			// },
		})
	}
	return responseList
}

func CoreToTourResponseList(tours []dashboard.Tour) []TourResponse {
	var responseList []TourResponse
	for _, tour := range tours {
		responseList = append(responseList, TourResponse{
			ID:        tour.ID,
			TourName:  tour.TourName,
			Image:     tour.Image,
			Thumbnail: tour.Thumbnail,
			City: CityResponse{
				CityName: tour.City.CityName,
			},
		})
	}
	return responseList
}

func CoreToDashboardResponse(core *dashboard.Dashboard) DashboardResponse {
	return DashboardResponse{
		TotalCustomer:    core.TotalCustomer,
		TotalPengelola:   core.TotalPengelola,
		TotalTransaction: core.TotalTransaction,
		TotalTour:        core.TotalTour,
		RecentBooking:    CoreToBookingResponseList(core.RecentBooking),
		TopTours:         CoreToTourResponseList(core.TopTours),
	}
}
