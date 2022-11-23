package database

import (
	"database/sql"
	"log"
	"net"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewConfig() *Config {
	return &Config{}
}

func Connect(cfg *Config) (*sql.DB, error) {
	dsn := cfg.formatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error open DB connection:", err)
	}

	log.Printf("Created DB connection: %s/%s\n", cfg.Host, cfg.Database)

	return db, err
}

func Close(db *sql.DB) error {
	if err := db.Close(); err != nil {
		log.Println("Can not close DB connection:", err)
		return err
	}

	log.Println("Closed DB connection")
	return nil
}

func (cfg Config) formatDSN() string {
	c := mysql.NewConfig()
	c.Net = "tcp"
	c.Addr = net.JoinHostPort(cfg.Host, cfg.Port)
	c.User = cfg.User
	c.Passwd = cfg.Password
	c.DBName = cfg.Database
	c.ParseTime = true
	c.MultiStatements = true

	return c.FormatDSN()
}
