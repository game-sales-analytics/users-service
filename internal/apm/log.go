package apm

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func SetSpanTagsFromLogEntry(span *sentry.Span, entry *logrus.Entry) {
	for field, value := range entry.Data {
		span.SetTag(field, fmt.Sprintf("%s", value))
	}
}
