package user

import (
	"errors"

	"github.com/LydiaTrack/ground/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserCommand struct {
	Username    string                 `json:"username"`
	Password    string                 `json:"password"`
	PersonInfo  *PersonInfo            `json:"personInfo"`
	ContactInfo ContactInfo            `json:"contactInfo"`
	Properties  map[string]interface{} `json:"properties"`
	Avatar      string                 `json:"avatar,omitempty"`
	OAuthInfo   *OAuthInfo             `json:"OAuthInfo,omitempty"`
}

type UpdateUserCommand struct {
	Username                 string                 `json:"username,omitempty"`
	Password                 string                 `json:"password,omitempty"`
	Avatar                   string                 `json:"avatar,omitempty"`
	PersonInfo               *PersonInfo            `json:"personInfo,omitempty"`
	ContactInfo              *ContactInfo           `json:"contactInfo,omitempty"`
	Properties               map[string]interface{} `json:"properties,omitempty"`
	LastSeenChangelogVersion string                 `json:"lastSeenChangelogVersion,omitempty"`
	OAuthProviders           map[string]OAuthInfo   `json:"oauthProviders,omitempty"`
}

func (cmd UpdateUserCommand) Validate() error {
	if cmd.Username == "" {
		return errors.New("username is required")
	}

	if cmd.PersonInfo != nil {
		if err := cmd.PersonInfo.Validate(); err != nil {
			return err
		}
	}

	if cmd.Avatar != "" {
		if err := utils.ValidateUserAvatar(cmd.Avatar); err != nil {
			return err
		}
	}

	return nil
}

type DeleteUserCommand struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}

type AddRoleToUserCommand struct {
	UserID primitive.ObjectID `json:"userID"`
	RoleID primitive.ObjectID `json:"roleID"`
}

type RemoveRoleFromUserCommand struct {
	UserID primitive.ObjectID `json:"userID"`
	RoleID primitive.ObjectID `json:"roleID"`
}

type UpdatePasswordCommand struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type ResetPasswordCommand struct {
	NewPassword string `json:"newPassword"`
}
