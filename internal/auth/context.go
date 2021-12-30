package auth

import (
	"context"

	"github.com/getsentry/sentry-go"
)

type Context struct {
	context.Context
	span *sentry.Span
}

func NewContext(ctx context.Context, span *sentry.Span) Context {
	return Context{
		ctx,
		span,
	}
}
