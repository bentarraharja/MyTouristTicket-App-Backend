package handler

import (
	"my-tourist-ticket/app/middlewares"
	"my-tourist-ticket/features/voucher"
	"my-tourist-ticket/utils/responses"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type VoucherHandler struct {
	voucherService voucher.VoucherServiceInterface
}

func New(vs voucher.VoucherServiceInterface) *VoucherHandler {
	return &VoucherHandler{
		voucherService: vs,
	}
}

func (handler *VoucherHandler) CreateVoucher(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	newVoucher := VoucherRequest{}
	errBind := c.Bind(&newVoucher)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data, data not valid", nil))
	}

	voucherCore := RequestToCore(newVoucher)
	errInsert := handler.voucherService.Create(voucherCore, userIdLogin)
	if errInsert != nil {
		if strings.Contains(errInsert.Error(), "Error 1062 (23000): Duplicate entry") {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
		} else if strings.Contains(errInsert.Error(), "maaf anda tidak memiliki akses") {
			return c.JSON(http.StatusForbidden, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
		} else if strings.Contains(errInsert.Error(), "nama voucher tidak boleh kosong") {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
		} else if strings.Contains(errInsert.Error(), "code voucher tidak boleh kosong") {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
		} else if strings.Contains(errInsert.Error(), "nominal voucher tidak boleh kosong") {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
		} else if strings.Contains(errInsert.Error(), "tanggal expired voucher tidak boleh kosong") {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
		} else {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
		}
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success insert data", nil))
}

func (handler *VoucherHandler) GetAllVoucher(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	vouchers, err := handler.voucherService.SelectAllVoucher(userIdLogin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error get data", nil))
	}

	vouchersResponses := CoreToResponseListGetAllVoucher(vouchers)

	return c.JSON(http.StatusOK, responses.WebResponse("success get data", vouchersResponses))
}

func (handler *VoucherHandler) UpdateVoucher(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	newVoucher := VoucherRequest{}
	errBind := c.Bind(&newVoucher)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data, data not valid", nil))
	}

	vocId, err := strconv.Atoi(c.Param("voucher_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Invalid voucher ID", nil))
	}

	voucherCore := RequestToCore(newVoucher)
	errUpdate := handler.voucherService.Update(vocId, voucherCore, userIdLogin)
	if errUpdate != nil {
		if strings.Contains(errUpdate.Error(), "Error 1062 (23000): Duplicate entry") {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("error update data. "+errUpdate.Error(), nil))
		} else if strings.Contains(errUpdate.Error(), "maaf anda tidak memiliki akses") {
			return c.JSON(http.StatusForbidden, responses.WebResponse("error update data. "+errUpdate.Error(), nil))
		} else if strings.Contains(errUpdate.Error(), "update failed, no rows affected") {
			return c.JSON(http.StatusNotFound, responses.WebResponse("error update data. "+errUpdate.Error(), nil))
		} else {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error update data. "+errUpdate.Error(), nil))
		}
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success update data", nil))
}

func (handler *VoucherHandler) DeleteVoucher(c echo.Context) error {
	voucherId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error. id should be number", nil))
	}

	errDelete := handler.voucherService.Delete(voucherId)
	if errDelete != nil {
		if strings.Contains(errDelete.Error(), "error record not found") {
			return c.JSON(http.StatusNotFound, responses.WebResponse("error delete data "+errDelete.Error(), nil))
		} else {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error delete data "+errDelete.Error(), nil))
		}
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success delete data", nil))
}
