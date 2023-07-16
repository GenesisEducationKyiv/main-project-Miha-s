package models

import "fmt"

type Rate struct {
	Value float64 `json:"value'"`
}

func (rate *Rate) String() string {
	return fmt.Sprintf("%v", rate.Value)
}
