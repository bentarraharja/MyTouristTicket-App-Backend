package handler

import (
	"my-tourist-ticket/features/tour"
)

type TourRequest struct {
	TourName    string  `json:"tour_name" form:"tour_name"`
	Description string  `json:"description" form:"description"`
	Image       string  `json:"image" form:"image"`
	Thumbnail   string  `json:"thumbnail" form:"thumbnail"`
	Address     string  `json:"address" form:"address"`
	Latitude    float64 `json:"latitude" form:"latitude"`
	Longitude   float64 `json:"longitude" form:"longitude"`
	CityId      uint    `json:"city_id" form:"city_id"`
	UserId      uint    `json:"user_id" form:"user_id"`
}

type ReportRequest struct {
	UserId     uint   `json:"user_id" form:"user_id"`
	TourId     uint   `json:"tour_id" form:"tour_id"`
	TextReport string `json:"text_report" form:"text_report"`
}

func RequestToCore(input TourRequest) tour.Core {
	return tour.Core{
		TourName:    input.TourName,
		Description: input.Description,
		Image:       input.Image,
		Thumbnail:   input.Thumbnail,
		Address:     input.Address,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
		CityId:      input.CityId,
		UserId:      input.UserId,
	}
}

func ReportRequestToCore(input ReportRequest) tour.ReportCore {
	return tour.ReportCore{
		UserId:     input.UserId,
		TourId:     input.TourId,
		TextReport: input.TextReport,
	}
}
