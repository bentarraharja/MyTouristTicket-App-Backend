package handler

import (
	"my-tourist-ticket/app/middlewares"
	"my-tourist-ticket/features/city"
	"my-tourist-ticket/utils/responses"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CityHandler struct {
	cityService city.CityServiceInterface
}

func NewCity(service city.CityServiceInterface) *CityHandler {
	return &CityHandler{
		cityService: service,
	}
}

func (handler *CityHandler) CreateCity(c echo.Context) error {
	userId := middlewares.ExtractTokenUserId(c)

	userRole, err := handler.cityService.GetUserRoleById(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Internal Server Error", nil))
	}
	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, responses.WebResponse("Forbidden - User is not an admin", nil))
	}

	var cityReq CityRequest
	if err := c.Bind(&cityReq); err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	cityCore := RequestToCore(cityReq)

	if cityCore.CityName == "" || cityCore.Description == "" {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("City name and description are required", nil))
	}

	_, imageHeader, err := c.Request().FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error retrieving the image file", nil))
	}

	_, thumbnailHeader, err := c.Request().FormFile("thumbnail")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error retrieving the thumbnail file", nil))
	}

	err = handler.cityService.Create(cityCore, imageHeader, thumbnailHeader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error creating city", nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("city created successfully", nil))
}

func (handler *CityHandler) UpdateCity(c echo.Context) error {
	userId := middlewares.ExtractTokenUserId(c)

	userRole, err := handler.cityService.GetUserRoleById(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Internal Server Error", nil))
	}
	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, responses.WebResponse("Forbidden - User is not an admin", nil))
	}
	cityID, err := strconv.Atoi(c.Param("city_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Invalid city ID", nil))
	}

	var cityReq CityRequest
	if err := c.Bind(&cityReq); err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Error binding data. Data not valid", nil))
	}

	// Ubah request menjadi core model
	cityCore := RequestToCore(cityReq)

	// Dapatkan file gambar dan thumbnail dari form
	_, imageHeader, err := c.Request().FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Error retrieving the image file", nil))
	}

	_, thumbnailHeader, err := c.Request().FormFile("thumbnail")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Error retrieving the thumbnail file", nil))
	}

	err = handler.cityService.Update(cityID, cityCore, imageHeader, thumbnailHeader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Error updating city", nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("City updated successfully", nil))
}

func (handler *CityHandler) GetCityById(c echo.Context) error {
	cityID, err := strconv.Atoi(c.Param("city_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Invalid city ID", nil))
	}

	cityData, err := handler.cityService.SelectCityById(cityID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Error retrieving city data", nil))
	}

	cityResponse := ModelToResponse(cityData)

	return c.JSON(http.StatusOK, responses.WebResponse("City data retrieved successfully", cityResponse))
}

func (handler *CityHandler) DeleteCity(c echo.Context) error {
	userId := middlewares.ExtractTokenUserId(c)

	userRole, err := handler.cityService.GetUserRoleById(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Internal Server Error", nil))
	}
	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, responses.WebResponse("Forbidden - User is not an admin", nil))
	}

	id := c.Param("city_id")
	idParam, errConv := strconv.Atoi(id)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error. id should be number", nil))
	}

	err = handler.cityService.Delete(idParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error delete city. delete failed "+err.Error(), nil))
	}
	return c.JSON(http.StatusOK, responses.WebResponse("success delete city", nil))
}

func (handler *CityHandler) GetAllCity(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	citys, totalPage, err := handler.cityService.SelectAllCity(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error get data", nil))
	}

	cityResponses := CoreToResponseListGetAllCity(citys)

	return c.JSON(http.StatusOK, responses.WebResponsePagination("success get data", cityResponses, totalPage))
}
