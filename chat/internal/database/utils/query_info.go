package utils

import "encoding/json"

type QueryInfo struct {
	Sql  string        `json:"sql"`
	Args []interface{} `json:"args"`
}

func NewQueryInfo(sql string, args []interface{}) QueryInfo {
	return QueryInfo{
		Sql:  sql,
		Args: args,
	}
}

func (q *QueryInfo) String() string {
	b, _ := json.MarshalIndent(struct {
		SQL  string        `json:"sql"`
		Args []interface{} `json:"args"`
	}{
		SQL:  q.Sql,
		Args: q.Args,
	}, "", "  ")

	return string(b)
}
