package handler

import packages "my-tourist-ticket/features/package"

type PackageResponse struct {
	ID          uint              `json:"id"`
	TourID      uint              `json:"tour_id"`
	PackageName string            `json:"package_name"`
	Price       int               `json:"price"`
	JumlahTiket int               `json:"jumlah_tiket"`
	Benefits    []BenefitResponse `json:"benefits"`
}

type BenefitResponse struct {
	ID        uint   `json:"id"`
	PackageID uint   `json:"package_id"`
	Benefit   string `json:"benefit"`
}

type PackageResponseName struct {
	PackageName string `json:"package_name"`
}

func CoresToResponses(data []packages.Core) []PackageResponse {
	var responses []PackageResponse
	for _, core := range data {
		var benefitsResponse []BenefitResponse
		for _, benefit := range core.Benefits {
			benefitsResponse = append(benefitsResponse, BenefitResponse{
				ID:        benefit.ID,
				PackageID: benefit.PackageID,
				Benefit:   benefit.Benefit,
			})
		}
		responses = append(responses, PackageResponse{
			ID:          core.ID,
			TourID:      core.TourID,
			PackageName: core.PackageName,
			Price:       core.Price,
			JumlahTiket: core.JumlahTiket,
			Benefits:    benefitsResponse,
		})
	}
	return responses
}
