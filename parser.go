package qparser

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

const structTAG = "qp"

// ParseURLQuery parses the URL query parameters from the HTTP request and sets the corresponding fields in the struct.
// The struct fields must have a `qp` tag with the query parameter name.
// Fields with a `qp` tag value of "-" will be ignored.
func ParseURLQuery(r *http.Request, st any) error {
	v, err := reflecter(st)
	if err != nil {
		return err
	}

	v = v.Elem()
	t := v.Type()

	// Iterate over the struct fields
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		// Retrieve the qp tag value
		tag := strings.TrimSpace(fieldType.Tag.Get(structTAG))
		if tag == "" || tag == "-" {
			continue
		}
		// Retrieve the query parameter value
		queryValue := r.URL.Query().Get(tag)

		if queryValue == "" {
			continue
		}

		f := registerField(fieldType, fieldValue, queryValue)

		if err := f.SetValue(); err != nil {
			return err
		}
	}
	return nil
}

// reflecter checks if the input is a pointer to a struct.
func reflecter(st any) (reflect.Value, error) {
	var err error
	v := reflect.ValueOf(st)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("st must be a pointer to a struct")
	}
	return v, err
}
