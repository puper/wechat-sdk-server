package helpers

import (
	"reflect"
	"sort"
)

type SortBy struct {
	rv       reflect.Value
	lessFunc func(i, j int) bool
}

func NewSortBy(data interface{}, lessFunc func(i, j int) bool) *SortBy {
	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Ptr && rv.Kind() != reflect.Interface {
		panic("sort type error")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		panic("sort type error")
	}
	return &SortBy{
		rv:       rv,
		lessFunc: lessFunc,
	}
}

func Sort(data interface{}, lessFunc func(i, j int) bool) {
	by := NewSortBy(data, lessFunc)
	sort.Sort(by)
}

func (this *SortBy) Len() int {
	return this.rv.Len()
}

func (this *SortBy) Less(i, j int) bool {
	return this.lessFunc(i, j)
}

func (this *SortBy) Swap(i, j int) {
	tmp := this.rv.Index(i).Interface()
	this.rv.Index(i).Set(this.rv.Index(j))
	this.rv.Index(j).Set(reflect.ValueOf(tmp))
}
