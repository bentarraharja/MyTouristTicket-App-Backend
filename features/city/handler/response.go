package handler

import "my-tourist-ticket/features/city"

type CityResponse struct {
	ID          uint   `json:"id"`
	CityName    string `json:"city_name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Thumbnail   string `json:"thumbnail"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func ModelToResponse(cityModel city.Core) CityResponse {
	return CityResponse{
		ID:          cityModel.ID,
		CityName:    cityModel.CityName,
		Description: cityModel.Description,
		Image:       cityModel.Image,
		Thumbnail:   cityModel.Thumbnail,
		CreatedAt:   cityModel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   cityModel.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func CoreToGetAllResponseCity(data city.Core) CityResponse {
	return CityResponse{
		ID:          data.ID,
		CityName:    data.CityName,
		Description: data.Description,
		Image:       data.Image,
		Thumbnail:   data.Thumbnail,
		CreatedAt:   data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func CoreToResponseListGetAllCity(data []city.Core) []CityResponse {
	var results []CityResponse
	for _, v := range data {
		results = append(results, CoreToGetAllResponseCity(v))
	}
	return results
}
