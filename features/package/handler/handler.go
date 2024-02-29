package handler

import (
	"my-tourist-ticket/app/middlewares"
	packages "my-tourist-ticket/features/package"
	"my-tourist-ticket/utils/responses"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type PackageHandler struct {
	packageService packages.PackageServiceInterface
}

func New(ps packages.PackageServiceInterface) *PackageHandler {
	return &PackageHandler{
		packageService: ps,
	}
}

func (handler *PackageHandler) CreatePackage(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	tourID, err := strconv.Atoi(c.Param("tour_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error parsing tour id", nil))
	}

	newPackage := PackageRequest{}
	errBind := c.Bind(&newPackage)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data, data not valid", nil))
	}

	packageCore := RequestToCore(newPackage, uint(tourID))
	errCreate := handler.packageService.Create(newPackage.Benefits, packageCore)
	if errCreate != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error insert data", nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success insert data", nil))
}

func (handler *PackageHandler) GetPackageByTourId(c echo.Context) error {
	tourID, err := strconv.Atoi(c.Param("tour_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error parsing tour id", nil))
	}

	packages, err := handler.packageService.GetByTourId(uint(tourID))
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return c.JSON(http.StatusNotFound, responses.WebResponse("record not found", nil))
		} else {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error read data", nil))
		}
	}

	packageResponses := CoresToResponses(packages)
	return c.JSON(http.StatusOK, responses.WebResponse("success read data", packageResponses))
}

func (handler *PackageHandler) DeletePackage(c echo.Context) error {
	packageId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error. id should be number", nil))
	}

	errDelete := handler.packageService.Delete(packageId)
	if errDelete != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error delete data "+errDelete.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success delete data", nil))
}
