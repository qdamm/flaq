package flaq

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

type errorHandling int

// These constants cause FlagSet.Parse to behave as described if the parse fails.
const (
	ContinueOnError errorHandling = iota // Return a descriptive error.
	ExitOnError                          // Print a descriptive error to os.Stderr and call os.Exit(2).
	PanicOnError                         // Call panic with a descriptive error.
)

var flags = &FlagSet{}

func init() {
	flags.ErrorHandling = ExitOnError
}

// String adds a string flag with specified long/short form and description.
func String(svar *string, long, short, description string) {
	flags.String(svar, long, short, description)
}

// Bool adds a bool flag with specified long/short form and description.
// hasArg indicates whether the flag accepts an optional argument (eg. --flag=false).
func Bool(bvar *bool, long, short, description string, hasArg bool) {
	flags.Bool(bvar, long, short, description, hasArg)
}

// Int adds an int flag with specified long/short form and description.
func Int(ivar *int, long, short, description string) {
	flags.Int(ivar, long, short, description)
}

// Float64 adds a float64 flag with specified long/short form and description.
func Float64(fvar *float64, long, short, description string) {
	flags.Float64(fvar, long, short, description)
}

// Count adds a count flag with specified long/short form and description.
func Count(cvar *int, long, short, description string) {
	flags.Count(cvar, long, short, description)
}

// Duration adds a duration flag with specified long/short form and description.
func Duration(dvar *time.Duration, long, short, description string) {
	flags.Duration(dvar, long, short, description)
}

// Help sets the help flag's long/short form and description.
func Help(long, short, description string) {
	flags.Help(long, short, description)
}

// Struct adds flags using reflection and struct field tags.
func Struct(svar interface{}) {
	flags.Struct(svar)
}

