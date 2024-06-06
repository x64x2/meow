package fflag

import (
	"os"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"

	"codeberg.org/gruf/go-bytesize"
	"codeberg.org/gruf/go-byteutil"
)

type FlagSet struct {
	flags []*Flag
	funcs []func()
	auto  func(string) string
}

// GetShort will fetch the flag with matching short name.
func (set *FlagSet) GetShort(name string) *Flag {
	for _, flag := range set.flags {
		if flag.Short == name {
			return flag
		}
	}
	return nil
}

// GetLong will fetch the flag with matching long name.
func (set *FlagSet) GetLong(name string) *Flag {
	for _, flag := range set.flags {
		if flag.Long == name {
			return flag
		}
	}
	return nil
}

// GetEnv will fetch the flag with matching env name.
func (set *FlagSet) GetEnv(name string) *Flag {
	for _, flag := range set.flags {
		if flag.Env == name {
			return flag
		}
	}
	return nil
}

// Sort will sort the internal slice of flags, first by short name, then long name.
func (set *FlagSet) Sort() { sortFlags(set.flags) }

// Usage returns a string showing typical terminal-formatted usage information for all of the registered flags.
func (set *FlagSet) Usage() string {
	// Sort all flags.
	sortFlags(set.flags)

	// Allocate string buffer
	l := len(set.flags) * 32
	b := make([]byte, 0, l)
	buf := byteutil.Buffer{B: b}

	// Append each flag usage info.
	for _, flag := range set.flags {
		buf.B = flag.AppendUsage(buf.B)
	}

	return buf.String()
}

// AutoEnv will automatically assign environment keys to flags with a long-name provided, if nil will use a default function.
func (set *FlagSet) AutoEnv(fn func(long string) (env string)) {
	if fn == nil {
		fn = func(long string) string {
			env := long
			env = strings.ReplaceAll(env, "-", "_")
			env = strings.ToUpper(env)
			return env
		}
	}
	set.auto = fn
}

// Hook allows registering a function hook to be called directly after parse.
func (set *FlagSet) Hook(fn func()) {
	set.funcs = append(set.funcs, fn)
}

// Config registers a function hook to open path and pass to given 'fn' on supplied '-c' or '--config' value.
func (set *FlagSet) Config(fn func(*os.File) error) {
	path := set.String("c", "config", "", "Configuration file path")
	set.Hook(func() {
		// Check path was set.
		if *path == "" {
			return
		}

		// Open file with name.
		file, err := os.Open(*path)
		if err != nil {
			panic(err)
		}

		// Ensure closed.
		defer file.Close()

		// Call the given file handler.
		if err := fn(file); err != nil {
			panic(err)
		}
	})
	// don't make accessible by env.
	set.GetShort("c").Env = ""
}

// Help registers a function hook to print usage string and exit with code = 0 on '-h' or '--help'.
func (set *FlagSet) Help() {
	set.Func("h", "help", "Print usage information", func() {
		str := "Usage: " + os.Args[0] + " ...\n" + set.Usage()
		os.Stdout.WriteString(str)
		syscall.Exit(0)
	})
	// don't make accessible by env.
	set.GetShort("h").Env = ""
}

// Version registers a function hook to print given version string and exit with code = 0 on '-v' or '--version'.
func (set *FlagSet) Version(version string) {
	set.Func("v", "version", "Print version information", func() {
		os.Stdout.WriteString(version + "\n")
		syscall.Exit(0)
	})
	// don't make accessible by env.
	set.GetShort("v").Env = ""
}

// Func registers a function hook to be called on boolean result = true of given short,long,env Flag.
func (set *FlagSet) Func(short string, long string, usage string, fn func()) {
	ptr := set.Bool(short, long, false, usage)
	set.Hook(func() {
		if *ptr {
			fn()
		}
	})
}

