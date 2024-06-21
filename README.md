## ⚠️ Warning ⚠️
This package is still under development and not ready for production use.\
\
`qparser` is a simple package that help parse query parameters into struct in Go. It is inspired by [gorilla/schema](https://github.com/gorilla/schema) with main focus on query parameters. Built on top of Go stdlib, it uses custom struct tag `qp` to define the query parameter key and Go `net/http` package to retrieve the URL query values and heavily relies on `strconv` package to convert string values into desired types.

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

### Multiple Values Query
To allow multiple values for a single query parameter, you can use a slice type in the struct. Here's an example:
```go
// Representing filter for menu
type MenuFilter struct {
	Categories []string `qp:"categories"`
}

func GetMenus(w http.ResponseWriter, r *http.Request) {
	var f MenuFilter
	if err := qparser.ParseRequest(r, &f); err != nil {
        // Handle Error
	}
    // Do something with f.Categories
}
```
The multiple values are separated by comma `,` in the query string. For example, `/menus?categories=desserts,beverages`.

## Supported field types
Currently, it only supports basic primitive types such as:
- String
- Boolean
- Integers (int, int8, int16, int32 and int64)
- Unsigned Integers (uint, uint8, uint16, uint32 and uint64)
- Floats (float64 and float32)
- Slice of above types
- A pointer to one of above

## Future plans
- Support for various types such as multidimensional slice and nested struct
- Default value mechanism
- Mapped Errors
- Custom multivalues separator

## Contribution
CONTRIBUTION GUIDE SOON!\
Contributions are welcome! If you have any improvements, bug fixes, or new features you'd like to add, please let me know.
