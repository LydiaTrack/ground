package repository

import (
	"context"
	"time"

	"github.com/LydiaTrack/lydia-base/pkg/domain/audit"
	"github.com/LydiaTrack/lydia-base/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuditMongoRepository struct {
	collection *mongo.Collection
}

var (
	auditRepository *AuditMongoRepository
)

func newAuditRepository() *AuditMongoRepository {
	collection, err := mongodb.GetCollection("audits")
	if err != nil {
		panic(err)
	}

	return &AuditMongoRepository{
		collection: collection,
	}
}

// GetUserRepository returns a UserRepository instance if it is not initialized yet or
// returns the existing one
func GetAuditRepository() *AuditMongoRepository {
	if auditRepository == nil {
		auditRepository = newAuditRepository()
	}
	return auditRepository
}

// SaveAudit saves an audit
func (r AuditMongoRepository) SaveAudit(audit audit.Model) (audit.Model, error) {
	_, err := r.collection.InsertOne(context.Background(), audit)
	if err != nil {
		return audit, err
	}

	return audit, nil
}

// GetAudit gets an audit by id
func (r AuditMongoRepository) GetAudit(id primitive.ObjectID) (audit.Model, error) {
	var audit audit.Model
	err := r.collection.FindOne(context.Background(), primitive.M{"_id": id}).Decode(&audit)
	if err != nil {
		return audit, err
	}

	return audit, nil
}

// ExistsAudit checks if an audit exists
func (r AuditMongoRepository) ExistsAudit(id primitive.ObjectID) (bool, error) {
	var audit audit.Model
	err := r.collection.FindOne(context.Background(), primitive.M{"_id": id}).Decode(&audit)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetAudits gets all audits
func (r AuditMongoRepository) GetAudits() ([]audit.Model, error) {
	var audits []audit.Model
	cursor, err := r.collection.Find(context.Background(), primitive.M{})
	if err != nil {
		return audits, err
	}

	err = cursor.All(context.Background(), &audits)
	if err != nil {
		return audits, err
	}

	return audits, nil
}

// DeleteOlderThan deletes all audits older than a date
func (r AuditMongoRepository) DeleteOlderThan(date time.Time) error {
	_, err := r.collection.DeleteMany(context.Background(), primitive.M{"instant": primitive.M{"$lt": date}})
	if err != nil {
		return err
	}

	return nil
}

// DeleteInterval deletes all audits between two dates
func (r AuditMongoRepository) DeleteInterval(from time.Time, to time.Time) error {
	_, err := r.collection.DeleteMany(context.Background(), primitive.M{"instant": primitive.M{"$gte": from, "$lte": to}})
	if err != nil {
		return err
	}

	return nil
}
