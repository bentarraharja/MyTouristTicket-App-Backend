package service

import (
	"errors"
	packages "my-tourist-ticket/features/package"
)

type packageService struct {
	packageData packages.PackageDataInterface
}

func New(repo packages.PackageDataInterface) packages.PackageServiceInterface {
	return &packageService{
		packageData: repo,
	}
}

// Create implements packages.PackageServiceInterface.
func (service *packageService) Create(benefits []string, input packages.Core) error {
	if input.JumlahTiket == 0 {
		input.JumlahTiket = 1
	}

	err := service.packageData.Insert(benefits, input)
	if err != nil {
		return err
	}

	return nil
}

// GetByTourId implements packages.PackageServiceInterface.
func (service *packageService) GetByTourId(tourId uint) ([]packages.Core, error) {
	packages, err := service.packageData.SelectByTourId(tourId)
	if err != nil {
		return nil, err
	}
	return packages, nil
}

// Delete implements packages.PackageServiceInterface.
func (service *packageService) Delete(packageId int) error {
	if packageId <= 0 {
		return errors.New("invalid id")
	}

	benefits, errGet := service.packageData.SelectAllBenefitsByPackageId(packageId)
	if errGet != nil {
		return errGet
	}

	// Delete each task associated with the project
	for _, benefit := range benefits {
		errDel := service.packageData.DeleteBenefits(int(benefit.ID))
		if errDel != nil {
			return errDel
		}
	}

	err := service.packageData.Delete(packageId)
	return err
}
