package apm

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
)

func RecoverUnaryWithSentry(hub *sentry.Hub, ctx context.Context, request interface{}) {
	if err := recover(); err != nil {
		eventID := hub.RecoverWithContext(
			context.WithValue(ctx, sentry.RequestContextKey, request),
			err,
		)
		if eventID != nil {
			hub.Flush(1 * time.Second)
		}
	}
}
