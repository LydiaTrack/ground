package service

import (
	"context"
	"github.com/LydiaTrack/ground/pkg/mongodb/repository"
	"github.com/LydiaTrack/ground/pkg/responses"
	"time"

	"github.com/LydiaTrack/ground/internal/permissions"
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/audit"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var auditSearchFields = []string{"source", "operation", "additionalData"}

type AuditService struct {
	auditRepository AuditRepository
}

func NewAuditService(auditRepository AuditRepository) AuditService {
	return AuditService{
		auditRepository: auditRepository,
	}
}

// AuditRepository defines the custom methods required in addition to the base repository methods.
type AuditRepository interface {
	repository.Repository[audit.Model]
	DeleteOlderThan(ctx context.Context, date time.Time) error
	DeleteInterval(ctx context.Context, from time.Time, to time.Time) error
}

// CreateAudit creates an audit record after permission validation.
func (s AuditService) CreateAudit(command audit.CreateAuditCommand, authContext auth.PermissionContext) (audit.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditCreatePermission) != nil {
		return audit.Model{}, constants.ErrorPermissionDenied
	}

	auditModel := audit.NewAudit(primitive.NewObjectID().Hex(), command.Source, command.Operation,
		time.Now(), audit.WithAdditionalData(command.AdditionalData), audit.WithRelatedPrincipal(command.RelatedPrincipal))
	createResult, err := s.auditRepository.Create(context.Background(), auditModel)
	if err != nil {
		return audit.Model{}, constants.ErrorInternalServerError
	}
	insertedId := createResult.InsertedID.(primitive.ObjectID)

	auditAfterSave, err := s.auditRepository.GetByID(context.Background(), insertedId)

	return auditAfterSave, nil
}

// GetAudit retrieves an audit by ID after permission validation.
func (s AuditService) GetAudit(id string, authContext auth.PermissionContext) (audit.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditReadPermission) != nil {
		return audit.Model{}, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return audit.Model{}, constants.ErrorBadRequest
	}

	auditModel, err := s.auditRepository.GetByID(context.Background(), objID)
	if err != nil {
		return audit.Model{}, constants.ErrorInternalServerError
	}
	return auditModel, nil
}

// ExistsAudit checks if an audit exists by ID after permission validation.
func (s AuditService) ExistsAudit(id string, authContext auth.PermissionContext) (bool, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditReadPermission) != nil {
		return false, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, constants.ErrorBadRequest
	}

	exists, err := s.auditRepository.ExistsByID(context.Background(), objID)
	if err != nil {
		return false, constants.ErrorInternalServerError
	}
	return exists, nil
}

// Query retrieves all audits after permission validation.
func (s AuditService) Query(searchText string, authContext auth.PermissionContext) (responses.QueryResult[audit.Model], error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditReadPermission) != nil {
		return responses.QueryResult[audit.Model]{}, constants.ErrorPermissionDenied
	}

	audits, err := s.auditRepository.Query(context.Background(), nil, auditSearchFields, searchText)
	if err != nil {
		return responses.QueryResult[audit.Model]{}, constants.ErrorInternalServerError
	}
	return audits, nil
}

// QueryPaginated retrieves all audits in a paginated manner after permission validation.
func (s AuditService) QueryPaginated(searchText string, page, limit int, authContext auth.PermissionContext) (repository.PaginatedResult[audit.Model], error) {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditReadPermission) != nil {
		return repository.PaginatedResult[audit.Model]{}, constants.ErrorPermissionDenied
	}

	result, err := s.auditRepository.QueryPaginate(context.Background(), nil, auditSearchFields, searchText, page, limit, primitive.M{"instant": 1})
	if err != nil {
		return repository.PaginatedResult[audit.Model]{}, constants.ErrorInternalServerError
	}
	return result, nil
}

// DeleteOlderThan deletes audits older than a given date after permission validation.
func (s AuditService) DeleteOlderThan(command audit.DeleteOlderThanAuditCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditDeletePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	err := command.Validate()
	if err != nil {
		return constants.ErrorBadRequest
	}

	err = s.auditRepository.DeleteOlderThan(context.Background(), command.Instant)
	return err
}

// DeleteInterval deletes audits in a specific time interval after permission validation.
func (s AuditService) DeleteInterval(command audit.DeleteIntervalAuditCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.AuditDeletePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	err := command.Validate()
	if err != nil {
		return constants.ErrorBadRequest
	}

	err = s.auditRepository.DeleteInterval(context.Background(), command.From, command.To)
	return err
}
