package kdto

import (
	"github.com/inter-hubly/pilot/domain/valueobject"
)

type Campaign struct {
	Name       string                             `json:"name"`
	TemplateId string                             `json:"templateId"`
	ContactsID []string                           `json:"contactsIds"`
	Variables  []valueobject.Pair[string, string] `json:"variables"`
}
