package data

import (
	"my-tourist-ticket/features/user"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName    string
	NoKtp       string `gorm:"unique;default:null"`
	Address     string `gorm:"default:null"`
	PhoneNumber string `gorm:"unique"`
	Email       string `gorm:"unique"`
	Password    string
	Image       string
	Role        string `gorm:"not null"`
	Status      string `gorm:"not null"`
}

func CoreToModel(input user.Core) User {
	return User{
		FullName:    input.FullName,
		NoKtp:       input.NoKtp,
		Address:     input.Address,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
		Password:    input.Password,
		Image:       input.Image,
		Role:        input.Role,
		Status:      input.Status,
	}
}

func CoreToModelPengelola(pengelolaStatus string) User {
	return User{
		Status: pengelolaStatus,
	}
}

func (u User) ModelToCore() user.Core {
	return user.Core{
		ID:          u.ID,
		FullName:    u.FullName,
		NoKtp:       u.NoKtp,
		Address:     u.Address,
		PhoneNumber: u.PhoneNumber,
		Email:       u.Email,
		Password:    u.Password,
		Image:       u.Image,
		Role:        u.Role,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func (u User) ModelToCoreAdmin() user.Core {
	return user.Core{
		ID:          u.ID,
		FullName:    u.FullName,
		NoKtp:       u.NoKtp,
		Address:     u.Address,
		PhoneNumber: u.PhoneNumber,
		Email:       u.Email,
		Password:    u.Password,
		Image:       u.Image,
		Role:        u.Role,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
