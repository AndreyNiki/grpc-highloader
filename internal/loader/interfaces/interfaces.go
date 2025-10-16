package interfaces

import (
	"context"

	"github.com/AndreyNiki/grpc-highloader/internal/entity"
	"github.com/AndreyNiki/grpc-highloader/internal/metrics"
)

// Requester interface for dynamic GRPC requests.
type Requester interface {
	SendUnaryRPCRequest(ctx context.Context) error
	Close()
}

// RequesterFactory interface for makes Requester.
type RequesterFactory interface {
	NewRequester(req *entity.RequestParams, metrics *metrics.Metrics) (Requester, error)
}
