package domain

import "github.com/inter-hubly/pilot/domain/base"

type TemplateType string

const (
	Header TemplateType = "HEADER"
	Body   TemplateType = "BODY"
	Footer TemplateType = "FOOTER"
)

type Template struct {
	Id              string       `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string       `json:"name" bson:"name"`
	Category        string       `json:"category" bson:"category"`
	ParameterFormat string       `json:"parameterFormat" bson:"parameter_format,omitempty"`
	Language        string       `json:"language" bson:"language"`
	Status          string       `json:"status" bson:"status"`
	ResponseId      string       `json:"responseId" bson:"response_id"`
	Components      []Components `json:"components" bson:"components"`
	base.Entity     `json:"-,inline" bson:"-,inline"`
}

type Components struct {
	Type    TemplateType             `json:"type" bson:"type"`
	Format  string                   `json:"format" bson:"format"`
	Text    string                   `json:"text" bson:"text"`
	Example map[string][]interface{} `json:"example" bson:"example"`
}
