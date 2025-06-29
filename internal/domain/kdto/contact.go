package kdto

import "github.com/inter-hubly/pilot/domain/valueobject"

type Contact struct {
	Name      string                             `json:"name"`
	Phone     string                             `json:"phone"`
	Variables []valueobject.Pair[string, string] `json:"variables"`
}
