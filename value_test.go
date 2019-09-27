package flaq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStringValue(t *testing.T) {
	var svar string
	val := (*stringValue)(&svar)

	require.NoError(t, val.Set("ok"))
	require.Equal(t, svar, "ok")
}

func TestBoolValue(t *testing.T) {
	var bvar bool
	val := (*boolValue)(&bvar)

	require.NoError(t, val.Set(""))
	require.Equal(t, bvar, true)
}

func TestCountValue(t *testing.T) {
	var cvar int
	val := (*countValue)(&cvar)

	require.NoError(t, val.Set(""))
	require.NoError(t, val.Set(""))
	require.Equal(t, cvar, 2)
}

func TestIntValue(t *testing.T) {
	var ivar int
	val := (*intValue)(&ivar)

	require.Error(t, val.Set("invalid"))
	require.NoError(t, val.Set("2"))
	require.Equal(t, ivar, 2)
}

func TestDurationValue(t *testing.T) {
	var dvar time.Duration
	val := (*durationValue)(&dvar)

	require.Error(t, val.Set("invalid"))
	require.NoError(t, val.Set("5s"))
	require.Equal(t, 5*time.Second, dvar)
}

func TestFloat64Value(t *testing.T) {
	var fvar float64
	val := (*float64Value)(&fvar)

	require.Error(t, val.Set("invalid"))
	require.NoError(t, val.Set("3.14159265358979323846264338327950288419716939937510"))
	require.Equal(t, 3.14159265358979323846264338327950288419716939937510, fvar)
}
