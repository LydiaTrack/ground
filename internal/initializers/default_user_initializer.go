package initializers

import (
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/internal/domain/user"
	"github.com/LydiaTrack/lydia-base/internal/repository"
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/LydiaTrack/lydia-base/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

// InitializeDefaultUser initializes the default user with default credentials
func InitializeDefaultUser() error {

	// While using remote connection for MongoDB instead of container, the user can be exist in the database.
	// In this case, the default user will not be created.
	isExists := repository.GetUserRepository().ExistsByUsername(os.Getenv("DEFAULT_USER_USERNAME"))
	if isExists {
		utils.Log("Default user already exists")
		return nil
	}
	userCreateCmd := user.CreateUserCommand{
		Username: os.Getenv("DEFAULT_USER_USERNAME"),
		Password: os.Getenv("DEFAULT_USER_PASSWORD"),
		PersonInfo: user.PersonInfo{
			FirstName: "Lydia",
			LastName:  "Admin",
			BirthDate: primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err := service.NewUserService(repository.GetUserRepository()).
		CreateUser(userCreateCmd, []auth.Permission{auth.AdminPermission})
	if err != nil {
		return err
	}

	utils.Log("Default user created successfully")

	return nil
}
