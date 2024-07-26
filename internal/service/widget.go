package service

import (
	"context"
	"datapointbackend/internal/entity"
	"fmt"
)

type widgetRepository interface {
	GetAll(ctx context.Context) ([]entity.Widget, error)
	GetOne(ctx context.Context, id string) (entity.Widget, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, w entity.Widget, parentId *string) (string, error)
	Edit(ctx context.Context, w entity.Widget) error
}

type WidgetService struct {
	wr widgetRepository
}

func NewWidgetService(wr widgetRepository) *WidgetService {
	return &WidgetService{wr: wr}
}

func (s *WidgetService) GetAll(ctx context.Context) ([]entity.Widget, error) {
	sl, err := s.wr.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить все виджеты: %s", err.Error())
	}
	return sl, nil
}

func (s *WidgetService) GetOne(ctx context.Context, id string) (entity.Widget, error) {
	widget, err := s.wr.GetOne(ctx, id)
	if err != nil {
		return entity.Widget{}, err
	}
	return widget, nil
}

func (s *WidgetService) Delete(ctx context.Context, id string) error {
	if err := s.wr.Delete(ctx, id); err != nil {
		return fmt.Errorf("не удалось удалить виджет: %s", err.Error())
	}
	return nil
}

func (s *WidgetService) Create(ctx context.Context, w entity.Widget) (string, error) {
	id, err := s.wr.Create(ctx, w, nil)
	if err != nil {
		return "", fmt.Errorf("не удалось сохранить виджет: %s", err.Error())
	}
	return id, nil
}

func (s *WidgetService) Edit(ctx context.Context, w entity.Widget) error {
	if err := s.wr.Edit(ctx, w); err != nil {
		return err
	}

	return nil
}
