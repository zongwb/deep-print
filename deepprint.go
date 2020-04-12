// Package deepprint recursively prints the content of an artibrary data structure.
// It achieves the same effect as calling json.MarshalIndent, but can print
// unexported fields.
package deepprint

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

const (
	INDENT         = "  "
	nilAngleString = "<nil>"
	maxDepth       = 5
)

// DeepPrint prints the content of s recursively into a string.
func DeepPrint(s interface{}) (string, error) {
	var buf bytes.Buffer
	err := deepPrint(&buf, reflect.ValueOf(s), "", INDENT, 1)
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}

// deepPrint prints the content of v recursively. prefix is prepended
// before each new line.
func deepPrint(w io.Writer, v reflect.Value, prefix, indent string, depth int) error {
	if depth > maxDepth {
		return nil
	}
	switch v.Kind() {
	case reflect.Bool:
		_, err := fmt.Fprintf(w, "%t", v.Bool())
		return err
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err := fmt.Fprintf(w, "%d", v.Int())
		return err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		_, err := fmt.Fprintf(w, "%d", v.Uint())
		return err
	case reflect.Float32, reflect.Float64:
		_, err := fmt.Fprintf(w, "%f", v.Float())
		return err
	case reflect.Complex64, reflect.Complex128:
		_, err := fmt.Fprintf(w, "%f", v.Complex())
		return err
	case reflect.String:
		_, err := fmt.Fprintf(w, "\"%s\"", v.String())
		return err
	case reflect.Array:
		return printSlice(w, v, prefix, indent, depth)
	case reflect.Slice:
		if v.IsNil() {
			fmt.Fprint(w, nilAngleString)
			return nil
		}
		return printSlice(w, v, prefix, indent, depth)
	case reflect.Map:
		if v.IsNil() {
			fmt.Fprint(w, nilAngleString)
			return nil
		}
		return printMap(w, v, prefix, indent, depth)
	case reflect.Struct:
		return printStruct(w, v, prefix, indent, depth)
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			fmt.Fprint(w, nilAngleString)
			return nil
		}
		return deepPrint(w, v.Elem(), prefix, indent, depth+1)
	default:
		fmt.Fprintf(w, "non-printable type %v", v.Kind())
	}
	return nil
}

func printSlice(w io.Writer, v reflect.Value, prefix, indent string, depth int) error {
	fmt.Fprint(w, "[ ")
	len := v.Len()
	for i := 0; i < len; i++ {
		vv := v.Index(i)
		err := deepPrint(w, vv, prefix, indent, depth+1)
		if err != nil {
			return err
		}
		fmt.Fprint(w, " ")
	}
	fmt.Fprint(w, "]")
	return nil
}

func printMap(w io.Writer, v reflect.Value, prefix, indent string, depth int) error {
	keys := v.MapKeys()
	fmt.Fprint(w, "{\n")
	for _, key := range keys {
		fmt.Fprintf(w, "%s", prefix)
		err := deepPrint(w, key, prefix, indent, depth+1)
		if err != nil {
			return err
		}
		fmt.Fprint(w, ": ")

		err = deepPrint(w, v.MapIndex(key), prefix, indent, depth+1)
		if err != nil {
			return err
		}
		fmt.Fprint(w, ",\n")
	}
	fmt.Fprintf(w, "%s}", prefix)
	return nil
}

func printStruct(w io.Writer, v reflect.Value, prefix, indent string, depth int) error {
	v = reflect.Indirect(v)
	numF := v.NumField()
	fmt.Fprint(w, "{\n")
	for i := 0; i < numF; i++ {
		key := prefix + indent + "\"" + v.Type().Field(i).Name + "\": "
		fmt.Fprint(w, key)
		newPrefix := padRight(prefix, " ", len(key))
		err := deepPrint(w, v.Field(i), newPrefix, indent, depth+1)
		if err != nil {
			return err
		}
		fmt.Fprint(w, ",\n")
	}
	fmt.Fprintf(w, "%s}", prefix)
	return nil
}

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}
