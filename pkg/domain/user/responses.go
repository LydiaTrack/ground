package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateResponse struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Username    string               `json:"username" bson:"username"`
	PersonInfo  *PersonInfo          `json:"personInfo"`
	ContactInfo ContactInfo          `json:"contactInfo"`
	CreatedDate time.Time            `json:"createdDate" bson:"createdDate"`
	Version     int                  `json:"version" bson:"version"`
	RoleIds     []primitive.ObjectID `json:"roleIds" bson:"roleIds"`
}
