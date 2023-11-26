package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/domain/audit"
	"lydia-track-base/internal/mongodb"
	"time"
)

type AuditMongoRepository struct {
	collection *mongo.Collection
}

var (
	auditRepository *AuditMongoRepository
)

func newAuditRepository() *AuditMongoRepository {
	ctx := context.Background()
	// FIXME: Burada ileride uzaktaki bir mongodb instance'ına bağlanmak gerekecek
	collection := mongodb.GetCollection("audits", ctx)

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
func (r AuditMongoRepository) GetAudit(id bson.ObjectId) (audit.Model, error) {
	var audit audit.Model
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&audit)
	if err != nil {
		return audit, err
	}

	return audit, nil
}

// ExistsAudit checks if an audit exists
func (r AuditMongoRepository) ExistsAudit(id bson.ObjectId) (bool, error) {
	var audit audit.Model
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&audit)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetAudits gets all audits
func (r AuditMongoRepository) GetAudits() ([]audit.Model, error) {
	var audits []audit.Model
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return audits, err
	}

	err = cursor.All(context.Background(), &audits)
	if err != nil {
		return audits, err
	}

	return audits, nil
}

// DeleteAudit deletes an audit by id
// TODO: Are you sure you want to delete an audit by id?
/*func (r AuditMongoRepository) DeleteAudit(id bson.ObjectId) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}*/

// DeleteOlderThan deletes all audits older than a date
func (r AuditMongoRepository) DeleteOlderThan(date time.Time) error {
	_, err := r.collection.DeleteMany(context.Background(), bson.M{"instant": bson.M{"$lt": date}})
	if err != nil {
		return err
	}

	return nil
}

// DeleteInterval deletes all audits between two dates
func (r AuditMongoRepository) DeleteInterval(from time.Time, to time.Time) error {
	_, err := r.collection.DeleteMany(context.Background(), bson.M{"instant": bson.M{"$gte": from, "$lte": to}})
	if err != nil {
		return err
	}

	return nil
}
