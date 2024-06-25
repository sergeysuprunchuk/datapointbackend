package entity

import "datapointbackend/pkg/database"

type Source struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
	database.Config
}
