package flag

import "errors"

// Value is the interface to the dynamic value stored in a flag.
type Value interface {
	Set(string) error
}

type stringVal string

func (s *stringVal) Set(val string) error {
	*s = stringVal(val)
	return nil
}

type boolVal bool

func (b *boolVal) Set(_ string) error {
	*b = true
	return nil
}

// ErrHelp is the error returned when a help flag is seen.
var ErrHelp = errors.New("help requested")

type helpVal struct{}

func (b *helpVal) Set(_ string) error {
	return ErrHelp
}
