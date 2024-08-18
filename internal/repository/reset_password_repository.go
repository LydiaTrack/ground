package repository

import (
	"context"
	"errors"
	"github.com/LydiaTrack/lydia-base/pkg/domain/resetPassword"
	"github.com/LydiaTrack/lydia-base/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// A ResetPasswordMongoRepository that implements ResetPasswordRepository
type ResetPasswordMongoRepository struct {
	collection *mongo.Collection
}

var (
	resetPasswordRepository *ResetPasswordMongoRepository
)

// NewResetPasswordMongoRepository creates a new ResetPasswordMongoRepository instance
// which implements ResetPasswordRepository
func newResetPasswordMongoRepository() *ResetPasswordMongoRepository {

	collection, err := mongodb.GetCollection("resetPwCodes")
	if err != nil {
		panic(err)
	}

	return &ResetPasswordMongoRepository{
		collection: collection,
	}
}

// GetResetPasswordRepository returns a ResetPasswordRepository
func GetResetPasswordRepository() *ResetPasswordMongoRepository {
	if resetPasswordRepository == nil || resetPasswordRepository.collection == nil {
		resetPasswordRepository = newResetPasswordMongoRepository()
	}
	return resetPasswordRepository
}

// SaveResetPassword saves a resetPassword
func (r *ResetPasswordMongoRepository) SaveResetPassword(resetPasswordModel resetPassword.Model) (resetPassword.Model, error) {
	_, err := r.collection.InsertOne(context.Background(), resetPasswordModel)
	if err != nil {
		return resetPassword.Model{}, err
	}
	return resetPasswordModel, nil
}

// GetResetPasswordByCode gets a resetPassword by code
func (r *ResetPasswordMongoRepository) GetResetPasswordByCode(code string) (resetPassword.Model, error) {
	var resetPasswordModel resetPassword.Model
	err := r.collection.FindOne(context.Background(), primitive.M{"code": code}).Decode(&resetPasswordModel)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return resetPassword.Model{}, resetPassword.ErrResetPasswordNotFound
		} else {
			return resetPassword.Model{}, err
		}
	}
	return resetPasswordModel, nil
}

// DeleteResetPassword deletes a resetPassword by id
func (r *ResetPasswordMongoRepository) DeleteResetPassword(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), primitive.M{"_id": id})
	return err
}

// DeleteResetPasswordByCode deletes a resetPassword by code
func (r *ResetPasswordMongoRepository) DeleteResetPasswordByCode(code string) error {
	_, err := r.collection.DeleteOne(context.Background(), primitive.M{"code": code})
	return err
}
