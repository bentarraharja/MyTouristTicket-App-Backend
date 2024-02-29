package handler

import (
	"my-tourist-ticket/app/middlewares"
	"my-tourist-ticket/features/dashboard"
	"my-tourist-ticket/utils/responses"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	dashboardService dashboard.DashboardServiceInterface
}

func New(ds dashboard.DashboardServiceInterface) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: ds,
	}
}

// func (handler *DashboardHandler) Dashboard(c echo.Context) error {
// 	// Call the service method to fetch dashboard data
// 	dash, err := handler.dashboardService.Dashboard()
// 	if err != nil {
// 		log.Printf("Error fetching dashboard data: %v", err)
// 		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Failed to fetch dashboard data", nil))
// 	}

// 	// Return the dashboard data as JSON response
// 	return c.JSON(http.StatusOK, responses.WebResponse("Success", dash))
// }

func (handler *DashboardHandler) Dashboard(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	dash, err := handler.dashboardService.Dashboard(userIdLogin)
	if err != nil {
		if strings.Contains(err.Error(), "akses") {
			return c.JSON(http.StatusForbidden, responses.WebResponse(err.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error get dashboard", nil))
	}

	dashboardResp := CoreToDashboardResponse(dash)

	return c.JSON(http.StatusOK, dashboardResp)
}
