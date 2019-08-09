package impl

import (
	"errors"
	"strings"
)

var ErrEmpty = errors.New("Empty String!")

type StringServiceImpl struct {
}

func (StringServiceImpl) UpperCase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (StringServiceImpl) Count(s string) int {
	return len(s)
}
