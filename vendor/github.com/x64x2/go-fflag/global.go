package fflag

import (
	"os"
	"syscall"
	"time"

	"codeberg.org/gruf/go-bytesize"
)

// Global provides access to a global
// FlagSet instance accessible via
// global FlagSet functions below.
var Global FlagSet

func GetShort(name string) *Flag {
	return Global.GetShort(name)
}

func GetLong(name string) *Flag {
	return Global.GetLong(name)
}

func GetEnv(name string) *Flag {
	return Global.GetEnv(name)
}

func Usage() string {
	return Global.Usage()
}

func AutoEnv(fn func(long string) (env string)) {
	Global.AutoEnv(fn)
}

func Config(fn func(*os.File) error) {
	Global.Config(fn)
}

func Help() {
	Global.Help()
}

func Version(version string) {
	Global.Version(version)
}

func Func(short string, long string, usage string, fn func()) {
	Global.Func(short, long, usage, fn)
}

func Bool(short string, long string, value bool, usage string) *bool {
	return Global.Bool(short, long, value, usage)
}

func BoolVar(ptr *bool, short string, long string, value bool, usage string) {
	Global.BoolVar(ptr, short, long, value, usage)
}

func Int(short string, long string, value int, usage string) *int {
	return Global.Int(short, long, value, usage)
}

func IntVar(ptr *int, short string, long string, value int, usage string) {
	Global.IntVar(ptr, short, long, value, usage)
}

func Int64(short string, long string, value int64, usage string) *int64 {
	return Global.Int64(short, long, value, usage)
}

func Int64Var(ptr *int64, short string, long string, value int64, usage string) {
	Global.Int64Var(ptr, short, long, value, usage)
}

func Uint(short string, long string, value uint, usage string) *uint {
	return Global.Uint(short, long, value, usage)
}

func UintVar(ptr *uint, short string, long string, value uint, usage string) {
	Global.UintVar(ptr, short, long, value, usage)
}

func Uint64(short string, long string, value uint64, usage string) *uint64 {
	return Global.Uint64(short, long, value, usage)
}

func Uint64Var(ptr *uint64, short string, long string, value uint64, usage string) {
	Global.Uint64Var(ptr, short, long, value, usage)
}

func Float(short string, long string, value float64, usage string) *float64 {
	return Global.Float(short, long, value, usage)
}

func FloatVar(ptr *float64, short string, long string, value float64, usage string) {
	Global.FloatVar(ptr, short, long, value, usage)
}

func String(short string, long string, value string, usage string) *string {
	return Global.String(short, long, value, usage)
}

func StringVar(ptr *string, short string, long string, value string, usage string) {
	Global.StringVar(ptr, short, long, value, usage)
}

func Size(short string, long string, value bytesize.Size, usage string) *bytesize.Size {
	return Global.Size(short, long, value, usage)
}

func SizeVar(ptr *bytesize.Size, short string, long string, value bytesize.Size, usage string) {
	Global.SizeVar(ptr, short, long, value, usage)
}

func Duration(short string, long string, value time.Duration, usage string) *time.Duration {
	return Global.Duration(short, long, value, usage)
}

func DurationVar(ptr *time.Duration, short string, long string, value time.Duration, usage string) {
	Global.DurationVar(ptr, short, long, value, usage)
}

func Bytes(short string, long string, value []byte, usage string) *[]byte {
	return Global.Bytes(short, long, value, usage)
}

func BytesVar(ptr *[]byte, short string, long string, value []byte, usage string) {
	Global.BytesVar(ptr, short, long, value, usage)
}

func IntSlice(short string, long string, value []int, usage string) *[]int {
	return Global.IntSlice(short, long, value, usage)
}

func IntSliceVar(ptr *[]int, short string, long string, value []int, usage string) {
	Global.IntSliceVar(ptr, short, long, value, usage)
}

func Int64Slice(short string, long string, value []int64, usage string) *[]int64 {
	return Global.Int64Slice(short, long, value, usage)
}

func Int64SliceVar(ptr *[]int64, short string, long string, value []int64, usage string) {
	Global.Int64SliceVar(ptr, short, long, value, usage)
}

func UintSlice(short string, long string, value []uint, usage string) *[]uint {
	return Global.UintSlice(short, long, value, usage)
}

func UintSliceVar(ptr *[]uint, short string, long string, value []uint, usage string) {
	Global.UintSliceVar(ptr, short, long, value, usage)
}

func Uint64Slice(short string, long string, value []uint64, usage string) *[]uint64 {
	return Global.Uint64Slice(short, long, value, usage)
}

func Uint64SliceVar(ptr *[]uint64, short string, long string, value []uint64, usage string) {
	Global.Uint64SliceVar(ptr, short, long, value, usage)
}

func FloatSlice(short string, long string, value []float64, usage string) *[]float64 {
	return Global.FloatSlice(short, long, value, usage)
}

func FloatSliceVar(ptr *[]float64, short string, long string, value []float64, usage string) {
	Global.FloatSliceVar(ptr, short, long, value, usage)
}

func StringSlice(short string, long string, value []string, usage string) *[]string {
	return Global.StringSlice(short, long, value, usage)
}

func StringSliceVar(ptr *[]string, short string, long string, value []string, usage string) {
	Global.StringSliceVar(ptr, short, long, value, usage)
}

func StructFields(dst interface{}) {
	Global.StructFields(dst)
}

func Var(ptr Value, short string, long string, value string, usage string) {
	Global.Var(ptr, short, long, value, usage)
}

func Add(flag Flag) {
	Global.Add(flag)
}

func Reset() {
	Global.Reset()
}

func ParseEnv() (err error) {
	env := syscall.Environ()
	return Global.ParseEnv(env)
}

func ParseArgs() (unused []string, err error) {
	args := os.Args[1:] // args minus bin name.
	return Global.ParseArgs(args)
}

func Parse() (unused []string, err error) {
	args := os.Args[1:] // args minus bin name.
	envs := syscall.Environ()
	return Global.Parse(args, envs)
}

func MustParse() (unused []string) {
	var err error
	unused, err = Parse()
	if err != nil {
		panic(err)
	}
	return unused
}
