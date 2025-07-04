package qparser

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name        string
		dst         any
		expectError error
	}{
		{"ValidStructPointer", &struct{}{}, nil},
		{"NilPointer", nil, ErrNotPtr},
		{"NonPointer", struct{}{}, ErrNotPtr},
		{"PointerToNonStruct", new(int), ErrNotPtr},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := parse(url.Values{}, tc.dst)
			require.Equal(t, tc.expectError, err)
		})
	}
}

// While creating individual tests for each type may seem redundant or excessive,
// But it's good to have them as they provide a more granular view of the parser's behavior.
// These tests help ensure that the parser handles all expected input variations correctly,
// providing confidence in its robustness and reliability.

type (
	strAlias string

	boolAlias bool

	intAlias   int
	int64Alias int64
	int32Alias int32
	int16Alias int16
	int8Alias  int8

	uintAlias   uint
	uint64Alias uint64
	uint32Alias uint32
	uint16Alias uint16
	uint8Alias  uint8

	float32Alias float32
	float64Alias float64
)

func TestStrings(t *testing.T) {
	type strs struct {
		F1 string    `qp:"f1"`
		F2 *string   `qp:"f2"`
		F3 strAlias  `qp:"f3"`
		F4 *strAlias `qp:"f4"`
		F5 string    `qp:"-"`
		F6 *string   `qp:"-"`
		F7 string
		F8 *string
		F9 *string `qp:"f9"`
	}

	queryParams := "f1=hello&f2=world&f3=foo&f4=bar&f5=ignored&f6=ignored&f7=hello&f8=world"

	var s strs
	expected := strs{
		F1: "hello",
		F2: ptr("world"),
		F3: strAlias("foo"),
		F4: ptr(strAlias("bar")),
		F5: "",
		F6: nil,
		F7: "",
		F8: nil,
		F9: nil, // Not provided in query params
	}

	values, err := url.ParseQuery(queryParams)
	require.Nil(t, err)

	err = parse(values, &s)
	require.Nil(t, err)
	require.Equal(t, expected, s)
}

func TestBooleans(t *testing.T) {
	type bools struct {
		F1 bool       `qp:"f1"`
		F2 *bool      `qp:"f2"`
		F3 boolAlias  `qp:"f3"`
		F4 *boolAlias `qp:"f4"`
		F5 bool       `qp:"-"`
		F6 *bool      `qp:"-"`
		F7 bool
		F8 *bool
		F9 *bool `qp:"f9"`
	}

	queryParams := "f1=true&f2=false&f3=true&f4=false&f5=ignored&f6=ignored&f7=ignored&f8=ignored"
	var b bools
	expected := bools{
		F1: true,
		F2: ptr(false),
		F3: boolAlias(true),
		F4: ptr(boolAlias(false)),
		F5: false,
		F6: nil,
		F7: false,
		F8: nil,
		F9: nil, // Not provided in query params
	}

	values, err := url.ParseQuery(queryParams)
	require.Nil(t, err)

	err = parse(values, &b)
	require.Nil(t, err)
	require.Equal(t, expected, b)

	t.Run("Invalid", func(t *testing.T) {
		queryParams := "f1=invalid-boolean"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		err = parse(values, &b)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid value")
	})
}

