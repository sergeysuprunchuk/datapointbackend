package service

import (
	"context"
	"datapointbackend/internal/entity"
	"fmt"
)

type dashboardRepository interface {
	GetAll(ctx context.Context) ([]entity.Dashboard, error)
	GetOne(ctx context.Context, id string) (entity.Dashboard, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, d entity.Dashboard) (string, error)
	Edit(ctx context.Context, d entity.Dashboard) error
}

type DashboardService struct {
	dr dashboardRepository
}

func NewDashboardService(dr dashboardRepository) *DashboardService {
	return &DashboardService{dr: dr}
}

func (s *DashboardService) GetAll(ctx context.Context) ([]entity.Dashboard, error) {
	sl, err := s.dr.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить все дашборды: %s", err.Error())
	}
	return sl, nil
}

func (s *DashboardService) GetOne(ctx context.Context, id string) (entity.Dashboard, error) {
	d, err := s.dr.GetOne(ctx, id)
	if err != nil {
		return entity.Dashboard{}, err
	}
	return d, nil
}

func (s *DashboardService) Delete(ctx context.Context, id string) error {
	if err := s.dr.Delete(ctx, id); err != nil {
		return fmt.Errorf("не удалось удалить дашборд: %s", err.Error())
	}
	return nil
}

func (s *DashboardService) Create(ctx context.Context, d entity.Dashboard) (string, error) {
	id, err := s.dr.Create(ctx, d)
	if err != nil {
		return "", fmt.Errorf("не удалось сохранить дашборд: %s", err.Error())
	}
	return id, nil
}

func (s *DashboardService) Edit(ctx context.Context, d entity.Dashboard) error {
	if err := s.dr.Edit(ctx, d); err != nil {
		return err
	}
	return nil
}
