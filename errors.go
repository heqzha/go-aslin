package aslin

type ErrorType uint64

type (
	Error struct {
		Err  error
		Type ErrorType
		Meta interface{}
	}

	errors []*Error
)
