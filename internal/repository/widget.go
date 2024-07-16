package repository

import (
	"context"
	"database/sql"
	"datapointbackend/internal/entity"
	"datapointbackend/pkg/database"
	"fmt"
)

type WidgetRepository struct {
	db *database.Database
}

func NewWidgetRepository(db *database.Database) *WidgetRepository {
	return &WidgetRepository{db: db}
}

type extended struct {
	*entity.Widget
	ParentId *string
}

func (r *WidgetRepository) GetAll(ctx context.Context) ([]entity.Widget, error) {
	rows, err := r.db.Builder.
		Select("id", "name", "type", "props", "query", "parent_id").
		From("widget").
		QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var (
		widgets  = make([]*extended, 0)
		pointers = make(map[string]*entity.Widget, 0)
	)

	for rows.Next() {
		e := extended{Widget: new(entity.Widget)}

		if err = rows.Scan(&e.Id, &e.Name, &e.Type, &e.Props, &e.Query, &e.ParentId); err != nil {
			return nil, err
		}

		pointers[e.Id] = e.Widget

		widgets = append(widgets, &e)
	}

	for _, w := range widgets {
		if w.ParentId != nil {
			pointers[*w.ParentId].Children = append(pointers[*w.ParentId].Children, w.Widget)
		}
	}

	var sl []entity.Widget

	for _, w := range widgets {
		if w.ParentId == nil {
			sl = append(sl, *w.Widget)
		}
	}

	return sl, nil
}

func (r *WidgetRepository) GetOne(ctx context.Context, id string) (entity.Widget, error) {
	anchorSql, anchorArgs, err := r.db.Builder.
		Select("id", "name", "type", "props", "query", "parent_id").
		From("widget").
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return entity.Widget{}, err
	}

	var recursiveSql string
	if recursiveSql, _, err = r.db.Builder.
		Select("widget.id", "widget.name", "widget.type", "widget.props", "widget.query", "widget.parent_id").
		From("widget").
		Join("recursive ON widget.parent_id = recursive.id").
		ToSql(); err != nil {
		return entity.Widget{}, err
	}

	tailSql, _, _ := r.db.Builder.
		Select("id", "name", "type", "props", "query", "parent_id").
		From("recursive").
		ToSql()

	var rows *sql.Rows

	if rows, err = r.db.Conn.QueryContext(
		ctx,
		fmt.Sprintf("WITH RECURSIVE recursive AS (%s UNION %s) %s", anchorSql, recursiveSql, tailSql),
		anchorArgs...,
	); err != nil {
		return entity.Widget{}, err
	}

	defer func() { _ = rows.Close() }()

	var (
		widgets  = make([]*extended, 0)
		pointers = make(map[string]*entity.Widget, 0)
	)

	for rows.Next() {
		e := extended{Widget: new(entity.Widget)}

		if err = rows.Scan(&e.Id, &e.Name, &e.Type, &e.Props, &e.Query, &e.ParentId); err != nil {
			return entity.Widget{}, err
		}

		pointers[e.Id] = e.Widget

		widgets = append(widgets, &e)
	}

	for _, w := range widgets {
		if w.ParentId != nil {
			pointers[*w.ParentId].Children = append(pointers[*w.ParentId].Children, w.Widget)
		}
	}

	widget, ok := pointers[id]
	if !ok {
		return entity.Widget{}, fmt.Errorf("виджет с идентификатором %s не существует", id)
	}

	return *widget, nil
}

func (r *WidgetRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Builder.
		Delete("widget").
		Where("id = ?", id).
		ExecContext(ctx)
	return err
}

func (r *WidgetRepository) Create(ctx context.Context, w entity.Widget, parentId *string) (string, error) {
	//в дальнейшем нужно будет все это выполнять в транзакции.
	err := r.db.Builder.
		Insert("widget").
		Columns("name", "type", "props", "query", "parent_id").
		Values(w.Name, w.Type, w.Props, w.Query, parentId).
		Suffix("RETURNING id").
		QueryRowContext(ctx).
		Scan(&w.Id)
	if err != nil {
		return "", err
	}

	for _, child := range w.Children {
		if child == nil {
			return "", fmt.Errorf("неверный дочерний виджет")
		}

		if _, err = r.Create(ctx, *child, &w.Id); err != nil {
			return "", err
		}
	}

	return w.Id, nil
}

func (r *WidgetRepository) Edit(ctx context.Context, w entity.Widget) error {
	//TODO implement me
	panic("implement me")
}
