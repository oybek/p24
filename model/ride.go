package model

import (
	"fmt"
	"strings"
)

type Ride struct {
	Path      []string `json:"path"`
	When      string   `json:"when"`
	Car       string   `json:"car"`
	Phone     string   `json:"phone"`
	MessageId int64
}

func (r *Ride) Text() string {
	return fmt.Sprintf(
		"%s\n"+
			"🕒 %s\n"+
			"🚘 %s\n"+
			"📞 %s",
		strings.Join(r.Path, " ➡️ "),
		strings.ToLower(r.When),
		strings.ToLower(r.Car),
		r.Phone,
	)
}
