package email

// TemplateContext represents the data required to render an email template
type TemplateContext struct {
	Data interface{} `json:"data"`
}

type EmailCredentials struct {
	Address  string
	Password string
}

type SupportedEmailType string

const (
	EmailTypeResetPassword SupportedEmailType = "RESET_PASSWORD"
	EmailTypeFeedback      SupportedEmailType = "FEEDBACK"
)
