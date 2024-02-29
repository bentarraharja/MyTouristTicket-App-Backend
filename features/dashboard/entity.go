package dashboard

import (
	"my-tourist-ticket/features/city"
	packages "my-tourist-ticket/features/package"
	"my-tourist-ticket/features/tour"
	"time"
)

type Dashboard struct {
	TotalCustomer    int
	TotalPengelola   int
	TotalTransaction int
	TotalTour        int

	RecentBooking []Booking

	TopTours []Tour
}

type Booking struct {
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
	PhoneNumber string
	Greeting    string
	FullName    string
	Email       string
	Quantity    int
	ExpiredAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// User        user.Core
	Tour    tour.Core
	Package packages.Core
	// Voucher     voucher.Core
}

type Tour struct {
	ID          uint
	CityId      uint
	UserId      uint
	TourName    string
	Description string
	Image       string
	Thumbnail   string
	Addres      string
	Latitude    float64
	Longitude   float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	City        city.Core
}

type DashboardDataInterface interface {
	GetTotalCustomer() (int, error)
	GetTotalPengelola() (int, error)
	GetTotalTransaction() (int, error)
	GetTotalTour() (int, error)
	GetRecentTransaction() ([]Booking, error)
	GetTopTour() ([]Tour, error)
	GetUserRoleById(userId int) (string, error)
}

type DashboardServiceInterface interface {
	Dashboard(userId int) (*Dashboard, error)
}
