package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
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
	Postgres
)

type DBConfig struct {
	Type     DBType
	File     string
	URI      string
	LogLevel string
}

var SQL *gorm.DB

func Init(config *DBConfig) (err error) {
	var d gorm.Dialector
	switch config.Type {
	case SQLite:
		d = getSQLiteDialector(config)
	case MySQL:
		d = getMySqlDialector(config)
	case Postgres:
		d = getPostgresDialector(config)
	default:
		return errors.New("unsupported database")
	}

	db, err := gorm.Open(d, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",                              // table name prefix, table for `User` would be `t_users`
			SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
		Logger: logger.New(
			log.New(os.Stdout, "[GORM] ", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,                  // Slow SQL threshold
				LogLevel:                  getLogLevel(config.LogLevel), // Log level
				IgnoreRecordNotFoundError: true,                         // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,                        // Disable color
			}),
	})
	if err != nil {
		return err
	}

	SQL = db
	return nil
}

func getLogLevel(level string) logger.LogLevel {
	l := strings.ToLower(level)
	switch l {
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Silent
	}
}

func getSQLiteDialector(config *DBConfig) gorm.Dialector {
	f := config.File
	if strings.TrimSpace(f) == "" {
		f = utils.DBFile()
	}

	log.Println("db file -->", f)
	return sqlite.Open(fmt.Sprintf("%s?loc=Asia/Shanghai", f))
}

func getMySqlDialector(config *DBConfig) gorm.Dialector {
	return mysql.Open(config.URI)
}

func getPostgresDialector(config *DBConfig) gorm.Dialector {
	return postgres.Open(config.URI)
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
