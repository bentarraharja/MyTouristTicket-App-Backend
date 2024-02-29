package handler

import (
	"my-tourist-ticket/features/booking"
	ph "my-tourist-ticket/features/package/handler"
	th "my-tourist-ticket/features/tour/handler"
)

type BookingResponse struct {
	ID          string `json:"booking_id"`
	UserID      uint   `json:"user_id"`
	TourID      uint   `json:"tour_id"`
	PackageID   uint   `json:"package_id"`
	VoucherID   *uint  `json:"voucher_id"`
	PaymentType string `json:"payment_type"`
	GrossAmount int    `json:"gross_amount"`
	Status      string `json:"status"`
	VaNumber    string `json:"va_number"`
	Bank        string `json:"bank"`
	BookingDate string `json:"booking_date"`
	PhoneNumber string `json:"phone_number"`
	Greeting    string `json:"greeting"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Quantity    int    `json:"quantity"`
	ExpiredAt   string `json:"payment_expired"`
	CreatedAt   string `json:"created_at"`
}

type BookingResponseAdmin struct {
	ID          string          `json:"booking_id"`
	UserID      uint            `json:"user_id"`
	TourID      uint            `json:"tour_id"`
	PackageID   uint            `json:"package_id"`
	VoucherID   *uint           `json:"voucher_id"`
	PaymentType string          `json:"payment_type"`
	GrossAmount int             `json:"gross_amount"`
	Status      string          `json:"status"`
	VaNumber    string          `json:"va_number"`
	Bank        string          `json:"bank"`
	BookingDate string          `json:"booking_date"`
	PhoneNumber string          `json:"phone_number"`
	Greeting    string          `json:"greeting"`
	FullName    string          `json:"full_name"`
	Email       string          `json:"email"`
	Quantity    int             `json:"quantity"`
	ExpiredAt   string          `json:"payment_expired"`
	CreatedAt   string          `json:"created_at"`
	Tour        TourResponse    `json:"tour"`
	Package     PackageResponse `json:"package"`
}

type BookingResponseUserDetail struct {
	ID          string                 `json:"booking_id"`
	FullName    string                 `json:"full_name"`
	Greeting    string                 `json:"greeting"`
	BokingDate  string                 `json:"booking_date"`
	VaNumber    string                 `json:"va_number"`
	Bank        string                 `json:"bank"`
	GrossAmount int                    `json:"gross_amount"`
	Quantity    int                    `json:"quantity"`
	Tour        th.TourResponseBooking `json:"tour"`
	Package     ph.PackageResponseName `json:"package"`
	Voucher     VoucherResponse        `json:"voucher"`
}

type VoucherResponse struct {
	Name string `json:"voucher_name"`
}

type TourResponse struct {
	ID       uint   `json:"id"`
	TourName string `json:"tour_name"`
}

type PackageResponse struct {
	ID          uint   `json:"id"`
	PackageName string `json:"package_name"`
	Price       int    `json:"price"`
}

type BookingResponseUser struct {
	ID          string              `json:"id" form:"id"`
	GrossAmount int                 `json:"gross_amount" form:"gross_amount"`
	Status      string              `json:"status" form:"status"`
	Tour        th.TourResponseName `json:"tour" form:"tour"`
}

func CoreToResponseBooking(core *booking.Core) BookingResponse {
	return BookingResponse{
		ID:          core.ID,
		UserID:      core.UserID,
		TourID:      core.TourID,
		PackageID:   core.PackageID,
		VoucherID:   core.VoucherID,
		PaymentType: core.PaymentType,
		GrossAmount: core.GrossAmount,
		Status:      core.Status,
		VaNumber:    core.VaNumber,
		Bank:        core.Bank,
		BookingDate: core.BookingDate,
		PhoneNumber: core.PhoneNumber,
		Greeting:    core.Greeting,
		FullName:    core.FullName,
		Email:       core.Email,
		Quantity:    core.Quantity,
		ExpiredAt:   core.ExpiredAt.Format("2006-01-02 15:04:05"),
		CreatedAt:   core.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func CoreToResponse(b booking.Core) BookingResponseAdmin {
	return BookingResponseAdmin{
		ID:          b.ID,
		UserID:      b.UserID,
		TourID:      b.TourID,
		PackageID:   b.PackageID,
		VoucherID:   b.VoucherID,
		PaymentType: b.PaymentType,
		GrossAmount: b.GrossAmount,
		Status:      b.Status,
		VaNumber:    b.VaNumber,
		Bank:        b.Bank,
		BookingDate: b.BookingDate,
		PhoneNumber: b.PhoneNumber,
		Greeting:    b.Greeting,
		FullName:    b.FullName,
		Email:       b.Email,
		Quantity:    b.Quantity,
		ExpiredAt:   b.ExpiredAt.Format("2006-01-02 15:04:05"),
		CreatedAt:   b.CreatedAt.Format("2006-01-02 15:04:05"),
		Tour: TourResponse{
			ID:       b.Tour.ID,
			TourName: b.Tour.TourName,
		},
		Package: PackageResponse{
			ID:          b.Package.ID,
			PackageName: b.Package.PackageName,
			Price:       b.Package.Price,
		},
	}
}

func CoreToResponseList(p []booking.Core) []BookingResponseAdmin {
	var results []BookingResponseAdmin
	for _, v := range p {
		results = append(results, CoreToResponse(v))
	}
	return results
}

func CoreToResponseBookingUser(data booking.Core) BookingResponseUser {
	tourResponse := th.TourResponseName{
		TourName: data.Tour.TourName,
	}

	return BookingResponseUser{
		ID:          data.ID,
		GrossAmount: data.GrossAmount,
		Status:      data.Status,
		Tour:        tourResponse,
	}
}

func CoreToResponseListUser(data []booking.Core) []BookingResponseUser {
	var results []BookingResponseUser
	for _, v := range data {
		results = append(results, CoreToResponseBookingUser(v))
	}
	return results
}

func CoreToResponseBookingUserDetail(data *booking.Core) BookingResponseUserDetail {
	var result = BookingResponseUserDetail{
		ID:          data.ID,
		FullName:    data.FullName,
		Greeting:    data.Greeting,
		BokingDate:  data.BookingDate,
		VaNumber:    data.VaNumber,
		Bank:        data.Bank,
		GrossAmount: data.GrossAmount,
		Quantity:    data.Quantity,
		Tour: th.TourResponseBooking{
			TourName: data.Tour.TourName,
			Address:  data.Tour.Address,
			Image:    data.Tour.Image,
		},
		Package: ph.PackageResponseName{
			PackageName: data.Package.PackageName,
		},
		Voucher: VoucherResponse{
			Name: data.Voucher.Name,
		},
	}
	return result
}

type UserReviewResponse struct {
	FullName string `json:"full_name"`
	Image    string `json:"image"`
}

type ReviewResponse struct {
	UserID     uint               `json:"user_id"`
	TextReview string             `json:"text_review"`
	StartRate  float64            `json:"start_rate"`
	CreatedAt  string             `json:"created_at"`
	User       UserReviewResponse `json:"user"`
}

type ReviewTourResponse struct {
	TotalReview   int     `json:"total_review"`
	AverageReview float64 `json:"average_review"`

	Review []ReviewResponse `json:"reviews"`
}

func ReviewCoreToResponse(review booking.ReviewCore) ReviewResponse {
	return ReviewResponse{
		UserID:     review.UserID,
		TextReview: review.TextReview,
		StartRate:  review.StartRate,
		CreatedAt:  review.CreatedAt.Format("2006-01-02 15:04:05"),
		User: UserReviewResponse{
			FullName: review.User.FullName,
			Image:    review.User.Image,
		},
	}
}

func ReviewCoreToResponseList(reviews []booking.ReviewCore) []ReviewResponse {
	var responseList []ReviewResponse
	for _, review := range reviews {
		responseList = append(responseList, ReviewCoreToResponse(review))
	}
	return responseList
}

func ReviewTourCoreToResponse(reviewTour *booking.ReviewTourCore) ReviewTourResponse {
	return ReviewTourResponse{
		TotalReview:   reviewTour.TotalReview,
		AverageReview: reviewTour.AverageReview,
		Review:        ReviewCoreToResponseList(reviewTour.ReviewCore),
	}
}
