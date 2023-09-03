package mysql

import (
	"fmt"
	"log"

	"github.com/lyujsally/MinTikTok-lyujsally/kitex_gen/relation"
	"github.com/lyujsally/MinTikTok-lyujsally/settings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *settings.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	//	创建表
	err = checkTableExist(&relation.User{})
	if err != nil {
		log.Printf("create table failed:%v", err)
		return
	}

	return nil
}
func checkTableExist(model interface{}) error {
	if DB.Migrator().HasTable(model) {
		return nil
	}
	return DB.AutoMigrate(model)
}

func Close() {
	if DB != nil {
		sqlDB, _ := DB.DB()
		sqlDB.Close()
	}
}
