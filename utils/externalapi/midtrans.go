package externalapi

import (
	"errors"
	"my-tourist-ticket/app/configs"
	"my-tourist-ticket/features/booking"
	"time"

	mid "github.com/midtrans/midtrans-go"

	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransInterface interface {
	NewBookingPayment(book booking.Core) (*booking.Core, error)
	CancelBookingPayment(bookingId string) error
}

type midtrans struct {
	client      coreapi.Client
	environment mid.EnvironmentType
}

func New() MidtransInterface {
	environment := mid.Sandbox
	var client coreapi.Client
	client.New(configs.MID_KEY, environment)

	return &midtrans{
		client: client,
	}
}

// NewBookingPayment implements Midtrans.
func (pay *midtrans) NewBookingPayment(book booking.Core) (*booking.Core, error) {
	req := new(coreapi.ChargeReq)
	req.TransactionDetails = mid.TransactionDetails{
		OrderID:  book.ID,
		GrossAmt: int64(book.GrossAmount),
	}

	switch book.Bank {
	case "bca":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankBca,
		}
	case "bni":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankBni,
		}
	case "bri":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankBri,
		}
	default:
		return nil, errors.New("sorry, payment not support")
	}

	res, err := pay.client.ChargeTransaction(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != "201" {
		return nil, errors.New(res.StatusMessage)
	}

	if len(res.VaNumbers) == 1 {
		book.VaNumber = res.VaNumbers[0].VANumber
	}

	if res.PaymentType != "" {
		book.PaymentType = res.PaymentType
	}

	if res.TransactionStatus != "" {
		book.Status = res.TransactionStatus
	}

	if expiredAt, err := time.Parse("2006-01-02 15:04:05", res.ExpiryTime); err != nil {
		return nil, err
	} else {
		book.ExpiredAt = expiredAt
	}

	return &book, nil
}

func (pay *midtrans) CancelBookingPayment(bookingId string) error {
	res, _ := pay.client.CancelTransaction(bookingId)
	if res.StatusCode != "200" && res.StatusCode != "412" {
		return errors.New(res.StatusMessage)
	}

	return nil
}
