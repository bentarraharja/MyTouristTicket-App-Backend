package data

import (
	"errors"
	"my-tourist-ticket/features/user"

	"gorm.io/gorm"
)

type userQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.UserDataInterface {
	return &userQuery{
		db: db,
	}
}

// Insert implements user.UserDataInterface.
func (repo *userQuery) Insert(input user.Core) error {
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

// Login implements user.UserDataInterface.
func (repo *userQuery) Login(email string, password string) (data *user.Core, err error) {
	var userGorm User
	tx := repo.db.Where("email = ?", email).First(&userGorm)
	if tx.Error != nil {
		return nil, tx.Error
	}
	result := userGorm.ModelToCore()
	return &result, nil
}

// SelectById implements user.UserDataInterface.
func (repo *userQuery) SelectById(userIdLogin int) (*user.Core, error) {
	var userDataGorm User
	tx := repo.db.First(&userDataGorm, userIdLogin)
	if tx.Error != nil {
		return nil, tx.Error
	}

	result := userDataGorm.ModelToCore()
	return &result, nil
}

// Update implements user.UserDataInterface.
func (repo *userQuery) Update(userIdLogin int, input user.Core) error {
	dataGorm := CoreToModel(input)
	tx := repo.db.Model(&User{}).Where("id = ?", userIdLogin).Updates(dataGorm)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found ")
	}
	return nil
}

// Delete implements user.UserDataInterface.
func (repo *userQuery) Delete(userIdLogin int) error {
	tx := repo.db.Delete(&User{}, userIdLogin)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found")
	}
	return nil
}

// SelectAdminUsers implements user.UserDataInterface.
func (repo *userQuery) SelectAdminUsers(page, limit int) ([]user.Core, error, int) {
	var usersDataGorm []User

	tx := repo.db.Unscoped().Where("role = 'pengelola' AND status = 'pending'").Limit(limit).Offset((page - 1) * limit).Find(&usersDataGorm)
	if tx.Error != nil {
		return nil, tx.Error, 0
	}

	var totalData int64
	tc := repo.db.Unscoped().Model(&usersDataGorm).Where("role = 'pengelola' AND status = 'pending'").Count(&totalData)
	if tc.Error != nil {
		return nil, tc.Error, 0
	}
	totalPage := int((totalData + int64(limit) - 1) / int64(limit))

	var usersDataCore []user.Core
	for _, value := range usersDataGorm {
		var usersCore = value.ModelToCoreAdmin()
		usersDataCore = append(usersDataCore, usersCore)
	}

	return usersDataCore, nil, int(totalPage)
}

// UpdatePengelola implements user.UserDataInterface.
func (repo *userQuery) UpdatePengelola(pengelolaStatus string, pengelolaId int) error {
	dataGorm := CoreToModelPengelola(pengelolaStatus)
	tx := repo.db.Model(&User{}).Where("id = ?", pengelolaId).Updates(dataGorm)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found")
	}
	return nil
}
