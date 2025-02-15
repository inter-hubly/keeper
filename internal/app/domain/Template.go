package domain

type TemplateType string

const (
	Header TemplateType = "HEADER"
	Body   TemplateType = "BODY"
	Footer TemplateType = "footer"
)

type Template struct {
	Id              string       `json:"id" bson:"_id"`
	Name            string       `json:"name" bson:"name"`
	Category        string       `json:"category" bson:"category"`
	ParameterFormat string       `json:"parameterFormat" bson:"parameter_format"`
	Language        string       `json:"language" bson:"language"`
	Components      []Components `json:"components" bson:"components"`
}

type Components struct {
	Type    TemplateType             `json:"type" bson:"type"`
	Format  string                   `json:"format" bson:"format"`
	Text    string                   `json:"text" bson:"text"`
	Example map[string][]interface{} `json:"example" bson:"example"`
}
