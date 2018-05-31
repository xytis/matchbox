package storage

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Opts contains registered possible values for storages
type Opts struct {
	values map[string]string
}

var validators map[string]func(string) error

func init() {
	validators = make(map[string]func(string) error)
}

func registerOpt(name string, validator func(string) error) {
	if _, contains := validators[name]; contains {
		panic("duplicate opt name")
	}
	validators[name] = validator
}

// NewOpts creates an empty Opts structure
func NewOpts(values map[string]string) *Opts {
	return &Opts{
		values: values,
	}
}

// String returns internal contents
func (s *Opts) String() string {
	return fmt.Sprintf("%v", s.values)
}

// Set inserts value into Opts map
func (s *Opts) Set(value string) error {
	vals := strings.SplitN(value, "=", 2)
	if len(vals) >= 1 {
		key := vals[0]
		if validate, contains := validators[key]; contains {
			val := ""
			if len(vals) > 1 {
				val = vals[1]
			}
			if err := validate(val); err == nil {
				s.values[key] = val
			} else {
				return errors.Wrap(err, fmt.Sprintf("option %s validation failed", key))
			}
		} else {
			return errors.New("option is not supported")
		}
	}
	return nil
}

// Type returns type property
func (s *Opts) Type() string {
	return "map"
}
