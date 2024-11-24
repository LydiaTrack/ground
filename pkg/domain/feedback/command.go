package feedback

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateFeedbackCommand represents the required fields for creating a new feedback instance
type CreateFeedbackCommand struct {
	UserID  primitive.ObjectID `json:"userID" bson:"userId"`
	Message string             `json:"message" bson:"message"`
	Type    FeedbackType       `json:"type" bson:"type"`
}
