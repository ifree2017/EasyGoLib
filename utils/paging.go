package utils

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
)

type PageForm struct {
	Start uint   `form:"start"`
	Limit uint   `form:"limit"`
	Q     string `form:"q"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

func (p PageForm) String() string {
	return fmt.Sprintf("PageForm[Start=%d, Limit=%d, Q=%s, Sort=%s, Order=%s]", p.Start, p.Limit, p.Q, p.Sort, p.Order)
}

func NewPageForm() *PageForm {
	return &PageForm{
		Start: 0,
		Limit: math.MaxUint32,
	}
}

type PageResult struct {
	Total uint          `json:"total"`
	Rows  []interface{} `json:"rows"`
}

func NewPageResult(rows []interface{}) *PageResult {
	return &PageResult{
		Total: uint(len(rows)),
		Rows:  rows,
	}
}

func (pr *PageResult) Sort(by, order string) {
	sort.Slice(pr.Rows, func(i, j int) (ret bool) {
		va := reflect.ValueOf(pr.Rows[i])
		vb := reflect.ValueOf(pr.Rows[j])
		for va.Kind() == reflect.Interface || va.Kind() == reflect.Ptr {
			va = va.Elem()
		}
		for vb.Kind() == reflect.Interface || vb.Kind() == reflect.Ptr {
			vb = vb.Elem()
		}
		if va.Kind() == reflect.Struct && vb.Kind() == reflect.Struct {
			ret = fmt.Sprintf("%v", va.FieldByName(by)) < fmt.Sprintf("%v", vb.FieldByName(by))
		}
		if va.Kind() == reflect.Map && vb.Kind() == reflect.Map {
			ret = fmt.Sprintf("%v", va.MapIndex(reflect.ValueOf(by))) < fmt.Sprintf("%v", vb.MapIndex(reflect.ValueOf(by)))
		}
		if strings.HasPrefix(strings.ToLower(order), "desc") {
			ret = !ret
		}
		return
	})
}
