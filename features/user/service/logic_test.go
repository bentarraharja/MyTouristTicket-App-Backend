package service

import (
	"errors"
	"my-tourist-ticket/features/user"
	"my-tourist-ticket/mocks"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	repo := new(mocks.UserData)
	hash := new(mocks.HashService)
	srv := New(repo, hash)

	inputData := user.Core{
		FullName:    "kitten",
		Email:       "kitten@gmail.com",
		Password:    "password",
		PhoneNumber: "3434144",
	}

	t.Run("Success Create Customer", func(t *testing.T) {
		hash.On("HashPassword", mock.AnythingOfType("string")).Return("hashedPassword", nil).Once()
		repo.On("Insert", mock.Anything).Return(nil).Once()
		inputData.Role = "customer"
		err := srv.Create(inputData)

		assert.NoError(t, err)
		hash.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Success Create Pengelola", func(t *testing.T) {
		hash.On("HashPassword", mock.AnythingOfType("string")).Return("hashedPassword", nil).Once()
		repo.On("Insert", mock.Anything).Return(nil).Once()
		inputData.Role = "pengelola"
		err := srv.Create(inputData)
		inputData.Status = "pending"

		assert.NoError(t, err)
		hash.AssertExpectations(t)
		repo.AssertExpectations(t)
		assert.Equal(t, "pending", inputData.Status)
	})

	t.Run("Validation Error", func(t *testing.T) {
		repo.On("Insert", mock.Anything).Return(nil).Once()
		invalidInput := user.Core{}
		err := srv.Create(invalidInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required")
		repo.AssertNotCalled(t, "Insert")
	})

	t.Run("Hash Password Error", func(t *testing.T) {
		hash.On("HashPassword", mock.Anything).Return("", errors.New("hash error")).Once()
		err := srv.Create(inputData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Error hash password.")
		repo.AssertNotCalled(t, "Insert")
	})
}

func TestLogin(t *testing.T) {
	repo := new(mocks.UserData)
	hash := new(mocks.HashService)
	userService := New(repo, hash)

	inputLogin := user.Core{
		Email:    "updated@gmail.com",
		Password: "newpassword",
	}

	t.Run("empty email and password", func(t *testing.T) {
		_, _, err := userService.Login("", "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email dan password wajib diisi")
	})

	t.Run("empty email", func(t *testing.T) {
		_, _, err := userService.Login("", "password")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email wajib diisi")
	})

	t.Run("empty password", func(t *testing.T) {
		_, _, err := userService.Login("email@gmail.com", "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password wajib diisi")
	})

	t.Run("password not match", func(t *testing.T) {
		repo.On("Login", inputLogin.Email, inputLogin.Password).Return(&inputLogin, nil).Once()
		hash.On("CheckPasswordHash", inputLogin.Password, inputLogin.Password).Return(false).Once()

		_, _, err := userService.Login(inputLogin.Email, inputLogin.Password)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password tidak sesuai.")
	})

	t.Run("error on userData.Login", func(t *testing.T) {
		repo.On("Login", inputLogin.Email, inputLogin.Password).Return(nil, errors.New("some error")).Once()

		_, _, err := userService.Login(inputLogin.Email, inputLogin.Password)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "some error")
	})

	t.Run("success", func(t *testing.T) {
		repo.On("Login", inputLogin.Email, inputLogin.Password).Return(&inputLogin, nil).Once()
		hash.On("CheckPasswordHash", inputLogin.Password, inputLogin.Password).Return(true).Once()

		_, _, err := userService.Login(inputLogin.Email, inputLogin.Password)

		assert.NoError(t, err)
	})
}
func TestGetById(t *testing.T) {
	repo := new(mocks.UserData)
	hash := new(mocks.HashService)
	userService := New(repo, hash)

	var userIdLogin = 1
	expectedUser := &user.Core{ID: uint(userIdLogin), Email: "test@example.com"}

	t.Run("Get user", func(t *testing.T) {
		repo.On("SelectById", userIdLogin).Return(expectedUser, nil).Once()

		user, err := userService.GetById(userIdLogin)

		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

}

func TestUpdateUser(t *testing.T) {
	repo := new(mocks.UserData)
	hash := new(mocks.HashService)
	userService := New(repo, hash)

	input := user.Core{
		ID:       1,
		FullName: "updatedUsername",
		Email:    "updated@gmail.com",
		Password: "newpassword",
	}

	t.Run("invalid user id", func(t *testing.T) {
		err := userService.Update(0, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid id")
	})

	t.Run("hash password error", func(t *testing.T) {
		hash.On("HashPassword", input.Password).Return("", errors.New("hash error")).Once()

		err := userService.Update(1, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Error hash password.")
	})

	t.Run("success", func(t *testing.T) {
		hashedPassword := "hashedPassword"
		hash.On("HashPassword", input.Password).Return(hashedPassword, nil).Once()
		repo.On("Update", 1, user.Core{ID: input.ID, FullName: input.FullName, Email: input.Email, Password: hashedPassword}).Return(nil).Once()

		err := userService.Update(1, input)

		assert.NoError(t, err)
	})
}

func TestDelete(t *testing.T) {
	repo := new(mocks.UserData)
	hash := new(mocks.HashService)
	userService := New(repo, hash)

	t.Run("invalid user id", func(t *testing.T) {
		err := userService.Delete(0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid id")
	})

	t.Run("success", func(t *testing.T) {
		repo.On("Delete", 1).Return(nil).Once()

		err := userService.Delete(1)

		assert.NoError(t, err)
	})
}

func TestGetAdminUsers(t *testing.T) {
	repo := new(mocks.UserData)
	hash := new(mocks.HashService)
	userService := New(repo, hash)

	adminUser := user.Core{
		ID:   1,
		Role: "admin",
	}

	customerUser := user.Core{
		ID:   2,
		Role: "customer",
	}

	pengelolaUser := user.Core{
		ID:   3,
		Role: "pengelola",
	}

	// Data hasil yang diharapkan dari pemanggilan SelectAdminUsers
	expectedResult := []user.Core{
		{ID: 1, FullName: "Kitten1"},
		{ID: 2, FullName: "Kitten2"},
	}

	// Jumlah halaman total yang diharapkan
	expectedTotalPage := 5

	// Tes jika pengguna adalah admin
	t.Run("Admin user", func(t *testing.T) {
		repo.On("SelectById", 1).Return(&adminUser, nil).Once()
		repo.On("SelectAdminUsers", 1, 10).Return(expectedResult, nil, expectedTotalPage).Once()

		result, err, totalPage := userService.GetAdminUsers(1, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	// Tes jika pengguna adalah customer
	t.Run("Customer user", func(t *testing.T) {
		repo.On("SelectById", 2).Return(&customerUser, nil).Once()

		result, err, _ := userService.GetAdminUsers(2, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "Sorry, your role does not have this access.", err.Error())

		repo.AssertExpectations(t)
	})

	// Tes jika pengguna adalah pengelola
	t.Run("Pengelola user", func(t *testing.T) {
		repo.On("SelectById", 3).Return(&pengelolaUser, nil).Once()

		result, err, _ := userService.GetAdminUsers(3, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "Sorry, your role does not have this access.", err.Error())

		repo.AssertExpectations(t)
	})

	// Tes jika pengguna tidak ditemukan
	t.Run("User not found", func(t *testing.T) {
		repo.On("SelectById", 4).Return(nil, errors.New("user not found")).Once()

		result, err, _ := userService.GetAdminUsers(4, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "user not found", err.Error())

		repo.AssertExpectations(t)
	})

	// Tes jika page dan limit bernilai 0
	t.Run("Page and Limit are 0", func(t *testing.T) {
		repo.On("SelectById", 1).Return(&adminUser, nil).Once()
		repo.On("SelectAdminUsers", 1, 10).Return(expectedResult, nil, expectedTotalPage).Once()

		result, err, totalPage := userService.GetAdminUsers(1, 0, 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})
}

func TestUpdatePengelola(t *testing.T) {
	// Buat mock untuk UserData
	repo := new(mocks.UserData)
	hash := new(mocks.HashService)
	userService := New(repo, hash)

	// Data pengguna yang valid dengan peran admin
	adminUser := user.Core{
		ID:   1,
		Role: "admin",
	}

	// Data pengguna dengan peran customer
	customerUser := user.Core{
		ID:   2,
		Role: "customer",
	}

	// Data pengguna dengan peran pengelola
	pengelolaUser := user.Core{
		ID:   3,
		Role: "pengelola",
	}

	expectedStatus := "accepted"
	expectedID := 1

	// Set mock behavior for SelectById to return an error
	expectedError := errors.New("some validation error")
	repo.On("SelectById", mock.AnythingOfType("int")).Return(nil, expectedError).Once()

	// Tes jika terjadi kesalahan saat memanggil SelectById
	t.Run("Error when calling SelectById", func(t *testing.T) {
		err := userService.UpdatePengelola(1, expectedStatus, expectedID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		repo.AssertExpectations(t)
	})

	// Tes jika pengguna adalah admin
	t.Run("Admin user", func(t *testing.T) {
		repo.On("SelectById", 1).Return(&adminUser, nil).Once()
		repo.On("UpdatePengelola", expectedStatus, expectedID).Return(nil).Once()

		err := userService.UpdatePengelola(1, expectedStatus, expectedID)

		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	// Tes jika pengguna adalah customer
	t.Run("Customer user", func(t *testing.T) {
		repo.On("SelectById", 2).Return(&customerUser, nil).Once()

		err := userService.UpdatePengelola(2, expectedStatus, expectedID)

		assert.Error(t, err)
		assert.Equal(t, "Sorry, your role does not have this access.", err.Error())

		repo.AssertExpectations(t)
	})

	// Tes jika pengguna adalah pengelola
	t.Run("Pengelola user", func(t *testing.T) {
		repo.On("SelectById", 3).Return(&pengelolaUser, nil).Once()

		err := userService.UpdatePengelola(3, expectedStatus, expectedID)

		assert.Error(t, err)
		assert.Equal(t, "Sorry, your role does not have this access.", err.Error())

		repo.AssertExpectations(t)
	})
}