func TestSignedIntegers(t *testing.T) {
	type ints struct {
		F1  int         `qp:"f1"`
		F2  *int        `qp:"f2"`
		F3  intAlias    `qp:"f3"`
		F4  *intAlias   `qp:"f4"`
		F5  int64       `qp:"f5"`
		F6  *int64      `qp:"f6"`
		F7  int64Alias  `qp:"f7"`
		F8  *int64Alias `qp:"f8"`
		F9  int32       `qp:"f9"`
		F10 *int32      `qp:"f10"`
		F11 int32Alias  `qp:"f11"`
		F12 *int32Alias `qp:"f12"`
		F13 int16       `qp:"f13"`
		F14 *int16      `qp:"f14"`
		F15 int16Alias  `qp:"f15"`
		F16 *int16Alias `qp:"f16"`
		F17 int8        `qp:"f17"`
		F18 *int8       `qp:"f18"`
		F19 int8Alias   `qp:"f19"`
		F20 *int8Alias  `qp:"f20"`
	}

	queryParams := "f1=42&f2=69&f3=42&f4=69&f5=42&f6=69&f7=42&f8=69&f9=42&f10=69&f11=42&f12=69&f13=42&f14=69&f15=42&f16=69&f17=42&f18=69&f19=42&f20=69"
	var i ints
	expected := ints{
		F1:  42,
		F2:  ptr(69),
		F3:  intAlias(42),
		F4:  ptr(intAlias(69)),
		F5:  42,
		F6:  ptr(int64(69)),
		F7:  int64Alias(42),
		F8:  ptr(int64Alias(69)),
		F9:  42,
		F10: ptr(int32(69)),
		F11: int32Alias(42),
		F12: ptr(int32Alias(69)),
		F13: 42,
		F14: ptr(int16(69)),
		F15: int16Alias(42),
		F16: ptr(int16Alias(69)),
		F17: 42,
		F18: ptr(int8(69)),
		F19: int8Alias(42),
		F20: ptr(int8Alias(69)),
	}

	values, err := url.ParseQuery(queryParams)
	require.Nil(t, err)

	err = parse(values, &i)
	assert.Nil(t, err)
	assert.Equal(t, expected, i)

	t.Run("Invalid", func(t *testing.T) {
		queryParams := "f1=non-integer"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		err = parse(values, &i)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid value")
	})

	t.Run("OutOfRange", func(t *testing.T) {
		t.Run("int8", func(t *testing.T) {
			// int8 range is -128 to 127
			queryParams := "f20=-129"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &i)
			assert.NotNil(t, err)
		})

		t.Run("int16", func(t *testing.T) {
			// int16 range is -32768 to 32767
			queryParams := "f16=-32769"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &i)
			assert.NotNil(t, err)
		})

		t.Run("int32", func(t *testing.T) {
			// int32 range is -2147483648 to 2147483647
			queryParams := "f12=-2147483649"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &i)
			assert.NotNil(t, err)
		})

		t.Run("int64", func(t *testing.T) {
			// int64 range is -9223372036854775808 to 9223372036854775807
			queryParams := "f8=-9223372036854775809"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &i)
			assert.NotNil(t, err)
		})
	})
}

func TestUnsignedIntegers(t *testing.T) {
	type uints struct {
		F1  uint         `qp:"f1"`
		F2  *uint        `qp:"f2"`
		F3  uintAlias    `qp:"f3"`
		F4  *uintAlias   `qp:"f4"`
		F5  uint64       `qp:"f5"`
		F6  *uint64      `qp:"f6"`
		F7  uint64Alias  `qp:"f7"`
		F8  *uint64Alias `qp:"f8"`
		F9  uint32       `qp:"f9"`
		F10 *uint32      `qp:"f10"`
		F11 uint32Alias  `qp:"f11"`
		F12 *uint32Alias `qp:"f12"`
		F13 uint16       `qp:"f13"`
		F14 *uint16      `qp:"f14"`
		F15 uint16Alias  `qp:"f15"`
		F16 *uint16Alias `qp:"f16"`
		F17 uint8        `qp:"f17"`
		F18 *uint8       `qp:"f18"`
		F19 uint8Alias   `qp:"f19"`
		F20 *uint8Alias  `qp:"f20"`
	}

	queryParams := "f1=42&f2=69&f3=42&f4=69&f5=42&f6=69&f7=42&f8=69&f9=42&f10=69&f11=42&f12=69&f13=42&f14=69&f15=42&f16=69&f17=42&f18=69&f19=42&f20=69"
	var u uints
	expected := uints{
		F1:  42,
		F2:  ptr(uint(69)),
		F3:  uintAlias(42),
		F4:  ptr(uintAlias(69)),
		F5:  42,
		F6:  ptr(uint64(69)),
		F7:  uint64Alias(42),
		F8:  ptr(uint64Alias(69)),
		F9:  42,
		F10: ptr(uint32(69)),
		F11: uint32Alias(42),
		F12: ptr(uint32Alias(69)),
		F13: 42,
		F14: ptr(uint16(69)),
		F15: uint16Alias(42),
		F16: ptr(uint16Alias(69)),
		F17: 42,
		F18: ptr(uint8(69)),
		F19: uint8Alias(42),
		F20: ptr(uint8Alias(69)),
	}

	values, err := url.ParseQuery(queryParams)
	require.Nil(t, err)

	err = parse(values, &u)
	assert.Nil(t, err)
	assert.Equal(t, expected, u)

	t.Run("Invalid", func(t *testing.T) {
		queryParams := "f1=non-integer"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		err = parse(values, &u)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid value")
	})

	t.Run("OutOfRange", func(t *testing.T) {
		t.Run("uint8", func(t *testing.T) {
			// uint8 range is 0 to 255
			queryParams := "f20=256"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &u)
			assert.NotNil(t, err)
		})

		t.Run("uint16", func(t *testing.T) {
			// uint16 range is 0 to 65535
			queryParams := "f16=65536"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &u)
			assert.NotNil(t, err)
		})

		t.Run("uint32", func(t *testing.T) {
			// uint32 range is 0 to 4294967295
			queryParams := "f12=4294967296"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &u)
			assert.NotNil(t, err)
		})

		t.Run("uint64", func(t *testing.T) {
			// uint64 range is 0 to 18446744073709551615
			queryParams := "f8=18446744073709551616"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &u)
			assert.NotNil(t, err)
		})
	})

	t.Run("Negative", func(t *testing.T) {
		t.Run("uint8", func(t *testing.T) {
			// uint8 range is 0 to 255
			queryParams := "f20=-1"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &u)
			assert.NotNil(t, err)
		})

		t.Run("uint16", func(t *testing.T) {
			// uint16 range is 0 to 65535
			queryParams := "f16=-1"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &u)
			assert.NotNil(t, err)
		})

		t.Run("uint32", func(t *testing.T) {
			// uint32 range is 0 to 4294967295
			queryParams := "f12=-1"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &u)
			assert.NotNil(t, err)
		})

		t.Run("uint64", func(t *testing.T) {
			// uint64 range is 0 to 18446744073709551615
			queryParams := "f8=-1"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &u)
			assert.NotNil(t, err)
		})
	})
}

