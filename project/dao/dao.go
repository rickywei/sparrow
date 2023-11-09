package dao

import (
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/rickywei/sparrow/project/conf"
	"github.com/rickywei/sparrow/project/logger"
	"github.com/rickywei/sparrow/project/po"
)

var (
	Q  *Query
	db *gorm.DB

	models = []any{&po.User{}}
)

func init() {
	var err error
	masters := conf.Strings("mysql.dsn.masters")
	if len(masters) <= 0 {
		logger.L().Fatal("length of mysql.dsn.masters <= 0")
	}
	if db, err = gorm.Open(mysql.Open(masters[0]), &gorm.Config{}); err != nil {
		logger.L().Fatal("connect db failed", zap.Error(err))
	}
	if err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:           lo.Map(masters, func(dsn string, _ int) gorm.Dialector { return mysql.Open(dsn) }),
		Replicas:          lo.Map(conf.Strings("mysql.dsn.slaves"), func(dsn string, _ int) gorm.Dialector { return mysql.Open(dsn) }),
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: true,
	})); err != nil {
		logger.L().Fatal("connect db failed", zap.Error(err))
	}

	if conf.Bool("mysql.migrate.enable") {
		if err = db.AutoMigrate(models...); err != nil {
			logger.L().Fatal("auto migrate failed", zap.Error(err))
		}
	}

	Q = Use(db)
}
