package qparser

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
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

func parse(queryValues url.Values, dst any) error {
	v, err := validateStructPointer(dst)
	if err != nil {
		return err
	}
	v = v.Elem()
	t := v.Type()

	// Iterate over the struct fields
	for i := 0; i < t.NumField(); i++ {
		fieldMetadata := t.Field(i)
		fieldValue := v.Field(i)

		if err := parseField(queryValues, fieldMetadata, fieldValue); err != nil {
			return err
		}
	}
	return nil
}

func validateStructPointer(st any) (reflect.Value, error) {
	v := reflect.ValueOf(st)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return v, ErrNotPtr
	}
	return v, nil
}

// parseField handles parsing a single field based on its type
func parseField(
	queryValues url.Values,
	fieldMetadata reflect.StructField,
	fieldValue reflect.Value,
) error {
	switch {
	// NOTE: Always Check time type before struct type, because time.Time is also struct but we handle it differently
	case isTimeType(fieldValue):
		return parseTimeField(queryValues, fieldMetadata, fieldValue)
	case isStructType(fieldValue):
		return parseStructField(queryValues, fieldValue)
	default:
		return parsePrimitiveField(queryValues, fieldMetadata, fieldValue)
	}
}

// isTimeType checks if the field is time.Time, *time.Time, or aliases of these
func isTimeType(fieldValue reflect.Value) bool {
	timeType := reflect.TypeOf(time.Time{})
	fieldType := fieldValue.Type()

	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	// Check if the type is time.Time or an alias of time.Time
	return fieldType == timeType ||
		fieldType.ConvertibleTo(timeType)
}

// isStructType checks if the field is a struct or pointer to struct (excluding time.Time)
func isStructType(fieldValue reflect.Value) bool {
	if isTimeType(fieldValue) {
		return false
	}

	return fieldValue.Kind() == reflect.Struct ||
		(fieldValue.Kind() == reflect.Ptr && fieldValue.Type().Elem().Kind() == reflect.Struct)
}

// parseTimeField handles time.Time, *time.Time, and any time-alias fields
func parseTimeField(
	queryValues url.Values,
	fieldMetadata reflect.StructField,
	fieldValue reflect.Value,
) error {
	tag := getStructTag(fieldMetadata)
	if tag == "" {
		return nil
	}

	queryValue, exists := queryValues[tag]
	if !exists || len(queryValue) == 0 {
		return nil
	}

	// Get the actual type we're working with (dereference if pointer)
	targetType := fieldValue.Type()
	isPtr := targetType.Kind() == reflect.Ptr
	if isPtr {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(targetType.Elem()))
		}
		targetType = targetType.Elem()
	}

	// Create a time.Time value to parse into
	var timeVal time.Time
	if err := setValue(reflect.ValueOf(&timeVal).Elem(), queryValue); err != nil {
		return fmt.Errorf("%s %s(%s)", err.Error(), fieldMetadata.Name, fieldMetadata.Type)
	}

	// Convert to the target type (handles both time.Time and aliases)
	var finalVal reflect.Value
	if targetType == reflect.TypeOf(time.Time{}) {
		finalVal = reflect.ValueOf(timeVal)
	} else {
		finalVal = reflect.ValueOf(timeVal).Convert(targetType)
	}

	// Set the value (handling pointer case)
	if isPtr {
		fieldValue.Elem().Set(finalVal)

		// Clean up if we got a zero time
		if timeVal.IsZero() {
			fieldValue.Set(reflect.Zero(fieldValue.Type()))
		}
	} else {
		fieldValue.Set(finalVal)
	}

	return nil
}

// parseStructField handles struct and pointer to struct fields
func parseStructField(queryValues url.Values, fieldValue reflect.Value) error {
	var target reflect.Value
	var shouldCleanupNilPointer bool

	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			shouldCleanupNilPointer = true
		}
		target = fieldValue
	} else {
		target = fieldValue.Addr()
	}

	if err := parse(queryValues, target.Interface()); err != nil {
		return err
	}

	// Clean up nil pointer if all fields are zero
	if shouldCleanupNilPointer && isAllFieldsZero(fieldValue.Elem()) {
		fieldValue.Set(reflect.Zero(fieldValue.Type()))
	}

	return nil
}

// parsePrimitiveField handles primitive types and other non-struct fields
func parsePrimitiveField(
	queryValues url.Values,
	fieldMetadata reflect.StructField,
	fieldValue reflect.Value,
) error {
	tag := getStructTag(fieldMetadata)
	if tag == "" {
		return nil
	}

	queryValue, exists := queryValues[tag]
	if !exists || len(queryValue) == 0 {
		return nil
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("cannot set field %s, be sure it is exported", fieldMetadata.Name)
	}

	if err := setValue(fieldValue, queryValue); err != nil {
		return fmt.Errorf("%s %s(%s)", err.Error(), fieldMetadata.Name, fieldMetadata.Type)
	}

	return nil
}

// getStructTag extracts and validates the struct tag
func getStructTag(fieldMetadata reflect.StructField) string {
	tag := strings.TrimSpace(fieldMetadata.Tag.Get(structTAG))
	if tag == "" || tag == "-" {
		return ""
	}
	return tag
}

// isAllFieldsZero checks if all fields in a struct are zero values
func isAllFieldsZero(structValue reflect.Value) bool {
	for i := 0; i < structValue.NumField(); i++ {
		if !structValue.Field(i).IsZero() {
			return false
		}
	}
	return true
}
