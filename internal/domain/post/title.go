package post

import "errors"

var (
	ErrShortTitle = errors.New("слишком короткий заголовок")
	ErrLongTitle  = errors.New("слишком длинный заголовок")
)

type Title struct {
	Value string
}

func NewTitle(value string) (*Title, error) {
	if len(value) < 3 {
		return nil, ErrShortTitle
	}

	if len(value) > 255 {
		return nil, ErrLongTitle
	}

	return &Title{
		Value: value,
	}, nil
}
