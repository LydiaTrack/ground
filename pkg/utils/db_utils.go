package utils

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToObjectID converts various ID formats to a MongoDB ObjectID.
func ToObjectID(id interface{}) (primitive.ObjectID, error) {
	switch v := id.(type) {
	case string:
		return primitive.ObjectIDFromHex(v)
	case primitive.ObjectID:
		return v, nil
	default:
		return primitive.NilObjectID, errors.New("invalid ID format")
	}
}
