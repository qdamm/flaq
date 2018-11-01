package flag

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringValue(t *testing.T) {
	var svar string
	val := (*StringValue)(&svar)

	require.NoError(t, val.Set("ok"))
	require.Equal(t, svar, "ok")
}

func TestBoolValue(t *testing.T) {
	var bvar bool
	val := (*BoolValue)(&bvar)

	require.NoError(t, val.Set(""))
	require.Equal(t, bvar, true)
}

func TestCountValue(t *testing.T) {
	var cvar int
	val := (*CountValue)(&cvar)

	require.NoError(t, val.Set(""))
	require.NoError(t, val.Set(""))
	require.Equal(t, cvar, 2)
}

func TestHelpValue(t *testing.T) {
	val := new(HelpValue)
	require.Equal(t, ErrHelp, val.Set(""))
}
