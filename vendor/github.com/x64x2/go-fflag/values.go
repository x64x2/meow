package fflag

import (
	"strconv"
	"time"
	"unsafe"

	"codeberg.org/gruf/go-bytesize"
	"codeberg.org/gruf/go-split"
)

// Value defines a pointer to usable flag value type.
type Value interface {
	Set(string) error
	Kind() string
	String() string
}

// below are default implementations of Value.

type boolVar bool

func (v *boolVar) Set(in string) error {
	b, err := strconv.ParseBool(in)
	if err != nil {
		return err
	}
	*v = boolVar(b)
	return nil
}

func (*boolVar) Kind() string {
	return ""
}

func (v *boolVar) String() string {
	return strconv.FormatBool(bool(*v))
}

type boolSliceVar []bool

func (v *boolSliceVar) Set(in string) error {
	s, err := split.SplitBools[bool](in)
	*v = append(*v, boolSliceVar(s)...)
	return err
}

func (*boolSliceVar) Kind() string {
	return "[]bool"
}

func (v *boolSliceVar) String() string {
	return split.JoinBools(*v)
}

type intVar int

func (v *intVar) Set(in string) error {
	i, err := strconv.ParseInt(in, 10, strconv.IntSize)
	if err != nil {
		return err
	}
	*v = intVar(i)
	return nil
}

func (*intVar) Kind() string {
	return "int"
}

func (v *intVar) String() string {
	return strconv.FormatInt(int64(*v), 10)
}

type intSliceVar []int

func (v *intSliceVar) Set(in string) error {
	s, err := split.SplitInts[int](in)
	*v = append(*v, intSliceVar(s)...)
	return err
}

func (*intSliceVar) Kind() string {
	return "[]int"
}

func (v *intSliceVar) String() string {
	return split.JoinInts(*v)
}

type int64Var int64

func (v *int64Var) Set(in string) error {
	i, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return err
	}
	*v = int64Var(i)
	return nil
}

func (*int64Var) Kind() string {
	return "int64"
}

func (v *int64Var) String() string {
	return strconv.FormatInt(int64(*v), 10)
}

type int64SliceVar []int64

func (v *int64SliceVar) Set(in string) error {
	s, err := split.SplitInts[int64](in)
	*v = append(*v, int64SliceVar(s)...)
	return err
}

func (*int64SliceVar) Kind() string {
	return "[]int64"
}

func (v *int64SliceVar) String() string {
	return split.JoinInts(*v)
}

type uintVar uint

func (v *uintVar) Set(in string) error {
	u, err := strconv.ParseUint(in, 10, strconv.IntSize)
	if err != nil {
		return err
	}
	*v = uintVar(u)
	return nil
}

func (*uintVar) Kind() string {
	return "uint"
}

func (v *uintVar) String() string {
	return strconv.FormatUint(uint64(*v), 10)
}

type uintSliceVar []uint

func (v *uintSliceVar) Set(in string) error {
	s, err := split.SplitUints[uint](in)
	*v = append(*v, uintSliceVar(s)...)
	return err
}

func (*uintSliceVar) Kind() string {
	return "[]uint"
}

func (v *uintSliceVar) String() string {
	return split.JoinUints(*v)
}

type uint64Var uint64

func (v *uint64Var) Set(in string) error {
	u, err := strconv.ParseUint(in, 10, 64)
	if err != nil {
		return err
	}
	*v = uint64Var(u)
	return nil
}

func (*uint64Var) Kind() string {
	return "uint64"
}

func (v *uint64Var) String() string {
	return strconv.FormatUint(uint64(*v), 10)
}

type uint64SliceVar []uint64

func (v *uint64SliceVar) Set(in string) error {
	s, err := split.SplitUints[uint64](in)
	*v = append(*v, uint64SliceVar(s)...)
	return err
}

func (*uint64SliceVar) Kind() string {
	return "[]uint64"
}

func (v *uint64SliceVar) String() string {
	return split.JoinUints(*v)
}

type float64Var float64

func (v *float64Var) Set(in string) error {
	f, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return err
	}
	*v = float64Var(f)
	return nil
}

func (*float64Var) Kind() string {
	return "float64"
}

func (v *float64Var) String() string {
	return strconv.FormatFloat(float64(*v), 'g', -1, 64)
}

type float64SliceVar []float64

func (v *float64SliceVar) Set(in string) error {
	s, err := split.SplitFloats[float64](in)
	*v = append(*v, float64SliceVar(s)...)
	return err
}

func (*float64SliceVar) Kind() string {
	return "[]float64"
}

func (v *float64SliceVar) String() string {
	return split.JoinFloats(*v)
}

type stringVar string

func (v *stringVar) Set(in string) error {
	*v += stringVar(in)
	return nil
}

func (*stringVar) Kind() string {
	return "string"
}

func (v *stringVar) String() string {
	return string(*v)
}

type stringSliceVar []string

func (v *stringSliceVar) Set(in string) error {
	s, err := split.SplitStrings[string](in)
	*v = append(*v, stringSliceVar(s)...)
	return err
}

func (*stringSliceVar) Kind() string {
	return "[]string"
}

func (v *stringSliceVar) String() string {
	return split.JoinStrings(*v)
}

type bytesVar []byte

func (v *bytesVar) Set(in string) error {
	*v = append(*v, in...)
	return nil
}

func (*bytesVar) Kind() string {
	return "[]byte"
}

func (v *bytesVar) String() string {
	return *(*string)(unsafe.Pointer(v))
}

type sizeVar bytesize.Size

func (v *sizeVar) Set(in string) error {
	sz, err := bytesize.ParseSize(in)
	if err != nil {
		return err
	}
	*v = sizeVar(sz)
	return nil
}

func (v *sizeVar) Kind() string {
	return "size"
}

func (v *sizeVar) String() string {
	return (*bytesize.Size)(v).StringIEC()
}

type sizeSliceVar []bytesize.Size

func (v *sizeSliceVar) Set(in string) error {
	s, err := split.SplitSizes(in)
	*v = append(*v, sizeSliceVar(s)...)
	return err
}

func (v *sizeSliceVar) Kind() string {
	return "[]size"
}

func (v *sizeSliceVar) String() string {
	return split.JoinSizes(*v)
}

type durationVar time.Duration

func (v *durationVar) Set(in string) error {
	d, err := time.ParseDuration(in)
	if err != nil {
		return err
	}
	*v = durationVar(d)
	return nil
}

func (v *durationVar) Kind() string {
	return "duration"
}

func (v *durationVar) String() string {
	return (*time.Duration)(v).String()
}

type durationSliceVar []time.Duration

func (v *durationSliceVar) Set(in string) error {
	s, err := split.SplitDurations(in)
	*v = append(*v, durationSliceVar(s)...)
	return err
}

func (v *durationSliceVar) Kind() string {
	return "[]duration"
}

func (v *durationSliceVar) String() string {
	return split.JoinDurations(*v)
}
