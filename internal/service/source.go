package service

import (
	"context"
	"datapointbackend/internal/entity"
	"datapointbackend/pkg/database"
	"fmt"
)

type sourceRepository interface {
	GetAll(ctx context.Context) ([]entity.Source, error)
	GetOne(ctx context.Context, id string) (entity.Source, error)
	Edit(ctx context.Context, s entity.Source) error
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, s entity.Source) (string, error)
}

type SourceService struct {
	sr      sourceRepository
	sources map[string]*database.Database
}

func NewSourceService(sr sourceRepository) *SourceService {
	s := SourceService{sources: make(map[string]*database.Database), sr: sr}

	sl, _ := s.sr.GetAll(context.Background())

	for _, source := range sl {
		s.sources[source.Id], _ = database.New(source.Config)
	}

	return &s
}

func (s *SourceService) GetAll(ctx context.Context) ([]entity.Source, error) {
	sl, err := s.sr.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить все источники: %s", err.Error())
	}

	for i := range sl {
		sl[i].Connected = s.IsConnected(sl[i].Id)
	}

	return sl, nil
}

func (s *SourceService) GetOne(ctx context.Context, id string) (entity.Source, error) {
	source, err := s.sr.GetOne(ctx, id)
	if err != nil {
		return entity.Source{}, fmt.Errorf("%s", err.Error())
	}

	source.Connected = s.IsConnected(source.Id)

	return source, nil
}

// Edit - предполагается, что схема базы данных не изменится.
func (s *SourceService) Edit(ctx context.Context, source entity.Source) error {
	db, err := s.GetDatabase(source.Id)
	if err != nil {
		return err
	}

	var newDb *database.Database

	if newDb, err = database.New(source.Config); err != nil {
		return fmt.Errorf("не удалось подключиться к источнику: %s", err.Error())
	}

	if err = s.sr.Edit(ctx, source); err != nil {
		_ = newDb.Conn.Close()
		return fmt.Errorf("не удалось отредактировать источник: %s", err.Error())
	}

	_ = db.Conn.Close()

	s.sources[source.Id] = newDb

	return nil
}

func (s *SourceService) Delete(ctx context.Context, id string) error {
	db, err := s.GetDatabase(id)
	if err != nil {
		return err
	}

	if err = s.sr.Delete(ctx, id); err != nil {
		return fmt.Errorf("не удалось удалить источник: %s", err.Error())
	}

	_ = db.Conn.Close()

	delete(s.sources, id)

	return nil
}

func (s *SourceService) Create(ctx context.Context, source entity.Source) (string, error) {
	db, err := database.New(source.Config)
	if err != nil {
		return "", fmt.Errorf("не удалось подключиться к источнику: %s", err.Error())
	}

	if source.Id, err = s.sr.Create(ctx, source); err != nil {
		_ = db.Conn.Close()
		return "", fmt.Errorf("не удалось сохранить конфигурацию источника: %s", err.Error())
	}

	s.sources[source.Id] = db

	return source.Id, nil
}

func (s *SourceService) GetDatabase(id string) (*database.Database, error) {
	source, ok := s.sources[id]
	if !ok {
		return nil, fmt.Errorf("отсутствует источник %s", id)
	}
	return source, nil
}

func (s *SourceService) IsConnected(id string) bool {
	db, err := s.GetDatabase(id)
	if err != nil {
		return false
	}

	if err = db.Conn.Ping(); err != nil {
		return false
	}

	return true
}
