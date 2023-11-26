package service

import (
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/domain/audit"
	"lydia-track-base/internal/domain/audit/command"
	"time"
)

type AuditService struct {
	auditRepository AuditRepository
}

func NewAuditService(auditRepository AuditRepository) AuditService {
	return AuditService{
		auditRepository: auditRepository,
	}
}

type AuditRepository interface {
	// SaveAudit saves an audit
	SaveAudit(audit audit.Model) (audit.Model, error)
	// GetAudit gets an audit by id
	GetAudit(id bson.ObjectId) (audit.Model, error)
	// ExistsAudit checks if an audit exists
	ExistsAudit(id bson.ObjectId) (bool, error)
	// GetAudits gets all audits
	GetAudits() ([]audit.Model, error)
	// DeleteOlderThan deletes all audits older than a date
	DeleteOlderThan(date time.Time) error
	// DeleteInterval deletes all audits between two dates
	DeleteInterval(from time.Time, to time.Time) error
}

// CreateAudit TODO: Add permission check
func (s AuditService) CreateAudit(command command.CreateAuditCommand) (audit.Model, error) {
	auditModel := audit.NewAudit(bson.NewObjectId().Hex(), command.Source, command.Operation,
		time.Now(), audit.WithAdditionalData(command.AdditionalData), audit.WithRelatedPrincipal(command.RelatedPrincipal))
	return s.auditRepository.SaveAudit(auditModel)
}

// GetAudit TODO: Add permission check
func (s AuditService) GetAudit(id string) (audit.Model, error) {
	auditModel, error := s.auditRepository.GetAudit(bson.ObjectIdHex(id))
	return auditModel, error
}

// ExistsAudit TODO: Add permission check
func (s AuditService) ExistsAudit(id string) (bool, error) {
	exists, error := s.auditRepository.ExistsAudit(bson.ObjectIdHex(id))
	return exists, error
}

// GetAudits TODO: Add permission check
func (s AuditService) GetAudits() ([]audit.Model, error) {
	audits, error := s.auditRepository.GetAudits()
	return audits, error
}

// DeleteOlderThan TODO: Add permission check
func (s AuditService) DeleteOlderThan(command command.DeleteOlderThanAuditCommand) error {
	error := command.Validate()
	if error != nil {
		return error
	}

	error = s.auditRepository.DeleteOlderThan(command.Instant)
	return error
}

// DeleteInterval TODO: Add permission check
func (s AuditService) DeleteInterval(command command.DeleteIntervalAuditCommand) error {
	error := command.Validate()
	if error != nil {
		return error
	}

	error = s.auditRepository.DeleteInterval(command.From, command.To)
	return error
}
