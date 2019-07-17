package flaq

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func defaultUsage(flags *FlagSet) string {
	usage := flags.UsageLine
	if usage == "" {
		usage = "Usage: " + filepath.Base(os.Args[0]) + " [options]\n"
	}
	usage += "\nOptions\n"

	usages, maxUsageLen := make(map[*Flag]string), 0
	sortedFlags := make([]*Flag, 0, len(flags.flags))

	flags.VisitAll(func(f *Flag) {
		if !f.Hidden {
			usages[f] = flagUsage(f)
			if len(usages[f]) > maxUsageLen {
				maxUsageLen = len(usages[f])
			}
			for i := range sortedFlags {
				if sortedFlags[i].Long+sortedFlags[i].Short > f.Long+f.Short {
					sortedFlags = append(sortedFlags, nil)
					copy(sortedFlags[i+1:], sortedFlags[i:])
					sortedFlags[i] = f
					return
				}
			}
			sortedFlags = append(sortedFlags, f)
		}
	})
	if maxUsageLen > 25 {
		maxUsageLen = 25
	}
	for _, f := range sortedFlags {
		usage += fmt.Sprintf("  %-"+strconv.Itoa(maxUsageLen)+"s   %s\n", usages[f], f.Description)
	}
	return usage
}

// flagUsage returns the help usage for a given flag.
func flagUsage(f *Flag) string {
	if f.Usage != "" {
		return f.Usage
	}
	usage := "    "
	if f.Short != "" {
		usage = "-" + f.Short
		if f.Long != "" {
			usage += ", "
		}
	}

	if f.Long != "" {
		usage += "--" + f.Long
	}

	if f.Arg != nil {
		argName := f.Arg.Name
		if argName == "" {
			argName = "arg"
		}
		switch {
		case f.Long != "" && f.Arg.Default != "":
			usage += fmt.Sprintf("=<%s>", argName)
		case f.Short != "" && f.Arg.Default != "":
			usage += fmt.Sprintf("<%s>", argName)
		default:
			usage += fmt.Sprintf(" <%s>", argName)
		}
	}
	return usage
}
