package flaq

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

			flags := &FlagSet{}
			flags.Bool(&foo, "foo", "f", "", false)
			flags.Bool(&fooBar, "foo-bar", "", "", false)
			flags.String(&bar, "bar", "b", "")
			flags.Bool(&car, "car", "c", "", false)
			flags.Count(&count, "count", "", "")

			err := flags.Parse(f.args)
			if f.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, f.operands, flags.Args())
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

	flags := &FlagSet{
		Ordered: true,
	}
	flags.String(&bar, "bar", "", "")

	err := flags.Parse([]string{"op1", "--bar", "ok", "op2"})
	require.NoError(t, err)
	assert.Equal(t, []string{"op1", "--bar", "ok", "op2"}, flags.Args())
	assert.Equal(t, "", bar)
}

func TestParseAbbreviation(t *testing.T) {
	var bar, foo bool

	flags := &FlagSet{
		Abbreviations: true,
	}
	flags.Bool(&bar, "foo-bar", "", "", false)
	flags.Bool(&foo, "foo-foo", "", "", false)

	err := flags.Parse([]string{"--foo-b"})
	require.NoError(t, err)
	assert.True(t, bar)
}

func TestParseAmbiguousAbbreviation(t *testing.T) {
	var bar, foo bool

	flags := &FlagSet{
		Abbreviations: true,
	}
	flags.Bool(&bar, "foo-bar", "", "", false)
	flags.Bool(&foo, "foo-foo", "", "", false)

	err := flags.Parse([]string{"--foo"})
	require.Error(t, err)
}

func TestParseStruct(t *testing.T) {
	var opts = struct {
		Name string `flaq:"-n, --name string    name of the person to greet"`
		Yell bool   `flaq:"    --yell           whether to yell or not"`
	}{}

	flags := &FlagSet{}
	flags.Struct(&opts)

	err := flags.Parse([]string{"--name=ok", "--yell"})
	require.NoError(t, err)

	require.Equal(t, "ok", opts.Name)
	require.True(t, opts.Yell)
}

func TestParseStructFieldTag(t *testing.T) {
	fixtures := []struct {
		tag         string
		short       string
		long        string
		fieldType   string
		description string
	}{
		{
			tag:         "--yell  whether to yell or not",
			long:        "yell",
			description: "whether to yell or not",
		},
		{
			tag:         "-y  whether to yell or not",
			short:       "y",
			description: "whether to yell or not",
		},
		{
			tag:         "-y, --yell bool whether to yell or not",
			short:       "y",
			long:        "yell",
			fieldType:   "bool",
			description: "whether to yell or not",
		},
		{
			tag:         "    --yell  whether to yell or not",
			long:        "yell",
			description: "whether to yell or not",
		},
	}

	for _, fixture := range fixtures {
		flag, fieldType := parseStructFieldTag(fixture.tag)
		require.Equal(t, fixture.short, flag.Short)
		require.Equal(t, fixture.long, flag.Long)
		require.Equal(t, fixture.fieldType, fieldType)
		require.Equal(t, fixture.description, flag.Description)
	}
}
