package repository

import (
	"context"
	"errors"
	"github.com/LydiaTrack/ground/pkg/responses"
	"github.com/LydiaTrack/ground/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseRepository provides default implementations for common CRUD operations.
type BaseRepository[T any] struct {
	Collection *mongo.Collection
}

// NewBaseRepository creates a new instance of BaseRepository.
func NewBaseRepository[T any](collection *mongo.Collection) *BaseRepository[T] {
	return &BaseRepository[T]{Collection: collection}
}

// Create inserts a new document into the collection.
func (r *BaseRepository[T]) Create(ctx context.Context, entity T) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, entity)
}

// Exists checks if a document matching the filter exists in the collection.
func (r *BaseRepository[T]) Exists(ctx context.Context, filter interface{}) (bool, error) {
	// Default to empty filter if nil is provided
	if filter == nil {
		filter = bson.M{}
	}

	// Perform the query
	err := r.Collection.FindOne(ctx, filter).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	return err == nil, err
}

// ExistsByID checks if a document with the provided ID exists in the collection.
func (r *BaseRepository[T]) ExistsByID(ctx context.Context, id interface{}) (bool, error) {
	objectID, err := utils.ToObjectID(id)
	if err != nil {
		return false, err
	}

	return r.Exists(ctx, bson.M{"_id": objectID})
}

// GetByID retrieves a document by its ID.
func (r *BaseRepository[T]) GetByID(ctx context.Context, id interface{}) (T, error) {
	var result T
	objectID, err := utils.ToObjectID(id)
	if err != nil {
		return result, err
	}
	filter := bson.M{"_id": objectID}
	err = r.Collection.FindOne(ctx, filter).Decode(&result)
	return result, err
}

// Update modifies an existing document identified by its ID.
func (r *BaseRepository[T]) Update(ctx context.Context, id interface{}, updateCommand interface{}) (*mongo.UpdateResult, error) {
	// Convert the ID to ObjectID
	objectID, err := utils.ToObjectID(id)
	if err != nil {
		return nil, err
	}

	// Generate the update document
	updateDoc, err := utils.GenerateUpdateDocument(updateCommand)
	if err != nil {
		return nil, err
	}

	// Perform the update
	return r.Collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updateDoc})
}

// Delete removes a document from the collection by its ID.
func (r *BaseRepository[T]) Delete(ctx context.Context, id interface{}) (*mongo.DeleteResult, error) {
	objectID, err := utils.ToObjectID(id)
	if err != nil {
		return nil, err
	}
	return r.Collection.DeleteOne(ctx, bson.M{"_id": objectID})
}

// Query retrieves documents matching the provided filter.
func (r *BaseRepository[T]) Query(ctx context.Context, filter interface{}, searchFields []string, searchText string) (responses.QueryResult[T], error) {
	var results []T

	// Ensure filter is not nil
	if filter == nil {
		filter = bson.M{}
	}

	// Add search criteria if searchText and searchFields are provided
	if searchText != "" && len(searchFields) > 0 {
		searchConditions := bson.A{}
		for _, field := range searchFields {
			searchConditions = append(searchConditions, bson.M{
				field: bson.M{"$regex": searchText, "$options": "i"}, // Case-insensitive search
			})
		}
		filter = bson.M{"$and": []bson.M{filter.(bson.M), {"$or": searchConditions}}}
	}

	// Execute the query
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return responses.QueryResult[T]{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return responses.QueryResult[T]{}, err
	}

	if results == nil {
		results = []T{}
	}

	return responses.QueryResult[T]{Data: results, TotalElements: len(results)}, nil
}

func (r *BaseRepository[T]) QueryPaginate(ctx context.Context, filter interface{}, searchFields []string, searchText string, page, limit int, sort interface{}) (PaginatedResult[T], error) {
	var results []T

	// Ensure filter is not nil
	if filter == nil {
		filter = bson.M{}
	}

	// Add search criteria if searchText and searchFields are provided
	if searchText != "" && len(searchFields) > 0 {
		searchConditions := bson.A{}
		for _, field := range searchFields {
			searchConditions = append(searchConditions, bson.M{
				field: bson.M{"$regex": searchText, "$options": "i"}, // Case-insensitive search
			})
		}
		filter = bson.M{"$and": []bson.M{filter.(bson.M), {"$or": searchConditions}}}
	}

	// Pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := int64((page - 1) * limit)

	// Default sort order
	if sort == nil {
		sort = bson.M{"_id": 1} // Default to ascending by _id
	}

	// MongoDB query options
	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(sort)

	// Execute the query
	cursor, err := r.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		return PaginatedResult[T]{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return PaginatedResult[T]{}, err
	}

	// Count total matching documents
	totalElements, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return PaginatedResult[T]{}, err
	}

	if results == nil {
		results = []T{}
	}

	return PaginatedResult[T]{
		Data:          results,
		TotalElements: totalElements,
		Page:          page,
		Limit:         limit,
	}, nil
}
