package qparser

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

var ErrNotPtr = errors.New("not a pointer to a struct")

const structTAG = "qp"

// ParseRequest parses the URL query parameters from the HTTP request and sets the corresponding fields in the struct.
// dst should be a pointer to a struct.
func ParseRequest(r *http.Request, dst any) error {
	queryValues := r.URL.Query()
	return parse(queryValues, dst)
}

// ParseURL parses query parameters from the URL string and sets the corresponding fields in the struct.
// The URL string should be in the format "http://example.com?qpkey=value1&qpkey2=value2".
// dst should be a pointer to a struct.
func ParseURL(address string, dst any) error {
	urlObj, err := url.Parse(address)
	if err != nil {
		return err
	}
	queryValues := urlObj.Query()
	return parse(queryValues, dst)
}

// reflecter checks if the input is a pointer to a struct.
func reflecter(st any) (reflect.Value, error) {
	v := reflect.ValueOf(st)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return v, ErrNotPtr
	}
	return v, nil
}

func parse(queryValues url.Values, dst any) error {
	v, err := reflecter(dst)
	if err != nil {
		return err
	}
	v = v.Elem()
	t := v.Type()

	// Iterate over the struct fields
	for i := 0; i < t.NumField(); i++ {
		fieldMetadata := t.Field(i)
		fieldValue := v.Field(i)

		isPtr := fieldValue.Kind() == reflect.Ptr

		isStruct := fieldValue.Kind() == reflect.Struct
		if isPtr {
			isStruct = fieldValue.Type().Elem().Kind() == reflect.Struct
		}

		// TODO: Maybe spawning goroutine for recursive calls is a good idea
		if isStruct {
			if !isPtr {
				if err := parse(queryValues, fieldValue.Addr().Interface()); err != nil {
					return err
				}
				continue
			}

			// We need to dereference the pointer to the struct before performing the recursive call
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}

			if err := parse(queryValues, fieldValue.Interface()); err != nil {
				return err
			}

			// Check if all fields inside the child struct are zero and set the parent pointer to nil if so
			// This helps to provide consistency on checking pointers
			isZero := true
			subValue := fieldValue.Elem()
			for j := 0; j < subValue.NumField(); j++ {
				if !subValue.Field(j).IsZero() {
					isZero = false
					break
				}
			}

			if isZero {
				fieldValue.Set(reflect.Zero(fieldValue.Type()))
			}

			continue
		}

		tag := strings.TrimSpace(fieldMetadata.Tag.Get(structTAG))
		if tag == "" || tag == "-" {
			continue
		}

		// Retrieve the query parameter value
		queryValue := queryValues.Get(tag)
		if queryValue == "" {
			continue
		}

		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set field %s, be sure it is exported", fieldMetadata.Name)
		}

		if err := setValue(fieldValue, queryValue); err != nil {
			return fmt.Errorf("%s %s(%s)", err.Error(), fieldMetadata.Name, fieldMetadata.Type)
		}
	}

	return nil
}
