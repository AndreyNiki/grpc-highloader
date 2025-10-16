package interfaces

import (
	"context"

	"github.com/AndreyNiki/grpc-highloader/internal/entity"
	"github.com/AndreyNiki/grpc-highloader/internal/metrics"
)

// Loader interface for running requests.
type Loader interface {
	Run(ctx context.Context) error
	Close()
}

// LoaderFactory interface for making Loader's. Using for each request.
type LoaderFactory interface {
	NewLoader(req *entity.RequestParams, metrics *metrics.Metrics) (Loader, error)
}

// Parser interface for parsing proto file.
type Parser interface {
	ParseProto(fp string) (*entity.ParsedProto, error)
}
