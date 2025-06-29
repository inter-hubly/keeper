package kdto

type Variable struct {
	Label string `json:"label"`
	Type  string `json:"type"`
	Slug  string `json:"slug,omitempty"`
}
