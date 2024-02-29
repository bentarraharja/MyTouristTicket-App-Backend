package data

import (
	"errors"
	packages "my-tourist-ticket/features/package"

	"gorm.io/gorm"
)

type packageQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) packages.PackageDataInterface {
	return &packageQuery{
		db: db,
	}
}

// Insert implements packages.PackageDataInterface.
func (repo *packageQuery) Insert(benefits []string, input packages.Core) error {
	dataGorm := CoreToModel(input)

	tx := repo.db.Create(&dataGorm)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("insert failed, row affected = 0")
	}
	lastInsertedID := dataGorm.ID

	for _, value := range benefits {
		benefitValue := Benefit{
			PackageID: lastInsertedID,
			Benefit:   value,
		}

		tb := repo.db.Create(&benefitValue)
		if tb.Error != nil {
			return tb.Error
		}
	}
	return nil
}

// SelectByTourId implements packages.PackageDataInterface.
func (repo *packageQuery) SelectByTourId(tourId uint) ([]packages.Core, error) {
	var packageDataGorms []Package
	tx := repo.db.Preload("Benefits").Where("tour_id = ?", tourId).Find(&packageDataGorms)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	var results []packages.Core
	for _, packageDataGorm := range packageDataGorms {
		result := packageDataGorm.ModelToCore()
		results = append(results, result)
	}
	return results, nil
}

// SelectAllBenefitsByPackageId implements packages.PackageDataInterface.
func (repo *packageQuery) SelectAllBenefitsByPackageId(packageId int) ([]packages.BenefitCore, error) {
	var benefitsDataGorm []Benefit
	tx := repo.db.Where("package_id = ?", packageId).Find(&benefitsDataGorm)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// proses mapping dari struct gorm model ke struct core
	var benefitsDataCore []packages.BenefitCore
	for _, value := range benefitsDataGorm {
		var benefitCore = packages.BenefitCore{
			ID: value.ID,
		}
		benefitsDataCore = append(benefitsDataCore, benefitCore)
	}

	return benefitsDataCore, nil
}

// Delete implements packages.PackageDataInterface.
func (repo *packageQuery) Delete(packageId int) error {
	tx := repo.db.Where("id = ?", packageId).Delete(&Package{})
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found")
	}
	return nil
}

// DeleteBenefits implements packages.PackageDataInterface.
func (repo *packageQuery) DeleteBenefits(packageId int) error {
	tx := repo.db.Delete(&Benefit{}, packageId)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found")
	}
	return nil
}
