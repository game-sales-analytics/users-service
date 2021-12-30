package repository

import (
	"context"

	"github.com/getsentry/sentry-go"
)

type DBOperationContext struct {
	context.Context
	span *sentry.Span
}

func NewDBOperationContext(ctx context.Context, span *sentry.Span) DBOperationContext {
	return DBOperationContext{
		ctx,
		span,
	}
}
