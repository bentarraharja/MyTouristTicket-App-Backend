package handler

import (
	"my-tourist-ticket/app/middlewares"
	"my-tourist-ticket/features/user"
	"my-tourist-ticket/utils/cloudinary"
	"my-tourist-ticket/utils/responses"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService user.UserServiceInterface
	cld         cloudinary.CloudinaryUploaderInterface
}

func New(service user.UserServiceInterface, cloudinaryUploader cloudinary.CloudinaryUploaderInterface) *UserHandler {
	return &UserHandler{
		userService: service,
		cld:         cloudinaryUploader,
	}
}

func (handler *UserHandler) RegisterUser(c echo.Context) error {
	newUser := UserRequest{}
	errBind := c.Bind(&newUser)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data, data not valid", nil))
	}

	userCore := RequestToCore(newUser)
	errInsert := handler.userService.Create(userCore)
	if errInsert != nil {
		if strings.Contains(errInsert.Error(), "Error 1062 (23000): Duplicate entry") {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("error register data. "+errInsert.Error(), nil))
		} else {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error register data. "+errInsert.Error(), nil))
		}
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success register data", nil))
}

func (handler *UserHandler) Login(c echo.Context) error {
	var reqData = LoginRequest{}
	errBind := c.Bind(&reqData)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data, data not valid", nil))
	}
	result, token, err := handler.userService.Login(reqData.Email, reqData.Password)
	if err != nil {
		if strings.Contains(err.Error(), "email wajib diisi.") {
			return c.JSON(http.StatusUnauthorized, responses.WebResponse("error login. "+err.Error(), nil))
		} else if strings.Contains(err.Error(), "password wajib diisi.") {
			return c.JSON(http.StatusUnauthorized, responses.WebResponse("error login. "+err.Error(), nil))
		} else if strings.Contains(err.Error(), "password tidak sesuai.") {
			return c.JSON(http.StatusUnauthorized, responses.WebResponse("error login. "+err.Error(), nil))
		} else {
			return c.JSON(http.StatusUnauthorized, responses.WebResponse("error login. "+err.Error(), nil))
		}
	}
	var responseData = CoreToResponseLogin(result, token)
	return c.JSON(http.StatusOK, responses.WebResponse("success login", responseData))
}

func (handler *UserHandler) GetUser(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	result, errSelect := handler.userService.GetById(userIdLogin)
	if errSelect != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error read data. "+errSelect.Error(), nil))
	}

	if result.NoKtp == "" {
		var userResult = CoreToResponseUser(result)
		return c.JSON(http.StatusOK, responses.WebResponse("success read data", userResult))
	} else {
		var userResult = CoreToResponsePengelola(result)
		return c.JSON(http.StatusOK, responses.WebResponse("success read data", userResult))
	}
}

func (handler *UserHandler) UpdateUser(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	var userData = UserRequestUpdate{}
	errBind := c.Bind(&userData)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	fileData, err := c.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error retrieving the file", nil))
	}

	var imageURL string
	if fileData != nil {
		imageURL, err = handler.cld.UploadImage(fileData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error uploading the image", nil))
		}
	}

	userCore := UpdateRequestToCore(userData, imageURL)
	errUpdate := handler.userService.Update(userIdLogin, userCore)
	if errUpdate != nil {
		if strings.Contains(errUpdate.Error(), "Error 1062 (23000): Duplicate entry") {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("error update data. "+errUpdate.Error(), nil))
		} else {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error update data. "+errUpdate.Error(), nil))
		}
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success update data", nil))
}

func (handler *UserHandler) DeleteUser(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	errDelete := handler.userService.Delete(userIdLogin)
	if errDelete != nil {
		if strings.Contains(errDelete.Error(), "error record not found") {
			return c.JSON(http.StatusNotFound, responses.WebResponse("error delete data. "+errDelete.Error(), nil))
		} else {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error delete data. "+errDelete.Error(), nil))
		}
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success delete data", nil))
}

func (handler *UserHandler) GetAdminUserData(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	result, errSelect, totalPage := handler.userService.GetAdminUsers(userIdLogin, page, limit)
	if errSelect != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error read data. "+errSelect.Error(), nil))
	}

	var userResult = CoreToResponseList(result)
	return c.JSON(http.StatusOK, responses.WebResponsePagination("success read data", userResult, totalPage))
}

func (handler *UserHandler) UpdateUserPengelolaById(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	pengelolaId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error parsing pengelola id", nil))
	}

	pengelolaStatus := c.QueryParam("status")

	// var userData = AdminPengelolaRequestUpdate{}
	// errBind := c.Bind(&userData)
	// if errBind != nil {
	// 	return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	// }

	// userCore := RequestToCoreAdminPengelola(userData)
	errUpdate := handler.userService.UpdatePengelola(userIdLogin, pengelolaStatus, pengelolaId)
	if errUpdate != nil {
		if strings.Contains(errUpdate.Error(), "error record not found") {
			return c.JSON(http.StatusNotFound, responses.WebResponse("error update data. "+errUpdate.Error(), nil))
		} else {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error update data. "+errUpdate.Error(), nil))
		}
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success update data", nil))
}
