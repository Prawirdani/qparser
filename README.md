[![Go Reference](https://pkg.go.dev/badge/github.com/prawirdani/qparser?status.svg)](https://pkg.go.dev/github.com/prawirdani/qparser?tab=doc)
[![codecov](https://codecov.io/github/Prawirdani/qparser/graph/badge.svg)](https://codecov.io/github/Prawirdani/qparser)
[![Go Report Card](https://goreportcard.com/badge/github.com/prawirdani/qparser)](https://goreportcard.com/report/github.com/prawirdani/qparser)
![Build Status](https://github.com/prawirdani/qparser/actions/workflows/ci.yml/badge.svg)

`qparser` is a lightweight Go package designed to parse URL query parameters directly into Go structs. It is built on the Go standard library, offering a simple and focused solution for handling incoming query strings.

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
<div align="center">

| Pattern               | URL                                                                    |
| :---------------------|:-----------------------------------------------------------------------|
| Comma-separated       | /menus?`categories`=desserts,beverages,sides                           |
| Repeated Keys         | /menus?`categories`=desserts&`categories`=beverages&`categories`=sides |
| Combination of both   | /menus?`categories`=desserts,beverages&`categories`=sides              |
</div>

Simply ensure that the qp tags are defined appropriately in your struct fields to map these parameters correctly.

### Time Handling
Supports time.Time, *time.Time, and type aliases. Handles a variety of standard time formats, both with and without timezone offsets, and supports nanosecond-level precision. Date formats follow the YYYY-MM-DD layout.
<div align="center">

| Format Description                         |Layout Example                        |
| :------------------------------------------|:-------------------------------------|
| RFC3339 (Z or offset)                      |`2006-01-02T15:04:05Z07:00`           |
| RFC3339Nano (Z or offset, nanosecond prec) |`2006-01-02T15:04:05.999999999Z07:00` |
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
| Space separator + TZ                       |`2006-01-02 15:04:05-07:00`           |
| Space separator + TZ (+ offset)            |`2006-01-02 15:04:05+07:00`           |
| Space separator + fractional + TZ          |`2006-01-02 15:04:05.123456789-07:00` |

</div>

- Fractional seconds (milliseconds, microseconds, nanoseconds) are supported with or without a timezone.
- Timezones may use Z, +HH:MM, or -HH:MM.
- No support for named timezones (e.g., PST, UTC)—only numeric offsets.

Example:
```go
type ReportFilter struct {
    From time.Time `qp:"from"`
    To   time.Time `qp:"to"`
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

- Empty query values are not validated by default. For custom validation (including empty value checks), implement your own validation method on the struct or use a third-party validator such as go-playground/validator.
- Missing query parameters:
  - Primitive fields keep their zero values (0, "", false, etc.).
  - Pointer-to-primitive fields (e.g., `*string`, `*int`) remain `nil` when the parameter is missing. They are only allocated when the parameter is provided.
  - Slice fields (`[]T`) remain `nil` when the parameter is missing. They are allocated only when at least one value is successfully decoded.
  - Pointer-to-slice fields (`*[]T`) remain `nil` when the parameter is missing. They are allocated only when the parameter is provided.
  - Pointer-to-struct fields are **always initialized**, even when the nested parameters are missing. They contain the zero value of the struct.
- For repeated query parameters, the value is appended to the slice every time. If you want deduplication or sanitization, implement a post-processing method on your struct.
- The `qp` tag is case-sensitive and must match the query parameter key exactly.
- Pointer-to-struct fields offer no practical benefit because they are always initialized and never `nil`, you cannot rely on `nil` checks to detect whether a nested parameter group was supplied. If you need that behavior, inspect field values or apply custom post-processing.



## Benchmarks

```text
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i5-8259U CPU @ 2.30GHz
=== Sequential ===========================================================================================
Benchmark/Seq/Minimal-8                  1204855          983.8 ns/op          224 B/op        1 allocs/op
Benchmark/Seq/1-date-8                    930016           1233 ns/op          248 B/op        2 allocs/op
Benchmark/Seq/2-dates-8                   626414           1599 ns/op          272 B/op        3 allocs/op
Benchmark/Seq/slices-string-1*50-8        348960           2924 ns/op         1144 B/op        3 allocs/op
Benchmark/Seq/slices-int-1*50-8           292065           3613 ns/op          664 B/op        3 allocs/op
Benchmark/Seq/slices-2*50-8               202550           5514 ns/op         1584 B/op        5 allocs/op
Benchmark/Seq/slices-2*100-8              123055           9351 ns/op         2960 B/op        5 allocs/op
Benchmark/Seq/2*25-slices-and-2-dates-8   257817           4158 ns/op          944 B/op        7 allocs/op
=== Parallel =============================================================================================
Benchmark/Par/Minimal-8                  4529499          261.3 ns/op          224 B/op        1 allocs/op
Benchmark/Par/1-date-8                   3590239          356.3 ns/op          248 B/op        2 allocs/op
Benchmark/Par/2-dates-8                  2865769          420.1 ns/op          272 B/op        3 allocs/op
Benchmark/Par/slices-string-1*50-8       1261774          963.2 ns/op         1144 B/op        3 allocs/op
Benchmark/Par/slices-int-1*50-8          1136395           1025 ns/op          664 B/op        3 allocs/op
Benchmark/Par/slices-2*50-8               684651           1758 ns/op         1584 B/op        5 allocs/op
Benchmark/Par/slices-2*100-8              410012           2982 ns/op         2960 B/op        5 allocs/op
Benchmark/Par/2*25-slices-and-2-dates-8   990036           1244 ns/op          944 B/op        7 allocs/op
```
### Noticeable Behaviors

#### Allocation Behavior
- **Base cost**: 1 allocation for parsing infrastructure
- **Dates**: Each successful parse creates exactly 1 allocation
- **Slices**: Each slice parsing creates exactly 2 allocations

#### Scaling Behavior
- **Linear growth**: Time and memory scale proportionally with item count
- **Predictable**: 2× items = ~2× time, exactly 2× memory
- **Consistent**: Same allocation count regardless of data size

#### Cache Behavior
- **First use**: Reflection builds struct metadata cache
- **Subsequent uses**: Zero-allocation cache lookups
- **Thread-safe**: sync.Map enables concurrent access
- **Persistent**: Cache lives for application lifetime

#### Concurrency
Parallel execution leverages sync.Map for effective caching, significantly improving performance under concurrent workloads like HTTP handlers by reducing per-operation time
