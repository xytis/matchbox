package storagepb

import (
	"errors"
)

var (
	ErrIdRequired      = errors.New("Id is required")
	ErrProfileRequired = errors.New("Profile Id is required")
	ErrTypeRequired    = errors.New("Type is required")
)
