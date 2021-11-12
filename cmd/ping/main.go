package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/game-sales-analytics/users-service/internal/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logger := logrus.New()

	// TODO: make level configurable
	logger.SetLevel(logrus.TraceLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})

	if len(os.Args) != 2 {
		logger.Fatal("missing server address argument")
	}

	serverAddr := os.Args[1]
	logger.WithField("address", serverAddr).Debug("ping provided server address")

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(serverAddr, opts...)
	if nil != err {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewUsersServiceClient(conn)
	reply, err := client.Ping(context.Background(), &pb.PingRequest{})
	if nil != err {
		logger.WithError(err).Fatal("failed pinging server")
	}

	logger.WithField("pong", reply.Pong).Printf("ping successful")
	logger.Printf("exiting")
}
