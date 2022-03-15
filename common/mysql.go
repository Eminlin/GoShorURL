package common

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = newMySQLCon()
	if err != nil {
		log.Fatalf("connect mysql err:%s", err.Error())
	}
}

func newMySQLCon() (*gorm.DB, error) {
	c := AppConf.MySQL
	url := c.User + ":" + c.Password + "@tcp" + "(" + c.HostPort + ")/" + c.Database
	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Fatalf("invalid database source err:%s", err.Error())
	}
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(10)
	return db, nil
}
