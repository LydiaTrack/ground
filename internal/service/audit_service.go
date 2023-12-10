package service

import (
	"errors"
	"lydia-track-base/internal/domain/audit"
	"lydia-track-base/internal/domain/audit/command"
	"lydia-track-base/internal/domain/auth"
	"time"

	"gopkg.in/mgo.v2/bson"
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
func (s AuditService) CreateAudit(command commands.CreateAuditCommand, permissions []auth.Permission) (audit.Model, error) {
	if !CheckPermission(permissions, commands.CreatePermission) {
		return audit.Model{}, errors.New("not permitted")
	}

	auditModel := audit.NewAudit(bson.NewObjectId().Hex(), command.Source, command.Operation,
		time.Now(), audit.WithAdditionalData(command.AdditionalData), audit.WithRelatedPrincipal(command.RelatedPrincipal))
	return s.auditRepository.SaveAudit(auditModel)
}

// GetAudit TODO: Add permission check
func (s AuditService) GetAudit(id string, permissions []auth.Permission) (audit.Model, error) {
	if !CheckPermission(permissions, commands.ReadPermission) {
		return audit.Model{}, errors.New("not permitted")
	}

	auditModel, error := s.auditRepository.GetAudit(bson.ObjectIdHex(id))
	return auditModel, error
}

// ExistsAudit TODO: Add permission check
func (s AuditService) ExistsAudit(id string, permissions []auth.Permission) (bool, error) {
	if !CheckPermission(permissions, commands.ReadPermission) {
		return false, errors.New("not permitted")
	}

	exists, error := s.auditRepository.ExistsAudit(bson.ObjectIdHex(id))
	return exists, error
}

// GetAudits TODO: Add permission check
func (s AuditService) GetAudits(permissions []auth.Permission) ([]audit.Model, error) {
	if !CheckPermission(permissions, commands.ReadPermission) {
		return nil, errors.New("not permitted")
	}

	audits, error := s.auditRepository.GetAudits()
	return audits, error
}

// DeleteOlderThan TODO: Add permission check
func (s AuditService) DeleteOlderThan(command commands.DeleteOlderThanAuditCommand, permissions []auth.Permission) error {
	if !CheckPermission(permissions, commands.DeletePermission) {
		return errors.New("not permitted")
	}

	error := command.Validate()
	if error != nil {
		return error
	}

	error = s.auditRepository.DeleteOlderThan(command.Instant)
	return error
}

// DeleteInterval TODO: Add permission check
func (s AuditService) DeleteInterval(command commands.DeleteIntervalAuditCommand, permissions []auth.Permission) error {
	if !CheckPermission(permissions, commands.DeletePermission) {
		return errors.New("not permitted")
	}

	error := command.Validate()
	if error != nil {
		return error
	}

	error = s.auditRepository.DeleteInterval(command.From, command.To)
	return error
}
