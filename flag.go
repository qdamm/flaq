package flag

import (
	"fmt"
	"strings"
)

// FlagArg represents a flag argument. An empty Default indicates that the argument is required.
type FlagArg struct {
	Default string
}

// StringFlag returns a string flag with specified long/short form and description string.
// The argument svar points to a string variable in which to store the value of the flag.
func StringFlag(svar *string, long, short, description string) *Flag {
	return &Flag{Long: long, Short: short, Description: description, Arg: &FlagArg{}, Value: (*stringVal)(svar)}
}

// BoolFlag returns a bool flag with specified long/short form and description string.
// The argument bvar points to a bool variable in which to store the value of the flag.
func BoolFlag(bvar *bool, long, short, description string) *Flag {
	return &Flag{Long: long, Short: short, Description: description, Value: (*boolVal)(bvar)}
}

// HelpFlag returns a help flag with specified long/short form and description string.
func HelpFlag(long, short, description string) *Flag {
	return &Flag{Long: long, Short: short, Description: description, Value: new(helpVal)}
}

// Flag is a representation for a command line option.
type Flag struct {
	Long        string
	Short       string
	Description string
	Arg         *FlagArg
	Value       Value
}

// FlagSetOpt represents a functional option for a FlagSet.
type FlagSetOpt func(f *FlagSet)

// ContinueOnArg indicates that options parsing should continue when an argument (aka. operand) is encountered.
func ContinueOnArg() FlagSetOpt {
	return func(f *FlagSet) {
		f.continueOnArg = true
	}
}

// FlagSet represents a set of flags.
type FlagSet struct {
	flags         []*Flag
	seenArgs      []string
	args          []string
	continueOnArg bool
}

// String adds a string flag to the FlagSet.
func (f *FlagSet) String(svar *string, long, short, description string) {
	f.Add(StringFlag(svar, long, short, description))
}

// Bool adds a bool flag to the FlagSet.
func (f *FlagSet) Bool(bvar *bool, long, short, description string) {
	f.Add(BoolFlag(bvar, long, short, description))
}

// Help adds a help flag to the FlagSet.
func (f *FlagSet) Help(long, short, description string) {
	f.Add(HelpFlag(long, short, description))
}

// Add adds a flag to the FlagSet.
func (f *FlagSet) Add(flag *Flag) {
	f.flags = append(f.flags, flag)
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
func (f *FlagSet) Parse(args []string, opts ...FlagSetOpt) error {
	for _, opt := range opts {
		opt(f)
	}
	f.args, f.seenArgs = args, nil
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		return err
	}
}

// parseOne parses one flag. It reports whether a flag was seen.
func (f *FlagSet) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}

	arg := f.args[0]
	if len(arg) < 2 || arg[0] != '-' {
		if f.continueOnArg {
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
func (f *FlagSet) parseLong(name string) (bool, error) {
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
			candidates = append(candidates, flag)
			if len(flag.Long) == len(name) {
				candidates = candidates[len(candidates)-1:]
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

func (f *FlagSet) parseShort(name string) (bool, error) {
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

// Args returns the non-flag arguments (also known as "operands").
func (f *FlagSet) Args() []string {
	return append(f.seenArgs, f.args...)
}
