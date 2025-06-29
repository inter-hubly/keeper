package kdto

import (
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/domain/valueobject"
)

type Campaign struct {
	Name       string                             `json:"name"`
	TemplateId string                             `json:"templateId"`
	IaContext  string                             `json:"iaContext"`
	ContactsID []string                           `json:"contactsIds"`
	Variables  []valueobject.Pair[string, string] `json:"variables"`
	Flows      map[string]*entity.Flow            `json:"flows"`
}
