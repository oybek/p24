package rest

type City struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ByValue implements sort.Interface for sorting by the Value field
type ByValue []City

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
