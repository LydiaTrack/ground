package repository

import (
	"context"

	"github.com/LydiaTrack/ground/pkg/domain/session"
	"github.com/LydiaTrack/ground/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionMongoRepository struct {
	collection *mongo.Collection
}

var (
	sessionRepository *SessionMongoRepository
)

func newSessionMongoRepository() *SessionMongoRepository {

	collection, err := mongodb.GetCollection("sessions")
	if err != nil {
		panic(err)
	}

	return &SessionMongoRepository{
		collection: collection,
	}
}

func GetSessionRepository() *SessionMongoRepository {
	if sessionRepository == nil {
		sessionRepository = newSessionMongoRepository()
	}
	return sessionRepository
}

// SaveSession is a function that creates a session
func (s SessionMongoRepository) SaveSession(sessionModel session.InfoModel) (session.InfoModel, error) {
	_, err := s.collection.InsertOne(context.Background(), sessionModel)
	if err != nil {
		return session.InfoModel{}, err
	}
	return sessionModel, nil
}

// GetUserSession is a function that gets a user session
func (s SessionMongoRepository) GetUserSession(id primitive.ObjectID) (session.InfoModel, error) {
	var sessionModel session.InfoModel
	err := s.collection.FindOne(context.Background(), primitive.M{"userId": id}).Decode(&sessionModel)
	if err != nil {
		return session.InfoModel{}, err
	}
	return sessionModel, nil
}

// DeleteSessionByUserID is a function that deletes all sessions of a user
func (s SessionMongoRepository) DeleteSessionByUserID(userID primitive.ObjectID) error {
	_, err := s.collection.DeleteMany(context.Background(), primitive.M{"userId": userID})
	if err != nil {
		return err
	}
	return nil
}

// DeleteSessionByID is a function that deletes a session by id
func (s SessionMongoRepository) DeleteSessionByID(sessionID primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(context.Background(), primitive.M{"_id": sessionID})
	if err != nil {
		return err
	}
	return nil
}

// GetSessionByRefreshToken is a function that gets a session by refresh token
func (s SessionMongoRepository) GetSessionByRefreshToken(refreshToken string) (session.InfoModel, error) {
	var sessionModel session.InfoModel
	err := s.collection.FindOne(context.Background(), primitive.M{"refreshToken": refreshToken}).Decode(&sessionModel)
	if err != nil {
		return session.InfoModel{}, err
	}
	return sessionModel, nil
}

// DeleteExpiredSessions deletes all sessions that have expired before the given time
func (s SessionMongoRepository) DeleteExpiredSessions(currentTime int64) error {
	filter := primitive.M{"expireTime": primitive.M{"$lt": currentTime}}
	_, err := s.collection.DeleteMany(context.Background(), filter)
	return err
}
