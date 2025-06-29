package kdto

type WhatsAppTemplateResponse struct {
	Id       string `json:"id"`
	Status   string `json:"status"`
	Category string `json:"category"`
}

type WhatsAppMessageTemplateResponse struct {
	Name                     string         `json:"name"`
	Language                 string         `json:"language"`
	Components               []ComponentDto `json:"components"`
	ParameterFormat          string         `json:"parameter_format"`
	WhatsAppTemplateResponse `json:",inline"`
}

type ComponentDto struct {
	Type    string         `json:"type" bson:"type"`
	Format  string         `json:"format" bson:"format"`
	Text    *string        `json:"text" bson:"text"`
	Example map[string]any `json:"example,omitempty" bson:"example,omitempty"`
}
