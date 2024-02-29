package handler

import "my-tourist-ticket/features/city"

type CityRequest struct {
	CityName    string `json:"city_name" form:"city_name"`
	Description string `json:"description" form:"description"`
	Image       string `json:"image" form:"image"`
	Thumbnail   string `json:"thumbnail" form:"thumbnail"`
}

func RequestToCore(input CityRequest) city.Core {
	return city.Core{
		CityName:    input.CityName,
		Description: input.Description,
		Image:       input.Image,
		Thumbnail:   input.Thumbnail,
	}
}
