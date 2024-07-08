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

func (s *QueryService) Execute(ctx context.Context, query entity.Query) database.QResponse {
	db, err := s.ss.GetDatabase(query.SourceId)
	if err != nil {
		return database.QResponse{}.Errorf(err.Error())
	}

	return db.Execute(ctx, query.Query)
}
