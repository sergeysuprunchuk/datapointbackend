package database

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"strconv"
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

func (db *Database) GetFunctions() (map[string][]string, error) {
	switch db.Config.Driver {
	case PostgreSQL:
		return map[string][]string{
			NumberJSON: {AvgSql, CountSql, SumSql, MaxSql, MinSql},
		}, nil
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
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	IsPKey   bool   `json:"isPKey"`
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

func (t *Table) GetPKey() (Column, bool) {
	for _, column := range t.Columns {
		if column.IsPKey {
			return column, true
		}
	}
	return Column{}, false
}

func (db *Database) GetTables(ctx context.Context) ([]*Table, error) {
	rows, err := db.Builder.
		Select(
			"c.table_name",
			"c.column_name",
			"c.data_type",
			"NOT (c.is_nullable::boolean OR c.column_default IS NOT NULL)",          //required
			"tc.constraint_type IS NOT NULL AND tc.constraint_type = 'PRIMARY KEY'", //is_pkey
		).
		From("information_schema.columns c").
		LeftJoin("information_schema.constraint_column_usage ccu USING (table_name, column_name)").
		LeftJoin("information_schema.table_constraints tc USING (constraint_name)").
		Where("c.table_schema = 'public' AND (tc.constraint_type = 'PRIMARY KEY' OR tc.constraint_type IS NULL)").
		OrderBy("c.table_name", "c.column_name").
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
			c     Column
			cType SQLType
		)
		if err = rows.Scan(&tName, &c.Name, &cType, &c.Required, &c.IsPKey); err != nil {
			return nil, err
		}

		_, ok := tableAcc[tName]
		if !ok {
			newT := &Table{Name: tName}
			tables = append(tables, newT)
			tableAcc[tName] = newT
		}

		c.Type = cType.ToJSON()

		tableAcc[tName].Columns = append(tableAcc[tName].Columns, c)
	}

	return tables, nil
}

func (db *Database) GetTable(ctx context.Context, name string) (*Table, error) {
	tables, err := db.GetTables(ctx)
	if err != nil {
		return nil, err
	}

	for _, table := range tables {
		if table.Name == name {
			return table, nil
		}
	}

	return nil, fmt.Errorf("таблицы с именем %s не существует", name)
}

const (
	Select = "select"
	Insert = "insert"
	Update = "update"
	Delete = "delete"
)

type Query struct {
	Type    string     `json:"type"`
	Table   *QTable    `json:"table"`
	Columns []*QColumn `json:"columns"`
	Where   []*QColumn `json:"where"`

	//используется только в select.
	OrderBy []*QColumn `json:"orderBy"`
	Limit   uint64     `json:"limit"`
	Offset  uint64     `json:"offset"`
}

type QTableKey struct {
	Name      string `json:"name"`      //имя таблицы.
	Increment uint8  `json:"increment"` //приращение имени для создания уникальных псевдонимов.
}

func (k *QTableKey) String() string {
	if k.Increment == 0 {
		return k.Name
	}
	return fmt.Sprintf("%s_%d", k.Name, k.Increment)
}

type QTable struct {
	QTableKey

	//используется только в select.
	Next []*QTable `json:"next"` //таблицы, которые объединены с этой таблицей.
	Rule *Rule     `json:"rule"` //правило объединения с предыдущей таблицей.
}

func (t *QTable) Partial() string {
	return fmt.Sprintf(`"%s"`, t.Name)
}

func (t *QTable) Full() string {
	return fmt.Sprintf(`%s "%s"`, t.Partial(), t.String())
}

func (t *QTable) Join(b sq.SelectBuilder) (sq.SelectBuilder, error) {
	for _, nextQt := range t.Next {
		join := []string{nextQt.Full(), "ON"}

		if nextQt.Rule == nil {
			//возможно, объект в этот момент уже был изменен, но это не важно.
			return b, fmt.Errorf("невозможно объединить таблицы, поскольку не указаны правила")
		}

		for i, c := range nextQt.Rule.Conditions {
			if i != 0 {
				join = append(join, "AND")
			}

			join = append(join, c.Columns[0].Partial(), c.Operator, c.Columns[1].Partial())
		}

		switch nextQt.Rule.Type {
		case Join:
			b = b.Join(strings.Join(join, " "))
		case Left:
			b = b.LeftJoin(strings.Join(join, " "))
		case Right:
			b = b.RightJoin(strings.Join(join, " "))
		default:
			return b, fmt.Errorf("неизвестный тип соединения: %s", nextQt.Rule.Type)
		}

		var err error

		if b, err = nextQt.Join(b); err != nil {
			return b, err
		}
	}

	return b, nil
}

