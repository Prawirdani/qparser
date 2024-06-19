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
			// If the value is empty, let the field value remain nil
			if val == "" {
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
			f.refVal.SetString(val)
		case reflect.Bool:
			err = f.setBool(val)
		case reflect.Float64, reflect.Float32:
			err = f.setFloat(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = f.setInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err = f.setUint(val)
		default:
			err = fmt.Errorf("qparser: unsupported kind %s for field %s", f.refVal.Kind(), f.refType.Name)
		}
	} else {
		err = fmt.Errorf("qparser: cannot set field %s", f.refType.Name)
	}
	return err
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
