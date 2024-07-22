package entity

type DashboardWidget struct {
	*Widget
	X uint `json:"x"`
	Y uint `json:"y"`
	W uint `json:"w"`
	H uint `json:"h"`
}

type Dashboard struct {
	Id      string             `json:"id"`
	Name    string             `json:"name"`
	Widgets []*DashboardWidget `json:"widgets"`
}
