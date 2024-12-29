package repository

import (
	"context"

	"github.com/LydiaTrack/ground/pkg/domain/role"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/mongodb"
	"github.com/LydiaTrack/ground/pkg/mongodb/repository"
	"github.com/LydiaTrack/ground/pkg/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserMongoRepository struct {
	*repository.BaseRepository[user.Model]
	roleRepository *RoleMongoRepository
}

func GetUserMongoRepository(roleRepo *RoleMongoRepository) *UserMongoRepository {
	collection, err := mongodb.GetCollection("users")
	if err != nil {
		panic(err)
	}
	return &UserMongoRepository{
		BaseRepository: repository.NewBaseRepository[user.Model](collection),
		roleRepository: roleRepo,
	}
}

func (r *UserMongoRepository) GetUsers(ctx context.Context, searchText string) ([]user.Model, error) {
	searchFields := []string{"username", "contactInfo.email"}
	return r.Query(ctx, nil, searchFields, searchText)
}

func (r *UserMongoRepository) GetUsersPaginated(ctx context.Context, searchText string, page, limit int) (repository.PaginatedResult[user.Model], error) {
	searchFields := []string{"username", "contactInfo.email"}
	result, err := r.QueryPaginate(ctx, nil, searchFields, searchText, page, limit, bson.M{"username": 1})
	if err != nil {
		return repository.PaginatedResult[user.Model]{}, err
	}
	return result, err
}

// ExistsByUsernameAndEmail checks if a user exists by username or email
func (r *UserMongoRepository) ExistsByUsernameAndEmail(username string, email string) bool {
	count, err := r.Collection.CountDocuments(context.Background(), bson.M{"$or": []bson.M{{"username": username}, {"contactInfo.email": email}}})
	if err != nil {
		panic(err)
	}
	return count > 0
}

// ExistsByUsername checks if a user exists by username
func (r *UserMongoRepository) ExistsByUsername(username string) bool {
	count, err := r.Collection.CountDocuments(context.Background(), bson.M{"username": username})
	return err == nil && count > 0
}

// ExistsByEmail checks if a user exists by email
func (r *UserMongoRepository) ExistsByEmail(email string) bool {
	count, err := r.Collection.CountDocuments(context.Background(), bson.M{"contactInfo.email": email})
	return err == nil && count > 0
}

// GetByUsername retrieves a user by username
func (r *UserMongoRepository) GetByUsername(username string) (user.Model, error) {
	var userModel user.Model
	err := r.Collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&userModel)
	return userModel, err
}

// GetByEmail retrieves a user by email
func (r *UserMongoRepository) GetByEmail(email string) (user.Model, error) {
	var userModel user.Model
	err := r.Collection.FindOne(context.Background(), bson.M{"contactInfo.email": email}).Decode(&userModel)
	return userModel, err
}

// AddRole adds a role to a user
func (r *UserMongoRepository) AddRole(userID, roleID primitive.ObjectID) error {
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$addToSet": bson.M{"roleIds": roleID}})
	return err
}

// RemoveRole removes a role from a user
func (r *UserMongoRepository) RemoveRole(userID, roleID primitive.ObjectID) error {
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$pull": bson.M{"roleIds": roleID}})
	return err
}

// GetUserRoles retrieves roles for a user
func (r *UserMongoRepository) GetUserRoles(userID primitive.ObjectID) (responses.QueryResult[role.Model], error) {
	userModel, err := r.GetByID(context.Background(), userID)
	if err != nil {
		return responses.QueryResult[role.Model]{}, err
	}

	var roles []role.Model
	for _, roleID := range *userModel.RoleIDs {
		roleModel, err := r.roleRepository.GetByID(context.Background(), roleID)
		if err != nil {
			return responses.QueryResult[role.Model]{}, err
		}
		roles = append(roles, roleModel)
	}

	return *responses.NewQueryResult(len(roles), roles), nil
}

// UpdateUserPassword updates a user's password
func (r *UserMongoRepository) UpdateUserPassword(userID primitive.ObjectID, password string) error {
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$set": bson.M{"password": password}})
	return err
}
