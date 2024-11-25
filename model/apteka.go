package model

type Apteka struct {
	Name    string
	Phone   string
	Address string
}

func (a Apteka) IsValid() bool {
	return a.Name != "" && a.Phone != "" && a.Address != ""
}
