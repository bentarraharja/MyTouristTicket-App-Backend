package handler

import (
	"my-tourist-ticket/app/middlewares"
	"my-tourist-ticket/features/tour"
	"my-tourist-ticket/utils/responses"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type TourHandler struct {
	tourService tour.TourServiceInterface
}

func NewTour(service tour.TourServiceInterface) *TourHandler {
	return &TourHandler{
		tourService: service,
	}
}

func (handler *TourHandler) CreateTour(c echo.Context) error {
	userId := middlewares.ExtractTokenUserId(c)

	userRole, err := handler.tourService.GetUserRoleById(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Internal Server Error", nil))
	}
	if userRole != "pengelola" {
		return c.JSON(http.StatusForbidden, responses.WebResponse("Forbidden - User is not a pengelola", nil))
	}

	var tourReq TourRequest

	if err := c.Bind(&tourReq); err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	// Mendapatkan file gambar dan thumbnail dari formulir
	_, imageHeader, err := c.Request().FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error retrieving the image file", nil))
	}

	_, thumbnailHeader, err := c.Request().FormFile("thumbnail")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error retrieving the thumbnail file", nil))
	}

	tourCore := RequestToCore(tourReq)
	tourCore.UserId = uint(userId)

	// Memanggil tourService.Insert dengan argumen yang sesuai, termasuk ID pengguna
	err = handler.tourService.Insert(uint(userId), tourCore, imageHeader, thumbnailHeader)
	if err != nil {
		if strings.Contains(err.Error(), "is required") {
			return c.JSON(http.StatusBadRequest, responses.WebResponse(err.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error creating tour", nil))
	}

	return c.JSON(http.StatusCreated, responses.WebResponse("tour created successfully", nil))
}

func (handler *TourHandler) UpdateTour(c echo.Context) error {
	userId := middlewares.ExtractTokenUserId(c)

	userRole, err := handler.tourService.GetUserRoleById(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Internal Server Error", nil))
	}
	if userRole != "pengelola" {
		return c.JSON(http.StatusForbidden, responses.WebResponse("Forbidden - User is not an pengelola", nil))
	}

	tourID, err := strconv.Atoi(c.Param("tour_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Invalid tour ID", nil))
	}

	var tourReq TourRequest
	if err := c.Bind(&tourReq); err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Error binding data. Data not valid", nil))
	}

	// Ubah request menjadi core model
	tourCore := RequestToCore(tourReq)

	// Dapatkan file gambar dan thumbnail dari form
	_, imageHeader, err := c.Request().FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Error retrieving the image file", nil))
	}

	_, thumbnailHeader, err := c.Request().FormFile("thumbnail")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Error retrieving the thumbnail file", nil))
	}

	err = handler.tourService.Update(tourID, tourCore, imageHeader, thumbnailHeader)
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.WebResponse("Error updating tour, "+err.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Tour updated successfully", nil))
}

func (handler *TourHandler) GetTourById(c echo.Context) error {
	tourID, err := strconv.Atoi(c.Param("tour_id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.WebResponse("Not found", nil))
	}

	tourData, err := handler.tourService.SelectTourById(tourID)
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.WebResponse("Error retrieving tour data", nil))
	}

	tourResponse := CoreToGetAllResponseTour(tourData)

	return c.JSON(http.StatusOK, responses.WebResponse("Tour data retrieved successfully", tourResponse))
}

func (handler *TourHandler) DeleteTour(c echo.Context) error {
	userId := middlewares.ExtractTokenUserId(c)

	userRole, err := handler.tourService.GetUserRoleById(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Internal Server Error", nil))
	}
	if userRole == "costumer" {
		return c.JSON(http.StatusForbidden, responses.WebResponse("Forbidden - Sorry you not have access", nil))
	}

	id := c.Param("tour_id")
	idParam, errConv := strconv.Atoi(id)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error. id should be number", nil))
	}

	err = handler.tourService.Delete(idParam)
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.WebResponse("error delete tour. delete failed "+err.Error(), nil))
	}
	return c.JSON(http.StatusOK, responses.WebResponse("success delete tour", nil))
}

