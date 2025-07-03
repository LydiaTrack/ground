package repository

import (
	"context"
	"time"

	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStatsMongoRepository struct {
	Collection *mongo.Collection
}

func GetUserStatsMongoRepository() *UserStatsMongoRepository {
	collection, err := mongodb.GetCollection("user_stats")
	if err != nil {
		panic(err)
	}
	return &UserStatsMongoRepository{
		Collection: collection,
	}
}

// GetStatsByUserID retrieves stats for a specific user
func (r *UserStatsMongoRepository) GetStatsByUserID(userID primitive.ObjectID) (user.StatsDocument, error) {
	var statsDoc user.StatsDocument
	err := r.Collection.FindOne(context.Background(), bson.M{"userId": userID}).Decode(&statsDoc)
	return statsDoc, err
}

// CreateStats creates a new stats document for a user
func (r *UserStatsMongoRepository) CreateStats(stats user.StatsDocument) error {
	_, err := r.Collection.InsertOne(context.Background(), stats)
	return err
}

// UpdateStats updates a user's stats (full document update)
func (r *UserStatsMongoRepository) UpdateStats(stats user.StatsDocument) error {
	stats.SetField("updatedDate", time.Now())
	core := stats.GetCoreFields()
	_, err := r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": core.ID},
		bson.M{"$set": stats},
	)
	return err
}

// IncrementField increments a numeric field by a specific value
func (r *UserStatsMongoRepository) IncrementField(statsID primitive.ObjectID, fieldName string, increment int) error {
	// Get the current stats to calculate fields
	var stats user.StatsDocument
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": statsID}).Decode(&stats)
	if err != nil {
		return err
	}

	// Calculate general stat fields
	stats.CalculateStatFields()

	// Perform the increment and update calculated fields in a single operation
	_, err = r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": statsID},
		bson.M{
			"$inc": bson.M{fieldName: increment},
			"$set": bson.M{
				"updatedDate":     stats.GetTime("updatedDate"),
				"lastActiveDate":  stats.GetTime("lastActiveDate"),
				"dayAge":          stats.GetInt("dayAge"),
				"activeDaysCount": stats.GetInt("activeDaysCount"),
			},
		},
	)
	return err
}

// IncrementInt64Field increments an int64 field by a specific value
func (r *UserStatsMongoRepository) IncrementInt64Field(statsID primitive.ObjectID, fieldName string, increment int64) error {
	// Get the current stats to calculate fields
	var stats user.StatsDocument
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": statsID}).Decode(&stats)
	if err != nil {
		return err
	}

	// Calculate general stat fields
	stats.CalculateStatFields()

	// Perform the increment and update calculated fields in a single operation
	_, err = r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": statsID},
		bson.M{
			"$inc": bson.M{fieldName: increment},
			"$set": bson.M{
				"updatedDate":     stats.GetTime("updatedDate"),
				"lastActiveDate":  stats.GetTime("lastActiveDate"),
				"dayAge":          stats.GetInt("dayAge"),
				"activeDaysCount": stats.GetInt("activeDaysCount"),
			},
		},
	)
	return err
}

// UpdateField updates a specific field value
func (r *UserStatsMongoRepository) UpdateField(statsID primitive.ObjectID, fieldName string, value interface{}) error {
	// Get the current stats to calculate fields
	var stats user.StatsDocument
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": statsID}).Decode(&stats)
	if err != nil {
		return err
	}

	// Calculate general stat fields
	stats.CalculateStatFields()

	// Add the updated date and calculated fields to the update
	updateFields := bson.M{
		fieldName:         value,
		"updatedDate":     stats.GetTime("updatedDate"),
		"lastActiveDate":  stats.GetTime("lastActiveDate"),
		"dayAge":          stats.GetInt("dayAge"),
		"activeDaysCount": stats.GetInt("activeDaysCount"),
	}

	_, err = r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": statsID},
		bson.M{"$set": updateFields},
	)
	return err
}

// UpdateFields updates multiple specific fields at once
func (r *UserStatsMongoRepository) UpdateFields(statsID primitive.ObjectID, fields map[string]interface{}) error {
	// Get the current stats to calculate fields
	var stats user.StatsDocument
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": statsID}).Decode(&stats)
	if err != nil {
		return err
	}

	// Calculate general stat fields
	stats.CalculateStatFields()

	// Add the calculated fields to the update
	fields["updatedDate"] = stats.GetTime("updatedDate")
	fields["lastActiveDate"] = stats.GetTime("lastActiveDate")
	fields["dayAge"] = stats.GetInt("dayAge")
	fields["activeDaysCount"] = stats.GetInt("activeDaysCount")

	_, err = r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": statsID},
		bson.M{"$set": fields},
	)
	return err
}
