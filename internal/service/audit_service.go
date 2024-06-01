package service

import (
	"errors"
	"github.com/LydiaTrack/lydia-base/internal/permissions"
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/domain/audit"
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
func (s AuditService) CreateAudit(command audit.CreateAuditCommand, authContext auth.PermissionContext) (audit.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditCreatePermission) != nil {
		return audit.Model{}, errors.New("not permitted")
	}

	auditModel := audit.NewAudit(bson.NewObjectId().Hex(), command.Source, command.Operation,
		time.Now(), audit.WithAdditionalData(command.AdditionalData), audit.WithRelatedPrincipal(command.RelatedPrincipal))
	return s.auditRepository.SaveAudit(auditModel)
}

// GetAudit TODO: Add permission check
func (s AuditService) GetAudit(id string, authContext auth.PermissionContext) (audit.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditReadPermission) != nil {
		return audit.Model{}, errors.New("not permitted")
	}

	auditModel, err := s.auditRepository.GetAudit(bson.ObjectIdHex(id))
	return auditModel, err
}

// ExistsAudit TODO: Add permission check
func (s AuditService) ExistsAudit(id string, authContext auth.PermissionContext) (bool, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditReadPermission) != nil {
		return false, errors.New("not permitted")
	}

	exists, err := s.auditRepository.ExistsAudit(bson.ObjectIdHex(id))
	return exists, err
}

// GetAudits TODO: Add permission check
func (s AuditService) GetAudits(authContext auth.PermissionContext) ([]audit.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditReadPermission) != nil {
		return nil, errors.New("not permitted")
	}

	audits, err := s.auditRepository.GetAudits()
	return audits, err
}

// DeleteOlderThan TODO: Add permission check
func (s AuditService) DeleteOlderThan(command audit.DeleteOlderThanAuditCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditDeletePermission) != nil {
		return errors.New("not permitted")
	}

	err := command.Validate()
	if err != nil {
		return err
	}

	err = s.auditRepository.DeleteOlderThan(command.Instant)
	return err
}

// DeleteInterval TODO: Add permission check
func (s AuditService) DeleteInterval(command audit.DeleteIntervalAuditCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditDeletePermission) != nil {
		return errors.New("not permitted")
	}

	err := command.Validate()
	if err != nil {
		return err
	}

	err = s.auditRepository.DeleteInterval(command.From, command.To)
	return err
}
