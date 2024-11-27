package model

type Medicine struct {
	Name   string
	Amount int64
}

func (a Medicine) IsValid() bool {
	return a.Name != "" && a.Amount > 0
}
