package service

import (
	"errors"
	"fmt"

	"github.com/autumnterror/volha-backend/internal/platform/format"
)

var (
	ErrBadServiceCheck = errors.New("bad service check")
)

func wrapServiceCheck(op string, err error) error {
	if err == nil {
		return nil
	}
	return format.Error(op, fmt.Errorf("%w: %v", ErrBadServiceCheck, err))
}
