package user

import "time"

type Core struct {
	ID          uint
	FullName    string `validate:"required"`
	NoKtp       string
	Address     string
	PhoneNumber string `validate:"required"`
	Email       string `validate:"required,email"`
	Password    string `validate:"required"`
	Image       string
	Role        string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// interface untuk Data Layer
type UserDataInterface interface {
	Login(email, password string) (data *Core, err error)
	Insert(input Core) error
	SelectById(userIdLogin int) (*Core, error)
	Update(userIdLogin int, input Core) error
	Delete(userIdLogin int) error
	SelectAdminUsers(page, limit int) ([]Core, error, int)
	UpdatePengelola(pengelolaStatus string, pengelolaId int) error
}

// interface untuk Service Layer
type UserServiceInterface interface {
	Login(email, password string) (data *Core, token string, err error)
	Create(input Core) error
	GetById(userIdLogin int) (*Core, error)
	Update(userIdLogin int, input Core) error
	Delete(userIdLogin int) error
	GetAdminUsers(userIdLogin, page, limit int) ([]Core, error, int)
	UpdatePengelola(userIdLogin int, pengelolaStatus string, pengelolaId int) error
}
