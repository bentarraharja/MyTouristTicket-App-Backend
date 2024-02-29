package data

import (
	"errors"
	"my-tourist-ticket/features/user"
	"my-tourist-ticket/features/voucher"
	"time"

	"gorm.io/gorm"
)

type voucherQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) voucher.VoucherDataInterface {
	return &voucherQuery{
		db: db,
	}
}

// Insert implements voucher.VoucherDataInterface.
func (repo *voucherQuery) Insert(input voucher.Core) error {
	dataGorm := CoreToModel(input)

	tx := repo.db.Create(&dataGorm)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("insert failed, row affected = 0")
	}

	return nil
}

// GetUserRoleById
func (repo *voucherQuery) GetUserRoleById(userIdLogin int) (string, error) {
	var user user.Core
	if err := repo.db.Table("users").Where("id = ?", userIdLogin).First(&user).Error; err != nil {
		return "", err
	}

	return user.Role, nil
}

// SelectAllVoucher implements voucher.VoucherDataInterface.
func (repo *voucherQuery) SelectAllVoucher(userRole string) ([]voucher.Core, error) {
	var voucherDataGorm []Voucher
	if userRole == "costumer" {
		currentDate := time.Now()
		tx := repo.db.Where("expired_voucher >= ?", currentDate).Find(&voucherDataGorm)
		if tx.Error != nil {
			return nil, tx.Error
		}
	} else {
		tx := repo.db.Find(&voucherDataGorm)
		if tx.Error != nil {
			return nil, tx.Error
		}
	}

	var voucherCores []voucher.Core
	for _, v := range voucherDataGorm {
		voucherCores = append(voucherCores, v.ModelToCore())
	}

	return voucherCores, nil
}

// Update implements voucher.VoucherDataInterface.
func (repo *voucherQuery) Update(voucherId int, input voucher.Core) error {
	dataGorm := CoreToModel(input)

	tx := repo.db.Model(&dataGorm).Where("id = ?", voucherId).Updates(dataGorm)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("update failed, no rows affected")
	}

	return nil
}

// Delete implements voucher.VoucherDataInterface.
func (repo *voucherQuery) Delete(voucherId int) error {
	tx := repo.db.Where("id = ?", voucherId).Delete(&Voucher{})
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found")
	}
	return nil
}
