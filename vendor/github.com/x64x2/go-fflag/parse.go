package fflag

import (
	"fmt"
	"strings"
)

// ParseEnv parses environment variables from slice into FlagSet.
func (set *FlagSet) ParseEnv(env []string) (err error) {
	err = set.readEnv(env)
	if err != nil {
		return err
	}
	return set.parse()
}

// ParseArgs parses CLI arguments from slices into FlagSet, returning unused arguments.
func (set *FlagSet) ParseArgs(args []string) (unused []string, err error) {
	unused, err = set.readArgs(args)
	if err != nil {
		return nil, err
	}
	if err = set.parse(); err != nil {
		return nil, err
	}
	return unused, nil
}

// Parse parses environment variables, followed by CLI arguments, into FlagSet, returning unused arguments.
func (set *FlagSet) Parse(args, env []string) (unused []string, err error) {
	err = set.readEnv(env)
	if err != nil {
		return nil, err
	}
	unused, err = set.readArgs(args)
	if err != nil {
		return nil, err
	}
	if err = set.parse(); err != nil {
		return nil, err
	}
	return unused, nil
}

// MustParse calls FlagSet.Parse(...), panicking on error.
func (set *FlagSet) MustParse(args, env []string) (unused []string) {
	var err error
	unused, err = set.Parse(args, env)
	if err != nil {
		panic(err)
	}
	return unused
}

func (set *FlagSet) readEnv(env []string) (err error) {
	for _, kv := range env {
		// Look for separating '='.
		i := strings.Index(kv, "=")
		if i == -1 {
			return fmt.Errorf("malformed env: %s", env)
		}

		// Separate into key + value.
		key, val := kv[:i], kv[i+1:]

		// Look for flag with key.
		flag := set.GetEnv(key)
		if flag == nil {
			// unknown
			continue
		}

		// Append next arg value to existing values.
		flag.RawValues = append(flag.RawValues, val)
	}
	return nil
}

func (set *FlagSet) readArgs(args []string) (unused []string, err error) {
	for i := 0; i < len(args); i++ {
		var (
			// Current arg.
			arg = args[i]

			// Matching flag for arg.
			flag *Flag

			// Current flag's matching
			// string value (if set).
			value *string
		)

		if len(arg) == 0 || len(arg) == 1 || arg[0] != '-' {
			// Unrecognizable argument, add to unused.
			unused = append(unused, arg)
			continue
		}

		// Check if compound value e.g. --verbose=true .
		if idx := strings.IndexByte(arg, '='); idx > 0 {
			v := arg[idx+1:]
			value = &v
			arg = arg[:idx]
		}

		if arg[1] == '-' {
			// Prefix='--' parse long flag.
			flag = set.GetLong(arg[2:])
		} else {
			// Prefix='-' parse short flag.
			flag = set.GetShort(arg[1:])
		}

		if flag == nil {
			// Unrecognized argument.
			unused = append(unused, arg)
			continue
		}

		if value == nil {
			if _, ok := flag.Value.(*boolVar); ok {
				// Only booleans allow no value.
				flag.RawValues = []string{"true"}
				continue
			}

			if i == len(args)-1 {
				// Value expected and none was found.
				return nil, fmt.Errorf("flag %q expects a value", arg)
			}

			// Grab next as value.
			value = &args[i+1]
			i++
		}

		// Append next arg value to existing values slice.
		flag.RawValues = append(flag.RawValues, *value)
	}
	return unused, nil
}

func (set *FlagSet) parse() (err error) {
loop:
	for _, flag := range set.flags {
		switch {
		// Values were provided for flag.
		case flag.RawValues != nil:

		// Flag is required but no values.
		case flag.Required:
			err = fmt.Errorf("flag %q is required", flag.Name())
			break loop

		// Flag is unset but has default.
		case flag.Default != "":
			flag.RawValues = []string{flag.Default}

		// Unset flag, ignore.
		default:
			continue loop
		}

		for _, value := range flag.RawValues {
			// Attempt to parse flag value from raw values.
			if err = flag.Value.Set(value); err != nil {
				err = fmt.Errorf("error parsing %q value: %w", flag.Name(), err)
				break loop
			}
		}
	}

	// Call hooks, even on error.
	for _, fn := range set.funcs {
		fn()
	}

	return err
}
