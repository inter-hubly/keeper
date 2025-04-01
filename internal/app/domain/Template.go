package domain

import "github.com/inter-hubly/pilot/domain/base"

type TemplateType string

const (
	Header TemplateType = "HEADER"
	Body   TemplateType = "BODY"
	Footer TemplateType = "FOOTER"
)

type Template struct {
	Id              string         `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string         `json:"name" bson:"name"`
	Category        string         `json:"category" bson:"category"`
	ParameterFormat string         `json:"parameterFormat" bson:"parameterFormat,omitempty"`
	Language        string         `json:"language" bson:"language"`
	Status          TemplateStatus `json:"status" bson:"status"`
	ResponseId      string         `json:"responseId" bson:"response_id"`
	Components      []Components   `json:"components" bson:"components"`
	base.Entity     `json:"-,inline" bson:"-,inline"`
}

type TemplateStatus string

const (
	Approved TemplateStatus = "APPROVED"
	Rejected TemplateStatus = "REJECTED"
)

type Components struct {
	Type    TemplateType          `json:"type" bson:"type"`
	Format  string                `json:"format" bson:"format"`
	Text    string                `json:"text" bson:"text"`
	Example map[string][][]string `json:"example,omitempty" bson:"example,omitempty"`
}

func (t *Template) GetComponentMessages() string {
	var resp string
	for i := range t.Components {
		resp += t.Components[i].Text + " "
	}
	return resp
}
