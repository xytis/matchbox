package testfakes

import (
	"errors"
)

var (
	errIntentional      = errors.New("store: error for testing purposes")
	errGroupNotFound    = errors.New("store: group not found")
	errProfileNotFound  = errors.New("store: profile not found")
	errTemplateNotFound = errors.New("store: template not found")
)
