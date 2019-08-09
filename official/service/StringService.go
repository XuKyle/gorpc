package service

import (
	"errors"
	"strings"
)

var ErrEmpty = errors.New("Empty String!")

// service
type StringService interface {
	UpperCase(string) (string, error)
	Count(string) int
}

// service 具体实现
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
