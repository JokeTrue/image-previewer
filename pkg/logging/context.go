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
		oldFields, ok := val.(Fields)
		if !ok {
			panic("fields expected to be type of 'Fields'")
		}

		newFields = make(Fields, len(oldFields)+len(fields))
		for k, v := range oldFields {
			newFields[k] = v
		}
		for k, v := range fields {
			newFields[k] = v
		}
	}

	return context.WithValue(parent, fieldsKey{}, newFields)
}
