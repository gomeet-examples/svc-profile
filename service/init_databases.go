// Code generated by protoc-gen-gomeet-service. DO NOT EDIT.
// source: pb/profile.proto
package service

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type GormLogger struct{}

func (*GormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		log.WithFields(log.Fields{"module": "gorm", "type": "sql"}).Debug(v[3])
	}
	if v[0] == "log" {
		log.WithFields(log.Fields{"module": "gorm", "type": "log"}).Warn(v[2])
	}
}

// initDatabaseHandle initializes the databases handles of the server.
func (s *profileServer) initDatabaseHandle() error {
	// establish the connection if it's not ready
	if s.mysqlHandle == nil {
		if strings.Contains(s.mysqlDataSourceName, "?") {
			return fmt.Errorf("database connection error: data source name cannot contain options")
		}
		dsn := fmt.Sprintf("%s?charset=utf8&parseTime=True", s.mysqlDataSourceName)

		mysqlHandle, err := gorm.Open("mysql", dsn)
		if err != nil {
			log.WithFields(log.Fields{
				"DSN": s.mysqlDataSourceName,
			}).Infof("database connection error: %v", err)
			return err
		}
		mysqlHandle.SetLogger(&GormLogger{})
		mysqlHandle.LogMode(true)
		s.mysqlHandle = mysqlHandle
	}

	// ping the database server
	err := s.mysqlHandle.DB().Ping()
	if err != nil {
		return err
	}
	return nil
}

func (s *profileServer) closeDatabaseHandle() error {
	// close the connection if it's not empty
	if s.mysqlHandle != nil {
		err := s.mysqlHandle.Close()
		if err != nil {
			return err
		}
		s.mysqlHandle = nil
	}

	return nil
}
