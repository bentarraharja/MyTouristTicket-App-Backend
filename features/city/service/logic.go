package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"my-tourist-ticket/features/city"
)

type cityService struct {
	cityData city.CityDataInterface
}

func NewCity(repo city.CityDataInterface) city.CityServiceInterface {
	return &cityService{
		cityData: repo,
	}
}

// GetUserRoleById implements admin.AdminServiceInterface.
func (service *cityService) GetUserRoleById(userId int) (string, error) {
	return service.cityData.GetUserRoleById(userId)
}

// Insert implements city.CityDataInterface.
func (service *cityService) Create(input city.Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error {

	err := service.cityData.Insert(input, image, thumbnail)
	if err != nil {
		return fmt.Errorf("error creating city: %w", err)
	}

	return nil
}

// Update implements city.CityServiceInterface.
func (service *cityService) Update(cityId int, input city.Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error {
	err := service.cityData.Update(cityId, input, image, thumbnail)
	if err != nil {
		return fmt.Errorf("error update city: %w", err)
	}

	return nil
}

// Delete implements city.CityDataInterface.
func (service *cityService) Delete(cityId int) error {
	if cityId <= 0 {
		return errors.New("invalid id")
	}
	err := service.cityData.Delete(cityId)
	return err
}

// SelectCityById implements city.CityServiceInterface.
func (service *cityService) SelectCityById(cityId int) (city.Core, error) {
	data, err := service.cityData.SelectCityById(cityId)
	if err != nil {
		return city.Core{}, err
	}

	return data, nil
}

// SelectAllCity implements city.CityServiceInterface.
func (service *cityService) SelectAllCity(page int, limit int) ([]city.Core, int, error) {
	if page == 0 {
		page = 1
	}

	citys, totalPage, err := service.cityData.SelectAllCity(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return citys, totalPage, nil
}
