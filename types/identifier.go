package types

type ID int64

func (h ID) Valid() bool {
	return h > 0
}

func (h ID) Equal(o ID) bool {
	return h == o
}
