package service

import (
	"fmt"
	"strings"
)

type ListCondition struct {
	Field string
	Op    string
	Value interface{}
}

type ListConditions []*ListCondition

func (listConds ListConditions) WhereExpr() (query interface{}, args []interface{}) {
	var querys []string
	var values []interface{}
	for _, cond := range listConds {
		querys = append(querys, fmt.Sprintf("%s %s ?", cond.Field, cond.Op))
		values = append(values, cond.Value)
	}
	return strings.Join(querys, " AND "), values
}
