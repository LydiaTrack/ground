package audit

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Option func(m Model) Model

type Model struct {
	ID     bson.ObjectId `bson:"_id"`
	Source string        `bson:"source"`
	// This field should be a struct, but I don't know how to do it yet.
	Operation        string                 `bson:"operation"`
	Instant          time.Time              `bson:"instant"`
	AdditionalData   map[string]interface{} `bson:"additional_data,omitempty"`
	RelatedPrincipal string                 `bson:"related_principal,omitempty"`
}

/**
 * NewAudit creates a new audit with the given parameters and returns it. It also accepts a variadic number of options
 * to customize the audit. (functional options pattern)
 */
func NewAudit(id string, source string, operation string, instant time.Time, opts ...Option) Model {
	m := Model{
		ID:        bson.ObjectIdHex(id),
		Source:    source,
		Operation: operation,
		Instant:   instant,
	}

	for _, opt := range opts {
		m = opt(m)
	}

	return m
}

func WithAdditionalData(additionalData map[string]interface{}) Option {
	return func(m Model) Model {
		m.AdditionalData = additionalData
		return m
	}
}

func WithRelatedPrincipal(relatedPrincipal string) Option {
	return func(m Model) Model {
		m.RelatedPrincipal = relatedPrincipal
		return m
	}
}
