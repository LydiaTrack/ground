package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/session"
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
	GetUserSession(id primitive.ObjectID) (session.InfoModel, error)
	// DeleteSessionByUserID is a function that deletes a session
	DeleteSessionByUserID(id primitive.ObjectID) error
	// DeleteSessionByID is a function that deletes a session by id
	DeleteSessionByID(sessionID primitive.ObjectID) error
	// GetSessionByRefreshToken is a function that gets a session by refresh token
	GetSessionByRefreshToken(refreshToken string) (session.InfoModel, error)
}

func NewSessionService(sessionRepository SessionRepository, userService UserService) *SessionService {
	return &SessionService{
		sessionRepository: sessionRepository,
		UserService:       userService,
	}
}

// CreateSession is a function that creates a session
func (s SessionService) CreateSession(cmd session.CreateSessionCommand) (session.InfoModel, error) {
	// Check if user exists
	userID, err := primitive.ObjectIDFromHex(cmd.UserID)
	if err != nil {
		return session.InfoModel{}, err
	}
	exists, err := s.UserService.ExistsUser(cmd.UserID, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      &userID,
	})
	if err != nil {
		return session.InfoModel{}, err
	}
	if !exists {
		return session.InfoModel{}, constants.ErrorNotFound
	}
	sessionInfo := session.InfoModel{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		ExpireTime:   cmd.ExpireTime,
		RefreshToken: cmd.RefreshToken,
	}
	// TODO: Permission check
	return s.sessionRepository.SaveSession(sessionInfo)
}

// GetUserSession is a function that gets a user session
func (s SessionService) GetUserSession(id string) (session.InfoModel, error) {
	// Check if user exists
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return session.InfoModel{}, err
	}
	exists, err := s.UserService.ExistsUser(id, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      &userID,
	})
	if err != nil {
		return session.InfoModel{}, err
	}
	if !exists {
		return session.InfoModel{}, constants.ErrorNotFound
	}

	return s.sessionRepository.GetUserSession(userID)
}

// DeleteSessionByUser DeleteSession is a function that deletes a session
func (s SessionService) DeleteSessionByUser(userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	return s.sessionRepository.DeleteSessionByUserID(objID)
}

// DeleteSessionByID is a function that deletes a session by id
func (s SessionService) DeleteSessionByID(sessionID string) error {
	objID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return err
	}
	return s.sessionRepository.DeleteSessionByID(objID)
}

// IsUserHasActiveSession is a function that checks if a user has an active session
func (s SessionService) IsUserHasActiveSession(userID string) bool {
	sessionModel, err := s.GetUserSession(userID)
	if err != nil {
		return false
	}

	// Check if session still valid by comparing the expire time with the current time
	currentTime := time.Now().Unix()
	return sessionModel.ExpireTime >= currentTime
}

// GetSessionByRefreshToken is a function that gets a session by refresh token
func (s SessionService) GetSessionByRefreshToken(refreshToken string) (session.InfoModel, error) {
	return s.sessionRepository.GetSessionByRefreshToken(refreshToken)
}
