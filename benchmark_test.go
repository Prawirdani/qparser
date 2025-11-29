package qparser

import (
	"maps"
	"net/url"
	"strconv"
	"testing"
	"time"
)

type QFilter struct {
	Page       int       `qp:"page"`
	Limit      int       `qp:"limit"`
	SortBy     string    `qp:"sort_by"`
	Order      string    `qp:"order"`
	Q          string    `qp:"q"`
	Active     bool      `qp:"active"`
	Threshold  float32   `qp:"threshold"`
	Foo        string    `qp:"foo"`
	Bar        string    `qp:"bar"`
	Buz        string    `qp:"buz"`
	From       time.Time `qp:"from"`
	Until      time.Time `qp:"until"`
	IDs        []int     `qp:"ids"`
	Categories []string  `qp:"categories"`
}

var base = url.Values{
	"page":      {"1"},
	"limit":     {"10"},
	"sort_by":   {"created_at"},
	"order":     {"ASC"},
	"q":         {"lorem ipsum dolor"},
	"active":    {"true"},
	"threshold": {"6.66"},
	"foo":       {"bar"},
	"bar":       {"buz"},
	"buz":       {"quux"},
}

var tests = []struct {
	name   string
	values func() url.Values
}{
	{
		// no Date, no slices
		name:   "Minimal",
		values: func() url.Values { return maps.Clone(base) },
	},
	{
		name: "1-date",
		values: func() url.Values {
			now := time.Now()
			from := now.Format(time.RFC3339)

			v := maps.Clone(base)
			v["from"] = []string{from}
			return v
		},
	},
	{
		name: "2-dates",
		values: func() url.Values {
			now := time.Now()
			from := now.Format(time.RFC3339)
			until := now.Add(7 * 24 * time.Hour).Format(time.RFC3339)

			v := maps.Clone(base)
			v["from"] = []string{from}
			v["until"] = []string{until}
			return v
		},
	},
	{
		name: "slices-string-1*50",
		values: func() url.Values {
			v := maps.Clone(base)
			v["categories"] = makeCategorySlice(50)["categories"]
			return v
		},
	},
	{
		name: "slices-int-1*50",
		values: func() url.Values {
			v := maps.Clone(base)
			v["ids"] = makeIDSlice(50)["ids"]
			return v
		},
	},
	{
		name: "slices-2*50",
		values: func() url.Values {
			v := maps.Clone(base)
			v["ids"] = makeIDSlice(50)["ids"]
			v["categories"] = makeCategorySlice(50)["categories"]
			return v
		},
	},
	{
		name: "slices-2*100",
		values: func() url.Values {
			v := maps.Clone(base)
			v["ids"] = makeIDSlice(100)["ids"]
			v["categories"] = makeCategorySlice(100)["categories"]
			return v
		},
	},
	{
		name: "2*25-slices-and-2-dates",
		values: func() url.Values {
			now := time.Now()
			from := now.Format(time.RFC3339)
			until := now.Add(7 * 24 * time.Hour).Format(time.RFC3339)

			v := maps.Clone(base)
			v["from"] = []string{from}
			v["until"] = []string{until}
			v["ids"] = makeIDSlice(25)["ids"]
			v["categories"] = makeCategorySlice(25)["categories"]
			return v
		},
	},
}

func Benchmark(b *testing.B) {
	for _, tt := range tests {
		b.Run("Seq/"+tt.name, func(b *testing.B) {
			b.ReportAllocs()
			vals := tt.values()
			b.ResetTimer()
			for b.Loop() {
				var f QFilter
				if err := Parse(vals, &f); err != nil {
					b.Fatal(err)
				}
			}
		})
	}

	// Benchmark Parallel measures cache behavior under concurrent load similar to HTTP handlers.
	for _, tt := range tests {
		b.Run("Par/"+tt.name, func(b *testing.B) {
			b.ReportAllocs()
			vals := tt.values()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					var f QFilter
					if err := Parse(vals, &f); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func makeIDSlice(n int) url.Values {
	v := url.Values{}
	for i := range n {
		v.Add("ids", strconv.Itoa(i+1))
	}

	return v
}

func makeCategorySlice(n int) url.Values {
	v := url.Values{}
	for i := range n {
		v.Add("categories", strconv.Itoa(i+1))
	}

	return v
}
