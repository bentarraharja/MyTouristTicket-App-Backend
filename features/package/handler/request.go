package handler

import packages "my-tourist-ticket/features/package"

type PackageRequest struct {
	TourID      uint
	PackageName string   `json:"package_name" form:"package_name"`
	Price       int      `json:"price" form:"price"`
	JumlahTiket int      `json:"jumlah_tiket" form:"jumlah_tiket"`
	Benefits    []string `json:"benefits" form:"benefits"`
}

func RequestToCore(input PackageRequest, tourID uint) packages.Core {
	return packages.Core{
		TourID:      tourID,
		PackageName: input.PackageName,
		Price:       input.Price,
		JumlahTiket: input.JumlahTiket,
	}
}
