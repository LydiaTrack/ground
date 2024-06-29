package audit

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Option func(m Model) Model

type Model struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Source           string             `json:"source" bson:"source"`
	Operation        `json:"operation" bson:"operation"`
	Instant          time.Time              `json:"instant" bson:"instant"`
	AdditionalData   map[string]interface{} `json:"additionalData,omitempty" bson:"additionalData,omitempty"`
	RelatedPrincipal string                 `json:"relatedPrincipal,omitempty" bson:"relatedPrincipal,omitempty"`
}

type Operation struct {
	Domain  string `json:"domain" bson:"domain"`
	Command string `json:"command" bson:"command"`
}

/**
 * NewAudit creates a new audit with the given parameters and returns it. It also accepts a variadic number of options
 * to customize the audit. (functional options pattern)
 */
func NewAudit(id string, source string, operation Operation, instant time.Time, opts ...Option) Model {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Model{}
	}
	m := Model{
		ID:        objID,
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
