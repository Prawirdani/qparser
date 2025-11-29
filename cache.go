package qparser

import (
	"reflect"
	"sync"
	"time"
)

var structCache sync.Map

type structInfo struct {
	name                 string
	fields               []fieldInfo
	hasUnexportedWithTag bool
}

type fieldInfo struct {
	name     string
	tag      string
	typ      reflect.Type
	index    []int
	isNested bool
}

func getStructCache(rt reflect.Type) *structInfo {
	// Try to load from cache
	if cached, ok := structCache.Load(rt); ok {
		return cached.(*structInfo)
	}

	// Build struct info
	info := &structInfo{name: rt.Name()}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get("qp")

		if !field.IsExported() {
			if tag != "" {
				info.hasUnexportedWithTag = true
			}
			continue
		}

		if tag != "" {
			info.fields = append(info.fields, fieldInfo{
				name:     field.Name,
				tag:      tag,
				typ:      field.Type,
				index:    field.Index,
				isNested: false,
			})
		} else {
			// Check if this field is a nested struct (struct or pointer to struct)
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			if fieldType.Kind() == reflect.Struct && fieldType != reflect.TypeOf(time.Time{}) {
				info.fields = append(info.fields, fieldInfo{
					name:     field.Name,
					tag:      "",
					typ:      field.Type, // Keep the original type (may be pointer)
					index:    field.Index,
					isNested: true,
				})
			}
		}
	}

	// LoadOrStore handles race conditions atomically
	// If another goroutine stored a value first, we return that instead
	actual, _ := structCache.LoadOrStore(rt, info)
	return actual.(*structInfo)
}
