// Package qparser provides a query parameter decoder for Go.
//
// It parses URL query parameters into user-defined structs using reflection,
// supporting nested structs, slices, pointer fields, numeric types, booleans,
// strings, and time.Time with multiple timestamp formats.
// The Parse, ParseRequest, and ParseURL functions all decode query parameters
// into a struct value provided by the caller.
package qparser

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
)

// Parse decodes the provided url.Values into the struct pointed to by dst.
//
// dst must be a pointer to a struct. Nested fields, slices, primitive types,
// pointer fields, and time.Time are supported.
//
// Example:
//
//	var f Filter
//	err := qparser.Parse(url.Values{"age": {"30"}}, &f)
func Parse(values url.Values, dst any) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("dst must be a pointer to struct")
	}
	rv = rv.Elem()
	rt := rv.Type()
	return parseStruct(values, rv, rt)
}

// ParseRequest extracts the query parameters from an http.Request and
// decodes them into the struct pointed to by dst.
//
// Equivalent to calling Parse(r.URL.Query(), dst).
func ParseRequest(r *http.Request, dst any) error {
	query := r.URL.Query()
	return Parse(query, dst)
}

// ParseURL parses the query parameters from the provided URL string and
// decodes them into the struct pointed to by dst.
//
// Returns an error if the URL cannot be parsed.
func ParseURL(addr string, dst any) error {
	urlObj, err := url.Parse(addr)
	if err != nil {
		return err
	}
	queryValues := urlObj.Query()
	return Parse(queryValues, dst)
}
