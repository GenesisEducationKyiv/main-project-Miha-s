package emailsStorage

type EmailsStorage interface {
	AddEmail(email string) error
	GetAllEmails() map[string]struct{}
	ValidateEmail(email string) bool
	Close()
}
