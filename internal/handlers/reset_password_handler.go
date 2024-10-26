package handlers

import (
	"errors"
	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/domain/resetPassword"
	"github.com/LydiaTrack/ground/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ResetPasswordHandler struct {
	userService          service.UserService
	resetPasswordService service.ResetPasswordService
}

func NewResetPasswordHandler(userService service.UserService, resetPasswordService service.ResetPasswordService) ResetPasswordHandler {
	return ResetPasswordHandler{
		userService:          userService,
		resetPasswordService: resetPasswordService,
	}
}

// SendResetPasswordCode godoc
// @Summary Send a reset password code
// @Description send a reset password code.
// @Tags reset-password
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Success 200 {object} map[string]interface{}
// @Router /reset-password/send-email [post]
func (h ResetPasswordHandler) SendResetPasswordCode(c *gin.Context) {
	//ip := c.ClientIP()
	//method := c.Request.Method
	//endpoint := c.FullPath()
	var cmd resetPassword.SendResetPasswordCodeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.resetPasswordService.SendResetPasswordEmail(c, cmd)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	// Add the IP to the block list for 5 minutes
	//blocker.GlobalBlocker.Add(ip, method, endpoint, 5*time.Minute)

	c.JSON(200, gin.H{"message": "Code sent successfully"})
}

// ResetPassword godoc
// @Summary Reset password
// @Description reset password.
// @Tags reset-password
// @Accept json
// @Produce json
// @Param code body string true "Code"
// @Param password body string true "Password"
// @Success 200 {object} map[string]interface{}
// @Router /reset-password [post]
func (h ResetPasswordHandler) ResetPassword(c *gin.Context) {
	var cmd resetPassword.DoResetPasswordCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.resetPasswordService.ResetPassword(c, cmd)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Password reset successfully"})
}

// VerifyResetPasswordCode godoc
// @Summary Verify reset password code
// @Description verify reset password code.
// @Tags reset-password
// @Accept json
// @Produce json
// @Param code body string true "Code"
// @Param email body string true "Email"
// @Success 200 {object} map[string]interface{}
// @Router /reset-password/verify-code [post]
func (h ResetPasswordHandler) VerifyResetPasswordCode(c *gin.Context) {
	var cmd resetPassword.VerifyResetPasswordCodeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.resetPasswordService.VerifyResetPasswordCode(c, cmd)
	if err != nil {
		if errors.Is(err, resetPassword.ErrResetPasswordCodeExpired) {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, resetPassword.ErrResetPasswordCodeInvalid) {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, resetPassword.ErrResetPasswordNotFound) {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		utils.EvaluateError(err, c)
		return
	}

	c.JSON(200, gin.H{"message": "Code verified successfully"})
}
