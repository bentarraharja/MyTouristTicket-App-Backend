package city

import (
	"mime/multipart"
	"time"
)

type Core struct {
	ID          uint
	CityName    string `validate:"required"`
	Description string `validate:"required"`
	Image       string
	Thumbnail   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// interface untuk Data Layer
type CityDataInterface interface {
	GetUserRoleById(userId int) (string, error)
	Insert(input Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error
	Update(cityId int, input Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error
	Delete(cityId int) error
	SelectCityById(cityId int) (Core, error)
	SelectAllCity(page, limit int) ([]Core, int, error)
}

// interface untuk Service Layer
type CityServiceInterface interface {
	GetUserRoleById(userId int) (string, error)
	Create(input Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error
	Update(cityId int, input Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error
	Delete(cityId int) error
	SelectCityById(cityId int) (Core, error)
	SelectAllCity(page, limit int) ([]Core, int, error)
}
