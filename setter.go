package qparser

import (
	"fmt"
	"reflect"
	"strconv"
)

// field represents each field in the target struct with its reflection metadata and string value from query parameter.
type field struct {
	refField reflect.StructField
	refVal   reflect.Value
	value    string
}

// registerField initializes a new field struct with the given reflection metadata and value.
func registerField(t reflect.StructField, v reflect.Value, queryVal string) *field {
	return &field{
		refField: t,
		refVal:   v,
		value:    queryVal,
	}
}

// SetValue assigns the value to the field and converts it to the correct type.
// It handles pointer fields by dereferencing them and sets the value only if the field is settable.
func (f *field) SetValue() error {
	// Handle pointer field, by dereferencing it
	if f.refVal.Kind() == reflect.Ptr {
		if f.refVal.IsNil() {
			// If the value is empty, let the field value remain nil
			if f.value == "" {
				return nil
			}
			f.refVal.Set(reflect.New(f.refVal.Type().Elem()))
		}
		f.refVal = f.refVal.Elem()
	}

	var err error
	// Check if the field is settable
	if f.refVal.CanSet() {
		switch f.refVal.Kind() {
		case reflect.String:
			f.refVal.SetString(f.value)
		case reflect.Bool:
			err = f.setBool()
		case reflect.Float64, reflect.Float32:
			err = f.setFloat()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = f.setInt()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err = f.setUint()
		default:
			err = fmt.Errorf("qparser: unsupported kind %s for field %s", f.refVal.Kind(), f.refField.Name)
		}
	} else {
		err = fmt.Errorf("qparser: cannot set field %s", f.refField.Name)
	}
	return err
}

// bitSize determines the bit size of the field for integer and float types.
func (f *field) bitSize() int {
	switch f.refField.Type.Kind() {
	case reflect.Int8, reflect.Uint8:
		return 8
	case reflect.Int16, reflect.Uint16:
		return 16
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 32
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return 64
	default:
		return 0
	}
}

func (f *field) setBool() error {
	val, err := strconv.ParseBool(f.value)
	if err != nil {
		return f.whichErr(err)
	}
	f.refVal.SetBool(val)
	return nil
}

func (f *field) setFloat() error {
	val, err := strconv.ParseFloat(f.value, f.bitSize())
	if err != nil {
		return f.whichErr(err)
	}
	f.refVal.SetFloat(val)
	return nil
}

func (f *field) setInt() error {
	val, err := strconv.ParseInt(f.value, 10, f.bitSize())
	if err != nil {
		return f.whichErr(err)
	}
	f.refVal.SetInt(val)
	return nil
}

func (f *field) setUint() error {
	val, err := strconv.ParseUint(f.value, 10, f.bitSize())
	if err != nil {
		return f.whichErr(err)
	}
	f.refVal.SetUint(val)
	return nil
}

// whichErr wraps the error with additional context about the field and value.
func (f *field) whichErr(err error) error {
	switch e := err.(type) {
	case *strconv.NumError:
		if e.Err == strconv.ErrRange {
			return fmt.Errorf("value %q out of range for %s(%s) field", f.value, f.refField.Name, f.refField.Type)
		}
		return fmt.Errorf("invalid value %q for %s(%s) field", f.value, f.refField.Name, f.refField.Type)
	default:
		return err
	}
}
