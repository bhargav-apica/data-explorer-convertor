package models

type Tab struct {
	Key         string      `json:"key"`
	Order       int         `json:"order"`
	QueriesList []QueryItem `json:"queriesList"`
	Title       string      `json:"title"`
	Type        string      `json:"type"`
}

type QueryItem struct {
	Name      string       `json:"name"`
	Options   QueryOptions `json:"options"`
	Query     string       `json:"query"`
	Schema    string       `json:"schema"`
	ChartType *string      `json:"chart_type,omitempty"`
}

type QueryOptions struct {
	Description string    `json:"description"`
	Order       int       `json:"order"`
	Plot        QueryPlot `json:"plot"`
	UpperLimit  string    `json:"upperLimit"`
}

type QueryPlot struct {
	ErrorColumn string   `json:"errorColumn"`
	GroupBy     string   `json:"groupBy"`
	X           string   `json:"x"`
	XLabel      string   `json:"xLabel"`
	Y           []string `json:"y"`
	YLabel      string   `json:"yLabel"`
}

func NewQueryPlot() QueryPlot {
	return QueryPlot{
		ErrorColumn: "",
		GroupBy:     "",
		X:           "Timestamp",
		XLabel:      "Timestamp",
		Y:           []string{"Timestamp"},
		YLabel:      "value",
	}
}
