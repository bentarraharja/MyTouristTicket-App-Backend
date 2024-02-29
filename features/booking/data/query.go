package data

import (
	"errors"
	"fmt"
	"math"
	"my-tourist-ticket/features/booking"
	pd "my-tourist-ticket/features/package/data"
	"my-tourist-ticket/features/user"
	vd "my-tourist-ticket/features/voucher/data"
	"my-tourist-ticket/utils/externalapi"

	"gorm.io/gorm"
)

type bookingQuery struct {
	db              *gorm.DB
	paymentMidtrans externalapi.MidtransInterface
}

func New(db *gorm.DB, mi externalapi.MidtransInterface) booking.BookingDataInterface {
	return &bookingQuery{
		db:              db,
		paymentMidtrans: mi,
	}
}

// GetUserRoleById implements booking.BookingDataInterface.
func (repo *bookingQuery) GetUserRoleById(userId int) (string, error) {
	var user user.Core
	if err := repo.db.Table("users").Where("id = ?", userId).First(&user).Error; err != nil {
		return "", err
	}

	return user.Role, nil
}

// InsertBooking implements booking.BookingDataInterface.
func (repo *bookingQuery) InsertBooking(userIdLogin int, inputBooking booking.Core) (*booking.Core, error) {

	var totalHargaKeseluruhan int
	var packageGorm pd.Package
	ts := repo.db.Where("tour_id = ? AND id = ?", inputBooking.TourID, inputBooking.PackageID).First(&packageGorm)
	if ts.Error != nil {
		return nil, ts.Error
	}

	if inputBooking.VoucherID != nil {
		var voucherGorm vd.Voucher
		ts := repo.db.Where("id = ?", inputBooking.VoucherID).First(&voucherGorm)
		if ts.Error != nil {
			return nil, ts.Error
		}

		var existingUseVoucher Booking
		if err := repo.db.Where("user_id = ? AND voucher_id = ?", userIdLogin, inputBooking.VoucherID).First(&existingUseVoucher).Error; err == nil {
			return nil, errors.New("user has already used this voucher")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		totalHargaAwal := packageGorm.Price * inputBooking.Quantity
		if totalHargaAwal < voucherGorm.DiscountValue {
			return nil, errors.New("maaf, anda tidak bisa menggunakan voucher ini karena total pembayaran anda terlalu rendah")
		} else {
			// totalHargaKeseluruhan = ((packageGorm.JumlahTiket * packageGorm.Price) * inputBooking.Quantity) - voucherGorm.DiscountValue
			totalHargaKeseluruhan = totalHargaAwal - voucherGorm.DiscountValue
		}
	} else {
		// totalHargaKeseluruhan = (packageGorm.JumlahTiket * packageGorm.Price) * inputBooking.Quantity
		totalHargaKeseluruhan = packageGorm.Price * inputBooking.Quantity
	}

	inputBooking.GrossAmount = totalHargaKeseluruhan

	payment, errPay := repo.paymentMidtrans.NewBookingPayment(inputBooking)
	if errPay != nil {
		return nil, errPay
	}

	bookingPaymentModel := CoreToModelBooking(inputBooking)
	bookingPaymentModel.PaymentType = payment.PaymentType
	bookingPaymentModel.Status = payment.Status
	bookingPaymentModel.VaNumber = payment.VaNumber
	bookingPaymentModel.ExpiredAt = payment.ExpiredAt

	tx := repo.db.Create(&bookingPaymentModel)
	if tx.Error != nil {
		return nil, tx.Error
	}

	bookingCore := ModelToCoreBooking(bookingPaymentModel)

	return &bookingCore, nil
}

func (repo *bookingQuery) CancelBooking(userIdLogin int, bookingId string, bookingCore booking.Core) error {
	if bookingCore.Status == "cancelled" {
		repo.paymentMidtrans.CancelBookingPayment(bookingId)
	}

	dataGorm := CoreToModelBookingCancel(bookingCore)
	tx := repo.db.Model(&Booking{}).Where("id = ? AND user_id = ?", bookingId, userIdLogin).Updates(dataGorm)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found")
	}
	return nil
}

// InsertBookingReview implements booking.BookingDataInterface.
func (repo *bookingQuery) InsertBookingReview(inputReview booking.ReviewCore) error {
	dataGorm := CoreReviewToModelReview(inputReview)

	tx := repo.db.Create(&dataGorm)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("insert failed, row affected = 0")

	}
	return nil
}

// Update implements booking.BookingDataInterface.
func (repo *bookingQuery) WebhoocksData(reqNotif booking.Core) error {
	dataGorm := CoreToModel(reqNotif)
	tx := repo.db.Model(&Booking{}).Where("id = ?", reqNotif.ID).Updates(dataGorm)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found ")
	}
	return nil
}

func (repo *bookingQuery) SelectBookingUser(userIdLogin int) ([]booking.Core, error) {
	var bookingDataGorms []Booking
	tx := repo.db.Preload("Tour").Where("user_id = ?", userIdLogin).Order("created_at DESC").Find(&bookingDataGorms)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var results []booking.Core
	for _, bookingDataGorm := range bookingDataGorms {
		result := bookingDataGorm.ModelToCoreBookingUser()
		results = append(results, result)
	}
	return results, nil
}

