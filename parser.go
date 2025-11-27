package qparser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// parseStruct traverses struct fields and maps query parameters to field values
func parseStruct(query map[string][]string, rv reflect.Value, rt reflect.Type) error {
	info := getStructCache(rt)
	if info.hasUnexportedWithTag {
		return ErrUnexportedStruct
	}

	for _, field := range info.fields {
		if field.isNested {
			if err := parseNestedField(query, rv, field, info.name); err != nil {
				return err
			}
			continue
		}

		vals, ok := query[field.tag]
		if !ok {
			continue
		}

		fv := rv.FieldByIndex(field.index)
		if err := setFieldValue(fv, field.typ, vals); err != nil {
			return wrapFieldError(fmt.Sprintf("%s.%s", info.name, field.name), err)
		}
	}
	return nil
}

// parseNestedField handles embedded or nested struct fields
func parseNestedField(query map[string][]string, rv reflect.Value, field fieldInfo, parentName string) error {
	fv := rv.FieldByIndex(field.index)
	ft := field.typ

	// Handle pointer to struct
	if ft.Kind() == reflect.Ptr {
		if fv.IsNil() {
			fv.Set(reflect.New(ft.Elem()))
		}
		fv = fv.Elem()
		ft = ft.Elem()
	}

	if err := parseStruct(query, fv, ft); err != nil {
		return wrapFieldError(fmt.Sprintf("%s.%s", parentName, field.name), err)
	}
	return nil
}

// setFieldValue routes to the appropriate handler based on field type
func setFieldValue(fv reflect.Value, ft reflect.Type, vals []string) error {
	switch ft.Kind() {
	case reflect.Ptr:
		return setPtrField(fv, ft.Elem(), vals)
	case reflect.Slice:
		return setSliceField(fv, ft, vals)
	default:
		if len(vals) == 0 {
			return nil
		}
		return setSingleValue(vals[0], fv, ft)
	}
}

// setPtrField handles pointer fields, including *[]T
func setPtrField(fv reflect.Value, elemType reflect.Type, vals []string) error {
	if elemType.Kind() == reflect.Slice {
		parts := splitAndTrim(vals)
		if len(parts) == 0 {
			return nil
		}
		slice := reflect.MakeSlice(elemType, len(parts), len(parts))
		if err := fillSlice(slice, elemType.Elem(), parts); err != nil {
			return err
		}
		ptr := reflect.New(elemType)
		ptr.Elem().Set(slice)
		fv.Set(ptr)
		return nil
	}

	if len(vals) == 0 || vals[0] == "" {
		return nil
	}

	elemVal := reflect.New(elemType)
	if err := setSingleValue(vals[0], elemVal.Elem(), elemType); err != nil {
		return err
	}
	fv.Set(elemVal)
	return nil
}

// setSliceField handles slice fields
func setSliceField(fv reflect.Value, ft reflect.Type, vals []string) error {
	parts := splitAndTrim(vals)
	if len(parts) == 0 {
		return nil
	}

	slice := reflect.MakeSlice(ft, len(parts), len(parts))
	if err := fillSlice(slice, ft.Elem(), parts); err != nil {
		return err
	}
	fv.Set(slice)
	return nil
}

// setSingleValue parses a single value and sets it on the reflect.Value
func setSingleValue(val string, fv reflect.Value, typ reflect.Type) error {
	// No look up table, just raw dog switch for maximum perf
	// WARN: mega switch for raw performance. Maintain with care.
	switch typ.Kind() {
	// ----- Pointers -----
	case reflect.Ptr:
		elemType := typ.Elem()
		elemVal := reflect.New(elemType)
		if err := setSingleValue(val, elemVal.Elem(), elemType); err != nil {
			return err
		}
		fv.Set(elemVal)

	// ----- Special structs -----
	case reflect.Struct:
		if typ == reflect.TypeOf(time.Time{}) {
			t, err := parseTime(val)
			if err != nil {
				return err
			}
			fv.Set(reflect.ValueOf(t))
			return nil
		}
		return fmt.Errorf("%w: %v", ErrUnsupportedKind, typ.Kind())

	// ----- Strings -----
	case reflect.String:
		fv.SetString(val)

	// ----- Booleans -----
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return ErrInvalidValue
		}
		fv.SetBool(b)

	// ----- Signed integers -----
	case reflect.Int:
		n, err := strconv.ParseInt(val, 10, strconv.IntSize)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetInt(n)

	case reflect.Int8:
		n, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetInt(n)

	case reflect.Int16:
		n, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetInt(n)

	case reflect.Int32:
		n, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetInt(n)

	case reflect.Int64:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetInt(n)

	// ----- Unsigned integers -----
	case reflect.Uint:
		n, err := strconv.ParseUint(val, 10, strconv.IntSize)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetUint(n)

	case reflect.Uint8:
		n, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetUint(n)

	case reflect.Uint16:
		n, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetUint(n)

	case reflect.Uint32:
		n, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetUint(n)

	case reflect.Uint64:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetUint(n)

	// ----- Floats -----
	case reflect.Float32:
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetFloat(f)

	case reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return strconvNumError(err, val)
		}
		fv.SetFloat(f)

	default:
		return fmt.Errorf("%w: %v", ErrUnsupportedKind, typ.Kind())
	}

	return nil
}

// fillSlice populates a slice with parsed values
func fillSlice(slice reflect.Value, elemType reflect.Type, parts []string) error {
	for i, part := range parts {
		if err := setSingleValue(part, slice.Index(i), elemType); err != nil {
			return fmt.Errorf("element [%d]: %w", i, err)
		}
	}
	return nil
}

// splitAndTrim splits comma-separated values and trims whitespace
func splitAndTrim(vals []string) []string {
	var parts []string
	for _, v := range vals {
		for p := range strings.SplitSeq(v, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				parts = append(parts, p)
			}
		}
	}
	return parts
}
