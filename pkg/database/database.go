package database

import (
	"context"
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

const (
	NumberJSON  = "number"
	StringJSON  = "string"
	BooleanJSON = "boolean"
	Unsupported = "unsupported"
)

type SQLType string

func (t SQLType) ToJSON() string {
	switch t {
	case Smallint, Integer, Bigint, Decimal, Numeric, Real, Double, Smallserial, Serial, Bigserial:
		return NumberJSON

	case Boolean:
		return BooleanJSON

	case CharacterVarying, Character, Text, Timestamp, Timestamptz, Time, Timetz, Interval:
		return StringJSON

	default:
		return Unsupported
	}
}

const (
	Smallint    SQLType = "smallint"
	Integer     SQLType = "integer"
	Bigint      SQLType = "bigint"
	Decimal     SQLType = "decimal"
	Numeric     SQLType = "numeric"
	Real        SQLType = "real"
	Double      SQLType = "double precision"
	Smallserial SQLType = "smallserial"
	Serial      SQLType = "serial"
	Bigserial   SQLType = "bigserial"

	CharacterVarying SQLType = "character varying"
	Character        SQLType = "character"
	Text             SQLType = "text"

	Boolean SQLType = "boolean"

	Timestamp   SQLType = "timestamp without time zone"
	Timestamptz SQLType = "timestamp with time zone"
	Time        SQLType = "time without time zone"
	Timetz      SQLType = "time with time zone"
	Interval    SQLType = "interval"
)

type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

func (db *Database) GetTables(ctx context.Context) ([]*Table, error) {
	rows, err := db.Builder.
		Select("table_name", "column_name", "data_type").
		From("information_schema.columns").
		Where("table_schema = ?", "public").
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var (
		tables   = make([]*Table, 0)
		tableAcc = make(map[string]*Table)
	)
	for rows.Next() {
		var (
			tName string
			cName string
			cType SQLType
		)
		if err = rows.Scan(&tName, &cName, &cType); err != nil {
			return nil, err
		}

		_, ok := tableAcc[tName]
		if !ok {
			newT := &Table{Name: tName}
			tables = append(tables, newT)
			tableAcc[tName] = newT
		}

		tableAcc[tName].Columns = append(tableAcc[tName].Columns, Column{
			Name: cName,
			Type: cType.ToJSON(),
		})
	}

	return tables, nil
}
