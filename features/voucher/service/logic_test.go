package service

import (
	"errors"
	"my-tourist-ticket/features/voucher"
	"my-tourist-ticket/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoucherCreate(t *testing.T) {
	mockRepo := new(mocks.VoucherData)
	service := New(mockRepo)

	input := voucher.Core{
		Name:           "Test Voucher",
		Code:           "TEST123",
		Description:    "test",
		DiscountValue:  10,
		ExpiredVoucher: "2024-12-31",
	}
	userIdLogin := 12

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetUserRoleById", userIdLogin).Return("admin", nil).Once()
		mockRepo.On("Insert", input).Return(nil).Once()

		err := service.Create(input, userIdLogin)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Access Denied", func(t *testing.T) {
		mockRepo.On("GetUserRoleById", userIdLogin).Return("costumer", nil).Once()

		err := service.Create(input, userIdLogin)

		assert.EqualError(t, err, "maaf anda tidak memiliki akses")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from GetUserRoleById", func(t *testing.T) {
		mockRepo.On("GetUserRoleById", userIdLogin).Return("", errors.New("error getting user role")).Once()

		err := service.Create(input, userIdLogin)

		assert.EqualError(t, err, "error getting user role")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation for Empty Fields", func(t *testing.T) {
		tests := []struct {
			name     string
			voucher  voucher.Core
			expected error
		}{
			{
				name:     "Empty Name",
				voucher:  voucher.Core{Code: "ketupat", Description: "senilai Rp.50.000", DiscountValue: 50000, ExpiredVoucher: "2024-03-21"},
				expected: errors.New("nama voucher tidak boleh kosong"),
			},
			{
				name:     "Empty Code",
				voucher:  voucher.Core{Name: "Voucher Lebaran", Description: "senilai Rp.50.000", DiscountValue: 50000, ExpiredVoucher: "2024-03-21"},
				expected: errors.New("code voucher tidak boleh kosong"),
			},
			{
				name:     "Empty Discount Value",
				voucher:  voucher.Core{Name: "Voucher Lebaran", Code: "ketupat", Description: "senilai Rp.50.000", ExpiredVoucher: "2024-03-21"},
				expected: errors.New("nominal voucher tidak boleh kosong"),
			},
			{
				name:     "Empty Expired Voucher",
				voucher:  voucher.Core{Name: "Voucher Lebaran", Code: "ketupat", Description: "senilai Rp.50.000", DiscountValue: 50000},
				expected: errors.New("tanggal expired voucher tidak boleh kosong"),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				mockRepo.On("GetUserRoleById", userIdLogin).Return("admin", nil).Once() // Panggilan ini harus ditambahkan di sini agar dipanggil sebelum memanggil service.Create
				err := service.Create(test.voucher, userIdLogin)
				assert.Equal(t, test.expected, err)
			})
		}
	})

	t.Run("Error from Insert", func(t *testing.T) {
		mockRepo.On("GetUserRoleById", userIdLogin).Return("admin", nil).Once() // Panggilan ini harus ditambahkan di sini agar dipanggil sebelum memanggil service.Create
		mockRepo.On("Insert", input).Return(errors.New("failed to insert voucher")).Once()

		err := service.Create(input, userIdLogin)
		assert.EqualError(t, err, "failed to insert voucher")
	})
}

func TestSelectAllVoucher(t *testing.T) {
	mockRepo := new(mocks.VoucherData)
	service := New(mockRepo)

	userIdLogin := 123

	t.Run("Success", func(t *testing.T) {
		expectedVouchers := []voucher.Core{
			{Name: "Voucher 1", Code: "VOUCHER1", DiscountValue: 20, ExpiredVoucher: "2024-12-31"},
			{Name: "Voucher 2", Code: "VOUCHER2", DiscountValue: 30, ExpiredVoucher: "2024-12-31"},
		}
		mockRepo.On("GetUserRoleById", userIdLogin).Return("admin", nil).Once()
		mockRepo.On("SelectAllVoucher", "admin").Return(expectedVouchers, nil).Once()

		vouchers, err := service.SelectAllVoucher(userIdLogin)

		assert.NoError(t, err)
		assert.Equal(t, expectedVouchers, vouchers)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from GetUserRoleById", func(t *testing.T) {
		mockRepo.On("GetUserRoleById", userIdLogin).Return("", errors.New("error getting user role")).Once()

		vouchers, err := service.SelectAllVoucher(userIdLogin)

		assert.EqualError(t, err, "error getting user role")
		assert.Nil(t, vouchers)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from SelectAllVoucher", func(t *testing.T) {
		mockRepo.On("GetUserRoleById", userIdLogin).Return("admin", nil).Once()
		mockRepo.On("SelectAllVoucher", "admin").Return(nil, errors.New("error selecting all vouchers")).Once()

		vouchers, err := service.SelectAllVoucher(userIdLogin)

		assert.EqualError(t, err, "error selecting all vouchers")
		assert.Nil(t, vouchers)
		mockRepo.AssertExpectations(t)
	})
}

func TestVoucherUpdate(t *testing.T) {
	mockRepo := new(mocks.VoucherData)
	service := New(mockRepo)

	voucherID := 5
	input := voucher.Core{
		Name:           "Updated Voucher",
		Code:           "UPDATED123",
		Description:    "updated test",
		DiscountValue:  15,
		ExpiredVoucher: "2024-12-31",
	}
	userIdLogin := 12

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetUserRoleById", userIdLogin).Return("admin", nil).Once()
		mockRepo.On("Update", voucherID, input).Return(nil).Once()

		err := service.Update(voucherID, input, userIdLogin)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Access Denied", func(t *testing.T) {
		// Mocking GetUserRoleById untuk mengembalikan peran "costumer"
		mockRepo.On("GetUserRoleById", userIdLogin).Return("costumer", nil).Once()

		// Pastikan bahwa metode Update tidak dipanggil dalam kasus ini
		mockRepo.AssertNotCalled(t, "Update")

		// Panggil fungsi Update
		err := service.Update(voucherID, input, userIdLogin)

		// Periksa bahwa pesan error yang diharapkan diterima
		assert.EqualError(t, err, "maaf anda tidak memiliki akses")
		// Periksa bahwa semua pemanggilan ekspektasi pada mockRepo terpenuhi
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from GetUserRoleById", func(t *testing.T) {
		mockRepo.On("GetUserRoleById", userIdLogin).Return("", errors.New("error getting user role")).Once()

		err := service.Update(voucherID, input, userIdLogin)

		assert.EqualError(t, err, "error getting user role")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from Update", func(t *testing.T) {
		mockRepo.On("GetUserRoleById", userIdLogin).Return("admin", nil).Once()
		mockRepo.On("Update", voucherID, input).Return(errors.New("error updating voucher")).Once()

		err := service.Update(voucherID, input, userIdLogin)

		assert.EqualError(t, err, "error updating voucher")
		mockRepo.AssertExpectations(t)
	})
}

func TestVoucherDelete(t *testing.T) {
	mockRepo := new(mocks.VoucherData)
	service := New(mockRepo)

	voucherID := 123

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Delete", voucherID).Return(nil).Once()

		err := service.Delete(voucherID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		err := service.Delete(0)

		assert.EqualError(t, err, "invalid id")
	})
}
