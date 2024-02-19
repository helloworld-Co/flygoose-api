package datasource

import (
	"flygoose/configs"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	masterDB *gorm.DB
)

func GetMasterDB() *gorm.DB {
	return masterDB
}

// 初始化mysql
func InitMySql(cfg *configs.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	masterDB = db
}

// 初始化postgresql
func InitPostgreSQL(cfg *configs.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	masterDB = db
}

// 创建数据库
func CheckAndCreateDatabase(cfg *configs.Config) (err error) {

	dbName := cfg.Database.Name
	if dbName == "" {
		println("databasenameerror:数据库名称不能为空")
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", cfg.Database.User, cfg.Database.Password, cfg.Database.Host)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
	})
	if err != nil {
		println("sql.Open error %w\n", err)
		return
	}

	println("sql open sucess:" + dsn)

	createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database.Name)
	db.Exec(createDB)
	if err != nil {
		println("db.Execerror %w\n", err)
	}
	println("createDB sucess :" + createDB)
	return
}
