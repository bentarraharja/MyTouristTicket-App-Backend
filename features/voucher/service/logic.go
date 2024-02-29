package service

import (
	"errors"
	"my-tourist-ticket/features/voucher"
)

type voucherService struct {
	voucherData voucher.VoucherDataInterface
}

func New(repo voucher.VoucherDataInterface) voucher.VoucherServiceInterface {
	return &voucherService{
		voucherData: repo,
	}
}

// Create implements voucher.VoucherServiceInterface.
func (service *voucherService) Create(input voucher.Core, userIdLogin int) error {
	userRole, errSelect := service.voucherData.GetUserRoleById(userIdLogin)
	if errSelect != nil {
		return errSelect
	}

	if userRole == "costumer" || userRole == "pengelola" {
		return errors.New("maaf anda tidak memiliki akses")
	}

	if input.Name == "" {
		return errors.New("nama voucher tidak boleh kosong")
	}
	if input.Code == "" {
		return errors.New("code voucher tidak boleh kosong")
	}
	if input.DiscountValue == 0 {
		return errors.New("nominal voucher tidak boleh kosong")
	}
	if input.ExpiredVoucher == "" {
		return errors.New("tanggal expired voucher tidak boleh kosong")
	}

	err := service.voucherData.Insert(input)
	if err != nil {
		return err
	}

	return nil
}

// SelectAllVoucher implements voucher.VoucherServiceInterface.
func (service *voucherService) SelectAllVoucher(userIdLogin int) ([]voucher.Core, error) {
	userRole, errSelect := service.voucherData.GetUserRoleById(userIdLogin)
	if errSelect != nil {
		return nil, errSelect
	}

	vouchers, err := service.voucherData.SelectAllVoucher(userRole)
	if err != nil {
		return nil, err
	}

	return vouchers, nil
}

// Update implements voucher.VoucherServiceInterface.
func (service *voucherService) Update(voucherId int, input voucher.Core, userIdLogin int) error {
	userRole, errSelect := service.voucherData.GetUserRoleById(userIdLogin)
	if errSelect != nil {
		return errSelect
	}

	if userRole == "costumer" || userRole == "pengelola" {
		return errors.New("maaf anda tidak memiliki akses")
	}

	err := service.voucherData.Update(voucherId, input)
	if err != nil {
		return err
	}

	return nil
}

// Delete implements voucher.VoucherServiceInterface.
func (service *voucherService) Delete(voucherId int) error {
	if voucherId <= 0 {
		return errors.New("invalid id")
	}

	err := service.voucherData.Delete(voucherId)
	return err
}
