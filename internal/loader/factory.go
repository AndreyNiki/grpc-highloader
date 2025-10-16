package loader

import (
	"github.com/AndreyNiki/grpc-highloader/internal/entity"
	guiinterfaces "github.com/AndreyNiki/grpc-highloader/internal/gui/interfaces"
	"github.com/AndreyNiki/grpc-highloader/internal/loader/interfaces"
	"github.com/AndreyNiki/grpc-highloader/internal/metrics"
)

// LoaderFactory implements factory interface for GUI.
type LoaderFactory struct {
	requesterFactory interfaces.RequesterFactory
}

// NewLoaderFactory create a new LoaderFactory.
func NewLoaderFactory(requesterFactory interfaces.RequesterFactory) *LoaderFactory {
	return &LoaderFactory{
		requesterFactory: requesterFactory,
	}
}

// NewLoader create a new Loader.
func (f *LoaderFactory) NewLoader(req *entity.RequestParams, metrics *metrics.Metrics) (guiinterfaces.Loader, error) {
	requester, err := f.requesterFactory.NewRequester(req, metrics)
	if err != nil {
		return nil, err
	}

	loader := NewRequestLoader(requester, req, metrics)
	return loader, nil
}
