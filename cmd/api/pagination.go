package main

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = o
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	return fq, nil
}
