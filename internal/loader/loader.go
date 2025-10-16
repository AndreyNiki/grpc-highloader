package loader

import (
	"context"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/AndreyNiki/grpc-highloader/internal/entity"
	"github.com/AndreyNiki/grpc-highloader/internal/loader/interfaces"
	"github.com/AndreyNiki/grpc-highloader/internal/logger"
	"github.com/AndreyNiki/grpc-highloader/internal/metrics"
)

// RequestLoader implements loader interface for GUI.
type RequestLoader struct {
	requester interfaces.Requester
	req       *entity.RequestParams
	metrics   *metrics.Metrics
}

// NewRequestLoader create a new loader.
func NewRequestLoader(
	requester interfaces.Requester,
	req *entity.RequestParams,
	metrics *metrics.Metrics,
) *RequestLoader {
	return &RequestLoader{
		requester: requester,
		req:       req,
		metrics:   metrics,
	}
}

// Run requests.
func (rl *RequestLoader) Run(ctx context.Context) error {
	md := metadata.New(rl.req.Metadata)
	ctx = metadata.NewOutgoingContext(ctx, md)

	rl.metrics.Reset()

	ticker := time.NewTicker(time.Second / time.Duration(rl.req.RPS))
	defer ticker.Stop()

	log := logger.LoggerFromContext(ctx)

	for {
		if ctx.Err() != nil {
			log.Info("context canceled")
			break
		}

		go func() {
			var cancel context.CancelFunc
			if rl.req.RequestDeadline != nil {
				ctx, cancel = context.WithDeadline(ctx, time.Now().Add(*rl.req.RequestDeadline))
			}

			switch rl.req.MethodType {
			case entity.MethodTypeUnaryRPC:
				err := rl.requester.SendUnaryRPCRequest(ctx)
				if err != nil {
					log.Error("Error send unary rpc request", "Error", err)
				}
			}
			if cancel != nil {
				cancel()
			}
		}()
		<-ticker.C
	}

	return nil
}

// Close loader.
func (rl *RequestLoader) Close() {
	rl.requester.Close()
}
