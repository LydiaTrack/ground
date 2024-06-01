package repository

import (
	"context"
	"github.com/LydiaTrack/lydia-base/pkg/domain/session"
	"github.com/LydiaTrack/lydia-base/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
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
func (s SessionMongoRepository) GetUserSession(id bson.ObjectId) (session.InfoModel, error) {
	var sessionModel session.InfoModel
	err := s.collection.FindOne(context.Background(), bson.M{"userId": id}).Decode(&sessionModel)
	if err != nil {
		return session.InfoModel{}, err
	}
	return sessionModel, nil
}

// DeleteSessionByUserId is a function that deletes all sessions of a user
func (s SessionMongoRepository) DeleteSessionByUserId(userId bson.ObjectId) error {
	_, err := s.collection.DeleteMany(context.Background(), bson.M{"userId": userId})
	if err != nil {
		return err
	}
	return nil
}

// DeleteSessionById is a function that deletes a session by id
func (s SessionMongoRepository) DeleteSessionById(sessionId bson.ObjectId) error {
	_, err := s.collection.DeleteOne(context.Background(), bson.M{"_id": sessionId})
	if err != nil {
		return err
	}
	return nil
}
