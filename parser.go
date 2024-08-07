package qparser

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const structTAG = "qp"

// ParseRequest parses the URL query parameters from the HTTP request and sets the corresponding fields in the struct.
func ParseRequest(r *http.Request, structPointer any) error {
	queryValues := r.URL.Query()
	return parse(queryValues, structPointer)
}

// ParseURL parses query parameters from the URL string and sets the corresponding fields in the struct.
func ParseURL(address string, structPointer any) error {
	urlObj, err := url.Parse(address)
	if err != nil {
		return err
	}
	queryValues := urlObj.Query()
	return parse(queryValues, structPointer)
}

// parse parses the provided URL query parameters into the given struct.
// The struct fields must have `qp` tags specifying the query parameter names.
// If a field's `qp` tag is "-" or empty, the field will be ignored.
func parse(querieValues url.Values, st any) error {
	v, err := reflecter(st)
	if err != nil {
		return err
	}

	v = v.Elem()
	t := v.Type()

	// Iterate over the struct fields
	for i := 0; i < t.NumField(); i++ {
		fieldMetadata := t.Field(i)
		fieldValue := deref(v.Field(i))

		// Handle nested struct by recursively calling parse
		if fieldValue.Kind() == reflect.Struct {
			err := parse(querieValues, fieldValue.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}

		// Retrieve the qp tag value
		tag := strings.TrimSpace(fieldMetadata.Tag.Get(structTAG))
		if tag == "" || tag == "-" {
			continue
		}

		// Retrieve the query parameter value
		queryValue := querieValues.Get(tag)

		if queryValue == "" {
			continue
		}

		if fieldValue.CanSet() {
			err := setValue(fieldValue, queryValue)
			if err != nil {
				return fmt.Errorf("%s %s(%s)", err.Error(), fieldMetadata.Name, fieldMetadata.Type)
			}
		} else {
			return fmt.Errorf("cannot set field %s, be sure it is exported", fieldMetadata.Name)
		}

	}
	return nil

}

var errNotPtr = errors.New("not a pointer to a struct")

// reflecter checks if the input is a pointer to a struct.
func reflecter(st any) (reflect.Value, error) {
	v := reflect.ValueOf(st)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return v, errNotPtr
	}
	return v, nil
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
