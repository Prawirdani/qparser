package qparser

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	Search     string    `qp:"q"`
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
		assert.Equal(t, "lorem", sp.Search)
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
		Search:     "lorem",
	}
	err := ParseURL(url, &sp)
	require.Nil(t, err)

	assert.Equal(t, expected, sp)

	t.Run("Invalid-URL", func(t *testing.T) {
		url := "ht@tp://example.com?page=1&limit=10&q=lorem&categories=foo,bar,baz"
		var sp SearchParams
		err := ParseURL(url, &sp)
		require.Error(t, err)
	})
}
