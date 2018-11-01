package flag

import "errors"

// Value is the interface to the dynamic value stored in a flag.
type Value interface {
	Set(string) error
}

type StringValue string

func (s *StringValue) Set(val string) error {
	*s = StringValue(val)
	return nil
}

type BoolValue bool

func (b *BoolValue) Set(_ string) error {
	*b = true
	return nil
}

type CountValue int

func (c *CountValue) Set(_ string) error {
	*c++
	return nil
}

// ErrHelp is the error returned when a help flag is seen.
var ErrHelp = errors.New("help requested")

type HelpValue struct{}

func (b *HelpValue) Set(_ string) error {
	return ErrHelp
}
