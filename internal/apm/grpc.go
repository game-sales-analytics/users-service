package apm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func uuid2traceID(id uuid.UUID) ([16]byte, error) {
	bin, err := id.MarshalBinary()
	if nil != err {
		return [16]byte{}, err
	}
	if len(bin) != 16 {
		return [16]byte{}, errors.New("unable to serialize generated uuid")
	}

	var out [16]byte
	for i := 0; i < 16; i++ {
		out[i] = bin[i]
	}

	return out, nil
}

func generateTraceID() ([16]byte, error) {
	id, err := uuid.NewRandom()
	if nil != err {
		return [16]byte{}, err
	}

	return uuid2traceID(id)
}

func parseTraceID(raw string) ([16]byte, error) {
	id, err := uuid.Parse(raw)
	if nil != err {
		return [16]byte{}, err
	}

	return uuid2traceID(id)
}

func ReadOrGenerateTraceID(ctx context.Context) ([16]byte, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return generateTraceID()
	}

	if traceIDs := meta.Get("sentry-trace-id"); len(traceIDs) > 0 {
		return parseTraceID(traceIDs[0])
	}

	return generateTraceID()
}
