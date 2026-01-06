package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxFieldsKey struct{}

func NewContext(parent context.Context, fields ...zap.Field) context.Context {
	val := parent.Value(ctxFieldsKey{})
	if val == nil {
		return context.WithValue(parent, ctxFieldsKey{}, fields)
	}
	return context.WithValue(parent, ctxFieldsKey{}, append(fields, val.([]zap.Field)...))
}
