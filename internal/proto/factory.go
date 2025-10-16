package proto

import (
	"github.com/AndreyNiki/grpc-highloader/internal/entity"
	"github.com/AndreyNiki/grpc-highloader/internal/loader/interfaces"
	"github.com/AndreyNiki/grpc-highloader/internal/metrics"
)

// RequesterFactory implements requester factory for Loader.
type RequesterFactory struct{}

// NewRequesterFactory create a new RequesterFactory.
func NewRequesterFactory() *RequesterFactory {
	return &RequesterFactory{}
}

// NewRequester create a new interfaces.Requester.
func (f *RequesterFactory) NewRequester(
	req *entity.RequestParams,
	metrics *metrics.Metrics,
) (interfaces.Requester, error) {
	return NewRequester(req, metrics)
}
