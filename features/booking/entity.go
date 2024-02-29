package booking

import (
	packages "my-tourist-ticket/features/package"
	"my-tourist-ticket/features/tour"
	"my-tourist-ticket/features/user"
	"my-tourist-ticket/features/voucher"
	"time"
)

type Core struct {
	ID          string
	UserID      uint
	TourID      uint
	PackageID   uint
	VoucherID   *uint
	PaymentType string
	GrossAmount int
	Status      string
	VaNumber    string
	Bank        string
	BookingDate string
	PhoneNumber string
	Greeting    string
	FullName    string
	Email       string
	Quantity    int
	ExpiredAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        user.Core
	Tour        tour.Core
	Package     packages.Core
	Voucher     voucher.Core
}

type ReviewCore struct {
	ID         uint
	BookingID  string
	UserID     uint
	TextReview string
	StartRate  float64
	CreatedAt  time.Time
	User       user.Core
}

type ReviewTourCore struct {
	TotalReview   int
	AverageReview float64
	ReviewCore    []ReviewCore
}

// interface untuk Data Layer
type BookingDataInterface interface {
	InsertBooking(userIdLogin int, inputBooking Core) (*Core, error)
	CancelBooking(userIdLogin int, orderId string, bookingCore Core) error
	InsertBookingReview(inputReview ReviewCore) error
	WebhoocksData(reqNotif Core) error
	SelectBookingUser(userIdLogin int) ([]Core, error)
	SelectBookingUserDetail(userIdLogin int, bookingId string) (*Core, error)
	SelectAllBooking(page, limit int) ([]Core, int, error)
	GetUserRoleById(userId int) (string, error)
	SelectAllBookingPengelola(pengelolaID int, page, limit int) ([]Core, int, error)
	GetAllBookingReview(tourId, limit int) ([]ReviewCore, error)
	GetTotalTourReview(tourId int) (int, error)
	GetAverageTourReview(tourId int) (float64, error)
}

// interface untuk Service Layer
type BookingServiceInterface interface {
	CreateBooking(userIdLogin int, inputBooking Core) (*Core, error)
	CancelBooking(userIdLogin int, orderId string, bookingCore Core) error
	CreateBookingReview(inputReview ReviewCore) error
	WebhoocksService(reqNotif Core) error
	GetBookingUser(userIdLogin int) ([]Core, error)
	GetBookingUserDetail(userIdLogin int, bookingId string) (*Core, error)
	SelectAllBooking(page, limit int) ([]Core, int, error)
	GetUserRoleById(userId int) (string, error)
	SelectAllBookingPengelola(pengelolaID int, page, limit int) ([]Core, int, error)
	GetAllBookingReview(tourId, limit int) (*ReviewTourCore, error)
}
