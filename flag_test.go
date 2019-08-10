package flaq

import (
	"fmt"
	"testing"
	"time"

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
		number      int
		duration    time.Duration
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
		{
			args:   []string{"--number=50"},
			number: 50,
		},
		{
			args:     []string{"--duration=5m"},
			duration: time.Duration(5 * time.Minute),
		},
	}

	for _, f := range fixtures {
		t.Run(fmt.Sprintf("%q", f.args), func(t *testing.T) {
			var car, foo, fooBar bool
			var bar string
			var count, number int
			var duration time.Duration

			flags := &FlagSet{}
			flags.Bool(&foo, "foo", "f", "", false)
			flags.Bool(&fooBar, "foo-bar", "", "", false)
			flags.String(&bar, "bar", "b", "")
			flags.Bool(&car, "car", "c", "", false)
			flags.Count(&count, "count", "", "")
			flags.Int(&number, "number", "", "")
			flags.Duration(&duration, "duration", "", "")

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
			assert.Equal(t, f.number, number)
			assert.Equal(t, f.duration, duration)
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
	require.Equal(t, []string{"op1", "--bar", "ok", "op2"}, flags.Args())
	require.Equal(t, "", bar)
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
	require.True(t, bar)
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
		Name     string        `flaq:"-n, --name string         name of the person to greet"`
		Yell     bool          `flaq:"    --yell                whether to yell or not"`
		Bool     bool          `flaq:"    --bool bool           whether to bool or not"`
		Int      int           `flaq:"    --int int             whether to int or not"`
		Count    int           `flaq:"-c, --count count         whether to count or not"`
		Duration time.Duration `flaq:"    --duration duration   whether to duration or not"`

		RandomJSONField string `json:"randomField"`
	}{}

	flags := &FlagSet{}
	flags.Struct(&opts)

	err := flags.Parse([]string{
		"--name=ok",
		"--yell",
		"--bool=true",
		"--duration=3s",
		"--int=100",
		"-ccc",
	})
	require.NoError(t, err)

	require.Equal(t, "ok", opts.Name)
	require.True(t, opts.Yell)
	require.True(t, opts.Bool)
	require.Equal(t, 3*time.Second, opts.Duration)
	require.Equal(t, 100, opts.Int)
	require.Equal(t, 3, opts.Count)
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

func TestUsage(t *testing.T) {
	var bar bool
	var foo string

	flags := &FlagSet{}
	flags.String(&foo, "foo", "", "Foo to the foo")
	flags.Bool(&bar, "bar", "b", "Foo to the bar", false)

	expectedUsage := `Usage: flaq.test [options]

Options
  -b, --bar            Foo to the bar
      --foo <string>   Foo to the foo
`
	require.Equal(t, expectedUsage, flags.Usage())
}

func TestCustomUsage(t *testing.T) {
	customUsage := "This is custooooom"

	flags := &FlagSet{
		UsageFunc: func(_ *FlagSet) string {
			return customUsage
		},
	}

	require.Equal(t, customUsage, flags.Usage())
}
