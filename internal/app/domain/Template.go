package domain

type TemplateType string

const (
	Header TemplateType = "HEADER"
	Body   TemplateType = "BODY"
	Footer TemplateType = "footer"
)

type Template struct {
	Id              string `bson:"id"`
	Name            string `bson:"name"`
	Category        string `bson:"category"`
	ParameterFormat string `bson:"parameter_format"`
	Language        string `bson:"language"`
}

type Components struct {
	Type    TemplateType             `bson:"type"`
	Format  string                   `bson:"format"`
	Text    string                   `bson:"text"`
	Example map[string][]interface{} `bson:"example"`
}
