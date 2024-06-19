package qparser

import (
	"fmt"
	"reflect"
	"strconv"
)

// Each field in the target struct is represented by a field struct
type field struct {
	refType reflect.StructField
	refVal  reflect.Value
}

func registerField(ft reflect.StructField, fv reflect.Value) *field {
	return &field{
		refType: ft,
		refVal:  fv,
	}
}

// Set assign the value to the field and converts it to the correct type.
func (f field) Set(val string) error {
	// Handle pointer field, by dereferencing it
	if f.refVal.Kind() == reflect.Ptr {
		if f.refVal.IsNil() {
			// Since pointer field is allowed to be nil, we should return nil if the value is empty.
			if val == "" {
				return nil
			}
			f.refVal.Set(reflect.New(f.refVal.Type().Elem()))
		}
		f.refVal = f.refVal.Elem()
	}

	if val == "" {
		return fmt.Errorf("empty value for field %s", f.refType.Name)
	}

	if f.refVal.CanSet() {
		if setter, ok := setters[f.refVal.Kind()]; ok {
			return setter(f, val)
		}
		return fmt.Errorf("unsupported kind %s for field %s", f.refVal.Kind(), f.refType.Name)
	}
	return fmt.Errorf("cannot set field %s", f.refType.Name)
}

var setters = map[reflect.Kind]setter{
	reflect.String:  (field).setString,
	reflect.Bool:    (field).setBool,
	reflect.Float64: (field).setFloat,
	reflect.Float32: (field).setFloat,
	reflect.Int:     (field).setInt,
	reflect.Int8:    (field).setInt,
	reflect.Int16:   (field).setInt,
	reflect.Int32:   (field).setInt,
	reflect.Int64:   (field).setInt,
	reflect.Uint:    (field).setUint,
	reflect.Uint8:   (field).setUint,
	reflect.Uint16:  (field).setUint,
	reflect.Uint32:  (field).setUint,
	reflect.Uint64:  (field).setUint,
}

func (f field) bitSize() int {
	switch f.refType.Type.Kind() {
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

type setter func(f field, v string) error

func (f field) setString(v string) error {
	f.refVal.SetString(v)
	return nil
}

func (f field) setBool(v string) error {
	val, err := strconv.ParseBool(v)
	if err != nil {
		return errorSet(err, f.refType.Name)
	}
	f.refVal.SetBool(val)
	return nil
}

func (f field) setFloat(v string) error {
	val, err := strconv.ParseFloat(v, f.bitSize())
	if err != nil {
		return errorSet(err, f.refType.Name)
	}
	f.refVal.SetFloat(val)
	return nil
}

func (f field) setInt(v string) error {
	val, err := strconv.ParseInt(v, 10, f.bitSize())
	if err != nil {
		return errorSet(err, f.refType.Name)
	}
	f.refVal.SetInt(val)
	return nil
}

func (f field) setUint(v string) error {
	val, err := strconv.ParseUint(v, 10, f.bitSize())
	if err != nil {
		return errorSet(err, f.refType.Name)
	}
	f.refVal.SetUint(val)
	return nil
}

func errorSet(err error, name string) error {
	return fmt.Errorf("%v at %v field", err, name)
}