func (handler *TourHandler) GetAllTour(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	tours, totalPage, err := handler.tourService.SelectAllTour(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error get data", nil))
	}

	tourResponses := CoreToResponseListGetAllTourIncludeReport(tours)

	return c.JSON(http.StatusOK, responses.WebResponsePagination("success get data", tourResponses, totalPage))
}

// GetTourByPengelola handles the request to get tours by pengelola.
func (handler *TourHandler) GetTourByPengelola(c echo.Context) error {
	userId := middlewares.ExtractTokenUserId(c)

	userRole, err := handler.tourService.GetUserRoleById(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Internal Server Error", nil))
	}
	if userRole != "pengelola" {
		return c.JSON(http.StatusForbidden, responses.WebResponse("Forbidden - User is not a pengelola", nil))
	}

	// Parse pagination parameters from the request
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	tours, totalPage, err := handler.tourService.SelectTourByPengelola(userId, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Error reading data", nil))
	}

	tourResponses := CoreToResponseListGetAllTour(tours)

	return c.JSON(http.StatusOK, responses.WebResponsePagination("success get data", tourResponses, totalPage))
}

func (handler *TourHandler) GetTourByCityID(c echo.Context) error {

	cityID, err := strconv.Atoi(c.Param("city_id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.WebResponse("Invalid city Id", nil))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	tours, totalPage, err := handler.tourService.GetTourByCityID(uint(cityID), page, limit)
	if err != nil {
		if err.Error() == "city not found" {
			return c.JSON(http.StatusNotFound, responses.WebResponse("City not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Error reading data", nil))
	}

	tourResponses := CoreToResponseListGetAllTour(tours)

	return c.JSON(http.StatusOK, responses.WebResponsePagination("success get data", tourResponses, totalPage))
}

func (handler *TourHandler) CreateReportTour(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}
	tourId, err := strconv.Atoi(c.Param("tour_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Invalid Tour Id", nil))
	}

	newReport := ReportRequest{}
	newReport.UserId = uint(userIdLogin)
	newReport.TourId = uint(tourId)
	errBind := c.Bind(&newReport)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	if newReport.TextReport == "" {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Text report are required", nil))
	}

	reportCore := ReportRequestToCore(newReport)
	errInsert := handler.tourService.InsertReportTour(userIdLogin, tourId, reportCore)
	if errInsert != nil {
		if strings.Contains(errInsert.Error(), "Error 1062 (23000): Duplicate entry") {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
		} else if strings.Contains(errInsert.Error(), "user has already reported this tour") {
			return c.JSON(http.StatusConflict, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
		} else {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
		}
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success insert data", nil))
}

func (handler *TourHandler) GetReportTour(c echo.Context) error {
	userId := middlewares.ExtractTokenUserId(c)

	userRole, err := handler.tourService.GetUserRoleById(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Internal Server Error", nil))
	}
	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, responses.WebResponse("Forbidden - User is not a admin", nil))
	}
	tourID, err := strconv.Atoi(c.Param("tour_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Invalid Tour Id", nil))
	}

	reports, err := handler.tourService.SelectReportTour(tourID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Error reading data", nil))
	}

	reportResponses := CoreReportToResponseListGetReportTour(reports)

	return c.JSON(http.StatusOK, responses.WebResponse("success get data", reportResponses))
}

func (handler *TourHandler) SearchTour(c echo.Context) error {
	query := c.QueryParam("tour_name")
	if query == "" {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("query parameter is required", nil))
	}

	tours, err := handler.tourService.SearchTour(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error read data", nil))
	}

	if len(tours) == 0 {
		return c.JSON(http.StatusNotFound, responses.WebResponse("The provided search query is not valid. Please provide a valid search term.", nil))
	}

	productResponses := CoreToResponseListGetAllTour(tours)
	return c.JSON(http.StatusOK, responses.WebResponse("success read data", productResponses))
}
