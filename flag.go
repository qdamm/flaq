package flag

import (
	"fmt"
	"strings"
)

// FlagArg represents a flag argument. An empty Default indicates that the argument is required.
type FlagArg struct {
	Default string
}

// String returns a string flag with specified long/short form and description string.
func String(svar *string, long, short, description string) *Flag {
	return &Flag{Long: long, Short: short, Description: description, Arg: &FlagArg{}, Value: (*StringValue)(svar)}
}

// Bool returns a bool flag with specified long/short form and description string.
func Bool(bvar *bool, long, short, description string) *Flag {
	return &Flag{Long: long, Short: short, Description: description, Value: (*BoolValue)(bvar)}
}

// Help returns a help flag with specified long/short form and description string.
func Help(long, short, description string) *Flag {
	return &Flag{Long: long, Short: short, Description: description, Value: new(HelpValue)}
}

// Count returns a count flag with specified long/short form and description string.
func Count(cvar *int, long, short, description string) *Flag {
	return &Flag{Long: long, Short: short, Description: description, Value: (*CountValue)(cvar)}
}

// Flag is a representation for a command line option.
type Flag struct {
	Long        string
	Short       string
	Description string
	Arg         *FlagArg
	Value       Value
}

// ParserOpt represents a functional option for a Parser.
type ParserOpt func(f *Parser)

// Abbreviations indicates that option abbreviations are supported.
func Abbreviations() ParserOpt {
	return func(f *Parser) {
		f.abbreviations = true
	}
}

// Ordered indicates that command-line options are expected before operands.
// This means that the parsing will stop when an operand is seen.
func Ordered() ParserOpt {
	return func(f *Parser) {
		f.ordered = true
	}
}

// NewParser returns a new flag parser.
func NewParser(args []string, opts ...ParserOpt) *Parser {
	parser := &Parser{args: args}
	for _, opt := range opts {
		opt(parser)
	}
	return parser
}

// Parser is a parser for command-line options.
type Parser struct {
	flags         []*Flag
	seenArgs      []string
	args          []string
	abbreviations bool
	ordered       bool
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
func (f *Parser) Parse(flags ...*Flag) error {
	f.flags, f.seenArgs = flags, nil
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		return err
	}
}

// parseOne parses one flag. It reports whether a flag was seen.
func (f *Parser) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}

	arg := f.args[0]
	if len(arg) < 2 || arg[0] != '-' {
		if !f.ordered {
			f.seenArgs, f.args = append(f.seenArgs, arg), f.args[1:]
			return f.parseOne()
		}
		return false, nil
	}
	f.args = f.args[1:]

	if arg[1] == '-' {
		if len(arg) == 2 {
			return false, nil
		}
		return f.parseLong(arg[2:])
	}
	return f.parseShort(arg[1:])
}

// parseLong parses a long flag. It reports whether a flag was seen.
func (f *Parser) parseLong(name string) (bool, error) {
	var flagArg string
	var hasFlagArg bool

	for i := 1; i < len(name); i++ {
		if name[i] == '=' {
			name, flagArg = name[0:i], name[i+1:]
			hasFlagArg = true
			break
		}
	}

	var candidates []*Flag
	for _, flag := range f.flags {
		if strings.HasPrefix(flag.Long, name) {
			if f.abbreviations {
				candidates = append(candidates, flag)
			}
			if len(flag.Long) == len(name) {
				candidates = append(candidates[:0], flag)
				break
			}
		}
	}

	switch len(candidates) {
	case 0:
		return false, fmt.Errorf("unknown option --%s", name)
	case 1:
		flag := candidates[0]

		if hasFlagArg {
			if flag.Arg == nil {
				return false, fmt.Errorf("unexpected argument '%s' for option --%s", flagArg, name)
			}
		} else if flag.Arg != nil {
			if flag.Arg.Default != "" {
				flagArg = flag.Arg.Default
			} else if len(f.args) == 0 {
				return false, fmt.Errorf("missing argument for option --%s", name)
			} else {
				flagArg, f.args = f.args[0], f.args[1:]
			}
		}

		if err := flag.Value.Set(flagArg); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, fmt.Errorf("multiple options matching --%s", name)
}

func (f *Parser) parseShort(name string) (bool, error) {
	for _, flag := range f.flags {
		switch {
		case flag.Short != string(name[0]):
			continue

		case len(name) > 1:
			if flag.Arg == nil {
				if err := flag.Value.Set(""); err != nil {
					return false, err
				}
				return f.parseShort(name[1:])
			}
			if err := flag.Value.Set(name[1:]); err != nil {
				return false, err
			}

		default:
			if flag.Arg != nil && flag.Arg.Default == "" {
				if len(f.args) == 0 {
					return false, fmt.Errorf("missing value for option -%c", name[0])
				}
				if err := flag.Value.Set(f.args[0]); err != nil {
					return false, err
				}
				f.args = f.args[1:]
			} else if err := flag.Value.Set(""); err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, fmt.Errorf("unknown option -%c", name[0])
}

// Next advances the parser to the next argument, the bool value indicates whether an argument was read.
func (f *Parser) Next() (string, bool) {
	if len(f.args) == 0 {
		return "", false
	}
	next := f.args[0]
	f.args = f.args[1:]
	return next, true
}

// Operands returns the operands (the non-option arguments).
func (f *Parser) Operands() []string {
	return append(f.seenArgs, f.args...)
}
