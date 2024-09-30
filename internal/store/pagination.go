package store

import (
	"net/http"
	"strconv"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset int    `json:"offset" validate:"omitempty,min=0"`
	Sort   string `json:"sort" validate:"omitempty,oneof=asc desc"`
	Search string `json:"search" validate:"omitempty,lte=100"`
	Since  string `json:"since"`
	Until  string `json:"until"`
}

func (q *PaginatedFeedQuery) Parse(r *http.Request) error {
	var err error
	qr := r.URL.Query()

	search := qr.Get("search")
	if search != "" {
		q.Search = search
	}

	since := qr.Get("since")
	if since != "" {
		q.Since = parseTime(since)
	}

	until := qr.Get("until")
	if until != "" {
		q.Until = parseTime(until)
	}

	limit := qr.Get("limit")
	if limit != "" {
		q.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return err
		}
	}

	offset := qr.Get("offset")
	if offset != "" {
		q.Offset, err = strconv.Atoi(offset)
		if err != nil {
			return err
		}
	}

	sort := qr.Get("sort")
	if sort != "" {
		q.Sort = sort
	}
	return nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
