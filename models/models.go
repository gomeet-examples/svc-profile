package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Profile is the GORM model definition.
type Profile struct {
	UUID      uuid.UUID `gorm:"primary_key;type:char(36);not null"`
	Gender    uint16
	Email     string    `gorm:"type:varchar(100);unique_index"`
	Name      string    `gorm:"size:255"`
	Birthday  time.Time `gorm:"type:date;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (*Profile) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UUID", uuid.New())
	return nil
}

// All Mysql models
func mysqlModels() (values []interface{}) {
	values = append(values, &Profile{})

	return
}
