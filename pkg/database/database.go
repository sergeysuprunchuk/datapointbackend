package database

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

const (
	PostgreSQL = "PostgreSQL"
)

type Config struct {
	Host         string `json:"host" yaml:"host"`
	Port         int    `json:"port" yaml:"port"`
	Username     string `json:"username" yaml:"username"`
	Password     string `json:"password" yaml:"password"`
	DatabaseName string `json:"databaseName" yaml:"database_name"`
	Driver       string `json:"driver" yaml:"driver"`
}

func (cfg *Config) Build() (
	driverName,
	dataSourceName string,
	placeholder sq.PlaceholderFormat,
	err error,
) {
	switch cfg.Driver {
	case PostgreSQL:
		return "postgres", fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName,
		), sq.Dollar, nil

	default:
		return "", "", nil, fmt.Errorf("неизвестный драйвер %s", cfg.Driver)
	}
}

type Database struct {
	Conn    *sql.DB
	Builder sq.StatementBuilderType
	Config  Config
}

func New(cfg Config) (*Database, error) {
	driverName, dataSourceName, placeholder, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	var conn *sql.DB

	if conn, err = sql.Open(driverName, dataSourceName); err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return &Database{
		Conn: conn,
		Builder: sq.StatementBuilder.
			PlaceholderFormat(placeholder).
			RunWith(conn),
		Config: cfg,
	}, nil
}
