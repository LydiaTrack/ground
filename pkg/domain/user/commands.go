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
}

type UpdateUserCommand struct {
	Username                 string                 `json:"username" bson:"username"`
	Avatar                   string                 `json:"avatar,omitempty" bson:"avatar,omitempty"`
	PersonInfo               *PersonInfo            `json:"personInfo" bson:"personInfo"`
	ContactInfo              ContactInfo            `json:"contactInfo" bson:"contactInfo"`
	LastSeenChangelogVersion string                 `json:"lastSeenChangelogVersion" bson:"lastSeenChangelogVersion"`
	Properties               map[string]interface{} `json:"properties" bson:"properties"`
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
		if err := utils.ValidateBase64Image(cmd.Avatar); err != nil {
			return err
		}
	}

	return nil
}

type DeleteUserCommand struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}

type AddRoleToUserCommand struct {
	UserID primitive.ObjectID `json:"userId"`
	RoleID primitive.ObjectID `json:"roleId"`
}

type RemoveRoleFromUserCommand struct {
	UserID primitive.ObjectID `json:"userId"`
	RoleID primitive.ObjectID `json:"roleId"`
}

type UpdatePasswordCommand struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type ResetPasswordCommand struct {
	NewPassword string `json:"newPassword"`
}
