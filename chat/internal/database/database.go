package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DatabaseConfig struct {
	Database string
	URL      string
	Username string
	Password string
}

type Database struct {
	cfg *DatabaseConfig
	db  *sql.DB
}

func NewDatabase(cfg *DatabaseConfig) (*Database, error) {
	d := &Database{cfg: cfg}

	if cfg == nil {
		return nil, fmt.Errorf("nil database config")
	}

	if cfg.URL == "" {
		return nil, fmt.Errorf("database url is empty")
	}

	parsed, err := mysql.ParseDSN(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid database dsn: %w", err)
	}

	parsed.User = cfg.Username
	parsed.Passwd = cfg.Password
	dsn := parsed.FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	d.db = db

	err = d.Ping()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Ping() error {
	if d.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := d.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (d *Database) DB() *sql.DB {
	return d.db
}
