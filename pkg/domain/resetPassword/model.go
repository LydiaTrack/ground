package resetPassword

import "go.mongodb.org/mongo-driver/bson/primitive"

type Model struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email"`
	Code      string             `json:"-"`
	ExpiresAt primitive.DateTime `json:"expiresAt" bson:"expiresAt"`
}

func NewModel(email, code string, expiresAt primitive.DateTime) Model {
	return Model{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Code:      code,
		ExpiresAt: expiresAt,
	}
}
