package qparser

import (
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

type SmallFilter struct {
	IDs []int `qp:"ids"`
}

type MediumFilter struct {
	IDs    []int  `qp:"ids"`
	Size   int    `qp:"size"`
	Name   string `qp:"name"`
	Active bool   `qp:"active"`
}

type LargeFilter struct {
	IDs       []int     `qp:"ids"`
	Size      int       `qp:"size"`
	Name      string    `qp:"name"`
	Active    bool      `qp:"active"`
	Category  string    `qp:"category"`
	Tags      []string  `qp:"tags"`
	Threshold float64   `qp:"threshold"`
	Start     time.Time `qp:"start"`
	End       time.Time `qp:"end"`
}

func Benchmark(b *testing.B) {
	tests := []struct {
		name    string
		makeVal func() url.Values
		fn      func(url.Values) error
	}{
		{
			name: "Small",
			makeVal: func() url.Values {
				return makeValues(50, nil)
			},
			fn: func(v url.Values) error {
				var f SmallFilter
				return parse(v, &f)
			},
		},
		{
			name: "Medium",
			makeVal: func() url.Values {
				fields := map[string]string{
					"size":   "42",
					"name":   "test",
					"active": "true",
				}
				return makeValues(50, fields)
			},
			fn: func(v url.Values) error {
				var f MediumFilter
				return parse(v, &f)
			},
		},
		{
			name: "Large",
			makeVal: func() url.Values {
				fields := map[string]string{
					"size":      "99",
					"name":      "rich",
					"active":    "true",
					"category":  "food",
					"tags":      "tag1,tag2,tag3",
					"threshold": "3.14",
				}
				return makeValues(50, fields)
			},
			fn: func(v url.Values) error {
				var f LargeFilter
				return parse(v, &f)
			},
		},
		{
			name: "LargeWithDate",
			makeVal: func() url.Values {
				fields := map[string]string{
					"size":      "99",
					"name":      "rich",
					"active":    "true",
					"category":  "food",
					"tags":      "tag1,tag2,tag3",
					"threshold": "3.14",
					"start":     "2025-07-04T17:12:32.123+07:00",
					"end":       "2025-10-04T17:12:32.123+07:00",
				}
				return makeValues(50, fields)
			},
			fn: func(v url.Values) error {
				var f LargeFilter
				return parse(v, &f)
			},
		},
	}

	for _, tt := range tests {
		vals := tt.makeVal()
		b.ResetTimer()
		b.Run(tt.name, func(b *testing.B) {
			for b.Loop() {
				if err := tt.fn(vals); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func makeValues(nIDs int, extra map[string]string) url.Values {
	v := url.Values{}

	ids := make([]string, nIDs)
	for i := range ids {
		ids[i] = strconv.Itoa(i + 1)
	}
	v.Set("ids", strings.Join(ids, ","))

	for k, val := range extra {
		v.Set(k, val)
	}

	return v
}
