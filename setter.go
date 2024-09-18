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
func setValue(v reflect.Value, queryValue []string) error {
	v = deref(v)

	if v.Kind() == reflect.Slice {
		values := func() (result []string) {
			for _, value := range queryValue {
				value = strings.TrimSpace(value)
				if value == "" {
					continue
				}

				// If the value contains the separator, split it and process each part
				if strings.Contains(value, SEPARATOR) {
					for _, v := range strings.Split(value, SEPARATOR) {
						if v = strings.TrimSpace(v); v != "" {
							result = append(result, v)
						}
					}
				} else {
					result = append(result, value)
				}
			}
			return result
		}()

		// Create a new slice with the same type, length and cap as the values
		cpSlice := reflect.MakeSlice(v.Type(), len(values), len(values))

		for i, value := range values {
			item := cpSlice.Index(i)

			if err := setValue(item, []string{value}); err != nil {
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

	result := set(queryValue[0])
	if !result.IsValid() || result == invalidReflect {
		return fmt.Errorf("invalid value '%s' for field", queryValue)
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

// deref helps to dereference a pointer value.
func deref(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}
