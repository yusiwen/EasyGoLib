package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
	Id      string `gorm:"type:varchar(256);primary_key;unique"`
	Name    string `gorm:"type:varchar(256)"`
	Address string `gorm:"type:varchar(256)"`
	Email   string `gorm:"type:varchar(256)"`
}

func TestSQLite(t *testing.T) {
	err := Init(&DBConfig{
		Type:     SQLite,
		URI:      "",
		LogLevel: "Info",
	})
	if err != nil {
		t.Error(err)
	}
	err = SQL.AutoMigrate(TestModel{})
	if err != nil {
		t.Error(err)
	}
	SQL.Where("1=1").Delete(&TestModel{})
	result := SQL.Create(&TestModel{
		Id:      "1",
		Name:    "Test_User",
		Address: "Test Address 1#",
		Email:   "test@email.com",
	})
	if result.Error != nil {
		t.Error(result.Error)
	}
	count := int64(-1)
	SQL.Model(TestModel{}).Where("name = ?", "Test_User").Count(&count)
	assert.Equal(t, int64(1), count, "count should be 1")
}

func TestMySQL(t *testing.T) {
	err := Init(&DBConfig{
		Type:     MySQL,
		URI:      "db_test_user:db_test_user_password@tcp(lattepanda:3306)/db_test?charset=utf8mb4&parseTime=True&loc=Local",
		LogLevel: "Info",
	})
	if err != nil {
		t.Error(err)
	}
	err = SQL.AutoMigrate(TestModel{})
	if err != nil {
		t.Error(err)
	}
	SQL.Where("1=1").Delete(&TestModel{})
	SQL.Create(&TestModel{
		Id:      "1",
		Name:    "Test_User",
		Address: "Test Address 1#",
		Email:   "test@email.com",
	})
	count := int64(-1)
	SQL.Model(TestModel{}).Where("name = ?", "Test_User").Count(&count)
	assert.Equal(t, int64(1), count, "count should be 1")
}