func (repo *bookingQuery) SelectBookingUserDetail(userIdLogin int, bookingId string) (*booking.Core, error) {
	var bookingDataGorm Booking
	tx := repo.db.Preload("Tour").Preload("Package").Preload("Voucher").Where("user_id = ? AND id = ?", userIdLogin, bookingId).Find(&bookingDataGorm)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, errors.New("booking id not found")
	}

	result := bookingDataGorm.ModelToCore()

	return &result, nil
}

// SelectAllBooking implements booking.BookingDataInterface.
func (repo *bookingQuery) SelectAllBooking(page int, limit int) ([]booking.Core, int, error) {
	var bookingGorm []Booking
	query := repo.db.Order("created_at desc")

	var totalData int64
	err := query.Model(&Booking{}).Count(&totalData).Error
	if err != nil {
		return nil, 0, err
	}

	totalPage := int((totalData + int64(limit) - 1) / int64(limit))

	// Retrieve booking data with associated user, tour, and package
	err = query.Limit(limit).Offset((page - 1) * limit).Preload("Package").Preload("Tour").Find(&bookingGorm).Error
	if err != nil {
		return nil, 0, err
	}

	// Convert booking data to booking.Core
	bookingCore, err := ModelToCoreList(bookingGorm)
	if err != nil {
		return nil, 0, err
	}

	// for i := range bookingCore {
	// 	bookingCore[i].Package.ID = bookingGorm[i].PackageID
	// }

	return bookingCore, totalPage, nil
}

// SelectAllBookingPengelola implements booking.BookingDataInterface.
func (repo *bookingQuery) SelectAllBookingPengelola(pengelolaID int, page int, limit int) ([]booking.Core, int, error) {
	var bookingGorm []Booking
	query := repo.db.Order("created_at desc")

	var totalData int64
	err := query.Model(&Booking{}).Where("tour_id IN (SELECT id FROM tours WHERE user_id = ?)", pengelolaID).Count(&totalData).Error
	if err != nil {
		return nil, 0, err
	}

	totalPage := int((totalData + int64(limit) - 1) / int64(limit))

	// Retrieve booking data with associated user, tour, and package
	err = query.Limit(limit).Offset((page - 1) * limit).Preload("Package").Preload("Tour").Find(&bookingGorm).Error
	if err != nil {
		return nil, 0, err
	}

	// Convert booking data to booking.Core
	bookingCore, err := ModelToCoreList(bookingGorm)
	if err != nil {
		return nil, 0, err
	}

	return bookingCore, totalPage, nil
}

// GetAllBookingReview implements booking.BookingDataInterface.
func (repo *bookingQuery) GetAllBookingReview(tourId, limit int) ([]booking.ReviewCore, error) {
	// Check if tour ID exists
	var count int64
	repo.db.Model(&Booking{}).Where("tour_id = ?", tourId).Count(&count)
	if count == 0 {
		return nil, fmt.Errorf("tour ID %d does not exist", tourId)
	}

	var reviews []Review
	query := repo.db.
		Preload("User").
		Joins("JOIN bookings ON reviews.booking_id = bookings.id").
		Where("bookings.tour_id = ?", tourId).
		Order("reviews.created_at desc").
		Limit(limit).
		Find(&reviews)

	if query.Error != nil {
		return nil, query.Error
	}

	var reviewCores []booking.ReviewCore
	for _, r := range reviews {
		reviewCores = append(reviewCores, ModelToReviewCore(r))
	}

	return reviewCores, nil
}

// GetAverageTourReview implements booking.BookingDataInterface.
func (repo *bookingQuery) GetAverageTourReview(tourId int) (float64, error) {
	var averageReview float64
	var totalReviews int64
	err := repo.db.Model(&Review{}).
		Joins("JOIN bookings ON reviews.booking_id = bookings.id").
		Where("bookings.tour_id = ?", tourId).
		Count(&totalReviews).Error
	if err != nil {
		return 0, err
	}

	if totalReviews == 0 {
		return 0, nil
	}

	var sumRating float64
	err = repo.db.Model(&Review{}).
		Select("COALESCE(SUM(start_rate), 0)").
		Joins("JOIN bookings ON reviews.booking_id = bookings.id").
		Where("bookings.tour_id = ?", tourId).
		Scan(&sumRating).Error
	if err != nil {
		return 0, err
	}

	averageReview = float64(sumRating) / float64(totalReviews)
	averageReviews := math.Round(averageReview*10) / 10
	return averageReviews, nil
}

// GetTotalTourReview implements booking.BookingDataInterface.
func (repo *bookingQuery) GetTotalTourReview(tourId int) (int, error) {
	var totalReviews int64
	err := repo.db.Model(&Review{}).
		Joins("JOIN bookings ON reviews.booking_id = bookings.id").
		Where("bookings.tour_id = ?", tourId).
		Count(&totalReviews).Error
	if err != nil {
		return 0, err
	}

	return int(totalReviews), nil
}
