package voucher

import "time"

type Core struct {
	ID             uint
	Name           string
	Code           string
	Description    string
	DiscountValue  int
	ExpiredVoucher string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// interface untuk Data Layer
type VoucherDataInterface interface {
	Insert(input Core) error
	GetUserRoleById(userId int) (string, error)
	SelectAllVoucher(userRole string) ([]Core, error)
	Update(voucherId int, input Core) error
	Delete(voucherId int) error
}

// interface untuk Service Layer
type VoucherServiceInterface interface {
	Create(input Core, userIdLogin int) error
	SelectAllVoucher(userIdLogin int) ([]Core, error)
	Update(voucherId int, input Core, userIdLogin int) error
	Delete(voucherId int) error
}
