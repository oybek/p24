package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
)

type Trip struct {
	ChatId    int64     `json:"chat_id"`
	Path      []string  `json:"path"`
	StartTime time.Time `json:"start_time"`
	Phone     string    `json:"phone"`
	Comment   string    `json:"comment"`
}

func (t Trip) IsValid() bool {
	return len(t.Path) > 1 && len(t.Phone) > 0
}

func (t *Trip) String() string {
	pathText := strings.Join(lo.Map(t.Path, func(s string, i int) string {
		return fmt.Sprintf("`%d│ %s`", i, s)
	}), "\n")

	timeText := fmt.Sprintf(
		"%02d:%02d %02d/%02d/%04d",
		t.StartTime.Hour(),
		t.StartTime.Minute(),
		t.StartTime.Day(),
		t.StartTime.Month(),
		t.StartTime.Year(),
	)

	return fmt.Sprintf(
		"%s\n\n🕖 %s\n🚙 %s\n📞 %s",
		pathText,
		timeText,
		t.Comment,
		t.Phone,
	)
}
