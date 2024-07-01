package service

import (
	"context"
	"datapointbackend/internal/entity"
	"datapointbackend/pkg/database"
)

type QueryService struct {
	ss *SourceService
}

func NewQueryService(ss *SourceService) *QueryService {
	return &QueryService{ss: ss}
}

func (s *QueryService) Execute(ctx context.Context, query entity.Query) (database.QueryResponse, error) {
	db, err := s.ss.GetDatabase(query.SourceId)
	if err != nil {
		return database.QueryResponse{}, err
	}

	var response database.QueryResponse

	if response, err = db.Execute(ctx, query.Query); err != nil {
		return database.QueryResponse{}, err
	}

	return response, nil
}
