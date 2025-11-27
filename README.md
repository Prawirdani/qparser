[![Go Reference](https://pkg.go.dev/badge/github.com/prawirdani/qparser?status.svg)](https://pkg.go.dev/github.com/prawirdani/qparser?tab=doc)
[![codecov](https://codecov.io/github/Prawirdani/qparser/graph/badge.svg)](https://codecov.io/github/Prawirdani/qparser)
[![Go Report Card](https://goreportcard.com/badge/github.com/prawirdani/qparser)](https://goreportcard.com/report/github.com/prawirdani/qparser)
![Build Status](https://github.com/prawirdani/qparser/actions/workflows/ci.yml/badge.svg)

`qparser` is a simple package that helps parse query parameters into structs in Go. It is inspired by [gorilla/schema](https://github.com/gorilla/schema) with a main focus on query parameters. Built on top of Go stdlib, it uses a custom struct tag `qp` to define the query parameter key.

## Table of Contents
- [Installation](#installation)
- [Examples](#examples)
- [Supported field types](#supported-field-types)
- [Error Handling](#error-handling)
- [Notes](#notes)
- [Benchmarks](#benchmarks)

## Installation
```bash
go get -u github.com/prawirdani/qparser@latest
```

## Examples

### Parse from `net/http` request.
Here's an example of how to use `qparser` to parse query parameters into struct from `net/http` request.
```go
// Representing basic pagination, /path?page=1&limit=5
type Pagination struct {
    Page    int `qp:"page"`
    Limit   int `qp:"limit"`
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
    var pagination Pagination

    err := qparser.ParseRequest(r, &pagination)
    if err != nil {
        // Handle Error
    }
        
    // Do something with pagination
}
```

### Parse from URL
You can also parse query parameters from URL string by calling the `ParseURL` function. Here's an example:
```go

func main() {
    var pagination Pagination

    url := "http://example.com/path?page=1&limit=5"

    err := qparser.ParseURL(url, &pagination)
    if err != nil {
        // Handle Error
    }

    // Do something with pagination
}
```

### Parse from url.Values
You can also parse query parameters directly from `url.Values` by calling the `Parse` function. Here's an example:
```go

func main() {
    var pagination Pagination

    values := url.Values{
        "page":  []string{"1"},
        "limit": []string{"5"},
    }

    err := qparser.Parse(values, &pagination)
    if err != nil {
        // Handle Error
    }

    // Do something with pagination
}
```
### Multiple Values Query & Nested Struct
To support multiple values for a single query parameter, use a slice type. For nested structs, utilize the qp tag within the fields of the nested struct to pass the query parameters. It's important to note that the parent struct containing the nested/child struct **should not have its own qp tag**. Here's an example:
```go
// Representing filter for menu
type MenuFilter struct {
    Categories []string `qp:"categories"`
    Available  bool     `qp:"available"`
}

type Pagination struct {
    Page    int `qp:"page"`
    Limit   int `qp:"limit"`
}

type MenuQueryParams struct {
    Filter     MenuFilter
    Pagination Pagination
}

func GetMenus(w http.ResponseWriter, r *http.Request) {
    var f MenuQueryParams
    if err := qparser.ParseRequest(r, &f); err != nil {
        // Handle Error
    }
    // Do something with f.Filter and f.Pagination
}
```
There are three ways for the parser to handle multiple values query parameters:
1. Comma-separated values: `/menus?categories=desserts,beverages,sides`
2. Repeated Keys: `/menus?categories=desserts&categories=beverages&categories=sides`
3. Combination of both: `/menus?categories=desserts,beverages&categories=sides`

Simply ensure that the qp tags are defined appropriately in your struct fields to map these parameters correctly.

### Time Handling
Supports time.Time, *time.Time, and type aliases. Handles a variety of standard time formats, both with and without timezone offsets, and supports nanosecond-level precision. Date formats follow the YYYY-MM-DD layout.
<div align="center">

| Format Description                         |Layout Example                        |
| :------------------------------------------|:-------------------------------------|
| Time only                                  |`15:04:05`                            |
| Date only                                  |`2006-01-02`                          |
| Date & time (space separated)              |`2006-01-02 15:04:05`                 |
| Date & time + milliseconds                 |`2006-01-02T15:04:05.000`             |
| Date & time + microseconds                 |`2006-01-02T15:04:05.000000`          |
| Date & time + nanoseconds                  |`2006-01-02T15:04:05.999999999`       |
| Date & time + TZ offset                    |`2006-01-02T15:04:05-07:00`           |
| Date & time + milliseconds + TZ            |`2006-01-02T15:04:05.000-07:00`       |
| Date & time + microseconds + TZ            |`2006-01-02T15:04:05.000000-07:00`    |
| Date & time + nanoseconds + TZ             |`2006-01-02T15:04:05.999999999-07:00` |
| RFC3339 (Z or offset)                      |`2006-01-02T15:04:05Z07:00`           |
| RFC3339Nano (Z or offset, nanosecond prec) |`2006-01-02T15:04:05.999999999Z07:00` |
| Space separator + TZ                       |`2006-01-02 15:04:05-07:00`           |
| Space separator + TZ (+ offset)            |`2006-01-02 15:04:05+07:00`           |
| Space separator + fractional + TZ          |`2006-01-02 15:04:05.123456789-07:00`|

</div>

- Fractional seconds (milliseconds, microseconds, nanoseconds) are supported with or without a timezone.
- Timezones may use Z, +HH:MM, or -HH:MM.
- No support for named timezones (e.g., PST, UTC)—only numeric offsets.

Example:
```go
type ReportFilter struct {
    From time.Time `qp:"from"`
    To time.Time `qp:"to"`
}
func main() {
    var filter ReportFilter 

    url := "http://example.com/reports?from=2025-07-01&to=2025-07-31"

    err := qparser.ParseURL(url, &filter)
    if err != nil {
        // Handle Error
    }
    // Do something with filter
}
```

## Supported field types
- String
- Boolean
- Integers (int, int8, int16, int32 and int64)
- Unsigned Integers (uint, uint8, uint16, uint32 and uint64)
- Floats (float64 and float32)
- Slice of above types
- Nested Struct
- time.Time
- A pointer to one of above


## Error Handling
qparser provides detailed error information with field-specific context. Errors are wrapped with field names to help identify exactly which parameter failed to parse.

```go
type UserFilter struct {
    Age    int    `qp:"age"`
    Active bool   `qp:"active"`
    Name   string `qp:"name"`
}

func main() {
    var filter UserFilter
    
    // Simulate invalid query parameters
    values := url.Values{
        "age":    []string{"invalid_age"},  // Invalid integer
        "active": []string{"maybe"},        // Invalid boolean
        "name":   []string{"John"},         // Valid
    }
    
    err := qparser.Parse(values, &filter)
    if err != nil {
        // Check for specific error types
        var fieldErr *qparser.FieldError
        if errors.As(err, &fieldErr) {
            fmt.Printf("Field error in %q: %v\n", fieldErr.FieldName, fieldErr.Err)
            
            // Handle specific error types
            switch {
            case errors.Is(fieldErr.Err, qparser.ErrInvalidValue):
                fmt.Printf("Invalid value format for field %s\n", fieldErr.FieldName)
            case errors.Is(fieldErr.Err, qparser.ErrOutOfRange):
                fmt.Printf("Value out of range for field %s\n", fieldErr.FieldName)
            case errors.Is(fieldErr.Err, qparser.ErrUnsupportedKind):
                fmt.Printf("Unsupported type for field %s\n", fieldErr.FieldName)
            }
        } else {
            fmt.Printf("General error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("Parsed filter: %+v\n", filter)
}
```

### Error Types

qparser defines several error types for different parsing scenarios:

- **`ErrInvalidValue`**: Value cannot be parsed as the target type (e.g., "abc" as integer)
- **`ErrOutOfRange`**: Value is too large for the target numeric type (e.g., "999" as int8)
- **`ErrUnsupportedKind`**: Target type is not supported by the parser
- **`ErrUnexportedStruct`**: Struct contains unexported fields with `qp` tags

### FieldError Structure

The `FieldError` type provides both the field name and the underlying error:

```go
type FieldError struct {
    FieldName string  // Name of the field that failed
    Err       error   // Underlying error
}
```

This allows you to:
- Identify exactly which field failed parsing
- Access the specific error type for targeted error handling
- Provide user-friendly error messages in your API responses


## Notes
- Empty query values are not validated by default. For custom validation (including empty value checks), implement your own validation method on the struct or use a third-party validator such as [go-playground/validator](https://github.com/go-playground/validator).
- Missing query parameters:
    - Primitive fields keep their zero values (0, "", false, etc.).
    - Pointer fields are always initialized, even when the parameter is missing. They will contain the zero value of the underlying type.
    - Slice fields become an empty slice ([]T{}), not nil.
    - Pointer nested structs are also always initialized. Missing fields inside them receive their zero values.
- For repeated query parameters, the value is appended to the slice every time. If you want deduplication or sanitization, implement a post-processing method on your struct.
- The qp tag is case-sensitive and must match the query parameter key exactly.


## Benchmarks
```text
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i5-8259U CPU @ 2.30GHz
Benchmark/Small-8              224376     4602 ns/op       2680 B/op      13 allocs/op
Benchmark/Medium-8             246818     4851 ns/op       2720 B/op      13 allocs/op
Benchmark/Large-8              173422     5867 ns/op       3056 B/op      21 allocs/op
Benchmark/LargeWithDate-8      160401     6543 ns/op       3104 B/op      23 allocs/op
```
