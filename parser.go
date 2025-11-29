package qparser

import (
	"fmt"
	"reflect"
	"strconv"
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
		slice, err := parseSliceFromStrings(vals, elemType)
		if err != nil {
			return err
		}
		if slice.Len() == 0 {
			return nil
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
	// Parse directly from comma-separated values without splitting
	slice, err := parseSliceFromStrings(vals, ft)
	if err != nil {
		return err
	}
	if slice.Len() == 0 {
		return nil
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

// parseSliceFromStrings parses comma-separated values directly into a slice without intermediate allocations
func parseSliceFromStrings(vals []string, sliceType reflect.Type) (reflect.Value, error) {
	if len(vals) == 0 {
		return reflect.Zero(sliceType), nil
	}

	// Count total elements needed for pre-allocation
	totalElements := 0
	for _, v := range vals {
		if v == "" {
			continue
		}
		// Count commas + 1 for number of elements
		for i := 0; i < len(v); i++ {
			if v[i] == ',' {
				totalElements++
			}
		}
		totalElements++ // +1 for the string itself
	}

	if totalElements == 0 {
		return reflect.Zero(sliceType), nil
	}

	slice := reflect.MakeSlice(sliceType, totalElements, totalElements)
	elemType := sliceType.Elem()
	elemIndex := 0

	// Parse directly without creating intermediate strings
	for _, v := range vals {
		vLen := len(v)
		if vLen == 0 {
			continue
		}

		start := 0
		for i := 0; i <= vLen; i++ {
			if i == vLen || v[i] == ',' {
				// Trim whitespace using indices directly
				trimStart := start
				trimEnd := i

				// Trim leading whitespace - optimized with single comparison
				for trimStart < trimEnd {
					c := v[trimStart]
					if c > ' ' && c != '\t' && c != '\n' && c != '\r' {
						break
					}
					trimStart++
				}

				// Trim trailing whitespace - optimized
				for trimStart < trimEnd {
					c := v[trimEnd-1]
					if c > ' ' && c != '\t' && c != '\n' && c != '\r' {
						break
					}
					trimEnd--
				}

				// Only process non-empty trimmed parts
				if trimStart < trimEnd {
					if err := setSingleValue(v[trimStart:trimEnd], slice.Index(elemIndex), elemType); err != nil {
						return reflect.Zero(sliceType), fmt.Errorf("element [%d]: %w", elemIndex, err)
					}
					elemIndex++
				}
				start = i + 1
			}
		}
	}

	// Resize slice if we skipped empty elements
	if elemIndex < totalElements {
		slice = slice.Slice(0, elemIndex)
	}

	return slice, nil
}
