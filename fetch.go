package speedrungql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type OrderDirection string

const (
	Ascending  OrderDirection = "ASC"
	Descending OrderDirection = "DESC"
)

type FetchOption func(*request)

type request struct {
	filters []requestFilter
	order   *requestOrder
	paging  requestPaging
}

type requestFilter struct {
	field string
	value interface{}
}

type requestOrder struct {
	field     *string
	direction *OrderDirection
}

type requestPaging struct {
	max    *int
	offset *int
}

func WithFilter(field string, value interface{}) FetchOption {
	return func(r *request) {
		r.filters = append(r.filters, requestFilter{
			field: field,
			value: value,
		})
	}
}

func WithOrder(field *string, direction *OrderDirection) FetchOption {
	return func(r *request) {
		r.order = &requestOrder{
			field:     field,
			direction: direction,
		}
	}
}

func WithLimit(limit int) FetchOption {
	return func(r *request) {
		r.paging.max = &limit
	}
}

func WithOffset(offset int) FetchOption {
	return func(r *request) {
		r.paging.offset = &offset
	}
}

func (c *Client) fetch(ctx context.Context, path string, result interface{}, opts ...FetchOption) error {
	var r request
	for _, opt := range opts {
		opt(&r)
	}

	u := c.BaseURL + path
	values := url.Values{}

	for _, filter := range r.filters {
		values.Set(filter.field, fmt.Sprintf("%s", filter.value))
	}

	if r.order != nil {
		if r.order.field != nil {
			values.Set("orderby", strings.ToLower(*r.order.field))
		}
		if r.order.direction != nil {
			values.Set("direction", strings.ToLower(string(*r.order.direction)))
		}
	}

	if r.paging.max != nil {
		values.Set("max", fmt.Sprintf("%d", *r.paging.max))
	}
	if r.paging.offset != nil {
		values.Set("offset", fmt.Sprintf("%d", *r.paging.offset))
	}

	if len(values) > 0 {
		u += "?" + values.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}

	return nil
}
