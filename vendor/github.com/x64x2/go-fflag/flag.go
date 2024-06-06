package fflag

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Flag struct {
	// Short is the short-form flag name, parsed from CLI flags prefixed by '-'.
	Short string

	// Long is the long-form flag name, parsed from CLI flags prefixed by '--'.
	Long string

	// Env ...
	Env string

	// Usage is the description for this CLI flag, used when generating usage string.
	Usage string

	// Required specifies whether this flag is required, if so returning error.
	Required bool

	// Default is the default value for this CLI flag, used when no value provided.
	Default string

	// RawValues contains raw value strings populated during parse stage.
	RawValues []string

	// Value stores the actual value pointer that is set upon argument parsing.
	Value Value
}

// Name returns a CLI readable name for this flag, preferring long. e.g. "--long".
func (f Flag) Name() string {
	switch {
	case f.Long != "":
		return "--" + f.Long
	case f.Env != "":
		return f.Env
	case f.Short != "":
		return "-" + f.Short
	default:
		return ""
	}
}

// AppendUsage will append (new-line terminated) usage information for flag to b.
func (f Flag) AppendUsage(b []byte) []byte {
	if f.Short != "" {
		// Append short-flag
		b = append(b, " -"...)
		b = append(b, f.Short...)
	}

	if f.Long != "" {
		if f.Short == "" {
			// Add prefix whitespace
			b = append(b, "   "...)
		}

		// Append long-flag
		b = append(b, " --"...)
		b = append(b, f.Long...)
	}

	if kind := f.Value.Kind(); kind != "" {
		// Append type information
		b = append(b, ' ')
		b = append(b, kind...)
	}

	if f.Env != "" {
		// Append env-flag key
		b = append(b, " (env: $"...)
		b = append(b, f.Env...)
		b = append(b, ")"...)
	}

	if f.Default != "" {
		// Append default information
		b = append(b, " (default: "...)
		b = append(b, f.Default...)
		b = append(b, ")"...)
	}

	if f.Usage != "" {
		// Append usage information
		b = append(b, "\n    \t"...)
		b = append(b, strings.ReplaceAll(
			// Replace new-lines with
			// prefixed new-lines.
			f.Usage,
			"\n",
			"\n    \t",
		)...)
	}

	b = append(b, '\n')
	return b
}

// Letter returns the first letter associated
// with this flag, preferring short over long.
func (f *Flag) Letter() rune {
	if f.Short != "" {
		r, _ := utf8.DecodeRuneInString(f.Short)
		return r
	}
	if f.Long != "" {
		r, _ := utf8.DecodeRuneInString(f.Long)
		return r
	}
	panic("no flag name")
}

func sortFlags(flags []*Flag) {
	slices.SortFunc(flags, func(i, j *Flag) int {
		ir := i.Letter()
		jr := j.Letter()
		if ir < jr {
			return -1
		} else if ir > jr {
			return +1
		}
		return 0
	})
}
