package database

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"strings"
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
	AvgSql   = "avg"
	CountSql = "count"
	SumSql   = "sum"
	MaxSql   = "max"
	MinSql   = "min"
)

func (db *Database) GetFunctions() ([]string, error) {
	switch db.Config.Driver {
	case PostgreSQL:
		return []string{AvgSql, CountSql, SumSql, MaxSql, MinSql}, nil
	default:
		return nil, fmt.Errorf("неизвестный драйвер %s", db.Config.Driver)
	}
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

type Query struct {
	Table QTable `json:"table"`
}

type QTable struct {
	Name      string    `json:"name"`
	Increment uint8     `json:"increment"`
	Columns   []QColumn `json:"columns"`
	Next      []QTable  `json:"next"`
	Rule      Rule      `json:"rule"`
}

func (t *QTable) Join(b sq.SelectBuilder) sq.SelectBuilder {
	for _, nextT := range t.Next {
		b = b.Columns(nextT.SqlColumns()...)
		b = nextT.Rule.Join(b, *t, nextT)
		b = nextT.Join(b)
	}
	return b
}

func (t *QTable) SqlAlias() string {
	if t.Increment == 0 {
		return t.Name
	}
	return fmt.Sprintf("%s_%d", t.Name, t.Increment)
}

func (t *QTable) SqlColumns() []string {
	columns := make([]string, 0)
	for _, c := range t.Columns {
		if len(c.Fun) != 0 {
			columns = append(columns, fmt.Sprintf(`%s("%s"."%s") "%s %s.%s"`,
				c.Fun, t.SqlAlias(), c.Name, c.Fun, t.SqlAlias(), c.Name))
			continue
		}
		columns = append(columns, fmt.Sprintf(`"%s"."%s" "%s.%s"`,
			t.SqlAlias(), c.Name, t.SqlAlias(), c.Name))
	}
	return columns
}

type QColumn struct {
	Name     string `json:"name"`
	Fun      string `json:"fun"`
	Key      string `json:"key"`
	KeyOrder uint8  `json:"keyOrder"`
}

const (
	Join  = "join"
	Left  = "left"
	Right = "right"
)

type Rule struct {
	Type       string      `json:"type"`
	Conditions []Condition `json:"conditions"`
}

func (r *Rule) Join(b sq.SelectBuilder, left, right QTable) sq.SelectBuilder {
	sl := []string{fmt.Sprintf(`"%s" "%s" ON`, right.Name, right.SqlAlias())}

	for i, cond := range r.Conditions {
		if i != 0 {
			sl = append(sl, "AND")
		}
		sl = append(sl, fmt.Sprintf(`"%s"."%s" %s "%s"."%s"`,
			left.SqlAlias(), cond.Left, cond.Op, right.SqlAlias(), cond.Right,
		))
	}

	join := strings.Join(sl, " ")

	switch r.Type {
	case Join:
		b = b.Join(join)
	case Left:
		b = b.LeftJoin(join)
	case Right:
		b = b.RightJoin(join)
	}

	return b
}

type Condition struct {
	Left  string `json:"left"`
	Right string `json:"right"`
	Op    string `json:"operator"`
}

func (db *Database) Parse(query Query) (sq.SelectBuilder, map[string][]string) {
	b := db.Builder.
		Select(query.Table.SqlColumns()...).
		From(fmt.Sprintf(`"%s" "%s"`, query.Table.Name, query.Table.SqlAlias()))

	b = query.Table.Join(b)

	var (
		groupBy []string
		hasFun  bool
		rules   = make(map[string][]string)
	)

	tables := []QTable{query.Table}

	for i := 0; i < len(tables); i++ {
		t := tables[i]
		for _, c := range t.Columns {
			if len(c.Fun) != 0 {
				rules[c.Key] = append(rules[c.Key], fmt.Sprintf("%s %s.%s",
					c.Fun, t.SqlAlias(), c.Name))
				hasFun = true
				continue
			}

			groupBy = append(groupBy, fmt.Sprintf(`"%s.%s"`, t.SqlAlias(), c.Name))

			rules[c.Key] = append(rules[c.Key], fmt.Sprintf("%s.%s",
				t.SqlAlias(), c.Name))
		}

		if len(t.Next) != 0 {
			tables = append(tables, t.Next...)
		}
	}

	if hasFun {
		b = b.GroupBy(groupBy...)
	}

	return b, rules
}

type QueryResponse struct {
	Rules  map[string][]string `json:"rules"`
	Data   []map[string]any    `json:"data"`
	RawSql string              `json:"rawSql"`
}

func (db *Database) Execute(ctx context.Context, query Query) (QueryResponse, error) {
	b, rules := db.Parse(query)

	rows, err := b.QueryContext(ctx)
	if err != nil {
		return QueryResponse{}, err
	}
	defer func() { _ = rows.Close() }()

	var columns []string

	if columns, err = rows.Columns(); err != nil {
		return QueryResponse{}, err
	}

	response := QueryResponse{Rules: rules}

	for rows.Next() {
		var (
			dest = make([]any, 0)
			item = make(map[string]any)
		)

		for _, c := range columns {
			item[c] = new(any)
			dest = append(dest, item[c])
		}

		if err = rows.Scan(dest...); err != nil {
			return QueryResponse{}, err
		}

		response.Data = append(response.Data, item)
	}

	response.RawSql, _, _ = b.ToSql()

	return response, nil
}
