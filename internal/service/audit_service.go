package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/LydiaTrack/lydia-base/internal/permissions"
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/constants"
	"github.com/LydiaTrack/lydia-base/pkg/domain/audit"
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
	GetAudit(id primitive.ObjectID) (audit.Model, error)
	// ExistsAudit checks if an audit exists
	ExistsAudit(id primitive.ObjectID) (bool, error)
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
		return audit.Model{}, constants.ErrorPermissionDenied
	}

	auditModel := audit.NewAudit(primitive.NewObjectID().Hex(), command.Source, command.Operation,
		time.Now(), audit.WithAdditionalData(command.AdditionalData), audit.WithRelatedPrincipal(command.RelatedPrincipal))
	auditModel, err := s.auditRepository.SaveAudit(auditModel)
	if err != nil {
		return audit.Model{}, constants.ErrorInternalServerError
	}
	return auditModel, nil
}

// GetAudit TODO: Add permission check
func (s AuditService) GetAudit(id string, authContext auth.PermissionContext) (audit.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditReadPermission) != nil {
		return audit.Model{}, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return audit.Model{}, constants.ErrorBadRequest
	}
	auditModel, err := s.auditRepository.GetAudit(objID)
	if err != nil {
		return audit.Model{}, constants.ErrorInternalServerError
	}
	return auditModel, nil
}

// ExistsAudit TODO: Add permission check
func (s AuditService) ExistsAudit(id string, authContext auth.PermissionContext) (bool, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditReadPermission) != nil {
		return false, constants.ErrorPermissionDenied
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, constants.ErrorBadRequest
	}
	exists, err := s.auditRepository.ExistsAudit(objID)
	if err != nil {
		return false, constants.ErrorInternalServerError
	}
	return exists, nil
}

// GetAudits TODO: Add permission check
func (s AuditService) GetAudits(authContext auth.PermissionContext) ([]audit.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditReadPermission) != nil {
		return nil, constants.ErrorPermissionDenied
	}

	audits, err := s.auditRepository.GetAudits()
	if err != nil {
		return nil, constants.ErrorInternalServerError
	}
	return audits, nil
}

// DeleteOlderThan TODO: Add permission check
func (s AuditService) DeleteOlderThan(command audit.DeleteOlderThanAuditCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditDeletePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	err := command.Validate()
	if err != nil {
		return constants.ErrorBadRequest
	}

	err = s.auditRepository.DeleteOlderThan(command.Instant)
	return err
}

// DeleteInterval TODO: Add permission check
func (s AuditService) DeleteInterval(command audit.DeleteIntervalAuditCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditDeletePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	err := command.Validate()
	if err != nil {
		return constants.ErrorBadRequest
	}

	err = s.auditRepository.DeleteInterval(command.From, command.To)
	return err
}
