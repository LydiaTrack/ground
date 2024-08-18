package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateUserCommand struct {
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	PersonInfo  *PersonInfo `json:"personInfo"`
	ContactInfo ContactInfo `json:"contactInfo"`
}

type UpdateUserCommand struct {
	Username                 string      `json:"username"`
	PersonInfo               *PersonInfo `json:"personInfo"`
	ContactInfo              ContactInfo `json:"contactInfo"`
	LastSeenChangelogVersion string      `json:"lastSeenChangelogVersion"`
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
