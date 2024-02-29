package routes

import (
	"my-tourist-ticket/app/cache"
	"my-tourist-ticket/app/middlewares"
	"my-tourist-ticket/utils/cloudinary"
	"my-tourist-ticket/utils/encrypts"
	"my-tourist-ticket/utils/externalapi"

	ud "my-tourist-ticket/features/user/data"
	uh "my-tourist-ticket/features/user/handler"
	us "my-tourist-ticket/features/user/service"

	cd "my-tourist-ticket/features/city/data"
	ch "my-tourist-ticket/features/city/handler"
	cs "my-tourist-ticket/features/city/service"

	td "my-tourist-ticket/features/tour/data"
	th "my-tourist-ticket/features/tour/handler"
	ts "my-tourist-ticket/features/tour/service"

	pd "my-tourist-ticket/features/package/data"
	ph "my-tourist-ticket/features/package/handler"
	ps "my-tourist-ticket/features/package/service"

	vd "my-tourist-ticket/features/voucher/data"
	vh "my-tourist-ticket/features/voucher/handler"
	vs "my-tourist-ticket/features/voucher/service"

	bd "my-tourist-ticket/features/booking/data"
	bh "my-tourist-ticket/features/booking/handler"
	bs "my-tourist-ticket/features/booking/service"

	dd "my-tourist-ticket/features/dashboard/data"
	dh "my-tourist-ticket/features/dashboard/handler"
	ds "my-tourist-ticket/features/dashboard/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, e *echo.Echo, redisClient cache.RedisInterface) {
	hash := encrypts.New()
	cloudinaryUploader := cloudinary.New()
	midtrans := externalapi.New()

	userData := ud.New(db)
	userService := us.New(userData, hash)
	userHandlerAPI := uh.New(userService, cloudinaryUploader)

	cityData := cd.NewCity(db, cloudinaryUploader)
	cityService := cs.NewCity(cityData)
	cityHandlerAPI := ch.NewCity(cityService)

	tourData := td.NewTour(db, redisClient, cloudinaryUploader)
	tourService := ts.NewTour(tourData)
	tourHandlerAPI := th.NewTour(tourService)

	packageData := pd.New(db)
	packageService := ps.New(packageData)
	packageHandlerAPI := ph.New(packageService)

	voucherData := vd.New(db)
	voucherService := vs.New(voucherData)
	voucherHandlerAPI := vh.New(voucherService)

	bookingData := bd.New(db, midtrans)
	bookingService := bs.New(bookingData)
	bookingHandlerAPI := bh.New(bookingService)

	dashboardData := dd.NewDashboard(db)
	dashboardService := ds.New(dashboardData)
	dashboardHandlerAPI := dh.New(dashboardService)

	// define routes/ endpoint USERS
	e.POST("/login", userHandlerAPI.Login)
	e.POST("/users", userHandlerAPI.RegisterUser)
	e.GET("/users", userHandlerAPI.GetUser, middlewares.JWTMiddleware())
	e.PUT("/users", userHandlerAPI.UpdateUser, middlewares.JWTMiddleware())
	e.DELETE("/users", userHandlerAPI.DeleteUser, middlewares.JWTMiddleware())
	e.GET("/users/admin", userHandlerAPI.GetAdminUserData, middlewares.JWTMiddleware())
	e.PUT("/users/admin/:id", userHandlerAPI.UpdateUserPengelolaById, middlewares.JWTMiddleware())

	//define routes/ endpoint CITY
	e.POST("/citys", cityHandlerAPI.CreateCity, middlewares.JWTMiddleware())
	e.PUT("/citys/:city_id", cityHandlerAPI.UpdateCity, middlewares.JWTMiddleware())
	e.GET("/citys/:city_id", cityHandlerAPI.GetCityById)
	e.DELETE("/citys/:city_id", cityHandlerAPI.DeleteCity, middlewares.JWTMiddleware())
	e.GET("/citys", cityHandlerAPI.GetAllCity)

	//define routes/ endpoint TOUR
	e.POST("/tours", tourHandlerAPI.CreateTour, middlewares.JWTMiddleware())
	e.PUT("/tours/:tour_id", tourHandlerAPI.UpdateTour, middlewares.JWTMiddleware())
	e.GET("/tours/:tour_id", tourHandlerAPI.GetTourById)
	e.DELETE("/tours/:tour_id", tourHandlerAPI.DeleteTour, middlewares.JWTMiddleware())
	e.GET("/tours", tourHandlerAPI.GetAllTour)
	e.GET("/tours/pengelola", tourHandlerAPI.GetTourByPengelola, middlewares.JWTMiddleware())
	e.GET("/tours/bycity/:city_id", tourHandlerAPI.GetTourByCityID)
	e.POST("/tours/:tour_id/report", tourHandlerAPI.CreateReportTour, middlewares.JWTMiddleware())
	e.GET("/tours/:tour_id/report", tourHandlerAPI.GetReportTour, middlewares.JWTMiddleware())
	e.GET("/tours/search", tourHandlerAPI.SearchTour)
	e.GET("/tours/:tour_id/reviews", bookingHandlerAPI.GetAllBookingTourReview)

	//define routes/ endpoint PACKAGE
	e.POST("/packages/:tour_id", packageHandlerAPI.CreatePackage, middlewares.JWTMiddleware())
	e.GET("/packages/:tour_id", packageHandlerAPI.GetPackageByTourId)
	e.DELETE("/packages/:id", packageHandlerAPI.DeletePackage, middlewares.JWTMiddleware())

	//define routes/ endpoint VOUCHER
	e.POST("/vouchers", voucherHandlerAPI.CreateVoucher, middlewares.JWTMiddleware())
	e.GET("/vouchers", voucherHandlerAPI.GetAllVoucher, middlewares.JWTMiddleware())
	e.PUT("/vouchers/:voucher_id", voucherHandlerAPI.UpdateVoucher, middlewares.JWTMiddleware())
	e.DELETE("/vouchers/:id", voucherHandlerAPI.DeleteVoucher, middlewares.JWTMiddleware())

	//define routes/ endpoint Booking
	e.POST("/bookings", bookingHandlerAPI.CreateBooking, middlewares.JWTMiddleware())
	e.PUT("/bookings/:id", bookingHandlerAPI.CancelBookingById, middlewares.JWTMiddleware())
	e.POST("/bookings/:booking_id/review", bookingHandlerAPI.CreateBookingReview, middlewares.JWTMiddleware())
	e.POST("/bookings/notification", bookingHandlerAPI.WebhoocksNotification)
	e.GET("/bookings/users", bookingHandlerAPI.GetBookingUser, middlewares.JWTMiddleware())
	e.GET("/bookings/users/:id", bookingHandlerAPI.GetBookingUserDetail, middlewares.JWTMiddleware())
	e.GET("/bookings/admin", bookingHandlerAPI.GetAllBooking, middlewares.JWTMiddleware())
	e.GET("/bookings/pengelola", bookingHandlerAPI.GetAllBookingPengelola, middlewares.JWTMiddleware())

	//define routes/ endpoint Dashboard
	e.GET("/admin/dashboard", dashboardHandlerAPI.Dashboard, middlewares.JWTMiddleware())
}
