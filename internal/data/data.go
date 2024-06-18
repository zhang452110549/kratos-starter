package data

import (
	"kratos-starter/internal/asset"
	"kratos-starter/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Data .
type Data struct {
	// wrapped database client
	db  *gorm.DB
	rdb redis.UniversalClient
}

// NewData .
func NewData(c *conf.Data, _logger log.Logger,
	_db *gorm.DB,
	_rdb redis.UniversalClient,
) (*Data, func(), error) {
	logger := log.NewHelper(log.With(_logger, "module", "repository/data"))

	if err := upDBByGoose(_db); err != nil {
		panic(err)
	}

	cleanup := func() {
		logger.Info("closing the data resources")
		if sd, err := _db.DB(); err != nil {
			logger.Error(err)
		} else {
			_ = sd.Close()
		}
		if err := _rdb.Close(); err != nil {
			logger.Error(err)
		}
	}
	return &Data{
		db:  _db,
		rdb: _rdb,
	}, cleanup, nil
}

func upDBByGoose(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	goose.SetBaseFS(asset.DataSql)
	goose.SetTableName("t_goose_db_version")

	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	return goose.Up(sqlDB, "data_sql")
}
