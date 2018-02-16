package models

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Profile is the GORM model definition.
type Profile struct {
	Uuid      string `gorm:"primary_key;type:char(36);not null"`
	Gender    uint16
	Email     string `gorm:"type:varchar(100);unique_index"`
	Name      string `gorm:"size:255"`
	Birthday  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// All Mysql models
func mysqlModels() (values []interface{}) {
	values = append(values, &Profile{})

	return
}
