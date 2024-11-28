package database

// Iterator - description of the iterator interface
type Iterator interface {
	Next() bool
	Scan(dest ...any) error
}

// ItToSlice convert Iterator to Slice
func itToSlice[T any](it Iterator, scan func(Iterator) (T, error)) ([]T, error) {
	ts := []T{}

	for it.Next() {
		t, err := scan(it)
		if err != nil {
			return ts, err
		}
		ts = append(ts, t)
	}

	return ts, nil
}
