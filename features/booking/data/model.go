package data

import (
	"my-tourist-ticket/features/booking"
	packages "my-tourist-ticket/features/package"
	pd "my-tourist-ticket/features/package/data"
	"my-tourist-ticket/features/tour"
	td "my-tourist-ticket/features/tour/data"
	"my-tourist-ticket/features/user"
	ud "my-tourist-ticket/features/user/data"
	"my-tourist-ticket/features/voucher"
	vd "my-tourist-ticket/features/voucher/data"
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID string `gorm:"type:varchar(36);primary_key" json:"id"`
	gorm.Model
	UserID      uint  `gorm:"not null"`
	TourID      uint  `gorm:"not null"`
	PackageID   uint  `gorm:"not null"`
	VoucherID   *uint `gorm:"default:null;omitempty"`
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
	User        ud.User
	Tour        td.Tour
	Package     pd.Package
	Voucher     vd.Voucher
}

type Review struct {
	gorm.Model
	BookingID  string `gorm:"not null;unique"`
	UserID     uint   `gorm:"not null"`
	TextReview string
	StartRate  float64
	Booking    Booking
	User       ud.User
}

func CoreToModelBooking(input booking.Core) Booking {
	return Booking{
		ID:          input.ID,
		UserID:      input.UserID,
		TourID:      input.TourID,
		PackageID:   input.PackageID,
		VoucherID:   input.VoucherID,
		PaymentType: input.PaymentType,
		GrossAmount: input.GrossAmount,
		Status:      input.Status,
		VaNumber:    input.VaNumber,
		Bank:        input.Bank,
		BookingDate: input.BookingDate,
		PhoneNumber: input.PhoneNumber,
		Greeting:    input.Greeting,
		FullName:    input.FullName,
		Email:       input.Email,
		Quantity:    input.Quantity,
		ExpiredAt:   input.ExpiredAt,
	}
}

func ModelToCoreBooking(model Booking) booking.Core {
	return booking.Core{
		ID:          model.ID,
		UserID:      model.UserID,
		TourID:      model.TourID,
		PackageID:   model.PackageID,
		VoucherID:   model.VoucherID,
		PaymentType: model.PaymentType,
		GrossAmount: model.GrossAmount,
		Status:      model.Status,
		VaNumber:    model.VaNumber,
		Bank:        model.Bank,
		BookingDate: model.BookingDate,
		PhoneNumber: model.PhoneNumber,
		Greeting:    model.Greeting,
		FullName:    model.FullName,
		Email:       model.Email,
		Quantity:    model.Quantity,
		ExpiredAt:   model.ExpiredAt,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func CoreReviewToModelReview(input booking.ReviewCore) Review {
	return Review{
		BookingID:  input.BookingID,
		UserID:     input.UserID,
		TextReview: input.TextReview,
		StartRate:  input.StartRate,
	}
}

func CoreToModel(reqNotif booking.Core) Booking {
	return Booking{
		Status: reqNotif.Status,
	}
}

func CoreToModelBookingCancel(input booking.Core) Booking {
	return Booking{
		Status: input.Status,
	}
}

func (b Booking) ModelToCore() booking.Core {
	return booking.Core{
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
		ExpiredAt:   b.ExpiredAt,
		CreatedAt:   b.CreatedAt,
		User: user.Core{
			ID:          b.User.ID,
			FullName:    b.User.FullName,
			NoKtp:       b.User.NoKtp,
			Address:     b.User.Address,
			PhoneNumber: b.User.PhoneNumber,
			Email:       b.User.Email,
			Image:       b.User.Image,
			Role:        b.User.Role,
			Status:      b.User.Status,
			CreatedAt:   b.User.CreatedAt,
			UpdatedAt:   b.User.UpdatedAt,
		},
		Tour: tour.Core{
			ID:       b.Tour.ID,
			TourName: b.Tour.TourName,
			Address:  b.Tour.Addres,
			Image:    b.Tour.Image,
		},
		Package: packages.Core{
			ID:          b.Package.ID,
			TourID:      b.Package.TourID,
			PackageName: b.Package.PackageName,
			Price:       b.Package.Price,
		},
		Voucher: voucher.Core{
			ID:   b.Voucher.ID,
			Name: b.Voucher.Name,
		},
	}
}

func ModelToCoreList(bookings []Booking) ([]booking.Core, error) {
	var bookingCores []booking.Core

	for _, b := range bookings {
		core := b.ModelToCore()
		core.Package = packages.Core{
			ID:          b.Package.ID,
			PackageName: b.Package.PackageName,
			Price:       b.Package.Price,
		}
		bookingCores = append(bookingCores, core)
	}

	return bookingCores, nil
}

func (b Booking) ModelToCoreBookingUser() booking.Core {
	return booking.Core{
		ID:          b.ID,
		UserID:      b.UserID,
		TourID:      b.TourID,
		PackageID:   b.PackageID,
		VoucherID:   b.VoucherID,
		GrossAmount: b.GrossAmount,
		Status:      b.Status,
		Tour:        b.Tour.ModelToCoreTourBooking(),
	}
}

func ModelToReviewCore(r Review) booking.ReviewCore {
	return booking.ReviewCore{
		UserID:     r.UserID,
		TextReview: r.TextReview,
		StartRate:  r.StartRate,
		CreatedAt:  r.CreatedAt,
		User: user.Core{
			FullName: r.User.FullName,
			Image:    r.User.Image,
		},
	}
}
