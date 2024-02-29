package handler

import "my-tourist-ticket/features/voucher"

type VoucherResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Description    string `json:"description"`
	DiscountValue  int    `json:"discount_value"`
	ExpiredVoucher string `json:"expired_voucher"`
}

func CoreToGetAllResponse(data voucher.Core) VoucherResponse {
	return VoucherResponse{
		ID:             data.ID,
		Name:           data.Name,
		Code:           data.Code,
		Description:    data.Description,
		DiscountValue:  data.DiscountValue,
		ExpiredVoucher: data.ExpiredVoucher,
	}
}

func CoreToResponseListGetAllVoucher(data []voucher.Core) []VoucherResponse {
	var results []VoucherResponse
	for _, v := range data {
		results = append(results, CoreToGetAllResponse(v))
	}
	return results
}
