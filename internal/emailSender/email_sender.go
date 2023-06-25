package emailSender

type EmailSender interface {
	BroadcastEmails(recipients *map[string]struct{}, body string)
	SendEmail(recipient, body string) error
}
