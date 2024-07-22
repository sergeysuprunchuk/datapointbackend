package repository

import (
	"context"
	"datapointbackend/internal/entity"
	"datapointbackend/pkg/database"
	"fmt"
)

type DashboardRepository struct {
	db *database.Database
}

func NewDashboardRepository(db *database.Database) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetAll(ctx context.Context) ([]entity.Dashboard, error) {
	nilUuid := "00000000-0000-0000-0000-000000000000"

	rows, err := r.db.Builder.
		Select(
			"d.id d_id", "d.name d_name",
			"w.id w_id", "w.name w_name", "w.type", "w.parent_id", "w.props", "w.query",
			"dw.x", "dw.y", "dw.w", "dw.h",
		).From("dashboard d").
		LeftJoin("dashboard_widget dw ON dw.dashboard_id = d.id").
		LeftJoin("widget w ON w.id = dw.widget_id").
		Suffix("UNION").
		SuffixExpr(r.db.Builder.
			Select(
				"r.d_id", "r.d_name",
				"w.id", "w.name", "w.type", "w.parent_id", "w.props", "w.query",
				"0::SMALLINT", "0::SMALLINT", "0::SMALLINT", "0::SMALLINT",
			).
			From("widget w").
			Join("r ON r.w_id = w.parent_id"),
		).
		Prefix("WITH RECURSIVE r AS (").
		Suffix(")").
		SuffixExpr(r.db.Builder.
			Select(
				"r.d_id", "r.d_name",
				fmt.Sprintf("COALESCE(r.w_id, '%s')", nilUuid), "COALESCE(r.w_name, '')", "COALESCE(r.type, '')", "r.parent_id", "r.props", "r.query",
				"COALESCE(r.x, 0)", "COALESCE(r.y, 0)", "COALESCE(r.w, 0)", "COALESCE(r.h, 0)",
			).
			From("r")).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var (
		dashboards    []*entity.Dashboard
		dashboardDict = make(map[string]*entity.Dashboard)
		parentId      *string
		widgetDict    = make(map[string]map[string]*entity.Widget)
		widgets       []*extended
	)

	for rows.Next() {
		var (
			d entity.Dashboard
			w = entity.DashboardWidget{Widget: new(entity.Widget)}
		)

		if err = rows.Scan(
			&d.Id, &d.Name,
			&w.Id, &w.Name, &w.Type, &parentId, &w.Props, &w.Query,
			&w.X, &w.Y, &w.W, &w.H,
		); err != nil {
			return nil, err
		}

		if _, ok := dashboardDict[d.Id]; !ok {
			dashboardDict[d.Id] = &d
			dashboards = append(dashboards, &d)
		}

		if w.Id == nilUuid {
			continue
		}

		if widgetDict[d.Id] == nil {
			widgetDict[d.Id] = make(map[string]*entity.Widget)
		}

		widgetDict[d.Id][w.Id] = w.Widget

		if parentId == nil {
			dashboardDict[d.Id].Widgets = append(dashboardDict[d.Id].Widgets, &w)
		} else {
			widgets = append(widgets, &extended{Widget: w.Widget, parentId: parentId, dashboardId: d.Id})
		}
	}

	for _, w := range widgets {
		widgetDict[w.dashboardId][*w.parentId].Children =
			append(widgetDict[w.dashboardId][*w.parentId].Children, w.Widget)
	}

	var result []entity.Dashboard

	for _, d := range dashboards {
		result = append(result, *d)
	}

	return result, nil
}

func (r *DashboardRepository) GetOne(ctx context.Context, id string) (entity.Dashboard, error) {
	//TODO implement me
	panic("implement me")
}

func (r *DashboardRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Builder.
		Delete("dashboard").
		Where("id = ?", id).
		ExecContext(ctx)
	return err
}

func (r *DashboardRepository) Create(ctx context.Context, d entity.Dashboard) (string, error) {
	//в дальнейшем все это должно быть выполнено в транзакции.
	//!не могу использовать оператор WITH,
	//поскольку squirrel не увеличивает placeholder.

	err := r.db.Builder.
		Insert("dashboard").
		Columns("name").
		Values(d.Name).
		Suffix("RETURNING id").
		QueryRowContext(ctx).
		Scan(&d.Id)
	if err != nil {
		return "", err
	}

	if len(d.Widgets) != 0 {
		b := r.db.Builder.
			Insert("dashboard_widget").
			Columns("dashboard_id", "widget_id", "x", "y", "w", "h")

		for _, w := range d.Widgets {
			b = b.
				Values(d.Id, w.Id, w.X, w.Y, w.W, w.H)
		}

		if _, err = b.ExecContext(ctx); err != nil {
			return "", err
		}
	}

	return d.Id, nil
}
