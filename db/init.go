package db

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/pajri/personal-backend/config"
)

func InitDB() (*sql.DB, error) {
	dbConfig := config.Config.DB
	cfg := mysql.NewConfig()
	cfg.User = dbConfig.Username
	cfg.Addr = dbConfig.Host
	cfg.Net = "tcp"
	cfg.Params = map[string]string{"parseTime": "true"}
	cfg.DBName = dbConfig.DbName
	cfg.Passwd = dbConfig.Password
	cfg.AllowCleartextPasswords = true

	return sql.Open(`mysql`, cfg.FormatDSN())
}
