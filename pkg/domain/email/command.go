package email

type SendEmailCommand struct {
	To      string
	Subject string
	Body    string
}
