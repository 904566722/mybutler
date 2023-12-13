package db

import (
	"fmt"
	"strconv"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"changeme/backend/pkg/configs"
	"changeme/backend/pkg/model/entity"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func InitMysql() {
	mysqlConf := configs.Default().Mysql
	dsn := genMysqlDSN(mysqlConf.User, mysqlConf.Passwd, mysqlConf.Host, strconv.Itoa(mysqlConf.Port), mysqlConf.DBName)
	db = NewDB(dsn)

	if err := db.AutoMigrate(
		&entity.Task{},
		&entity.SubTask{},
		&entity.TaskCheckIn{},
	); err != nil {
		panic(err)
	}
}

func NewDB(dsn string) *gorm.DB {
	dbOnce.Do(func() {
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	})
	return db
}

func Default() *gorm.DB {
	return db
}

func genMysqlDSN(user, passwd, host, port, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, passwd, host, port, dbname)
}
