package kdto

import "github.com/inter-hubly/pilot/domain/valueobject"

type Campaign struct {
	Name       string                             `json:"name"`
	Template   Template                           `json:"template"`
	ContactsID []string                           `json:"contactsId"`
	Variables  []valueobject.Pair[string, string] `json:"variables"`
}

type Template struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}
