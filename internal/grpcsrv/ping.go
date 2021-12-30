package grpcsrv

import (
	"context"

	"github.com/getsentry/sentry-go"

	"github.com/game-sales-analytics/users-service/internal/pb"
)

func (s server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	span := sentry.StartSpan(ctx, "ping", sentry.TransactionName("handle-ping-request"))
	span.Status = sentry.SpanStatusOK
	defer span.Finish()

	reply := &pb.PingReply{
		Pong: true,
	}

	return reply, nil
}
