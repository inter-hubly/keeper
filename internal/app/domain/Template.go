package domain

type Template struct {
	Name            string `json:"name"`
	Category        string `json:"category"`
	ParameterFormat string `json:"parameter_format"`
	Language        string `json:"language"`
}

type Components struct {
	Type    string                   `json:"type"`
	Format  string                   `json:"format"`
	Text    string                   `json:"text"`
	Example map[string][]interface{} `json:"example"`
}
