package flaq

import (
	"strconv"
	"time"
)

// Value is the interface to the dynamic value stored in a flag.
type Value interface {
	Set(string) error
}

type stringValue string

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

type boolValue bool

func (b *boolValue) Set(val string) error {
	if val == "" {
		*b = true
		return nil
	}
	v, err := strconv.ParseBool(val)
	*b = boolValue(v)
	return err
}

type countValue int

func (c *countValue) Set(_ string) error {
	*c++
	return nil
}

type intValue int

func (i *intValue) Set(val string) error {
	v, err := strconv.Atoi(val)
	*i = intValue(v)
	return err
}

type durationValue time.Duration

func (d *durationValue) Set(val string) error {
	v, err := time.ParseDuration(val)
	*d = durationValue(v)
	return err
}
