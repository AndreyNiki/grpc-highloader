package metrics

import (
	"sync/atomic"

	"google.golang.org/grpc/codes"
)

// Metric stat with value.
type Metric struct {
	Value *atomic.Int64
}

// Metrics struct with needed metrics.
type Metrics struct {
	RequestCounter                          *Metric
	RequestPerSecondGauge                   *Metric
	ResponseStatusOKCounter                 *Metric
	ResponseStatusUnknownCounter            *Metric
	ResponseStatusCancelledCounter          *Metric
	ResponseStatusInvalidArgumentCounter    *Metric
	ResponseStatusDeadlineExceededCounter   *Metric
	ResponseStatusNotFoundCounter           *Metric
	ResponseStatusAlreadyExistsCounter      *Metric
	ResponseStatusPermissionDeniedCounter   *Metric
	ResponseStatusResourceExhaustedCounter  *Metric
	ResponseStatusFailedPreconditionCounter *Metric
	ResponseStatusAbortedCounter            *Metric
	ResponseStatusOutOfRangeCounter         *Metric
	ResponseStatusUnimplementedCounter      *Metric
	ResponseStatusUnavailableCounter        *Metric
	ResponseStatusDataLossCounter           *Metric
	ResponseStatusUnauthenticatedCounter    *Metric
}

// InitMetrics initialize metrics.
func InitMetrics() *Metrics {
	return &Metrics{
		RequestCounter:                          &Metric{Value: &atomic.Int64{}},
		RequestPerSecondGauge:                   &Metric{Value: &atomic.Int64{}},
		ResponseStatusOKCounter:                 &Metric{Value: &atomic.Int64{}},
		ResponseStatusUnknownCounter:            &Metric{Value: &atomic.Int64{}},
		ResponseStatusCancelledCounter:          &Metric{Value: &atomic.Int64{}},
		ResponseStatusInvalidArgumentCounter:    &Metric{Value: &atomic.Int64{}},
		ResponseStatusDeadlineExceededCounter:   &Metric{Value: &atomic.Int64{}},
		ResponseStatusNotFoundCounter:           &Metric{Value: &atomic.Int64{}},
		ResponseStatusAlreadyExistsCounter:      &Metric{Value: &atomic.Int64{}},
		ResponseStatusPermissionDeniedCounter:   &Metric{Value: &atomic.Int64{}},
		ResponseStatusResourceExhaustedCounter:  &Metric{Value: &atomic.Int64{}},
		ResponseStatusFailedPreconditionCounter: &Metric{Value: &atomic.Int64{}},
		ResponseStatusAbortedCounter:            &Metric{Value: &atomic.Int64{}},
		ResponseStatusOutOfRangeCounter:         &Metric{Value: &atomic.Int64{}},
		ResponseStatusUnimplementedCounter:      &Metric{Value: &atomic.Int64{}},
		ResponseStatusUnavailableCounter:        &Metric{Value: &atomic.Int64{}},
		ResponseStatusDataLossCounter:           &Metric{Value: &atomic.Int64{}},
		ResponseStatusUnauthenticatedCounter:    &Metric{Value: &atomic.Int64{}},
	}
}

// SetRequestPerSecond set value for RequestPerSecondGauge.
func (m *Metrics) SetRequestPerSecond(value int64) {
	m.RequestPerSecondGauge.Value.Store(value)
}

// IncrementRequestCount increment value for RequestCounter.
func (m *Metrics) IncrementRequestCount() {
	m.RequestCounter.Value.Add(1)
}

// IncrementResponseStatus increment response status by code.
func (m *Metrics) IncrementResponseStatus(code codes.Code) {
	switch code {
	case codes.OK:
		m.ResponseStatusOKCounter.Value.Add(1)
	case codes.Unknown:
		m.ResponseStatusUnknownCounter.Value.Add(1)
	case codes.Canceled:
		m.ResponseStatusCancelledCounter.Value.Add(1)
	case codes.InvalidArgument:
		m.ResponseStatusInvalidArgumentCounter.Value.Add(1)
	case codes.DeadlineExceeded:
		m.ResponseStatusDeadlineExceededCounter.Value.Add(1)
	case codes.NotFound:
		m.ResponseStatusNotFoundCounter.Value.Add(1)
	case codes.AlreadyExists:
		m.ResponseStatusAlreadyExistsCounter.Value.Add(1)
	case codes.PermissionDenied:
		m.ResponseStatusPermissionDeniedCounter.Value.Add(1)
	case codes.ResourceExhausted:
		m.ResponseStatusResourceExhaustedCounter.Value.Add(1)
	case codes.FailedPrecondition:
		m.ResponseStatusFailedPreconditionCounter.Value.Add(1)
	case codes.Aborted:
		m.ResponseStatusAbortedCounter.Value.Add(1)
	case codes.OutOfRange:
		m.ResponseStatusOutOfRangeCounter.Value.Add(1)
	case codes.Unimplemented:
		m.ResponseStatusUnimplementedCounter.Value.Add(1)
	case codes.Unavailable:
		m.ResponseStatusUnavailableCounter.Value.Add(1)
	case codes.DataLoss:
		m.ResponseStatusDataLossCounter.Value.Add(1)
	case codes.Unauthenticated:
		m.ResponseStatusUnauthenticatedCounter.Value.Add(1)
	}
}

// Reset all metrics.
func (m *Metrics) Reset() {
	m.RequestCounter.Value.Store(0)
	m.ResponseStatusOKCounter.Value.Store(0)
	m.ResponseStatusUnknownCounter.Value.Store(0)
	m.ResponseStatusCancelledCounter.Value.Store(0)
	m.ResponseStatusInvalidArgumentCounter.Value.Store(0)
	m.ResponseStatusDeadlineExceededCounter.Value.Store(0)
	m.ResponseStatusNotFoundCounter.Value.Store(0)
	m.ResponseStatusAlreadyExistsCounter.Value.Store(0)
	m.ResponseStatusPermissionDeniedCounter.Value.Store(0)
	m.ResponseStatusResourceExhaustedCounter.Value.Store(0)
	m.ResponseStatusFailedPreconditionCounter.Value.Store(0)
	m.ResponseStatusAbortedCounter.Value.Store(0)
	m.ResponseStatusOutOfRangeCounter.Value.Store(0)
	m.ResponseStatusUnimplementedCounter.Value.Store(0)
	m.ResponseStatusUnavailableCounter.Value.Store(0)
	m.ResponseStatusDataLossCounter.Value.Store(0)
	m.ResponseStatusUnauthenticatedCounter.Value.Store(0)
}
