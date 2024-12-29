package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Paginate retrieves a paginated list of documents based on the filter, page, limit, and sort criteria.
func (r *BaseRepository[T]) Paginate(ctx context.Context, filter interface{}, page, limit int, sort interface{}) (PaginatedResult[T], error) {
	var results []T

	// Ensure page and limit have valid values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Default to an empty filter if none is provided
	if filter == nil {
		filter = bson.M{}
	}

	// Default sort order if none is provided
	if sort == nil {
		sort = bson.M{"_id": 1} // Sort by `_id` ascending
	}

	// Calculate the number of documents to skip
	skip := int64((page - 1) * limit)

	// Set up find options with pagination and sorting
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(skip)
	findOptions.SetSort(sort)

	// Execute the query
	cursor, err := r.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		return PaginatedResult[T]{}, err
	}
	defer cursor.Close(ctx)

	// Decode the results
	if err := cursor.All(ctx, &results); err != nil {
		return PaginatedResult[T]{}, err
	}

	// Get the total count of documents matching the filter
	totalCount, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return PaginatedResult[T]{}, err
	}

	// Construct and return the paginated result
	return PaginatedResult[T]{
		Data:       results,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}, nil
}
