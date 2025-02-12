package kdto

import (
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/valueobject"
)

type Campaign struct {
	Name       string                             `json:"name"`
	Template   base.TemplateInfo                  `json:"template"`
	ContactsID []string                           `json:"contactsId"`
	Variables  []valueobject.Pair[string, string] `json:"variables"`
}
