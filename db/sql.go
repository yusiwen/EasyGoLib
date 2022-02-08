package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MeloQi/EasyGoLib/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Model struct {
	ID        string         `structs:"id" gorm:"primary_key" form:"id" json:"id"`
	CreatedAt utils.DateTime `structs:"-" json:"createdAt" gorm:"type:datetime"`
	UpdatedAt utils.DateTime `structs:"-" json:"updatedAt" gorm:"type:datetime"`
	// DeletedAt *time.Time `sql:"index" structs:"-"`
}

type DBType int

const (
	SQLite DBType = iota
	MySQL
)

type DBConfig struct {
	Type     DBType
	URI      string
	LogLevel logger.LogLevel
}

var SQL *gorm.DB

func Init(config *DBConfig) (err error) {
	switch config.Type {
	case SQLite:
		SQL, err = createSQLite(config)
		if err != nil {
			return err
		}
	case MySQL:
		SQL, err = createMySQL(config)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported database")
	}
	return nil
}

func createSQLite(config *DBConfig) (*gorm.DB, error) {
	f := config.URI
	if strings.TrimSpace(f) == "" {
		f = utils.DBFile()
	}

	log.Println("db file -->", f)
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s?loc=Asia/Shanghai", f)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",                              // table name prefix, table for `User` would be `t_users`
			SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
		Logger: logger.New(
			log.New(os.Stdout, "[GORM] ", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,     // Slow SQL threshold
				LogLevel:                  config.LogLevel, // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,           // Disable color
			}),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createMySQL(config *DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.URI), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",                              // table name prefix, table for `User` would be `t_users`
			SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
		Logger: logger.New(
			log.New(os.Stdout, "[GORM] ", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,     // Slow SQL threshold
				LogLevel:                  config.LogLevel, // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,           // Disable color
			}),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Close() {
	if SQL != nil {
		db, err := SQL.DB()
		if err != nil {
			log.Println("cannot get DB")
			return
		}
		db.Close()
		SQL = nil
	}
}
