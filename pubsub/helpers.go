package pubsub

import (
	"bytes"
	"reflect"
	"unicode"
	"unicode/utf8"
)

var (
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
)

type buffer struct {
	*bytes.Buffer
}

func (b *buffer) Close() error {
	b.Buffer.Reset()
	return nil
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}