func TestFloats(t *testing.T) {
	type floats struct {
		F1 float32       `qp:"f1"`
		F2 *float32      `qp:"f2"`
		F3 float32Alias  `qp:"f3"`
		F4 *float32Alias `qp:"f4"`
		F5 float64       `qp:"f5"`
		F6 *float64      `qp:"f6"`
		F7 float64Alias  `qp:"f7"`
		F8 *float64Alias `qp:"f8"`
	}

	queryParams := "f1=3.14&f2=2.718&f3=3.14&f4=2.718&f5=3.14&f6=2.718&f7=3.14&f8=2.718"
	var f floats
	expected := floats{
		F1: 3.14,
		F2: ptr(float32(2.718)),
		F3: float32Alias(3.14),
		F4: ptr(float32Alias(2.718)),
		F5: 3.14,
		F6: ptr(float64(2.718)),
		F7: float64Alias(3.14),
		F8: ptr(float64Alias(2.718)),
	}

	values, err := url.ParseQuery(queryParams)
	require.Nil(t, err)

	err = parse(values, &f)
	assert.Nil(t, err)
	assert.Equal(t, expected, f)

	t.Run("Invalid", func(t *testing.T) {
		queryParams := "f1=non-float"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		err = parse(values, &f)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid value")
	})

	t.Run("OutOfRange", func(t *testing.T) {
		t.Run("float32", func(t *testing.T) {
			queryParams := "f4=3.40282346638528859811704183484516925440e+38"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &f)
			assert.NotNil(t, err)
		})

		t.Run("float64", func(t *testing.T) {
			queryParams := "f8=1.797693134862315708145274237317043567981e+308"
			values, err := url.ParseQuery(queryParams)
			require.Nil(t, err)

			err = parse(values, &f)
			assert.NotNil(t, err)
		})
	})
}

func TestTime(t *testing.T) {
	type times struct {
		F1 time.Time  `qp:"f1"`
		F2 *time.Time `qp:"f2"`
		F3 time.Time  `qp:"-"`
	}

	currDate := time.Now().Truncate(time.Nanosecond)

	dateStr := currDate.Format(time.RFC3339Nano)
	queryParams := fmt.Sprintf(
		"f1=%s&f2=%s&f3=ignored",
		dateStr,
		dateStr,
	)
	expected := times{
		F1: currDate,
		F2: ptr(currDate),
		F3: time.Time{},
	}

	values, err := url.ParseQuery(queryParams)
	require.Nil(t, err)

	var result times
	err = parse(values, &result)

	require.Nil(t, err)
	require.True(t, expected.F1.Equal(result.F1))
	require.True(t, expected.F2.Equal(*result.F2))
	require.True(t, expected.F3.Equal(result.F3))

	t.Run("Invalid", func(t *testing.T) {
		queryParams := "f1=invalid-date"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		err = parse(values, &result)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid value")
	})
}