func (set *FlagSet) Bool(short string, long string, value bool, usage string) *bool {
	ptr := new(bool)
	set.BoolVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) BoolVar(ptr *bool, short string, long string, value bool, usage string) {
	var vstr string
	if value {
		vstr = (*boolVar)(&value).String()
	}
	set.Var((*boolVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) Int(short string, long string, value int, usage string) *int {
	ptr := new(int)
	set.IntVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) IntVar(ptr *int, short string, long string, value int, usage string) {
	var vstr string
	if value != 0 {
		vstr = (*intVar)(&value).String()
	}
	set.Var((*intVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) Int64(short string, long string, value int64, usage string) *int64 {
	ptr := new(int64)
	set.Int64Var(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) Int64Var(ptr *int64, short string, long string, value int64, usage string) {
	var vstr string
	if value != 0 {
		vstr = (*int64Var)(&value).String()
	}
	set.Var((*int64Var)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) Uint(short string, long string, value uint, usage string) *uint {
	ptr := new(uint)
	set.UintVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) UintVar(ptr *uint, short string, long string, value uint, usage string) {
	var vstr string
	if value != 0 {
		vstr = (*uintVar)(&value).String()
	}
	set.Var((*uintVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) Uint64(short string, long string, value uint64, usage string) *uint64 {
	ptr := new(uint64)
	set.Uint64Var(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) Uint64Var(ptr *uint64, short string, long string, value uint64, usage string) {
	var vstr string
	if value != 0 {
		vstr = (*uint64Var)(&value).String()
	}
	set.Var((*uint64Var)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) Float(short string, long string, value float64, usage string) *float64 {
	ptr := new(float64)
	set.FloatVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) FloatVar(ptr *float64, short string, long string, value float64, usage string) {
	var vstr string
	if value != 0 {
		vstr = (*float64Var)(&value).String()
	}
	set.Var((*float64Var)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) String(short string, long string, value string, usage string) *string {
	ptr := new(string)
	set.StringVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) StringVar(ptr *string, short string, long string, value string, usage string) {
	set.Var((*stringVar)(ptr), short, long, value, usage)
}

func (set *FlagSet) Size(short string, long string, value bytesize.Size, usage string) *bytesize.Size {
	ptr := new(bytesize.Size)
	set.SizeVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) SizeVar(ptr *bytesize.Size, short string, long string, value bytesize.Size, usage string) {
	var vstr string
	if value != 0 {
		vstr = value.StringIEC()
	}
	set.Var((*sizeVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) Duration(short string, long string, value time.Duration, usage string) *time.Duration {
	ptr := new(time.Duration)
	set.DurationVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) DurationVar(ptr *time.Duration, short string, long string, value time.Duration, usage string) {
	var vstr string
	if value != 0 {
		vstr = value.String()
	}
	set.Var((*durationVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) Bytes(short string, long string, value []byte, usage string) *[]byte {
	ptr := new([]byte)
	set.BytesVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) BytesVar(ptr *[]byte, short string, long string, value []byte, usage string) {
	set.Var((*bytesVar)(ptr), short, long, *(*string)(unsafe.Pointer(&value)), usage)
}

func (set *FlagSet) BoolSlice(short string, long string, value []bool, usage string) *[]bool {
	ptr := new([]bool)
	set.BoolSliceVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) BoolSliceVar(ptr *[]bool, short string, long string, value []bool, usage string) {
	var vstr string
	if value != nil {
		vstr = (*boolSliceVar)(&value).String()
	}
	set.Var((*boolSliceVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) IntSlice(short string, long string, value []int, usage string) *[]int {
	ptr := new([]int)
	set.IntSliceVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) IntSliceVar(ptr *[]int, short string, long string, value []int, usage string) {
	var vstr string
	if value != nil {
		vstr = (*intSliceVar)(&value).String()
	}
	set.Var((*intSliceVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) Int64Slice(short string, long string, value []int64, usage string) *[]int64 {
	ptr := new([]int64)
	set.Int64SliceVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) Int64SliceVar(ptr *[]int64, short string, long string, value []int64, usage string) {
	var vstr string
	if value != nil {
		vstr = (*int64SliceVar)(&value).String()
	}
	set.Var((*int64SliceVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) UintSlice(short string, long string, value []uint, usage string) *[]uint {
	ptr := new([]uint)
	set.UintSliceVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) UintSliceVar(ptr *[]uint, short string, long string, value []uint, usage string) {
	var vstr string
	if value != nil {
		vstr = (*uintSliceVar)(&value).String()
	}
	set.Var((*uintSliceVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) Uint64Slice(short string, long string, value []uint64, usage string) *[]uint64 {
	ptr := new([]uint64)
	set.Uint64SliceVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) Uint64SliceVar(ptr *[]uint64, short string, long string, value []uint64, usage string) {
	var vstr string
	if value != nil {
		vstr = (*uint64SliceVar)(&value).String()
	}
	set.Var((*uint64SliceVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) FloatSlice(short string, long string, value []float64, usage string) *[]float64 {
	ptr := new([]float64)
	set.FloatSliceVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) FloatSliceVar(ptr *[]float64, short string, long string, value []float64, usage string) {
	var vstr string
	if value != nil {
		vstr = (*float64SliceVar)(&value).String()
	}
	set.Var((*float64SliceVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) StringSlice(short string, long string, value []string, usage string) *[]string {
	ptr := new([]string)
	set.StringSliceVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) StringSliceVar(ptr *[]string, short string, long string, value []string, usage string) {
	var vstr string
	if value != nil {
		vstr = (*stringSliceVar)(&value).String()
	}
	set.Var((*stringSliceVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) SizeSlice(short string, long string, value []bytesize.Size, usage string) *[]bytesize.Size {
	ptr := new([]bytesize.Size)
	set.SizeSliceVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) SizeSliceVar(ptr *[]bytesize.Size, short string, long string, value []bytesize.Size, usage string) {
	var vstr string
	if value != nil {
		vstr = (*sizeSliceVar)(&value).String()
	}
	set.Var((*sizeSliceVar)(ptr), short, long, vstr, usage)
}

func (set *FlagSet) DurationSlice(short string, long string, value []time.Duration, usage string) *[]time.Duration {
	ptr := new([]time.Duration)
	set.DurationSliceVar(ptr, short, long, value, usage)
	return ptr
}

func (set *FlagSet) DurationSliceVar(ptr *[]time.Duration, short string, long string, value []time.Duration, usage string) {
	var vstr string
	if value != nil {
		vstr = (*durationSliceVar)(ptr).String()
	}
	set.Var((*durationSliceVar)(ptr), short, long, vstr, usage)
}

// Var will add given Value implementation to the Flagset, with given short / long / env names and usage, and encoded default 'value' string.
func (set *FlagSet) Var(ptr Value, short string, long string, value string, usage string) {
	set.Add(Flag{
		Short:    short,
		Long:     long,
		Usage:    usage,
		Required: false,
		Default:  value,
		Value:    ptr,
	})
}

// Add will check the validity of, and add given flag to this FlagSet.
func (set *FlagSet) Add(flag Flag) {
	if set.auto != nil &&
		flag.Long != "" && flag.Env == "" {
		// Generate environment key name.
		flag.Env = set.auto(flag.Long)
	}

	switch {
	// Flag contains no matchable name
	case flag.Short == "" && flag.Long == "" && flag.Env == "":
		panic("no flag short, long or env name")

	// Flag short name must be a single char
	case utf8.RuneCountInString(flag.Short) > 1:
		panic("flag short name must be single rune")

	// Flag names cannot start with '-'
	case strings.HasPrefix(flag.Short, "-"):
		panic("flag short name starts with hypen")
	case strings.HasPrefix(flag.Long, "-"):
		panic("flag long name starts with hypen")

	// Flag names must not contain any illegal chars
	case strings.ContainsAny(flag.Short, illegalChars):
		panic("flag short name contains illegal char")
	case strings.ContainsAny(flag.Long, illegalChars):
		panic("flag long name contains illegal char(s)")
	case strings.ContainsAny(flag.Env, illegalChars):
		panic("flag env name contains illegal char(s)")

	// Flag names where provided must be unique
	case flag.Short != "" && set.GetShort(flag.Short) != nil:
		panic("flag short name conflict")
	case flag.Long != "" && set.GetLong(flag.Long) != nil:
		panic("flag long name conflict")
	case flag.Env != "" && set.GetEnv(flag.Env) != nil:
		panic("flag env name conflict")

	// Check for flag dst
	case flag.Value == nil:
		panic("nil flag value destination")
	}

	// Append flag to internal slice
	set.flags = append(set.flags, &flag)
}

// Reset will reset all Flags in FlagSet.
func (set *FlagSet) Reset() {
	set.flags = set.flags[:0]
	set.funcs = set.funcs[:0]
}
