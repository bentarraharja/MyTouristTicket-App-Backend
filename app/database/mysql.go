package database

import (
	"fmt"
	"my-tourist-ticket/app/configs"
	bd "my-tourist-ticket/features/booking/data"
	cd "my-tourist-ticket/features/city/data"
	pd "my-tourist-ticket/features/package/data"
	td "my-tourist-ticket/features/tour/data"
	ud "my-tourist-ticket/features/user/data"
	vd "my-tourist-ticket/features/voucher/data"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDBMysql(cfg *configs.AppConfig) *gorm.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.DB_USERNAME, cfg.DB_PASSWORD, cfg.DB_HOSTNAME, cfg.DB_PORT, cfg.DB_NAME)

	DB, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return DB
}

func InitMigrate(db *gorm.DB) {
	db.AutoMigrate(&ud.User{}, &cd.City{}, &td.Tour{}, &td.Report{}, &pd.Package{}, &pd.Benefit{}, &vd.Voucher{}, &bd.Booking{}, &bd.Review{})
}
