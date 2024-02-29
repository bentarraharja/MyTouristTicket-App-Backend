package data

import (
	bd "my-tourist-ticket/features/booking/data"
	"my-tourist-ticket/features/dashboard"
	td "my-tourist-ticket/features/tour/data"
	"my-tourist-ticket/features/user"
	ud "my-tourist-ticket/features/user/data"

	"gorm.io/gorm"
)

type dashboardQuery struct {
	db *gorm.DB
}

func NewDashboard(db *gorm.DB) dashboard.DashboardDataInterface {
	return &dashboardQuery{
		db: db,
	}
}

// GetUserRoleById implements dashboard.DashboardDataInterface.
func (repo *dashboardQuery) GetUserRoleById(userId int) (string, error) {
	var user user.Core
	if err := repo.db.Table("users").Where("id = ?", userId).First(&user).Error; err != nil {
		return "", err
	}

	return user.Role, nil
}

// GetTotalCustomer implements dashboard.DashboardDataInterface.
func (repo *dashboardQuery) GetTotalCustomer() (int, error) {
	var totalUser int64
	err := repo.db.Model(&ud.User{}).Where(&ud.User{Role: "costumer"}).Count(&totalUser).Error
	return int(totalUser), err
}

// GetTotalPengelola implements dashboard.DashboardDataInterface.
func (repo *dashboardQuery) GetTotalPengelola() (int, error) {
	var totalUser int64
	err := repo.db.Model(&ud.User{}).Where(&ud.User{Role: "pengelola"}).Count(&totalUser).Error
	return int(totalUser), err
}

// GetTotalTransaction implements dashboard.DashboardDataInterface.
func (repo *dashboardQuery) GetTotalTransaction() (int, error) {
	var totalBooking int64
	if err := repo.db.Model(&bd.Booking{}).Count(&totalBooking).Error; err != nil {
		return 0, err
	}
	return int(totalBooking), nil
}

// GetRecentTransaction implements dashboard.DashboardDataInterface.
func (repo *dashboardQuery) GetRecentTransaction() ([]dashboard.Booking, error) {
	var gormBookings []bd.Booking
	if err := repo.db.Preload("Tour").Preload("Package").Model(&bd.Booking{}).Order("created_at desc").Limit(5).Find(&gormBookings).Error; err != nil {
		return nil, err
	}

	var data []dashboard.Booking
	for _, gormBooking := range gormBookings {
		data = append(data, BookingModelToDashboard(gormBooking))
	}

	return data, nil
}

// GetTopTour implements dashboard.DashboardDataInterface.
func (repo *dashboardQuery) GetTopTour() ([]dashboard.Tour, error) {
	var topTours []td.Tour

	// Subquery untuk menghitung jumlah pemesanan untuk setiap tur
	subquery := repo.db.Model(&bd.Booking{}).
		Select("tour_id, COUNT(*) as booking_count").
		Group("tour_id").
		Order("booking_count DESC").
		Limit(5)

	// Lakukan penggabungan (join) dengan tabel td.Tour untuk mendapatkan data tur
	query := repo.db.Preload("City").Model(&td.Tour{}).
		Joins("JOIN (?) as bookings ON tours.id = bookings.tour_id", subquery).
		Order("bookings.booking_count DESC").
		Find(&topTours)

	if query.Error != nil {
		return nil, query.Error
	}

	var dashboardTopTours []dashboard.Tour
	for _, tour := range topTours {
		dashboardTopTours = append(dashboardTopTours, TourModelToDashboard(tour))
	}

	return dashboardTopTours, nil

}

// GetTotalTour implements dashboard.DashboardDataInterface.
func (repo *dashboardQuery) GetTotalTour() (int, error) {
	var totalTour int64
	if err := repo.db.Model(&td.Tour{}).Count(&totalTour).Error; err != nil {
		return 0, err
	}
	return int(totalTour), nil
}
