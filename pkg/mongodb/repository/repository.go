package repository

import (
	"context"
	"github.com/LydiaTrack/ground/pkg/responses"

	"go.mongodb.org/mongo-driver/mongo"
)

// Repository defines the interface for generic CRUD operations and pagination.
type Repository[T any] interface {
	Create(ctx context.Context, entity T) (*mongo.InsertOneResult, error)
	GetByID(ctx context.Context, id interface{}) (T, error)
	Update(ctx context.Context, id interface{}, update interface{}) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, id interface{}) (*mongo.DeleteResult, error)
	Exists(ctx context.Context, filter interface{}) (bool, error)
	ExistsByID(ctx context.Context, id interface{}) (bool, error)
	Query(ctx context.Context, filter interface{}, searchFields []string, searchText string) (responses.QueryResult[T], error)
	QueryPaginate(ctx context.Context, filter interface{}, searchFields []string, searchText string, page, limit int, sort interface{}) (responses.PaginatedResult[T], error)
}
