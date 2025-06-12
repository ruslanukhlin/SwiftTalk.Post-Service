package post

import "errors"

var (
	ErrShortContent = errors.New("слишком короткое содержание")
	ErrLongContent = errors.New("слишком длинное содержание")
)

type Content struct {
	Value string
}

func NewContent(value string) (*Content, error) {
	if(len(value) < 3) {
		return nil, ErrShortContent
	}

	if(len(value) > 100000) {
		return nil, ErrLongContent
	}

	return &Content{
		Value: value,
	}, nil
}