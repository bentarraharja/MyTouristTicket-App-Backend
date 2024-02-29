package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"my-tourist-ticket/features/tour"
)

type tourService struct {
	tourData tour.TourDataInterface
}

func NewTour(repo tour.TourDataInterface) tour.TourServiceInterface {
	return &tourService{
		tourData: repo,
	}
}

// GetUserRoleById implements tour.TourServiceInterface.
func (service *tourService) GetUserRoleById(userId int) (string, error) {
	return service.tourData.GetUserRoleById(userId)
}

// Insert implements tour.TourServiceInterface.
func (service *tourService) Insert(userId uint, input tour.Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error {
	if input.TourName == "" {
		return errors.New("tour name is required")
	}
	if input.Description == "" {
		return errors.New("description is required")
	}
	if input.Address == "" {
		return errors.New("address is required")
	}
	if image == nil {
		return errors.New("image is required")
	}
	if thumbnail == nil {
		return errors.New("thumbnail is required")
	}

	if input.Latitude == 0 {
		return errors.New("latitude is required")
	}
	if input.Longitude == 0 {
		return errors.New("longitude is required")
	}
	if input.CityId == 0 {
		return errors.New("city id is required")
	}

	err := service.tourData.Insert(userId, input, image, thumbnail)
	if err != nil {
		return fmt.Errorf("error creating tour: %w", err)
	}

	return nil
}

// Update implements tour.TourServiceInterface.
func (service *tourService) Update(tourId int, input tour.Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error {
	return service.tourData.Update(tourId, input, image, thumbnail)

}

// SelectTourById implements tour.TourServiceInterface.
func (service *tourService) SelectTourById(tourId int) (tour.Core, error) {
	data, err := service.tourData.SelectTourById(tourId)
	if err != nil {
		return tour.Core{}, err
	}

	return data, nil
}

// Delete implements tour.TourServiceInterface.
func (service *tourService) Delete(tourId int) error {
	if tourId <= 0 {
		return errors.New("invalid id")
	}
	err := service.tourData.Delete(tourId)
	return err
}

// SelectAllTour implements tour.TourServiceInterface.
func (service *tourService) SelectAllTour(page int, limit int) ([]tour.Core, int, error) {

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 12
	}

	// Panggil SelectAllTour dari tourData
	tours, totalPage, err := service.tourData.SelectAllTour(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return tours, totalPage, nil
}

// SelectTourByPengelola implements tour.TourServiceInterface.
func (service *tourService) SelectTourByPengelola(userId int, page, limit int) ([]tour.Core, int, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 12
	}

	tours, totalPage, err := service.tourData.SelectTourByPengelola(userId, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return tours, totalPage, nil
}

// GetTourByCityID implements tour.TourServiceInterface.
func (service *tourService) GetTourByCityID(cityID uint, page, limit int) ([]tour.Core, int, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 6
	}

	tours, totalPage, err := service.tourData.GetTourByCityID(cityID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return tours, totalPage, nil
}

// InsertReportTour implements tour.TourServiceInterface.
func (service *tourService) InsertReportTour(userId int, tourId int, input tour.ReportCore) error {
	err := service.tourData.InsertReportTour(userId, tourId, input)
	if err != nil {
		return fmt.Errorf("error creating report: %w", err)
	}

	return nil
}

// SelectReportTour implements tour.TourServiceInterface.
func (service *tourService) SelectReportTour(tourId int) ([]tour.ReportCore, error) {
	reports, err := service.tourData.SelectReportTour(tourId)
	if err != nil {
		return nil, fmt.Errorf("error get report: %w", err)
	}

	return reports, nil
}

// SearchTour implements tour.TourServiceInterface.
func (service *tourService) SearchTour(query string) ([]tour.Core, error) {
	tours, err := service.tourData.SearchTour(query)
	if err != nil {
		return tours, err
	}
	return tours, nil
}
