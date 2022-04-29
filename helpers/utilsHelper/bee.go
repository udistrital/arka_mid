package utilsHelper

import (
	"fmt"
	"net/url"
	"strings"
)

func EncodeUrl(query string, fields string, sortby string, order string, offset string, limit string) string {
	params := url.Values{}

	if len(query) > 0 {
		params.Add("query", query)
	}

	if len(fields) > 0 {
		params.Add("fields", fields)
	}

	if len(sortby) > 0 {
		params.Add("sortby", sortby)

	}

	if len(order) > 0 {
		params.Add("order", order)

	}

	if len(offset) > 0 {
		params.Add("offset", offset)

	}

	if len(limit) > 0 {
		params.Add("limit", limit)

	}

	return params.Encode()
}

type Sorting struct {
	By  string
	Asc bool
}

func (c *Sorting) OrderStr() string {
	if c.Asc {
		return "asc"
	}
	return "desc"
}

type Query struct {
	Query  map[string]string
	Limit  int `default:"10"`
	Offset int `default:"0"`
	Fields []string
	Sort   []Sorting
}

func (opts *Query) Encode() string {

	sortby := make([]string, 0)
	order := make([]string, 0)
	for _, v := range opts.Sort {
		sortby = append(sortby, v.By)
		order = append(order, v.OrderStr())
	}

	joinedQuery := ""
	periods := len(opts.Query)
	for k, v := range opts.Query {
		joinedQuery += k + ":" + v
		if periods > 1 {
			joinedQuery += ","
			periods--
		}
	}

	return EncodeUrl(joinedQuery,
		strings.Join(opts.Fields, ","),
		strings.Join(sortby, ","),
		strings.Join(order, ","),
		fmt.Sprint(opts.Offset),
		fmt.Sprint(opts.Limit),
	)
}