// специальные ключи к QColumn Payload для создания правил разбора.
const (
	MetaKey = "metaKey"
	Order   = "order"
)

type QColumn struct {
	Column

	TableKey QTableKey `json:"tableKey"` //ключ таблицы, которой принадлежит столбец.

	Payload map[string]any `json:"payload"` //специальные данные, привязанные к этому столбцу.

	//используется только в select.
	Func string `json:"func"`

	//используется в insert, update, delete и where.
	Value any `json:"value"`
}

func (c *QColumn) MetaKey() string {
	key, ok := c.Payload[MetaKey].(string)
	if !ok {
		return ""
	}
	return key
}

func (c *QColumn) String() string {
	result := fmt.Sprintf("%s.%s", c.TableKey.String(), c.Name)
	if len(c.Func) == 0 {
		return result
	}
	return fmt.Sprintf("%s %s", c.Func, result)
}

func (c *QColumn) Partial() string {
	if len(c.Func) != 0 {
		return fmt.Sprintf(`%s("%s"."%s")`, c.Func, c.TableKey.String(), c.Name)
	}
	return fmt.Sprintf(`"%s"."%s"`, c.TableKey.String(), c.Name)
}

func (c *QColumn) Full() string {
	return fmt.Sprintf(`%s "%s"`, c.Partial(), c.String())
}

const (
	Join  = "join"
	Left  = "left"
	Right = "right"
)

type Rule struct {
	Type       string       `json:"type"`
	Conditions []*Condition `json:"conditions"`
}

type Condition struct {
	Columns  [2]*QColumn `json:"columns"` //предыдущий и текущий столбец таблицы.
	Operator string      `json:"operator"`
}

type QResponse struct {
	Data   any    `json:"data"`
	Err    string `json:"err"`
	RawSql string `json:"rawSql"`
}

func (r QResponse) Errorf(format string, a ...any) QResponse {
	r.Err = fmt.Sprintf(format, a...)
	return r
}

func (r QResponse) errExecute(err error) QResponse {
	return r.Errorf("не удалось исполнить запрос: %s", err.Error())
}

func (r QResponse) errParse(err error) QResponse {
	return r.Errorf("не удалось разобрать запрос: %s", err.Error())
}

func (db *Database) Execute(ctx context.Context, query Query) QResponse {
	switch query.Type {
	case Select:
		return db.executeSelect(ctx, query)
	case Insert:
		return db.executeInsert(ctx, query)
	case Update:
		return db.executeUpdate(ctx, query)
	case Delete:
		return db.executeDelete(ctx, query)

	default:
		return QResponse{}.Errorf("неизвестный тип команды %s", query.Type)
	}
}

const (
	specialPrefix     = "$"
	specialRootPrefix = specialPrefix + specialPrefix
	specialRootId     = specialPrefix + "root_id"
)

func (db *Database) executeSelect(ctx context.Context, query Query) QResponse {
	b, rules, err := db.parseSelect(query)
	if err != nil {
		return QResponse{}.errParse(err)
	}

	b = b.Columns(fmt.Sprintf(`COUNT(*) OVER() "%stotal"`, specialRootPrefix))

	var rows *sql.Rows

	if rows, err = b.QueryContext(ctx); err != nil {
		return QResponse{}.errExecute(err)
	}
	defer func() { _ = rows.Close() }()

	var (
		data    []map[string]any
		columns []string
		total   uint64
	)

	if columns, err = rows.Columns(); err != nil {
		return QResponse{}.errExecute(err)
	}

	for rows.Next() {
		dest, item := make([]any, 0), make(map[string]any)

		for _, c := range columns {
			if !strings.HasPrefix(c, specialRootPrefix) {
				item[c] = new(any)
				dest = append(dest, item[c])
			}
		}

		if err = rows.Scan(append(dest, &total)...); err != nil {
			return QResponse{}.errExecute(err)
		}

		data = append(data, item)
	}

	rawSql, _ := b.MustSql()

	return QResponse{
		Data: struct {
			Rules map[string][]string `json:"rules"`
			Data  []map[string]any    `json:"data"`
			Total uint64              `json:"total"`
		}{
			Rules: rules,
			Data:  data,
			Total: total,
		},
		RawSql: rawSql,
	}
}

