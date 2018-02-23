package service

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/gomeet-examples/svc-profile/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.PanicLevel)
	return
}

func newProfileServerTest(t *testing.T) (*profileServer, *gorm.DB, error) {
	// FIXME: do something better with mock? or with ENV VAR or with flags
	// see. https://siongui.github.io/2017/04/28/command-line-argument-in-golang-test/
	dsn := "gomeet:totomysql@tcp(localhost:3306)/svc_profile_test"
	models.MigrateSchema(dsn)
	server := &profileServer{
		mysqlDataSourceName: dsn,
	}
	err := server.initDatabaseHandle()
	assert.Nil(t, err, "Create: init database fail")
	assert.NotNil(t, server.mysqlHandle, "Create: init database fail nil handle")
	db := server.mysqlHandle.Exec("TRUNCATE TABLE profiles")

	return server, db, err
}
