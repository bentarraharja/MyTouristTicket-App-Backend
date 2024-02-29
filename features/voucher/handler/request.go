package handler

import (
	"my-tourist-ticket/features/voucher"
)

type VoucherRequest struct {
	Name           string `json:"name" form:"name"`
	Code           string `json:"code" form:"code"`
	Description    string `json:"description" form:"description"`
	DiscountValue  int    `json:"discount_value" form:"discount_value"`
	ExpiredVoucher string `json:"expired_voucher" form:"expired_voucher"`
}

func RequestToCore(input VoucherRequest) voucher.Core {
	return voucher.Core{
		Name:           input.Name,
		Code:           input.Code,
		Description:    input.Description,
		DiscountValue:  input.DiscountValue,
		ExpiredVoucher: input.ExpiredVoucher,
	}
}
