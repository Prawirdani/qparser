package qparser

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommonTypes(t *testing.T) {
	type commonTypes struct {
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

	t.Run("success", func(t *testing.T) {
		queryParams := "f1=hello&f2=world&f3=-42&f4=69&f5=true&f6=false&f7=3.14&f8=42.0&f9=42"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		f2 := "world"
		f4 := 69
		f6 := false
		f8 := 42.0
		expected := commonTypes{
			F1: "hello",
			F2: &f2,
			F3: -42,
			F4: &f4,
			F5: true,
			F6: &f6,
			F7: 3.14,
			F8: &f8,
		}

		var parsed commonTypes
		err = parse(values, &parsed)
		require.Nil(t, err)
		require.Equal(t, expected.F1, parsed.F1)
		require.Equal(t, expected.F2, parsed.F2)
		require.Equal(t, expected.F3, parsed.F3)
		require.Equal(t, expected.F4, parsed.F4)
		require.Equal(t, expected.F5, parsed.F5)
		require.Equal(t, expected.F6, parsed.F6)
		require.Equal(t, expected.F7, parsed.F7)
		require.Equal(t, expected.F8, parsed.F8)
		require.Zero(t, parsed.F9)

	})

	t.Run("invalid-values", func(t *testing.T) {
		queryParams := "f3=abc&f4=invalid&f5=invalid&f6=invalid&f7=invalid"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		var parsed commonTypes
		err = parse(values, &parsed)
		require.NotNil(t, err, "Error expected")
	})
}

func TestParseFixedIntTypes(t *testing.T) {
	type fixedIntTypes struct {
		F1 uint8  `qp:"f1"`
		F2 int8   `qp:"f2"`
		F3 uint16 `qp:"f3"`
		F4 int16  `qp:"f4"`
		F5 uint32 `qp:"f5"`
		F6 int32  `qp:"f6"`
		F7 uint64 `qp:"f7"`
		F8 int64  `qp:"f8"`
	}
	t.Run("success", func(t *testing.T) {
		queryParams := "f1=1&f2=-1&f3=2&f4=-2&f5=3&f6=-3&f7=4&f8=-4"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		expected := fixedIntTypes{
			F1: 1,
			F2: -1,
			F3: 2,
			F4: -2,
			F5: 3,
			F6: -3,
			F7: 4,
			F8: -4,
		}

		var parsed fixedIntTypes
		err = parse(values, &parsed)
		require.Nil(t, err)
		require.Equal(t, expected.F1, parsed.F1)
		require.Equal(t, expected.F2, parsed.F2)
		require.Equal(t, expected.F3, parsed.F3)
		require.Equal(t, expected.F4, parsed.F4)
		require.Equal(t, expected.F5, parsed.F5)
		require.Equal(t, expected.F6, parsed.F6)
		require.Equal(t, expected.F7, parsed.F7)
		require.Equal(t, expected.F8, parsed.F8)

	})

	t.Run("out-of-range", func(t *testing.T) {
		queryParams := "f2=128"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		var parsed fixedIntTypes
		err = parse(values, &parsed)
		require.NotNil(t, err, "Error expected")
	})
}

func TestParseSliceTypes(t *testing.T) {
	type sliceTypes struct {
		F1 []string  `qp:"f1"`
		F2 []*int    `qp:"f2"`
		F3 []float64 `qp:"f3"`
		F4 []int8    `qp:"f4"`
		F5 []int     `qp:"f5"`
	}
	t.Run("success", func(t *testing.T) {
		queryParams := "f1=foo,bar,baz,qux&f2=1,2,3&f3=1.1,2.2,3.3&f4=1,2,3&f5=69,420"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		f2_1 := 1
		f2_2 := 2
		f2_3 := 3
		expected := sliceTypes{
			F1: []string{"foo", "bar", "baz", "qux"},
			F2: []*int{&f2_1, &f2_2, &f2_3},
			F3: []float64{1.1, 2.2, 3.3},
			F4: []int8{1, 2, 3},
			F5: []int{69, 420},
		}

		var parsed sliceTypes
		err = parse(values, &parsed)
		require.Nil(t, err)
		require.Equal(t, expected.F1, parsed.F1)
		require.Equal(t, expected.F2, parsed.F2)
		require.Equal(t, expected.F3, parsed.F3)
		require.Equal(t, expected.F4, parsed.F4)
		require.Equal(t, expected.F5, parsed.F5)
	})

	t.Run("invalid", func(t *testing.T) {
		queryParams := "f1=foo,bar,baz,qux&f4=1,2,255"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		var parsed sliceTypes
		err = parse(values, &parsed)
		require.NotNil(t, err, "Error expected")

	})
}

func TestParseNested(t *testing.T) {
	type Child1 struct {
		F1 string `qp:"f1"`
		F2 int    `qp:"f2"`
	}
	type Child2 struct {
		F3 string `qp:"f3"`
		F4 []int  `qp:"f4"`
	}
	type nestedSt struct {
		C1 Child1
		C2 *Child2
	}

	t.Run("success", func(t *testing.T) {
		queryParams := "f1=hello&f2=42&f3=world&f4=69,420"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		c1 := Child1{F1: "hello", F2: 42}
		c2 := &Child2{F3: "world", F4: []int{69, 420}}
		expected := nestedSt{C1: c1, C2: c2}

		var parsed nestedSt
		err = parse(values, &parsed)
		require.Nil(t, err)

		require.Equal(t, expected.C1.F1, parsed.C1.F1)
		require.Equal(t, expected.C1.F2, parsed.C1.F2)
		require.Equal(t, expected.C2.F3, parsed.C2.F3)
		require.Equal(t, expected.C2.F4, parsed.C2.F4)
	})
}

func TestParseRequest(t *testing.T) {
	type request struct {
		F1 string `qp:"f1"`
		F2 int    `qp:"f2"`
	}

	t.Run("success", func(t *testing.T) {
		queryParams := "f1=hello&f2=42"
		req, err := http.NewRequest(http.MethodGet, "http://example.com?"+queryParams, nil)
		require.Nil(t, err, fmt.Sprintf("Error creating request: %s", err))

		rr := httptest.NewRecorder()

		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			var parsed request
			err := ParseRequest(r, &parsed)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// require.Nil(t, err)
			require.Equal(t, "hello", parsed.F1)
			require.Equal(t, 42, parsed.F2)
			w.WriteHeader(http.StatusOK)
		}

		http.HandlerFunc(handlerFn).ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid", func(t *testing.T) {
		queryParams := "f1=hello&f2=world"
		req, err := http.NewRequest(http.MethodGet, "http://example.com?"+queryParams, nil)
		require.Nil(t, err, fmt.Sprintf("Error creating request: %s", err))

		rr := httptest.NewRecorder()

		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			var parsed request
			err := ParseRequest(r, &parsed)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		}

		http.HandlerFunc(handlerFn).ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

}

func TestParseURL(t *testing.T) {
	type request struct {
		F1 string `qp:"f1"`
		F2 int    `qp:"f2"`
	}
	t.Run("success", func(t *testing.T) {

		queryParams := "f1=hello&f2=42"
		address := "http://example.com?" + queryParams

		var parsed request
		err := ParseURL(address, &parsed)
		require.Nil(t, err)
		require.Equal(t, "hello", parsed.F1)
		require.Equal(t, 42, parsed.F2)
	})

	t.Run("invalid-url", func(t *testing.T) {
		address := "ht@tp://example.com?f1=hello&f2=42"
		var parsed request
		err := ParseURL(address, &parsed)
		require.NotNil(t, err)
	})
}
