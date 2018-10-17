package flag

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	fixtures := []struct {
		args        []string
		operands    []string
		expectError bool
		foo         bool
		fooBar      bool
		bar         string
		car         bool
	}{
		{
			args:     []string{"bar", "op1", "op2"},
			operands: []string{"bar", "op1", "op2"},
		},
		{
			args:     []string{"--bar", "ok", "op1", "op2"},
			operands: []string{"op1", "op2"},
			bar:      "ok",
		},
		{
			args:     []string{"--bar=ok", "op1", "op2"},
			operands: []string{"op1", "op2"},
			bar:      "ok",
		},
		{
			args:     []string{"--bar=ok", "op1", "op2"},
			operands: []string{"op1", "op2"},
			bar:      "ok",
		},
		{
			args:     []string{"--ba=ok", "op1", "op2"},
			operands: []string{"op1", "op2"},
			bar:      "ok",
		},
		{
			args:     []string{"-b", "ok", "op1", "op2"},
			operands: []string{"op1", "op2"},
			bar:      "ok",
		},
		{
			args:     []string{"-bok", "op1", "op2"},
			operands: []string{"op1", "op2"},
			bar:      "ok",
		},
		{
			args: []string{"-f"},
			foo:  true,
		},
		{
			args:     []string{"-fc", "op1"},
			operands: []string{"op1"},
			foo:      true,
			car:      true,
		},
		{
			args:        []string{"-b"},
			expectError: true,
		},
		{
			args: []string{"-fcb", "ok"},
			foo:  true,
			bar:  "ok",
			car:  true,
		},
		{
			args:     []string{"--foo", "--", "op1"},
			operands: []string{"op1"},
			foo:      true,
		},
		{
			args:   []string{"--foo-b"},
			fooBar: true,
		},
		{
			args: []string{"-cbbarval"},
			car:  true,
			bar:  "barval",
		},
	}

	for _, f := range fixtures {
		t.Run(fmt.Sprintf("%q", f.args), func(t *testing.T) {
			var car, foo, fooBar bool
			var bar string

			opts := &FlagSet{}
			opts.Bool(&foo, "foo", "f", "")
			opts.Bool(&fooBar, "foo-bar", "", "")
			opts.String(&bar, "bar", "b", "")
			opts.Bool(&car, "car", "c", "")

			err := opts.Parse(f.args)
			if f.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, f.operands, opts.Args())
			}
			assert.Equal(t, f.foo, foo)
			assert.Equal(t, f.fooBar, fooBar)
			assert.Equal(t, f.bar, bar)
			assert.Equal(t, f.car, car)
		})
	}
}

func TestParseContinuesOnArg(t *testing.T) {
	var bar string

	opts := &FlagSet{}
	opts.String(&bar, "bar", "", "")

	err := opts.Parse([]string{"op1", "--bar", "ok", "op2"}, ContinueOnArg())
	require.NoError(t, err)
	assert.Equal(t, []string{"op1", "op2"}, opts.Args())
	assert.Equal(t, "ok", bar)
}

func TestParseHelp(t *testing.T) {
	opts := &FlagSet{}
	opts.Help("help", "", "")

	err := opts.Parse([]string{"--help"})
	require.Equal(t, ErrHelp, err)
}
