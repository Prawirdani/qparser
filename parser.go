package qparser

import (
	"fmt"
	"net/http"
	"reflect"
)

const structTAG = "qp"

// ParseURLQuery parses the URL query parameters and assigns them to the struct fields
// Using the qp struct tag to map the query parameter to the struct field.
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
		tag := fieldType.Tag.Get(structTAG)
		if tag == "" {
			continue
		}
		// Retrieve the query parameter value
		queryValue := r.URL.Query().Get(tag)
		f := registerField(fieldType, fieldValue)

		if err := f.Set(queryValue); err != nil {
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
