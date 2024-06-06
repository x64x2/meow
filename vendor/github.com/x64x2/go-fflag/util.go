package fflag

import (
	"reflect"
)

const (
	// illegalChars specifies illegal characters that flag CLI / env flags may not contain.
	illegalChars = "\x00" + // zero byte

		// '=' used for compound key=value
		"=" +

		// various "whitespace" chars
		"\t\n\v\f\r " + // ASCII space chars
		"\u0009" + // char tab, '\t'
		"\u000A" + // line feed, '\n'
		"\u000B" + // line tab, '\v'
		"\u000C" + // form feed, '\f'
		"\u000D" + // carriage return, '\r'
		"\u0020" + // space
		"\u0085" + // next line, NEL
		"\u00a0" + // no break space
		"\u1680" + // ogham space (though often a line)
		"\u2000" + // en quad
		"\u2001" + // em quad
		"\u2002" + // en space
		"\u2003" + // em space
		"\u2004" + // three-per-em space
		"\u2005" + // four-per-em space
		"\u2006" + // six-per-em space
		"\u2007" + // figure space
		"\u2008" + // punctuation space
		"\u2009" + // thin space
		"\u200A" + // hair space
		"\u2028" + // line separator
		"\u2029" + // paragraph separator
		"\u202F" + // narrow no-break space
		"\u205F" + // medium math space
		"\u2300" // ideographic space
)

// isValidTag returns whether tag is a valid struct tag (i.e. not empty or '-' null).
func isValidTag(tag string) bool {
	return len(tag) > 0 && tag != "-"
}

// structFieldsForType returns a slice of struct fields for given struct type.
func structFieldsForType(t reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, t.NumField())
	for i := 0; i < len(fields); i++ {
		fields[i] = t.Field(i)
	}
	return fields
}
