package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToObjectID converts various ID formats to a MongoDB ObjectID.
func ToObjectID(id interface{}) (primitive.ObjectID, error) {
	switch v := id.(type) {
	case string:
		return primitive.ObjectIDFromHex(v)
	case primitive.ObjectID:
		return v, nil
	// If it is a pointer to a string, dereference it and convert to ObjectID
	case *string:
		if v == nil {
			return primitive.NilObjectID, nil
		}
		return ToObjectID(*v)
	// If it is a pointer to an ObjectID, dereference it
	case *primitive.ObjectID:
		if v == nil {
			return primitive.NilObjectID, nil
		}
		return *v, nil
	default:
		return primitive.NilObjectID, errors.New("invalid ID format")
	}
}

// GenerateUpdateDocument generates a BSON update document from a struct.
// Only non-zero and non-nil fields are included in the document. If the bson tag
// is empty, the json tag is considered.
func GenerateUpdateDocument(updateCommand interface{}) (bson.M, error) {
	updateDoc := bson.M{}

	val := reflect.ValueOf(updateCommand)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct type, got %T", updateCommand)
	}

	typ := reflect.TypeOf(updateCommand)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Prefer bson tag
		bsonKey := field.Tag.Get("bson")
		// If bson tag is empty, fallback to json tag
		if bsonKey == "" {
			bsonKey = field.Tag.Get("json")
		}

		// Skip invalid or excluded tags
		if bsonKey == "" || bsonKey == "-" {
			continue
		}

		// Extract the actual key name (omit everything after ',')
		bsonKey = strings.Split(bsonKey, ",")[0]

		// Check if the field has a non-zero value
		if !fieldValue.IsZero() {
			updateDoc[bsonKey] = fieldValue.Interface()
		}
		// If field has empty value, check if it is a pointer, and if it is nil just assign it
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			updateDoc[bsonKey] = nil
		}
	}

	return updateDoc, nil
}
