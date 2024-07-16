
## ⚠️ Warning ⚠️
This package is still under development and not ready for production use.\
\
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
By default, it does not validate empty query values. To perform empty value validation or implement some business rules validation, you can create your own validator or use a validator package like [validator](https://github.com/go-playground/validator).

### Multiple Values Query & Nested Struct
To allow multiple values for a single query parameter, you can use a slice type.
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
The multiple values are separated by comma `,` in the query string. For example, `/menus?categories=desserts,beverages`. For the nested struct, simply just pass the `qp` tag definition in the nested struct field. So the final query string will look like `/menus?categories=desserts,beverages&available=true&page=1&limit=5`.

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

## Supported field types
Currently, it only supports basic primitive types such as:
- String
- Boolean
- Integers (int, int8, int16, int32 and int64)
- Unsigned Integers (uint, uint8, uint16, uint32 and uint64)
- Floats (float64 and float32)
- Slice of above types
- Nested Struct
- A pointer to one of above

## Future plans
- Support for various types such as multidimensional slice and slice of Struct
- Default value mechanism
- Mapped Errors
- Custom multivalues separator

## Contribution
Contributions are welcome! If you have any improvements, bug fixes, or new features you'd like to add.

