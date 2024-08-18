package resetPassword

type SendResetPasswordCodeCommand struct {
	Email string `json:"email"`
}

type DoResetPasswordCommand struct {
	Code        string `json:"code"`
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
}

type VerifyResetPasswordCodeCommand struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}
