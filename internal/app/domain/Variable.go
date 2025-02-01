package domain

import "github.com/inter-hubly/pilot/domain/base"

type Variables struct {
	Id          string           `bson:"_id,omitempty"`
	Variable    []SingleVariable `bson:"variable"`
	base.Entity `bson:",inline"`
}

type SingleVariable struct {
	Slug  string `bson:"slug"`
	Label string `bson:"label"`
	Type  string `bson:"type"`
}
