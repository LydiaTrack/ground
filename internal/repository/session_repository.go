package repository

import (
	"context"

	"github.com/LydiaTrack/lydia-base/pkg/domain/session"
	"github.com/LydiaTrack/lydia-base/pkg/mongodb"
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
	ctx := context.Background()
	// FIXME: Burada ileride uzaktaki bir mongodb instance'ına bağlanmak gerekecek

	collection := mongodb.GetCollection("sessions", ctx)

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

// DeleteSessionByUserId is a function that deletes all sessions of a user
func (s SessionMongoRepository) DeleteSessionByUserId(userId primitive.ObjectID) error {
	_, err := s.collection.DeleteMany(context.Background(), primitive.M{"userId": userId})
	if err != nil {
		return err
	}
	return nil
}

// DeleteSessionById is a function that deletes a session by id
func (s SessionMongoRepository) DeleteSessionById(sessionId primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(context.Background(), primitive.M{"_id": sessionId})
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
