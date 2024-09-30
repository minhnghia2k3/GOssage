package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset int    `json:"offset" validate:"omitempty,min=0"`
	Sort   string `json:"sort" validate:"omitempty,oneof=asc desc"`
}

func (q *PaginatedFeedQuery) Parse(r *http.Request) error {
	var err error
	qr := r.URL.Query()

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
