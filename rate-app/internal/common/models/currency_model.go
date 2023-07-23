package models

import "fmt"

type Currency struct {
	From string
	To   string
}

func (curr *Currency) String() string {
	return fmt.Sprintf("%s->%s", curr.From, curr.To)
}
