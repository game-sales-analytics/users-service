package grpcsrv

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"

	"github.com/game-sales-analytics/users-service/internal/pb"
)

func (s server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	traceContextFields := apmlogrus.TraceContext(ctx)
	logrus.WithFields(traceContextFields).Debug("handling ping request")

	return &pb.PingReply{
		Pong: true,
	}, nil
}
