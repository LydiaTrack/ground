package repository

import (
	"context"
	"github.com/LydiaTrack/ground/pkg/mongodb"
	"time"

	"github.com/LydiaTrack/ground/pkg/domain/audit"
	"github.com/LydiaTrack/ground/pkg/mongodb/repository"
	"go.mongodb.org/mongo-driver/bson"
)

// AuditMongoRepository extends the generic BaseRepository with custom methods for the Audit model.
type AuditMongoRepository struct {
	*repository.BaseRepository[audit.Model]
}

// GetAuditRepository creates a new instance of AuditMongoRepository.
func GetAuditRepository() *AuditMongoRepository {
	collection, err := mongodb.GetCollection("audits")
	if err != nil {
		panic(err)
	}
	return &AuditMongoRepository{
		BaseRepository: repository.NewBaseRepository[audit.Model](collection),
	}
}

// DeleteOlderThan deletes all audits older than a specific date.
func (r *AuditMongoRepository) DeleteOlderThan(ctx context.Context, date time.Time) error {
	_, err := r.Collection.DeleteMany(ctx, bson.M{"instant": bson.M{"$lt": date}})
	return err
}

// DeleteInterval deletes all audits between two dates.
func (r *AuditMongoRepository) DeleteInterval(ctx context.Context, from, to time.Time) error {
	_, err := r.Collection.DeleteMany(ctx, bson.M{"instant": bson.M{"$gte": from, "$lte": to}})
	return err
}
