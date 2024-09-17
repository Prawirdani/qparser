`qparser` is a simple package that help parse query parameters into struct in Go. It is inspired by [gorilla/schema](https://github.com/gorilla/schema) with main focus on query parameters. Built on top of Go stdlib, it uses custom struct tag `qp` to define the query parameter key .

## Installation
```bash
go get -u github.com/prawirdani/qparser
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
To allow multiple values for a single query parameter, you can use a slice type. For **nested structs**, you can pass the query parameters using the qp tag within the nested struct fields. This allows you to build more complex query strings that include nested parameters.
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
Multiple values in the query string should be **separated by commas**. For example: `/menus?categories=desserts,beverages`.
A complete query string based on the above example might look like this: `/menus?categories=desserts,beverages&available=true&page=1&limit=5`.

Simply ensure that the qp tags are defined appropriately in your nested struct fields to map these parameters correctly.

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
    - Regular fields keep their zero values (e.g., `0` for int, `""` for string, `false` for bool)
    - Pointer and slice fields remain `nil`

## Supported field types
- String
- Boolean
- Integers (int, int8, int16, int32 and int64)
- Unsigned Integers (uint, uint8, uint16, uint32 and uint64)
- Floats (float64 and float32)
- Slice of above types
- Nested Struct
- A pointer to one of above

## Future plans
- Introduce benchmarking to assess and optimize performance.
- Enhance overall performance and efficiency.
- Expand support to include various types, such as multidimensional slices and slices of structs, complex types, and more.
- Implement a default value mechanism for enhanced flexibility.
- Provide mapped errors for clearer error handling.
- Add support for custom multi-value separators.

## Undiscovered Scenarios and Behavior
- What happens if fields have the same query parameter tag?


