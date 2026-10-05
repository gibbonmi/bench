package jsonfile

import (
	"fmt"
	"reflect"
	"strings"
)

func indirectType(schema reflect.Type) reflect.Type {
	for schema != nil && schema.Kind() == reflect.Pointer {
		schema = schema.Elem()
	}
	return schema
}

func elementType(schema reflect.Type) reflect.Type {
	schema = indirectType(schema)
	if schema != nil && (schema.Kind() == reflect.Slice || schema.Kind() == reflect.Array || schema.Kind() == reflect.Map) {
		return schema.Elem()
	}
	return nil
}

func objectField(schema reflect.Type, name string) (reflect.Type, error) {
	schema = indirectType(schema)
	if schema == nil || schema.Kind() != reflect.Struct {
		return elementType(schema), nil
	}
	var match reflect.Type
	for i := 0; i < schema.NumField(); i++ {
		field := schema.Field(i)
		if !field.IsExported() {
			continue
		}
		key, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if key == "-" {
			continue
		}
		if field.Anonymous && key == "" {
			return nil, fmt.Errorf("exact JSON schema requires a name for anonymous field %s", field.Name)
		}
		if key == "" {
			key = field.Name
		}
		if key == name {
			if match != nil {
				return nil, fmt.Errorf("ambiguous JSON schema field %q", name)
			}
			match = field.Type
		}
	}
	if match == nil {
		return nil, fmt.Errorf("unknown exact JSON field %q", name)
	}
	return match, nil
}
