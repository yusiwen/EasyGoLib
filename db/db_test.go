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
	count, err := doTest()
	if err != nil {
		t.Error(err)
	}
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
	count, err := doTest()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(1), count, "count should be 1")
}

func TestPostgres(t *testing.T) {
	err := Init(&DBConfig{
		Type:     Postgres,
		URI:      "host=lattepanda user=db_test_user password=db_test_user_password dbname=db_test port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		LogLevel: "Info",
	})
	if err != nil {
		t.Error(err)
	}
	count, err := doTest()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(1), count, "count should be 1")
}

func doTest() (int64, error) {
	err := SQL.AutoMigrate(TestModel{})
	if err != nil {
		return 0, err
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
	return count, nil
}
