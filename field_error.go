package qparser

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	// ErrUnexportedStruct is returned when a struct contains unexported fields with qp tags.
	// This prevents unintentional parsing of private fields.
	ErrUnexportedStruct = errors.New("struct has unexported fields with qp tags")

	// ErrInvalidValue indicates that a value cannot be parsed as the target type.
	// For example, parsing "abc" as an integer would return this error.
	ErrInvalidValue = errors.New("invalid value")

	// ErrOutOfRange indicates that a value is too large for the target numeric type.
	// For example, parsing "9999999999" as an int8 would return this error.
	ErrOutOfRange = errors.New("out of range")

	// ErrUnsupportedKind indicates that the target type is not supported by the parser.
	// This typically occurs with complex types like maps, channels, or unsupported structs.
	ErrUnsupportedKind = errors.New("unsupported kind")
)

type FieldError struct {
	FieldName string
	Err       error
}

func (e *FieldError) Error() string {
	if e.FieldName != "" {
		return fmt.Sprintf("failed to parse %q: %v", e.FieldName, e.Err)
	}
	return e.Err.Error()
}

func (e *FieldError) Unwrap() error {
	return e.Err
}

func wrapFieldError(fieldName string, err error) error {
	if err == nil {
		return nil
	}
	return &FieldError{FieldName: fieldName, Err: err}
}

// strconvNumError normalize strconv number parser (int, uint, floats) parsing error
func strconvNumError(err error, value string) error {
	e := err
	if numErr, ok := err.(*strconv.NumError); ok {
		switch numErr.Err {
		case strconv.ErrRange:
			e = ErrOutOfRange
		case strconv.ErrSyntax:
			e = ErrInvalidValue
		}
	}

	return fmt.Errorf("%w: %v", e, value)
}
