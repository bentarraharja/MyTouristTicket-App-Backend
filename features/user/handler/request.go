package handler

import "my-tourist-ticket/features/user"

type UserRequest struct {
	FullName    string `json:"full_name" form:"full_name"`
	NoKtp       string `json:"no_ktp" form:"no_ktp"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Email       string `json:"email" form:"email"`
	Password    string `json:"password" form:"password"`
	Image       string `json:"image" form:"image"`
	Role        string `json:"role" form:"role"`
	Status      string `json:"status" form:"status"`
}

type UserRequestUpdate struct {
	FullName    string `json:"full_name" form:"full_name"`
	NoKtp       string `json:"no_ktp" form:"no_ktp"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Email       string `json:"email" form:"email"`
	Password    string `json:"password" form:"password"`
	Image       string `json:"image" form:"image"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// type AdminPengelolaRequestUpdate struct {
// 	Status string `json:"status" form:"status"`
// }

func RequestToCore(input UserRequest) user.Core {
	return user.Core{
		FullName:    input.FullName,
		NoKtp:       input.NoKtp,
		Address:     input.Address,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
		Password:    input.Password,
		Image:       input.Image,
		Role:        input.Role,
		Status:      input.Status,
	}
}

func UpdateRequestToCore(input UserRequestUpdate, imageURL string) user.Core {
	return user.Core{
		FullName:    input.FullName,
		NoKtp:       input.NoKtp,
		Address:     input.Address,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
		Password:    input.Password,
		Image:       imageURL,
	}
}

// func RequestToCoreAdminPengelola(input AdminPengelolaRequestUpdate) user.Core {
// 	return user.Core{
// 		Status: input.Status,
// 	}
// }
