package abstraction

import (
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

type Filter struct {
	Pagination
	Search string        `query:"search"`
	SortBy []string      `example:"asc_column,-dsc_column"`
	Query  []FilterQuery `json:"-"` // skip swagger
}

type FilterQuery struct {
	Field    string
	Value    string
	Operator string
}

type FilterBuilder[T any] struct {
	ectx    echo.Context
	name    string
	entity  *T
	Payload *Filter
}

func NewFilterBuiler[T any](ectx echo.Context, name string) FilterBuilder[T] {
	return FilterBuilder[T]{
		ectx:    ectx,
		name:    name,
		entity:  new(T),
		Payload: &Filter{},
	}
}

func (a *FilterBuilder[T]) Bind() {
	req := a.ectx.Request()
	queries := req.URL.Query()
	structType := reflect.TypeOf(*a.entity)

	// filter
	ignores := map[string]bool{
		"page":   true,
		"limit":  true,
		"search": true,
	}
	for field, values := range queries {
		if ignores[field] {
			continue
		}

		for _, value := range values {
			if field == "sort_by" {
				a.ConstructSort(structType, value)
			} else {
				a.ConstructFilter(structType, field, value)
			}
		}
	}
}

func (a *FilterBuilder[T]) ConstructFilter(structType reflect.Type, field string, value string) {
	for i := 0; i < structType.NumField(); i++ {
		if structType.Field(i).Type.Kind() == reflect.Struct {
			a.ConstructFilter(structType.Field(i).Type, field, value)
		}

		jsonTag := structType.Field(i).Tag.Get("json")
		if jsonTag == field {
			fieldQuery := ""
			valueQuery := value
			operatorValue := "=" //will be the default condition

			filterTag := structType.Field(i).Tag.Get("filter")
			if filterTag != "" {
				// !TODO prepare for multiple value in filter tag, even for now we just use 1 value (column)
				filterTagValues := strings.Split(filterTag, ";")
				for _, v := range filterTagValues {
					if v != "" {
						filterTagValue := strings.Split(v, ":")
						if filterTagValue[0] == "column" {
							fieldQuery += filterTagValue[1]
						}
						if filterTagValue[0] == "operator" {
							operatorValue = filterTagValue[1]
						}
						if filterTagValue[0] == "ignore" {
							valueQuery = ""
							fieldQuery = ""
							operatorValue = ""
						}
					}
				}
			} else {
				fieldQuery = a.name + "." + field
			}

			a.Payload.Query = append(a.Payload.Query, FilterQuery{
				Field:    fieldQuery,
				Value:    valueQuery,
				Operator: operatorValue,
			})
		}
	}
}

func (a *FilterBuilder[T]) ConstructSort(structType reflect.Type, value string) {
	for i := 0; i < structType.NumField(); i++ {
		if structType.Field(i).Type.Kind() == reflect.Struct {
			a.ConstructSort(structType.Field(i).Type, value)
		}

		isAsc := true
		field := value
		if strings.Contains(value, "-") {
			isAsc = false
			field = value[1:]
		}

		jsonTag := structType.Field(i).Tag.Get("json")

		if jsonTag == field {
			fieldQuery := ""

			filterTag := structType.Field(i).Tag.Get("filter")
			if filterTag != "" {
				// !TODO prepare for multiple value in filter tag, even for now we just use 1 value (column)
				filterTagValues := strings.Split(filterTag, ";")
				for _, v := range filterTagValues {
					if v != "" {
						filterTagValue := strings.Split(v, ":")
						if filterTagValue[0] == "column" {
							fieldQuery += filterTagValue[1]
						}
					}
				}
			} else {
				fieldQuery = field
			}

			if !isAsc {
				fieldQuery = "-" + fieldQuery
			}

			a.Payload.SortBy = append(a.Payload.SortBy, fieldQuery)
		}
	}
}
