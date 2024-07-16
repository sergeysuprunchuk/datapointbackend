package entity

import "encoding/json"

type Widget struct {
	Id       string           `json:"id"`
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	Children []*Widget        `json:"children"`
	Props    *json.RawMessage `json:"props"`
	//чтобы ускорить разработку пока так, а дальше посмотрим :/.
	Query *json.RawMessage `json:"query"`
}
