package logging

import (
	"context"
)

type fieldsKey struct{}

// ContextWithFields adds logger fields to fields in context.
func ContextWithFields(parent context.Context, fields Fields) context.Context {
	var newFields Fields
	val := parent.Value(fieldsKey{})
	if val == nil {
		newFields = fields
	} else {
		newFields = make(Fields)
		oldFields, _ := val.(Fields)
		for k, v := range oldFields {
			newFields[k] = v
		}
		for k, v := range fields {
			newFields[k] = v
		}
	}

	return context.WithValue(parent, fieldsKey{}, newFields)
}
