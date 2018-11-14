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
		count       int
	}{
		{
			args:     []string{"bar", "op1", "op2"},
			operands: []string{"bar", "op1", "op2"},
		},
		{
			args:     []string{"op1", "--bar", "ok", "op2"},
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
			args:        []string{"--foo-b"},
			expectError: true,
		},
		{
			args: []string{"-cbbarval"},
			car:  true,
			bar:  "barval",
		},
		{
			args:  []string{"--count", "--count"},
			count: 2,
		},
		{
			args:        []string{"--foo=ok"},
			expectError: true,
		},
		{
			args:        []string{"--bar"},
			expectError: true,
		},
	}

	for _, f := range fixtures {
		t.Run(fmt.Sprintf("%q", f.args), func(t *testing.T) {
			var car, foo, fooBar bool
			var bar string
			var count int

			parser := NewParser(f.args)
			err := parser.Parse(
				Bool(&foo, "foo", "f", ""),
				Bool(&fooBar, "foo-bar", "", ""),
				String(&bar, "bar", "b", ""),
				Bool(&car, "car", "c", ""),
				Count(&count, "count", "", ""),
			)
			if f.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, f.operands, parser.Operands())
			}
			assert.Equal(t, f.foo, foo)
			assert.Equal(t, f.fooBar, fooBar)
			assert.Equal(t, f.bar, bar)
			assert.Equal(t, f.car, car)
			assert.Equal(t, f.count, count)
		})
	}
}

func TestParseOrderedArgs(t *testing.T) {
	var bar string
	parser := NewParser([]string{"op1", "--bar", "ok", "op2"}, Ordered())

	err := parser.Parse(String(&bar, "bar", "", ""))
	require.NoError(t, err)
	assert.Equal(t, []string{"op1", "--bar", "ok", "op2"}, parser.Operands())
	assert.Equal(t, "", bar)
}

func TestParseAbbreviation(t *testing.T) {
	var bar, foo bool
	opts := NewParser([]string{"--foo-b"}, Abbreviations())

	err := opts.Parse(
		Bool(&bar, "foo-bar", "", ""),
		Bool(&foo, "foo-foo", "", ""),
	)
	require.NoError(t, err)
	assert.True(t, bar)
}

func TestParseAmbiguousAbbreviation(t *testing.T) {
	var bar, foo bool
	opts := NewParser([]string{"--foo"}, Abbreviations())

	err := opts.Parse(
		Bool(&bar, "foo-bar", "", ""),
		Bool(&foo, "foo-foo", "", ""),
	)
	require.Error(t, err)
}

func TestParseHelp(t *testing.T) {
	err := NewParser([]string{"--help"}).Parse(Help("help", "", ""))
	require.Equal(t, ErrHelp, err)
}
