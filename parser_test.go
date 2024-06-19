package qparser

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type s1 struct {
	F1 string   `qp:"f1"`
	F2 *string  `qp:" f2"`
	F3 int      `qp:"f3 "`
	F4 *int     `qp:" f4 "`
	F5 bool     `qp:"f5"`
	F6 *bool    `qp:"f6"`
	F7 float64  `qp:"f7"`
	F8 *float64 `qp:"f8"`
	F9 int      `qp:"-"` // will be ignored
}

func TestParseS1(t *testing.T) {
	t.Run("filled", func(t *testing.T) {
		queryParams := "f1=hello&f2=world&f3=-42&f4=69&f5=true&f6=false&f7=3.14&f8=42.0&f9=42"
		f2 := "world"
		f4 := 69
		f6 := false
		f8 := 42.0
		expected := s1{
			F1: "hello",
			F2: &f2,
			F3: -42,
			F4: &f4,
			F5: true,
			F6: &f6,
			F7: 3.14,
			F8: &f8,
		}

		req, err := http.NewRequest("GET", "http://example.com?"+queryParams, nil)
		require.Nil(t, err, fmt.Sprintf("Error creating request: %s", err))

		rr := httptest.NewRecorder()

		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			var parsed s1
			err := ParseURLQuery(r, &parsed)
			require.Nil(t, err, fmt.Sprintf("Error parsing query: %s", err))
			require.Equal(t, expected.F1, parsed.F1)
			require.Equal(t, expected.F2, parsed.F2)
			require.Equal(t, expected.F3, parsed.F3)
			require.Equal(t, expected.F4, parsed.F4)
			require.Equal(t, expected.F5, parsed.F5)
			require.Equal(t, expected.F6, parsed.F6)
			require.Equal(t, expected.F7, parsed.F7)
			require.Equal(t, expected.F8, parsed.F8)
			require.Zero(t, parsed.F9)

			w.Write([]byte("Pass"))
		}

		http.HandlerFunc(handlerFn).ServeHTTP(rr, req)
	})

	t.Run("invalid-values", func(t *testing.T) {
		queryParams := "f3=abc&f4=invalid&f5=invalid&f6=invalid&f7=invalid"

		req, err := http.NewRequest("GET", "http://example.com?"+queryParams, nil)
		require.Nil(t, err, fmt.Sprintf("Error creating request: %s", err))
		rr := httptest.NewRecorder()
		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			var parsed s1
			err := ParseURLQuery(r, &parsed)
			t.Log(err)
			require.NotNil(t, err, "Error expected")
			w.Write([]byte("Pass"))
		}

		http.HandlerFunc(handlerFn).ServeHTTP(rr, req)
	})
}

type s2 struct {
	F1 uint8  `qp:"f1"`
	F2 int8   `qp:"f2"`
	F3 uint16 `qp:"f3"`
	F4 int16  `qp:"f4"`
	F5 uint32 `qp:"f5"`
	F6 int32  `qp:"f6"`
	F7 uint64 `qp:"f7"`
	F8 int64  `qp:"f8"`
}

func TestParseS2(t *testing.T) {
	t.Run("filled", func(t *testing.T) {
		queryParams := "f1=1&f2=-1&f3=2&f4=-2&f5=3&f6=-3&f7=4&f8=-4"
		expected := s2{
			F1: 1,
			F2: -1,
			F3: 2,
			F4: -2,
			F5: 3,
			F6: -3,
			F7: 4,
			F8: -4,
		}

		req, err := http.NewRequest("GET", "http://example.com?"+queryParams, nil)
		require.Nil(t, err, fmt.Sprintf("Error creating request: %s", err))

		rr := httptest.NewRecorder()

		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			var parsed s2
			err := ParseURLQuery(r, &parsed)
			require.Nil(t, err, fmt.Sprintf("Error parsing query: %s", err))
			require.Equal(t, expected.F1, parsed.F1)
			require.Equal(t, expected.F2, parsed.F2)
			require.Equal(t, expected.F3, parsed.F3)
			require.Equal(t, expected.F4, parsed.F4)
			require.Equal(t, expected.F5, parsed.F5)
			require.Equal(t, expected.F6, parsed.F6)
			require.Equal(t, expected.F7, parsed.F7)
			require.Equal(t, expected.F8, parsed.F8)

			w.Write([]byte("Pass"))
		}

		http.HandlerFunc(handlerFn).ServeHTTP(rr, req)
	})

	t.Run("out-of-range", func(t *testing.T) {
		queryParams := "f2=128"
		req, err := http.NewRequest("GET", "http://example.com?"+queryParams, nil)
		require.Nil(t, err, fmt.Sprintf("Error creating request: %s", err))

		rr := httptest.NewRecorder()

		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			var parsed s2
			err := ParseURLQuery(r, &parsed)
			t.Log(err)
			require.NotNil(t, err, "Error expected")
			w.Write([]byte("Pass"))
		}
		http.HandlerFunc(handlerFn).ServeHTTP(rr, req)
	})
}
