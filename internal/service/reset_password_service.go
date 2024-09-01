package service

import (
	"errors"
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/constants"
	"github.com/LydiaTrack/lydia-base/pkg/domain/email"
	"github.com/LydiaTrack/lydia-base/pkg/domain/resetPassword"
	"github.com/LydiaTrack/lydia-base/pkg/domain/user"
	"github.com/LydiaTrack/lydia-base/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ResetPasswordService struct {
	userService             UserService
	emailService            SimpleEmailService
	resetPasswordRepository ResetPasswordRepository
}

func NewResetPasswordService(resetPasswordRepository ResetPasswordRepository, emailService SimpleEmailService, userService UserService) *ResetPasswordService {
	return &ResetPasswordService{
		userService:             userService,
		emailService:            emailService,
		resetPasswordRepository: resetPasswordRepository,
	}
}

type ResetPasswordRepository interface {
	// SaveResetPassword saves a resetPassword
	SaveResetPassword(resetPasswordModel resetPassword.Model) (resetPassword.Model, error)
	// GetResetPasswordByCode gets a resetPassword by code
	GetResetPasswordByCode(code string) (resetPassword.Model, error)
	// DeleteResetPassword deletes a resetPassword by code
	DeleteResetPassword(id primitive.ObjectID) error
	// DeleteResetPasswordByCode deletes a resetPassword by code
	DeleteResetPasswordByCode(code string) error
}

func (s ResetPasswordService) createResetPassword(cmd resetPassword.SendResetPasswordCodeCommand) (resetPassword.Model, error) {
	authContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	}

	// Check if the user exists
	_, err := s.userService.GetUserByEmail(cmd.Email, authContext)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return resetPassword.Model{}, constants.ErrorNotFound
		}
		return resetPassword.Model{}, err
	}

	threeMinutesLater := primitive.DateTime(time.Now().Add(3*time.Minute).UnixNano() / int64(time.Millisecond))
	code, err := utils.Generate6DigitCode(false)
	if err != nil {
		return resetPassword.Model{}, err
	}

	resetPasswordModel := resetPassword.NewModel(cmd.Email, code, threeMinutesLater)
	resetPasswordModel, err = s.resetPasswordRepository.SaveResetPassword(resetPasswordModel)
	if err != nil {
		return resetPassword.Model{}, err
	}

	return resetPasswordModel, nil
}

// SendResetPasswordEmail sends a reset password email
func (s ResetPasswordService) SendResetPasswordEmail(c *gin.Context, cmd resetPassword.SendResetPasswordCodeCommand) error {
	resetPasswordModel, err := s.createResetPassword(cmd)
	if err != nil {
		return err
	}

	userModel, err := s.userService.GetUserByEmail(cmd.Email, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})
	if err != nil {
		return err
	}

	templateData := email.EmailTemplateData{
		Code:     resetPasswordModel.Code,
		Username: userModel.Username,
	}

	err = s.emailService.SendEmail(email.SendEmailCommand{
		To:      cmd.Email,
		Subject: "Renoten Reset Password",
	}, "RESET_PASSWORD", templateData)
	if err != nil {
		return err
	}

	return nil
}

// VerifyResetPasswordCode verifies a reset password code
func (s ResetPasswordService) VerifyResetPasswordCode(c *gin.Context, cmd resetPassword.VerifyResetPasswordCodeCommand) error {
	resetPasswordModel, err := s.resetPasswordRepository.GetResetPasswordByCode(cmd.Code)
	if err != nil {
		return constants.ErrorNotFound
	}

	if resetPasswordModel.Email != cmd.Email {
		return resetPassword.ErrResetPasswordCodeInvalid
	}

	if resetPasswordModel.ExpiresAt < primitive.DateTime(time.Now().UnixNano()/int64(time.Millisecond)) {
		return resetPassword.ErrResetPasswordCodeExpired
	}

	return nil
}

// ResetPassword resets a password by communicating with the user service
func (s ResetPasswordService) ResetPassword(c *gin.Context, cmd resetPassword.DoResetPasswordCommand) error {
	verifyCmd := resetPassword.VerifyResetPasswordCodeCommand{
		Code:  cmd.Code,
		Email: cmd.Email,
	}
	err := s.VerifyResetPasswordCode(c, verifyCmd)
	if err != nil {
		utils.EvaluateError(err, c)
	}

	authCtx := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	}

	userToUpdate, err := s.userService.GetUserByEmail(cmd.Email, authCtx)
	if err != nil {
		return err
	}
	updatePasswordCmd := user.ResetPasswordCommand{
		NewPassword: cmd.NewPassword,
	}

	err = s.userService.ResetUserPassword(userToUpdate.ID.Hex(), updatePasswordCmd)
	if err != nil {
		return err
	}

	// Delete the reset password code after the password has been reset
	// TODO: There should be a scheduled job to delete expired reset password codes
	err = s.resetPasswordRepository.DeleteResetPasswordByCode(cmd.Code)
	if err != nil {
		return err
	}

	return nil
}
