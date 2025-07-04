[![Go Reference](https://pkg.go.dev/badge/github.com/prawirdani/qparser?status.svg)](https://pkg.go.dev/github.com/prawirdani/qparser?tab=doc)
[![codecov](https://codecov.io/github/Prawirdani/qparser/graph/badge.svg)](https://codecov.io/github/Prawirdani/qparser)
[![Go Report Card](https://goreportcard.com/badge/github.com/prawirdani/qparser)](https://goreportcard.com/report/github.com/prawirdani/qparser)
![Build Status](https://github.com/prawirdani/qparser/actions/workflows/ci.yml/badge.svg)

`qparser` is a simple package that help parse query parameters into struct in Go. It is inspired by [gorilla/schema](https://github.com/gorilla/schema) with main focus on query parameters. Built on top of Go stdlib, it uses custom struct tag `qp` to define the query parameter key .

## Installation
```bash
go get -u github.com/prawirdani/qparser@latest
```

## Example
Here's an example of how to use `qparser` to parse query parameters into struct.
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

## Notes
- Empty query values are not validated by default. For custom validation (including empty value checks), you can create your own validator by creating a pointer/value receiver method on the struct or with help of a third party validator package like [go-playground/validator](https://github.com/go-playground/validator).
- Missing query parameters:
    - Primitive type fields keep their zero values (e.g., `0` for int, `""` for string, `false` for bool)
    - Pointer fields are remain **nil** and slice are set to nil slice ([]).
    - A pointer nested struct will remain nil, **only if all the fields are missing**. If any field is present, the struct will be initialized and the missing fields will be set to their zero values.
- For multiple values query parameters, same value will be appended to the slice. If you want to make sure each value in the slice is unique, you can create a pointer receiver method on the struct to remove the duplicates or sanitize the values.
- The `qp` tag value is **case-sensitive** and must match the query parameter key exactly.

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
| Space separator + fractional + TZ          |`2006-01-02 15:04:05.123456789 -07:00`|

</div>

- Fractional seconds (milliseconds, microseconds, nanoseconds) are supported with or without a timezone.
- Timezones may use Z, +HH:MM, or -HH:MM.
- No support for named timezones (e.g., PST, UTC)—only numeric offsets.

Example:
```go
type ReportSearchQuery struct {
    From time.Time `qp:"from"`
    To time.Time `qp:"to"`
}
func main() {
    var search ReportSearchQuery

    url := "http://example.com/reports?from=2025-07-01&to=2025-07-31"

    err := qparser.ParseURL(url, &search)
    if err != nil {
        // Handle Error
    }
    // Do something with search 
}
```

## TODO 
- [ ] Introduce benchmarking to evaluate and optimize performance
- [ ] Add support for default values via struct tags for better flexibility
- [ ] Explore parsing into interface fields without knowing the concrete type (may not be feasible due to reflection limits)
- [ ] Provide structured, mapped errors for clearer debugging and validation feedback
- [ ] Support custom multi-value separators to improve slice parsing flexibility
