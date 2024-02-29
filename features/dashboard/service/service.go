package service

import (
	"errors"
	"log"
	"my-tourist-ticket/features/dashboard"
)

type dashboardService struct {
	dashboardData dashboard.DashboardDataInterface
}

func New(repo dashboard.DashboardDataInterface) dashboard.DashboardServiceInterface {
	return &dashboardService{
		dashboardData: repo,
	}
}

// Dashboard implements dashboard.DashboardServiceInterface.
func (service *dashboardService) Dashboard(userId int) (*dashboard.Dashboard, error) {
	userRole, errSelect := service.dashboardData.GetUserRoleById(userId)
	if errSelect != nil {
		return nil, errSelect
	}
	if userRole != "admin" {
		return nil, errors.New("maaf anda tidak memiliki akses")
	}

	dash := &dashboard.Dashboard{}

	// Get total customers
	totalCustomers, err := service.dashboardData.GetTotalCustomer()
	if err != nil {
		log.Printf("Error fetching total customers: %v", err)
		return nil, err
	}
	dash.TotalCustomer = totalCustomers

	// Get total pengelola
	totalPengelola, err := service.dashboardData.GetTotalPengelola()
	if err != nil {
		log.Printf("Error fetching total pengelola: %v", err)
		return nil, err
	}
	dash.TotalPengelola = totalPengelola

	// Get total transactions
	totalTransactions, err := service.dashboardData.GetTotalTransaction()
	if err != nil {
		log.Printf("Error fetching total transactions: %v", err)
		return nil, err
	}
	dash.TotalTransaction = totalTransactions

	// Get recent transactions
	recentTransactions, err := service.dashboardData.GetRecentTransaction()
	if err != nil {
		log.Printf("Error fetching recent transactions: %v", err)
		return nil, err
	}
	dash.RecentBooking = recentTransactions

	// Get top tours
	topTours, err := service.dashboardData.GetTopTour()
	if err != nil {
		log.Printf("Error fetching top tours: %v", err)
		return nil, err
	}
	dash.TopTours = topTours

	// Get total tours
	totalTours, err := service.dashboardData.GetTotalTour()
	if err != nil {
		log.Printf("Error fetching total tours: %v", err)
		return nil, err
	}
	dash.TotalTour = totalTours

	return dash, nil
}
