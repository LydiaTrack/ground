package command

import (
	"errors"
	"time"
)

type CreateAuditCommand struct {
	Source           string                 `json:"source"`
	Operation        string                 `json:"operation"`
	AdditionalData   map[string]interface{} `json:"additional_data,omitempty"`
	RelatedPrincipal string                 `json:"related_principal,omitempty"`
}

type DeleteAuditCommand struct {
	ID string `json:"id"`
}

type DeleteOlderThanAuditCommand struct {
	Instant time.Time `json:"instant"`
}

func (dot DeleteOlderThanAuditCommand) Validate() error {
	// Instant cannot be in the future
	if dot.Instant.After(time.Now()) {
		return errors.New("instant cannot be in the future")
	}

	return nil
}

type DeleteIntervalAuditCommand struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func (di DeleteIntervalAuditCommand) Validate() error {
	// From cannot be in the future
	if di.From.After(time.Now()) {
		return errors.New("from cannot be in the future")
	}

	// From cannot be after to
	if di.From.After(di.To) {
		return errors.New("from cannot be after to")
	}

	return nil
}
