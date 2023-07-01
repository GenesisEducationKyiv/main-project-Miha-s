package models

type Email struct {
	Value string `json:"value"`
}

func (email *Email) String() string {
	return email.Value
}
