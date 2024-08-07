package qparser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// A Function that set query parameter value to the reflected value from the struct field.
type Setter func(value string) reflect.Value

// Basic types setters
var setters = map[reflect.Kind]Setter{
	reflect.String:  setStr,
	reflect.Bool:    setBool,
	reflect.Int:     setInt(0),
	reflect.Int8:    setInt(8),
	reflect.Int16:   setInt(16),
	reflect.Int32:   setInt(32),
	reflect.Int64:   setInt(64),
	reflect.Uint:    setUint(0),
	reflect.Uint8:   setUint(8),
	reflect.Uint16:  setUint(16),
	reflect.Uint32:  setUint(32),
	reflect.Uint64:  setUint(64),
	reflect.Float32: setFloat(32),
	reflect.Float64: setFloat(64),
}

var invalidReflect = reflect.Value{}

// Multiple query parameter values separator
const SEPARATOR = ","

// setValue sets the query parameter value to the reflected value from the struct field.
func setValue(v reflect.Value, value string) error {
	// BUG: We already deref the value in the parse function, so we should not deref it again here. But somehow it still reflecting as a pointer, so we need to deref it again.
	v = deref(v)

	if v.Kind() == reflect.Slice {
		values := func() (result []string) {
			for _, value := range strings.Split(value, SEPARATOR) {
				if value = strings.TrimSpace(value); value != "" {
					result = append(result, value)
				}
			}
			return result
		}()

		// Create a new slice with the same type and length as the values
		cpSlice := reflect.MakeSlice(v.Type(), len(values), len(values))

		for i, value := range values {
			item := cpSlice.Index(i)

			if err := setValue(item, value); err != nil {
				return err
			}
		}
		v.Set(cpSlice)
		return nil
	}

	set, ok := setters[v.Kind()]
	if !ok {
		return fmt.Errorf("unsupported kind %s", v.Kind())
	}

	result := set(value)
	if !result.IsValid() || result == invalidReflect {
		return fmt.Errorf("invalid value '%s' for field", value)
	}

	v.Set(result.Convert(v.Type()))

	return nil
}

func setStr(value string) reflect.Value {
	return reflect.ValueOf(value)
}

func setBool(value string) reflect.Value {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return invalidReflect
	}
	return reflect.ValueOf(v)
}

func setInt(bitSize int) Setter {
	return func(value string) reflect.Value {
		v, err := strconv.ParseInt(value, 10, bitSize)
		if err != nil {
			return invalidReflect
		}
		return reflect.ValueOf(v)
	}
}

func setUint(bitSize int) Setter {
	return func(value string) reflect.Value {
		v, err := strconv.ParseUint(value, 10, bitSize)
		if err != nil {
			return invalidReflect
		}
		return reflect.ValueOf(v)
	}
}

func setFloat(bitSize int) Setter {
	return func(value string) reflect.Value {
		v, err := strconv.ParseFloat(value, bitSize)
		if err != nil {
			return invalidReflect
		}
		return reflect.ValueOf(v)
	}
}
