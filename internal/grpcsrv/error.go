package grpcsrv

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errorInternal = status.Error(codes.Internal, "internal error occurred. try again later.")
)
