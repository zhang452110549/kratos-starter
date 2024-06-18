package data

import (
	"fmt"

	"kratos-starter/internal/conf"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(cfg *conf.Data) *gorm.DB {
	mysqlCfg := cfg.GetMysql()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlCfg.GetUserName(),
		mysqlCfg.GetPassword(),
		mysqlCfg.GetAddress(),
		mysqlCfg.GetDbname())

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	return db
}