func TestSlices(t *testing.T) {
	type slices struct {
		F1  []string     `qp:"f1"`
		F2  []*string    `qp:"f2"`
		F3  *[]string    `qp:"f3"`
		F4  []int        `qp:"f4"`
		F5  []int64      `qp:"f5"`
		F6  []int32Alias `qp:"f6"`
		F7  []uint16     `qp:"f7"`
		F8  []uint8Alias `qp:"f8"`
		F9  []float32    `qp:"f9"`
		F10 []float64    `qp:"f10"`
		F11 []bool       `qp:"f11"`
		F12 []boolAlias  `qp:"f12"`
		F13 []string     `qp:"f13"`
		F14 *[]string    `qp:"f14"`
		F15 []string     `qp:"-"`
		F16 []string
		F17 []string `qp:"f17"`
	}

	queryParams := "f1=foo,bar,baz,qux&f2=1,2,3&f3=foo,bar,baz,qux&f4=1,2,3&f5=1,2,3&f6=1,2,3&f7=1,2,3&f8=1,2,3&f9=1.1,2.2,3.3&f10=1.1,2.2,3.3&f11=true,false,true&f12=true,false,true&f17=foo.bar.baz.qux"
	var s slices
	expected := slices{
		F1:  []string{"foo", "bar", "baz", "qux"},
		F2:  []*string{ptr("1"), ptr("2"), ptr("3")},
		F3:  &[]string{"foo", "bar", "baz", "qux"},
		F4:  []int{1, 2, 3},
		F5:  []int64{1, 2, 3},
		F6:  []int32Alias{1, 2, 3},
		F7:  []uint16{1, 2, 3},
		F8:  []uint8Alias{1, 2, 3},
		F9:  []float32{1.1, 2.2, 3.3},
		F10: []float64{1.1, 2.2, 3.3},
		F11: []bool{true, false, true},
		F12: []boolAlias{true, false, true},
		F13: nil, // Not provided in query params
		F14: nil,
		F15: nil,
		F16: nil,
		F17: []string{"foo.bar.baz.qux"},
	}

	values, err := url.ParseQuery(queryParams)
	require.Nil(t, err)

	err = parse(values, &s)
	assert.Nil(t, err)
	assert.Equal(t, expected, s)

	t.Run("Repeated-Keys", func(t *testing.T) {
		queryParams := "f1=foo&f1=bar&f1=baz&f1=qux"
		var s slices
		expected := slices{
			F1: []string{"foo", "bar", "baz", "qux"},
		}

		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		err = parse(values, &s)
		assert.Nil(t, err)
		assert.Nil(t, s.F2)
		assert.Equal(t, expected, s)
	})

	// Mixed multiple values and single values comma-separated
	t.Run("Mixed-Keys", func(t *testing.T) {
		queryParams := "f1=foo,bar&f1=baz&f1=qux&f1="
		var s slices
		expected := slices{
			F1: []string{"foo", "bar", "baz", "qux"},
		}

		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		err = parse(values, &s)
		assert.Nil(t, err)
		assert.Equal(t, expected, s)
	})

	t.Run("Invalid", func(t *testing.T) {
		queryParams := "f4=1,2,invalid"
		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		err = parse(values, &s)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid value")
	})
}

func TestUnsupportedKind(t *testing.T) {
	type unsupported struct {
		F1 complex64 `qp:"f1"`
	}

	queryParams := "f1=1+2i"
	var u unsupported
	values, err := url.ParseQuery(queryParams)
	require.Nil(t, err)

	err = parse(values, &u)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported kind")
}

func TestUnexportedFields(t *testing.T) {
	queryParams := "f1=foo&f2=bar"

	type unexportedFields struct {
		F1 string `qp:"f1"`
		f2 string `qp:"f2"`
	}
	var st unexportedFields

	values, err := url.ParseQuery(queryParams)
	require.Nil(t, err)

	err = parse(values, &st)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "be sure it is exported")

	assert.Equal(t, "", st.f2) // No need, only to satisfy the linter
}

