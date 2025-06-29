package domain

import (
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/valueobject"
)

type Contact struct {
	Id          string                             `json:"id" bson:"_id,omitempty"`
	Name        string                             `json:"name" bson:"name"`
	Phone       string                             `json:"phone" bson:"phone"`
	Variables   []valueobject.Pair[string, string] `json:"variables" bson:"variables"`
	base.Entity `bson:",inline"`
}
