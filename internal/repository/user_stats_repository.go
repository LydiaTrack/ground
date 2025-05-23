package repository

import (
	"context"
	"time"

	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/mongodb"
	"github.com/LydiaTrack/ground/pkg/mongodb/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatsMongoRepository struct {
	*repository.BaseRepository[user.StatsModel]
}

func GetUserStatsMongoRepository() *UserStatsMongoRepository {
	collection, err := mongodb.GetCollection("user_stats")
	if err != nil {
		panic(err)
	}
	return &UserStatsMongoRepository{
		BaseRepository: repository.NewBaseRepository[user.StatsModel](collection),
	}
}

// GetStatsByUserID retrieves stats for a specific user
func (r *UserStatsMongoRepository) GetStatsByUserID(userID primitive.ObjectID) (user.StatsModel, error) {
	var statsModel user.StatsModel
	err := r.Collection.FindOne(context.Background(), bson.M{"userId": userID}).Decode(&statsModel)
	return statsModel, err
}

// CreateStats creates a new stats document for a user
func (r *UserStatsMongoRepository) CreateStats(stats *user.StatsModel) error {
	_, err := r.Collection.InsertOne(context.Background(), stats)
	return err
}

// UpdateStats updates a user's stats (full document update)
func (r *UserStatsMongoRepository) UpdateStats(stats *user.StatsModel) error {
	stats.UpdatedDate = time.Now()
	_, err := r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": stats.ID},
		bson.M{"$set": stats},
	)
	return err
}

// IncrementField increments a numeric field by a specific value
func (r *UserStatsMongoRepository) IncrementField(statsID primitive.ObjectID, fieldName string, increment int) error {
	_, err := r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": statsID},
		bson.M{
			"$inc": bson.M{fieldName: increment},
			"$set": bson.M{"updatedDate": time.Now()},
		},
	)
	return err
}

// IncrementInt64Field increments an int64 field by a specific value
func (r *UserStatsMongoRepository) IncrementInt64Field(statsID primitive.ObjectID, fieldName string, increment int64) error {
	_, err := r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": statsID},
		bson.M{
			"$inc": bson.M{fieldName: increment},
			"$set": bson.M{"updatedDate": time.Now()},
		},
	)
	return err
}

// UpdateField updates a specific field value
func (r *UserStatsMongoRepository) UpdateField(statsID primitive.ObjectID, fieldName string, value interface{}) error {
	// Get the current stats to calculate fields
	var stats user.StatsModel
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": statsID}).Decode(&stats)
	if err != nil {
		return err
	}

	// Calculate general stat fields
	stats.CalculateStatFields()

	// Add the updated date and calculated fields to the update
	updateFields := bson.M{
		fieldName:         value,
		"updatedDate":     stats.UpdatedDate,
		"lastActiveDate":  stats.LastActiveDate,
		"dayAge":          stats.DayAge,
		"activeDaysCount": stats.ActiveDaysCount,
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
	var stats user.StatsModel
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": statsID}).Decode(&stats)
	if err != nil {
		return err
	}

	// Calculate general stat fields
	stats.CalculateStatFields()

	// Add the calculated fields to the update
	fields["updatedDate"] = stats.UpdatedDate
	fields["lastActiveDate"] = stats.LastActiveDate
	fields["dayAge"] = stats.DayAge
	fields["activeDaysCount"] = stats.ActiveDaysCount

	_, err = r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": statsID},
		bson.M{"$set": fields},
	)
	return err
}
