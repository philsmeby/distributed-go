package datamanager

import (
	"github.com/distributed-go/dto"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db  *gorm.DB
	err error
)

type SensorMessage struct {
	gorm.Model
	dto.SensorMessage
}

func init() {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect to database")
	}
	//defer db.Close()

	//Migrate the schema
	db.AutoMigrate(&SensorMessage{})
	//
	//// Create
	//db.Create(&Product{Code: "L1212", Price: 1000})
	//
	//// Read
	//var product Product
	//db.First(&product, 1) // find product with id 1
	//db.First(&product, "code = ?", "L1212") // find product with code l1212
	//
	//// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 2000)
	//
	//// Delete - delete product
	//db.Delete(&product)
}