// Add adds a new flag.
func Add(flag *Flag) {
	flags.Add(flag)
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
func Parse() error {
	return flags.Parse(os.Args[1:])
}

// Args returns the remaining arguments once options have been parsed.
func Args() []string {
	return flags.Args()
}

// Flag is a representation for a command line option.
type Flag struct {
	// Long option name.
	Long string

	// Short option name.
	Short string

	// Description for the option, as it will appear in help usage.
	Description string

	// Specifies whether or not the option accepts an argument:
	//   - when set to nil, the option doesn't accept any argument
	//   - when set to &FlagArg{}, the option requires an argument
	//   - when set to &FlagArg{"default", ""}, the option accepts
	//     an optional argument which defaults to "default".
	Arg *FlagArg

	// Value is the interface to the dynamic value stored in the flag.
	Value Value

	// Hidden indicates that the option should be hidden in help usage.
	Hidden bool

	// Usage can be set to overwrite the option usage line in the help usage.
	Usage string

	// Final indicates that option parsing should stop when the option is seen.
	// This is generally not to be used, except for some options such as --help or --version.
	Final bool
}

// FlagArg represents a flag argument. An empty Default indicates that the argument is required.
type FlagArg struct {
	Default string
	Name    string
}

// A FlagSet represents a set of defined flags. The zero value of a FlagSet
// has ContinueOnError error handling and doesn't add a help flag.
type FlagSet struct {
	// Abbreviations indicates that option abbreviations are supported.
	Abbreviations bool

	// DisableHelp disables help usage.
	DisableHelp bool

	// ErrorHandling defines how FlagSet.Parse behaves when parsing fails.
	ErrorHandling errorHandling

	// Ordered indicates that command-line options are expected before operands.
	// This means that the parsing will stop when an operand is seen.
	Ordered bool

	// UsageLine can be set to overwrite the flag help usage.
	UsageLine string

	// UsageFunc
	UsageFunc func(*FlagSet) string

	flags    []*Flag
	seenArgs []string
	args     []string
	help     bool
	helpFlag *Flag
}

// String adds a string flag with specified long/short form and description.
func (f *FlagSet) String(svar *string, long, short, description string) {
	f.Add(&Flag{
		Long:        long,
		Short:       short,
		Description: description,
		Value:       (*stringValue)(svar),
		Arg:         &FlagArg{Name: "string"},
	})
}

// Bool adds a bool flag with specified long/short form and description.
// hasArg indicates whether the flag accepts an optional argument (eg. --flag=false).
func (f *FlagSet) Bool(bvar *bool, long, short, description string, hasArg bool) {
	flag := &Flag{
		Long:        long,
		Short:       short,
		Description: description,
		Value:       (*boolValue)(bvar),
	}
	if hasArg {
		flag.Arg = &FlagArg{Default: "true", Name: "bool"}
	}
	f.Add(flag)
}

// Int adds an int flag with specified long/short form and description.
func (f *FlagSet) Int(ivar *int, long, short, description string) {
	f.Add(&Flag{
		Long:        long,
		Short:       short,
		Description: description,
		Value:       (*intValue)(ivar),
		Arg:         &FlagArg{Name: "int"},
	})
}

// Float64 adds a float64 flag with specified long/short form and description.
func (f *FlagSet) Float64(fvar *float64, long, short, description string) {
	f.Add(&Flag{
		Long:        long,
		Short:       short,
		Description: description,
		Value:       (*float64Value)(fvar),
		Arg:         &FlagArg{Name: "float"},
	})
}

// Count adds a count flag with specified long/short form and description.
func (f *FlagSet) Count(cvar *int, long, short, description string) {
	f.Add(&Flag{
		Long:        long,
		Short:       short,
		Description: description,
		Value:       (*countValue)(cvar),
	})
}

// Duration adds a duration flag with specified long/short form and description.
func (f *FlagSet) Duration(dvar *time.Duration, long, short, description string) {
	f.Add(&Flag{
		Long:        long,
		Short:       short,
		Description: description,
		Value:       (*durationValue)(dvar),
		Arg:         &FlagArg{},
	})
}

// Help sets the help flag's long/short form and description.
func (f *FlagSet) Help(long, short, description string) {
	helpFlag := &Flag{
		Long:        long,
		Short:       short,
		Description: description,
		Value:       (*boolValue)(&f.help),
		Final:       true,
	}
	if f.helpFlag == nil {
		flags.Add(helpFlag)
		f.helpFlag = helpFlag
	} else {
		*f.helpFlag = *helpFlag
	}
}

// Struct adds flags using reflection and struct field tags.
func (f *FlagSet) Struct(svar interface{}) {
	sval := reflect.ValueOf(svar)
	if svar == nil || sval.Kind() != reflect.Ptr || sval.Elem().Kind() != reflect.Struct {
		panic("can only accept a pointer to a struct")
	}
	sval = sval.Elem()

	for i := 0; i < sval.NumField(); i++ {
		field := sval.Field(i)
		val := field.Addr().Interface()
		tag, ok := sval.Type().Field(i).Tag.Lookup("flaq")
		if !ok {
			continue
		}
		flag, fieldType := parseStructFieldTag(tag)
		switch fieldType {
		case "count":
			flag.Value = (*countValue)(val.(*int))
		case "duration":
			flag.Value = (*durationValue)(val.(*time.Duration))
			flag.Arg = &FlagArg{Name: "duration"}
		case "float64":
			flag.Value = (*float64Value)(val.(*float64))
			flag.Arg = &FlagArg{Name: "float"}
		case "int":
			flag.Value = (*intValue)(val.(*int))
			flag.Arg = &FlagArg{Name: "int"}
		case "string":
			flag.Value = (*stringValue)(val.(*string))
			flag.Arg = &FlagArg{Name: "string"}
		case "bool":
			flag.Value = (*boolValue)(val.(*bool))
			flag.Arg = &FlagArg{Default: "true", Name: "bool"}
		case "":
			flag.Value = (*boolValue)(val.(*bool))
		default:
			panic(fmt.Sprintf(`unknown struct field type "%s"`, fieldType))
		}
		f.Add(flag)
	}
}

// Add adds a flag to the flagset.
func (f *FlagSet) Add(flag *Flag) {
	f.flags = append(f.flags, flag)
}

// Usage returns help usage.
func (f *FlagSet) Usage() string {
	if f.UsageFunc != nil {
		return f.UsageFunc(f)
	}
	return defaultUsage(f)
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
func (f *FlagSet) Parse(args []string) error {
	if f.helpFlag == nil && !f.DisableHelp {
		flags.Help("help", "h", "show usage help")
	}
	f.args = args
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			if f.help {
				// When the default --help flag is encountered, print help usage
				// to stdout and exit. To change this behaviour, one should set
				// DisableHelp and implement their own help flag instead.
				fmt.Print(f.Usage())
				os.Exit(0)
			}
			break
		}
		switch f.ErrorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}

// parseOne parses one flag. It reports whether a flag was seen.
func (f *FlagSet) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}

	arg := f.args[0]
	if len(arg) < 2 || arg[0] != '-' {
		if !f.Ordered {
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
			if f.Abbreviations {
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
		return !flag.Final, nil
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
		return !flag.Final, nil
	}
	return false, fmt.Errorf("unknown option -%c", name[0])
}

// Args returns the remaining arguments once options have been parsed.
func (f *FlagSet) Args() []string {
	return append(f.seenArgs, f.args...)
}

// VisitAll visits all the flags, calling fn for each. It visits all flags, even those not set.
func (f *FlagSet) VisitAll(fn func(*Flag)) {
	for _, flag := range f.flags {
		fn(flag)
	}
}

func parseStructFieldTag(tag string) (*Flag, string) {
	flag := &Flag{}
	var fieldType string
	var parsing *string
	for j := 0; j < len(tag); j++ {
		switch parsing {
		case nil:
			if tag[j] == '-' {
				parsing = &flag.Short
			}
		case &flag.Short:
			switch tag[j] {
			case '-':
				parsing = &flag.Long
			case ',':
				parsing = &flag.Long
				j += 3
			case ' ':
				parsing = &fieldType
			default:
				flag.Short = string(tag[j])
			}
		case &flag.Long:
			if tag[j] == ' ' {
				parsing = &fieldType
			} else {
				flag.Long += string(tag[j])
			}
		case &fieldType:
			if tag[j] == ' ' {
				parsing = &flag.Description
			} else {
				fieldType += string(tag[j])
			}
		case &flag.Description:
			if flag.Description != "" || tag[j] != ' ' {
				flag.Description += string(tag[j])
			}
		}
	}
	return flag, fieldType
}
