package fflag

import (
	"reflect"
	"time"
	"unsafe"

	"codeberg.org/gruf/go-bytesize"
)

var (
	// reflected runtime type pointers.
	tDur       = reflect.TypeOf(time.Duration(0))
	tSize      = reflect.TypeOf(bytesize.Size(0))
	tDurSlice  = reflect.TypeOf([]time.Duration(nil))
	tSizeSlice = reflect.TypeOf([]bytesize.Size(nil))
	tValue     = reflect.TypeOf((*Value)(nil)).Elem()
)

// StructFields registers appropriately tagged struct fields as flags.
//
// Tags:
// "short" - the flag short name
// "long" - the flag long name
// "env" - flag env key name
// "usage" - the flag usage information
// "default" - the default flag value (in string form)
// "required" - whether flag MUST be provided
func (set *FlagSet) StructFields(dst any) {
	// Get reflect value + type
	rtype := reflect.TypeOf(dst)
	rvalue := reflect.ValueOf(dst)

	for rtype.Kind() == reflect.Pointer {
		// We need a non-nil dst
		if rvalue.IsNil() {
			panic("dst must be non-nil")
		}

		// Dereference pointer
		rtype = rtype.Elem()
		rvalue = rvalue.Elem()
	}

	if rtype.Kind() != reflect.Struct {
		panic("dst must be struct pointer")
	}

	// Generate fields for reflected type
	fields := structFieldsForType(rtype)
	set.generateStructFlags(rvalue, fields)
}

// generateStructFlags contains the main logic of GenerateFlags(), setup to allow recursion.
func (set *FlagSet) generateStructFlags(rvalue reflect.Value, fields []reflect.StructField) {
outer:
	for i := 0; i < len(fields); i++ {
		var (
			// configured CLI flag
			flag Flag

			// does struct field require flag
			hasFlag bool
		)

		if tag := fields[i].Tag.Get("short"); isValidTag(tag) {
			// Set flag's short-name
			flag.Short = tag
			hasFlag = true
		}

		if tag := fields[i].Tag.Get("long"); isValidTag(tag) {
			// Set flag's long-name
			flag.Long = tag
			hasFlag = true
		}

		if tag := fields[i].Tag.Get("env"); isValidTag(tag) {
			// Set flag's env-name
			flag.Env = tag
			hasFlag = true
		}

		if tag := fields[i].Tag.Get("usage"); isValidTag(tag) {
			// Set flag's usage information
			flag.Usage = tag
		}

		if tag := fields[i].Tag.Get("default"); isValidTag(tag) {
			// Set flag's default value
			flag.Default = tag
		}

		if tag := fields[i].Tag.Get("required"); isValidTag(tag) {
			// Set flag as required
			flag.Required = true
		}

		if !hasFlag {
			// no flag.
			continue
		}

		var (
			t = fields[i].Type
			c = 0
		)

		for t.Kind() == reflect.Pointer {
			// Dereference ptr
			t = t.Elem()
			c++
		}

		// Get struct field value
		vfield := rvalue.Field(i)

		if !hasFlag && t.Kind() == reflect.Struct {
			for i := 0; i < c; i++ {
				if vfield.IsNil() {
					// ignore nil structs
					continue outer
				}

				// Dereference pointer
				vfield = vfield.Elem()
			}

			// Recursively generate flags for structs
			subfields := structFieldsForType(fields[i].Type)
			set.generateStructFlags(vfield, subfields)
			continue
		}

		if !assignValueType(&flag, vfield, fields[i].Type) {
			// If no flag value was set, then this is not a kind we support by default
			panic("no Value implementation for struct field type: " + t.String())
		}

		// Append flag
		set.Add(flag)
	}
}

// assignValueType attempts to assign a flag's value pointer for given reflect type and value.
func assignValueType(flag *Flag, v reflect.Value, t reflect.Type) bool {
	// Try direct type ptr comparison
	if findByType(flag, v, t) {
		return true
	}

	// Check if implements Value{}
	if t.Implements(tValue) {
		i := v.Interface()
		flag.Value = i.(Value)
		return true
	} else if reflect.PointerTo(t).Implements(tValue) {
		i := v.Addr().Interface()
		flag.Value = i.(Value)
		return true
	}

	// Check by known primitive types
	if findByKind(flag, v, t) {
		return true
	}

	return false
}

func findByType(flag *Flag, v reflect.Value, t reflect.Type) bool {
	switch t {
	// bytesize.Size
	case tSize:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*sizeVar)(p)
		return true

	// time.Duration
	case tDur:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*durationVar)(p)
		return true

	// []bytesize.Size
	case tSizeSlice:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*sizeSliceVar)(p)
		return true

	// []time.Duration
	case tDurSlice:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*durationSliceVar)(p)
		return true

	default:
		return false
	}
}

func findByKind(flag *Flag, v reflect.Value, t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*boolVar)(p)
		return true

	case reflect.Int:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*intVar)(p)
		return true

	case reflect.Int64:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*int64Var)(p)
		return true

	case reflect.Uint, reflect.Uintptr:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*uintVar)(p)
		return true

	case reflect.Uint64:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*uint64Var)(p)
		return true

	case reflect.Float64:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*float64Var)(p)
		return true

	case reflect.String:
		p := unsafe.Pointer(v.UnsafeAddr())
		flag.Value = (*stringVar)(p)
		return true

	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8:
			p := unsafe.Pointer(v.UnsafeAddr())
			flag.Value = (*bytesVar)(p)
			return true

		case reflect.Bool:
			p := unsafe.Pointer(v.UnsafeAddr())
			flag.Value = (*boolSliceVar)(p)
			return true

		case reflect.Int:
			p := unsafe.Pointer(v.UnsafeAddr())
			flag.Value = (*intSliceVar)(p)
			return true

		case reflect.Int64:
			p := unsafe.Pointer(v.UnsafeAddr())
			flag.Value = (*int64SliceVar)(p)
			return true

		case reflect.Uint, reflect.Uintptr:
			p := unsafe.Pointer(v.UnsafeAddr())
			flag.Value = (*uintSliceVar)(p)
			return true

		case reflect.Uint64:
			p := unsafe.Pointer(v.UnsafeAddr())
			flag.Value = (*uint64SliceVar)(p)
			return true

		case reflect.Float64:
			p := unsafe.Pointer(v.UnsafeAddr())
			flag.Value = (*float64SliceVar)(p)
			return true

		case reflect.String:
			p := unsafe.Pointer(v.UnsafeAddr())
			flag.Value = (*stringSliceVar)(p)
			return true

		default:
			return false
		}

	default:
		return false
	}
}
