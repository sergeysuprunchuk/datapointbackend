package entity

import "datapointbackend/pkg/database"

type Query struct {
	SourceId string `json:"sourceId"`
	database.Query
}
