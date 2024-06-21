package qparser

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

const structTAG = "qp"

// ParseRequest parses the URL query parameters from the HTTP request and sets the corresponding fields in the struct.
// The struct fields must have a `qp` tag with the query parameter name.
// Fields with a `qp` tag value of "-" will be ignored.
func ParseRequest(r *http.Request, st any) error {
	v, err := reflecter(st)
	if err != nil {
		return err
	}

	v = v.Elem()
	t := v.Type()

	queryValues := r.URL.Query()
	// Iterate over the struct fields
	for i := 0; i < t.NumField(); i++ {
		fieldMetadata := t.Field(i)
		fieldValue := v.Field(i)

		// Retrieve the qp tag value
		tag := strings.TrimSpace(fieldMetadata.Tag.Get(structTAG))
		if tag == "" || tag == "-" {
			continue
		}

		// Retrieve the query parameter value
		queryValue := queryValues.Get(tag)

		if queryValue == "" {
			continue
		}

		if fieldValue.CanSet() {
			err := SetValue(fieldValue, queryValue)
			if err != nil {
				return fmt.Errorf("%s %s(%s)", err.Error(), fieldMetadata.Name, fieldMetadata.Type)
			}
		} else {
			return fmt.Errorf("qparser: cannot set field %s", fieldMetadata.Name)
		}

	}
	return nil
}

// reflecter checks if the input is a pointer to a struct.
func reflecter(st any) (reflect.Value, error) {
	var err error
	v := reflect.ValueOf(st)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("qparser: st must be a pointer to a struct")
	}
	return v, err
}
