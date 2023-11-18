package models

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"Gin-Example/src/gin-blog/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func init() {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)

	// 读取配置
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatalf("Fail to get section 'database' : %v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	db, err := gorm.Open(dbType, fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName,
	))

	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	// 配置Gorm在转义时 不会使用复数
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)  // 连接池最大连接数
	db.DB().SetMaxOpenConns(100) // 最大打开连接数

}

func closeDB() {
	defer db.Close()
}