func (db *Database) parseSelect(query Query) (sq.SelectBuilder, map[string][]string, error) {
	b := db.Builder.
		Select().
		From(query.Table.Full())

	var err error

	if b, err = query.Table.Join(b); err != nil {
		return b, nil, err
	}

	var (
		hasFunc bool
		pKey    *QColumn
		groupBy = make([]string, 0)
		rules   = make(map[string][]string)
	)

	for _, column := range query.Columns {
		if column.IsPKey && query.Table.QTableKey == column.TableKey {
			pKey = column
		}

		b = b.Columns(column.Full())

		rules[column.MetaKey()] = append(rules[column.MetaKey()], column.String())

		if len(column.Func) != 0 {
			hasFunc = true
			continue
		}

		groupBy = append(groupBy, column.Partial())
	}

	if pKey == nil && !hasFunc {
		var table *Table

		if table, err = db.GetTable(context.Background(), query.Table.Name); err != nil {
			return sq.SelectBuilder{}, nil, err
		}

		if newPKey, ok := table.GetPKey(); ok {
			pKey = &QColumn{Column: newPKey, TableKey: query.Table.QTableKey}
		}
	}

	if pKey != nil && !hasFunc {
		b = b.Columns(fmt.Sprintf(
			`%s "%s"`, pKey.Partial(), specialRootId,
		))
	}

	if hasFunc {
		b = b.GroupBy(groupBy...)
	}

	for _, column := range query.OrderBy {
		order, ok := column.Payload[Order]
		if !ok || (order != "ASC" && order != "DESC") {
			order = "ASC"
		}

		b = b.OrderBy(fmt.Sprintf("%s %s", column.Partial(), order))
	}

	for _, column := range query.Where {
		if column.Value != nil {
			b = b.Where(sq.Eq{column.Partial(): column.Value})
		}
	}

	if query.Limit != 0 {
		b = b.
			Limit(query.Limit).
			Offset(query.Offset)
	}

	return b, rules, nil
}

func (db *Database) executeInsert(ctx context.Context, query Query) QResponse {
	b, err := db.parseInsert(query)
	if err != nil {
		return QResponse{}.errParse(err)
	}

	if _, err = b.ExecContext(ctx); err != nil {
		return QResponse{}.errExecute(err)
	}

	rawSql, _ := b.MustSql()

	return QResponse{RawSql: rawSql}
}

func (db *Database) parseInsert(query Query) (sq.InsertBuilder, error) {
	b := db.Builder.Insert(query.Table.Partial())

	var values []any

	for _, column := range query.Columns {
		b = b.Columns(strconv.Quote(column.Name))
		values = append(values, column.Value)
	}

	return b.Values(values...), nil
}

func (db *Database) executeUpdate(ctx context.Context, query Query) QResponse {
	b, err := db.parseUpdate(query)
	if err != nil {
		return QResponse{}.errParse(err)
	}

	if _, err = b.ExecContext(ctx); err != nil {
		return QResponse{}.errExecute(err)
	}

	rawSql, _ := b.MustSql()

	return QResponse{RawSql: rawSql}
}

func (db *Database) parseUpdate(query Query) (sq.UpdateBuilder, error) {
	b := db.Builder.Update(query.Table.Partial())

	for _, column := range query.Columns {
		b = b.Set(strconv.Quote(column.Name), column.Value)
	}

	for _, column := range query.Where {
		if column.Value != nil {
			b = b.Where(sq.Eq{strconv.Quote(column.Name): column.Value})
		}
	}

	return b, nil
}

func (db *Database) executeDelete(ctx context.Context, query Query) QResponse {
	b, err := db.parseDelete(query)
	if err != nil {
		return QResponse{}.errParse(err)
	}

	if _, err = b.ExecContext(ctx); err != nil {
		return QResponse{}.errExecute(err)
	}

	rawSql, _ := b.MustSql()

	return QResponse{RawSql: rawSql}
}

func (db *Database) parseDelete(query Query) (sq.DeleteBuilder, error) {
	b := db.Builder.Delete(query.Table.Partial())

	for _, column := range query.Where {
		if column.Value != nil {
			b = b.Where(sq.Eq{strconv.Quote(column.Name): column.Value})
		}
	}

	return b, nil
}
