package email

type EmailTemplateData struct {
	Username string
	Code     string
}

type EmailCredentials struct {
	Address  string
	Password string
}
