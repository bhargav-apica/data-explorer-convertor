package models

type GrafanaDashboard struct {
	Panels []Panel `json:"panels"`
}

type Panel struct {
	Title   string   `json:"title"`
	Targets []Target `json:"targets"`
	IsLines *bool    `json:"lines"`
	Panels  *[]Panel `json:"panels"`
}

type Target struct {
	QueryExpr *string `json:"expr"`
	Legend    *string `json:"legendFormat"`
}
