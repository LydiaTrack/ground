package service

import (
	"errors"
	"github.com/LydiaTrack/lydia-base/internal/domain/auth"
	"github.com/LydiaTrack/lydia-base/internal/domain/session"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// SessionService is an interface that contains the methods for the session service
type SessionService struct {
	sessionRepository SessionRepository
	UserService
}

// SessionRepository is an interface that contains the methods for the session repository
type SessionRepository interface {
	// SaveSession is a function that creates a session
	SaveSession(model session.InfoModel) (session.InfoModel, error)
	// GetUserSession is a function that gets a user session
	GetUserSession(id bson.ObjectId) (session.InfoModel, error)
	// DeleteSessionByUserId is a function that deletes a session
	DeleteSessionByUserId(id bson.ObjectId) error
	// DeleteSessionById is a function that deletes a session by id
	DeleteSessionById(sessionId bson.ObjectId) error
}

func NewSessionService(sessionRepository SessionRepository, userService UserService) SessionService {
	return SessionService{
		sessionRepository: sessionRepository,
		UserService:       userService,
	}
}

// CreateSession is a function that creates a session
func (s SessionService) CreateSession(cmd session.CreateSessionCommand) (session.InfoModel, error) {
	// Check if user exists
	exists, err := s.UserService.ExistsUser(cmd.UserId, []auth.Permission{auth.AdminPermission})
	if err != nil {
		return session.InfoModel{}, err
	}
	if !exists {
		return session.InfoModel{}, errors.New("user does not exist")
	}
	sessionInfo := session.InfoModel{
		ID:           bson.NewObjectId(),
		UserId:       bson.ObjectIdHex(cmd.UserId),
		ExpireTime:   cmd.ExpireTime,
		RefreshToken: cmd.RefreshToken,
	}
	// TODO: Permission check
	return s.sessionRepository.SaveSession(sessionInfo)
}

// GetUserSession is a function that gets a user session
func (s SessionService) GetUserSession(id string) (session.InfoModel, error) {
	// Check if user exists
	exists, err := s.UserService.ExistsUser(id, []auth.Permission{auth.AdminPermission})
	if err != nil {
		return session.InfoModel{}, err
	}
	if !exists {
		return session.InfoModel{}, errors.New("user does not exist")
	}

	return s.sessionRepository.GetUserSession(bson.ObjectIdHex(id))
}

// DeleteSession is a function that deletes a session
func (s SessionService) DeleteSessionByUser(userId string) error {
	return s.sessionRepository.DeleteSessionByUserId(bson.ObjectIdHex(userId))
}

// DeleteSessionById is a function that deletes a session by id
func (s SessionService) DeleteSessionById(sessionId string) error {
	return s.sessionRepository.DeleteSessionById(bson.ObjectIdHex(sessionId))
}

// IsUserHasActiveSession is a function that checks if a user has an active session
func (s SessionService) IsUserHasActiveSession(userId string) bool {
	session, err := s.GetUserSession(userId)
	if err != nil {
		return false
	}

	// Check if session still valid by comparing the expire time with the current time
	currentTime := time.Now().Unix()
	return session.ExpireTime >= currentTime
}
