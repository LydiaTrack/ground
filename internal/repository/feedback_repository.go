package repository

import (
	"context"
	"time"

	"github.com/LydiaTrack/ground/pkg/domain/feedback"
	"github.com/LydiaTrack/ground/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbackMongoRepository struct {
	collection *mongo.Collection
}

var (
	feedbackRepository *FeedbackMongoRepository
)

func newFeedbackRepository() *FeedbackMongoRepository {
	collection, err := mongodb.GetCollection("feedbacks")
	if err != nil {
		panic(err)
	}

	return &FeedbackMongoRepository{
		collection: collection,
	}
}

// GetFeedbackRepository returns a FeedbackMongoRepository instance if it is not initialized yet or
// returns the existing one
func GetFeedbackRepository() *FeedbackMongoRepository {
	if feedbackRepository == nil {
		feedbackRepository = newFeedbackRepository()
	}
	return feedbackRepository
}

// SaveFeedback saves a feedback record in the database
func (r *FeedbackMongoRepository) SaveFeedback(f feedback.Model) (feedback.Model, error) {
	_, err := r.collection.InsertOne(context.Background(), f)
	if err != nil {
		return f, err
	}
	return f, nil
}

// GetFeedback retrieves a feedback record by ID
func (r *FeedbackMongoRepository) GetFeedback(id primitive.ObjectID) (feedback.Model, error) {
	var f feedback.Model
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&f)
	if err != nil {
		return f, err
	}
	return f, nil
}

// ExistsFeedback checks if a feedback record exists by ID
func (r *FeedbackMongoRepository) ExistsFeedback(id primitive.ObjectID) (bool, error) {
	var f feedback.Model
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&f)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetFeedbacks retrieves all feedback records
func (r *FeedbackMongoRepository) GetFeedbacks() ([]feedback.Model, error) {
	var feedbacks []feedback.Model
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return feedbacks, err
	}

	err = cursor.All(context.Background(), &feedbacks)
	if err != nil {
		return feedbacks, err
	}
	return feedbacks, nil
}

// DeleteFeedback deletes a feedback record by ID
func (r *FeedbackMongoRepository) DeleteFeedback(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

// DeleteOlderThan deletes all feedback records older than a specified date
func (r *FeedbackMongoRepository) DeleteOlderThan(date time.Time) error {
	_, err := r.collection.DeleteMany(context.Background(), bson.M{"createdAt": bson.M{"$lt": date}})
	return err
}

// UpdateFeedbackStatus updates the status of a feedback record by ID
func (r *FeedbackMongoRepository) UpdateFeedbackStatus(id primitive.ObjectID, status feedback.FeedbackStatus) error {
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	return err
}

// GetFeedbacksByUser retrieves all feedback records submitted by a specific user
func (r *FeedbackMongoRepository) GetFeedbacksByUser(userID primitive.ObjectID) ([]feedback.Model, error) {
	var feedbacks []feedback.Model
	cursor, err := r.collection.Find(context.Background(), bson.M{"userId": userID})
	if err != nil {
		return feedbacks, err
	}

	err = cursor.All(context.Background(), &feedbacks)
	if err != nil {
		return feedbacks, err
	}
	return feedbacks, nil
}
