package tour

import (
	"mime/multipart"
	"my-tourist-ticket/features/city"
	packages "my-tourist-ticket/features/package"
	"my-tourist-ticket/features/user"
	"time"
)

type Core struct {
	ID          uint
	CityId      uint
	UserId      uint
	TourName    string
	Description string
	Image       string
	Thumbnail   string
	Address     string
	Latitude    float64
	Longitude   float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	City        city.Core
	Package     packages.Core
	ReportCount int64
}

type ReportCore struct {
	ID         uint
	UserId     uint
	TourId     uint
	TextReport string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       user.Core
}

type TourDataInterface interface {
	GetUserRoleById(userId int) (string, error)
	Insert(userId uint, input Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error
	Update(tourId int, input Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error
	SelectTourById(tourId int) (Core, error)
	Delete(tourId int) error
	SelectAllTour(page, limit int) ([]Core, int, error)
	SelectTourByPengelola(userId int, page, limit int) ([]Core, int, error)
	GetTourByCityID(cityID uint, page, limit int) ([]Core, int, error)
	InsertReportTour(userId int, tourId int, input ReportCore) error
	SelectReportTour(tourId int) ([]ReportCore, error)
	SearchTour(query string) ([]Core, error)
}

type TourServiceInterface interface {
	GetUserRoleById(userId int) (string, error)
	Insert(userId uint, input Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error
	Update(tourId int, input Core, image *multipart.FileHeader, thumbnail *multipart.FileHeader) error
	SelectTourById(tourId int) (Core, error)
	Delete(tourId int) error
	SelectAllTour(page, limit int) ([]Core, int, error)
	SelectTourByPengelola(userId int, page, limit int) ([]Core, int, error)
	GetTourByCityID(cityID uint, page, limit int) ([]Core, int, error)
	InsertReportTour(userId int, tourId int, input ReportCore) error
	SelectReportTour(tourId int) ([]ReportCore, error)
	SearchTour(query string) ([]Core, error)
}