func TestNestedStruct(t *testing.T) {
	type childchild struct {
		F3 string `qp:"f3"`
	}

	type child1 struct {
		F1 string `qp:"f1"`
		F2 int    `qp:"f2"`
		CC childchild
	}

	type child2 struct {
		F4 []int   `qp:"f4"`
		F5 *string `qp:"f5"`
	}

	type ns struct {
		C1 child1
		C2 *child2
		F6 string `qp:"f6"`
	}

	queryParams := "f1=foo&f2=42&f3=bar&f4=69,420&f5=baz&f6=qux"
	var n ns
	expected := ns{
		C1: child1{
			F1: "foo",
			F2: 42,
			CC: childchild{F3: "bar"},
		},
		C2: &child2{
			F4: []int{69, 420},
			F5: ptr("baz"),
		},
		F6: "qux",
	}

	values, err := url.ParseQuery(queryParams)
	require.Nil(t, err)

	err = parse(values, &n)
	assert.Nil(t, err)
	assert.Equal(t, expected, n)

	t.Run("Empty-Nested-Pointer", func(t *testing.T) {
		type child struct {
			F1 *string `qp:"f1"`
		}
		type parent struct {
			C *child
		}

		queryParams := ""
		var s parent
		expected := parent{
			C: nil,
		}

		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		err = parse(values, &s)
		assert.Nil(t, err)
		assert.Equal(t, expected, s)
	})

	t.Run("Invalid-Nested-Pointer", func(t *testing.T) {
		type child struct {
			F1 complex64 `qp:"f1"`
		}

		type parent struct {
			C *child
		}

		queryParams := "f1=1+2i"
		var s parent

		values, err := url.ParseQuery(queryParams)
		require.Nil(t, err)

		err = parse(values, &s)
		assert.NotNil(t, err)
	})
}

type Pagination struct {
	Page  int `qp:"page"`
	Limit int `qp:"limit"`
}

type Filter struct {
	Categories []string `qp:"categories"`
}

type SearchParams struct {
	Pagination Pagination
	Filters    Filter
	Q          string    `qp:"q"`
	Date       time.Time `qp:"date"`
}

func TestParseRequest(t *testing.T) {
	currDate := time.Now().Truncate(time.Nanosecond)

	queryParams := "page=1&limit=10&q=lorem&categories=foo,bar,baz&date=" + currDate.Format(
		time.RFC3339Nano,
	)
	url := "http://example.com?" + queryParams

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.Nil(t, err, fmt.Sprintf("Error creating request: %s", err))

	rr := httptest.NewRecorder()

	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		var sp SearchParams
		err := ParseRequest(r, &sp)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		require.Nil(t, err)
		assert.Equal(t, 1, sp.Pagination.Page)
		assert.Equal(t, 10, sp.Pagination.Limit)
		assert.Equal(t, "lorem", sp.Q)
		assert.Equal(t, []string{"foo", "bar", "baz"}, sp.Filters.Categories)
		assert.True(t, currDate.Equal(sp.Date))
		w.WriteHeader(http.StatusOK)
	}

	http.HandlerFunc(handlerFn).ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	t.Run("Invalid", func(t *testing.T) {
		queryParams := "page=1&limit=invalid-limit&q=lorem"
		url := "http://example.com?" + queryParams
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.Nil(t, err, fmt.Sprintf("Error creating request: %s", err))

		rr := httptest.NewRecorder()

		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			var parsed SearchParams
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
	url := "http://example.com?page=1&limit=10&q=lorem&categories=foo,bar,baz"
	var sp SearchParams
	expected := SearchParams{
		Pagination: Pagination{Page: 1, Limit: 10},
		Filters:    Filter{Categories: []string{"foo", "bar", "baz"}},
		Q:          "lorem",
	}
	err := ParseURL(url, &sp)
	require.Nil(t, err)

	assert.Equal(t, expected, sp)

	t.Run("Invalid-URL", func(t *testing.T) {
		url := "ht@tp://example.com?page=1&limit=10&q=lorem&categories=foo,bar,baz"
		var sp SearchParams
		err := ParseURL(url, &sp)
		require.NotNil(t, err)
	})
}

func ptr[T any](v T) *T {
	return &v
}
