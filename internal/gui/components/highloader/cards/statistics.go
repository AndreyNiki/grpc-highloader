package cards

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/AndreyNiki/grpc-highloader/internal/gui/utils"
	"github.com/AndreyNiki/grpc-highloader/internal/metrics"
)

const (
	zeroValue                             = "0"
	scrapeInterval                        = 20 * time.Millisecond
	labelStatisticsReqs                   = "~Req/s"
	labelStatisticsTotalReqs              = "Total Requests"
	labelStatisticsReqsOK                 = "OK"
	labelStatisticsReqsUnknown            = "Unknown"
	labelStatisticsReqsCancelled          = "Cancelled"
	labelStatisticsReqsInvalidArgument    = "InvalidArgument"
	labelStatisticsReqsDeadlineExceeded   = "DeadlineExceeded"
	labelStatisticsReqsNotFound           = "NotFound"
	labelStatisticsReqsAlreadyExists      = "AlreadyExists"
	labelStatisticsReqsPermissionDenied   = "PermissionDenied"
	labelStatisticsReqsResourceExhausted  = "ResourceExhausted"
	labelStatisticsReqsFailedPrecondition = "FailedPrecondition"
	labelStatisticsReqsAborted            = "Aborted"
	labelStatisticsOutOfRange             = "OutOfRange"
	labelStatisticsUnimplemented          = "Unimplemented"
	labelStatisticsUnavailable            = "Unavailable"
	labelStatisticsDataLoss               = "DataLoss"
	labelStatisticsUnauthenticated        = "Unauthenticated"
)

// statistics struct with metrics.
type statistics struct {
	box     *fyne.Container
	Metrics *metrics.Metrics
	stats   []*metricStat
	info    *infoStat
}

// infoStat stat for showing in GUI.
type infoStat struct {
	reqPerSecond *widget.Label
}

// metricStat metric for showing in GUI.
type metricStat struct {
	value  *widget.Label
	metric *metrics.Metric
}

// newStatistics create a new statistics.
func newStatistics(metrics *metrics.Metrics) *statistics {
	s := &statistics{
		Metrics: metrics,
	}

	labelsBox := s.makeLabels()
	s.box = labelsBox
	return s
}

// makeLabels make labels with metrics for GUI.
func (s *statistics) makeLabels() *fyne.Container {
	valueTotalReqs := widget.NewLabel(zeroValue)
	labelTotalReqs := container.NewHBox(widget.NewLabel(labelStatisticsTotalReqs), valueTotalReqs)
	statsTotalReqs := &metricStat{
		value:  valueTotalReqs,
		metric: s.Metrics.RequestCounter,
	}

	valueReqsPerSecond := widget.NewLabel(zeroValue)
	labelReqsPerSecond := container.NewHBox(valueReqsPerSecond, widget.NewLabel(labelStatisticsReqs))
	info := &infoStat{
		reqPerSecond: valueReqsPerSecond,
	}

	mainLabel := container.NewHBox(labelTotalReqs, utils.NewLine(), labelReqsPerSecond)

	valueOK := widget.NewLabel(zeroValue)
	labelOK := container.NewHBox(widget.NewLabel(labelStatisticsReqsOK+":"), valueOK)
	statsOK := &metricStat{
		value:  valueOK,
		metric: s.Metrics.ResponseStatusOKCounter,
	}
	valueUnknown := widget.NewLabel(zeroValue)
	labelUnknown := container.NewHBox(widget.NewLabel(labelStatisticsReqsUnknown+":"), valueUnknown)
	statsUnknown := &metricStat{
		value:  valueUnknown,
		metric: s.Metrics.ResponseStatusUnknownCounter,
	}
	valueCancelled := widget.NewLabel(zeroValue)
	labelCancelled := container.NewHBox(widget.NewLabel(labelStatisticsReqsCancelled+":"), valueCancelled)
	statsCancelled := &metricStat{
		value:  valueCancelled,
		metric: s.Metrics.ResponseStatusCancelledCounter,
	}
	valueInvalidArgument := widget.NewLabel(zeroValue)
	labelInvalidArgument := container.NewHBox(widget.NewLabel(labelStatisticsReqsInvalidArgument+":"),
		valueInvalidArgument)
	statsInvalidArgument := &metricStat{
		value:  valueInvalidArgument,
		metric: s.Metrics.ResponseStatusInvalidArgumentCounter,
	}
	valueDeadlineExceeded := widget.NewLabel(zeroValue)
	labelDeadlineExceeded := container.NewHBox(widget.NewLabel(labelStatisticsReqsDeadlineExceeded+":"),
		valueDeadlineExceeded)
	statsDeadlineExceeded := &metricStat{
		value:  valueDeadlineExceeded,
		metric: s.Metrics.ResponseStatusDeadlineExceededCounter,
	}
	valueNotFound := widget.NewLabel(zeroValue)
	labelNotFound := container.NewHBox(widget.NewLabel(labelStatisticsReqsNotFound+":"), valueNotFound)
	statsNotFound := &metricStat{
		value:  valueNotFound,
		metric: s.Metrics.ResponseStatusNotFoundCounter,
	}
	rowOne := container.NewHBox(labelOK, labelUnknown, labelCancelled, labelInvalidArgument, labelDeadlineExceeded,
		labelNotFound)

	valueAlreadyExists := widget.NewLabel(zeroValue)
	labelAlreadyExists := container.NewHBox(widget.NewLabel(labelStatisticsReqsAlreadyExists+":"), valueAlreadyExists)
	statsAlreadyExists := &metricStat{
		value:  valueAlreadyExists,
		metric: s.Metrics.ResponseStatusAlreadyExistsCounter,
	}
	valuePermissionDenied := widget.NewLabel(zeroValue)
	labelPermissionDenied := container.NewHBox(widget.NewLabel(labelStatisticsReqsPermissionDenied+":"),
		valuePermissionDenied)
	statsPermissionDenied := &metricStat{
		value:  valuePermissionDenied,
		metric: s.Metrics.ResponseStatusPermissionDeniedCounter,
	}
	valueResourceExhausted := widget.NewLabel(zeroValue)
	labelResourceExhausted := container.NewHBox(widget.NewLabel(labelStatisticsReqsResourceExhausted+":"),
		valueResourceExhausted)
	statsResourceExhausted := &metricStat{
		value:  valueResourceExhausted,
		metric: s.Metrics.ResponseStatusResourceExhaustedCounter,
	}
	valueFailedPrecondition := widget.NewLabel(zeroValue)
	labelFailedPrecondition := container.NewHBox(widget.NewLabel(labelStatisticsReqsFailedPrecondition+":"),
		valueFailedPrecondition)
	statsFailedPrecondition := &metricStat{
		value:  valueFailedPrecondition,
		metric: s.Metrics.ResponseStatusFailedPreconditionCounter,
	}
	valueAborted := widget.NewLabel("0")
	labelAborted := container.NewHBox(widget.NewLabel(labelStatisticsReqsAborted+":"), valueAborted)
	statsAborted := &metricStat{
		value:  valueAborted,
		metric: s.Metrics.ResponseStatusAbortedCounter,
	}
	rowTwo := container.NewHBox(labelAlreadyExists, labelPermissionDenied, labelResourceExhausted,
		labelFailedPrecondition, labelAborted)

	valueOutOfRange := widget.NewLabel(zeroValue)
	labelOutOfRange := container.NewHBox(widget.NewLabel(labelStatisticsOutOfRange+":"), valueOutOfRange)
	statsOutOfRange := &metricStat{
		value:  valueOutOfRange,
		metric: s.Metrics.ResponseStatusOutOfRangeCounter,
	}
	valueUnimplemented := widget.NewLabel(zeroValue)
	labelUnimplemented := container.NewHBox(widget.NewLabel(labelStatisticsUnimplemented+":"), valueUnimplemented)
	statsUnimplemented := &metricStat{
		value:  valueUnimplemented,
		metric: s.Metrics.ResponseStatusUnimplementedCounter,
	}
	valueUnavailable := widget.NewLabel(zeroValue)
	labelUnavailable := container.NewHBox(widget.NewLabel(labelStatisticsUnavailable+":"), valueUnavailable)
	statsUnavailable := &metricStat{
		value:  valueUnavailable,
		metric: s.Metrics.ResponseStatusUnavailableCounter,
	}
	valueDataLoss := widget.NewLabel(zeroValue)
	labelDataLoss := container.NewHBox(widget.NewLabel(labelStatisticsDataLoss+":"), valueDataLoss)
	statsDataLoss := &metricStat{
		value:  valueDataLoss,
		metric: s.Metrics.ResponseStatusDataLossCounter,
	}
	valueUnauthenticated := widget.NewLabel(zeroValue)
	labelUnauthenticated := container.NewHBox(widget.NewLabel(labelStatisticsUnauthenticated+":"),
		valueUnauthenticated)
	statsUnauthenticated := &metricStat{
		value:  valueUnauthenticated,
		metric: s.Metrics.ResponseStatusUnauthenticatedCounter,
	}
	rowThree := container.NewHBox(labelOutOfRange, labelUnimplemented,
		labelUnavailable, labelDataLoss, labelUnauthenticated)

	stats := []*metricStat{statsTotalReqs, statsOK, statsUnknown, statsCancelled, statsInvalidArgument,
		statsDeadlineExceeded, statsNotFound, statsAlreadyExists, statsPermissionDenied, statsResourceExhausted,
		statsFailedPrecondition, statsAborted, statsOutOfRange, statsUnimplemented, statsUnavailable, statsDataLoss,
		statsUnauthenticated}
	s.stats = stats
	s.info = info
	box := container.NewVBox(mainLabel, utils.NewLine(), rowOne, rowTwo, rowThree)
	return box
}

// showStats show stats in GUI.
func (s *statistics) showStats(ctx context.Context) {
	fyne.Do(func() {
		s.resetValues()
		s.box.Refresh()
	})
	for {
		if ctx.Err() != nil {
			fyne.Do(func() {
				s.setValues()
				s.box.Refresh()
			})
			return
		}
		fyne.Do(func() {
			s.setValues()
			s.box.Refresh()
		})
		time.Sleep(scrapeInterval)
	}
}

// showStats show info in GUI.
func (s *statistics) showInfo(ctx context.Context) {
	startTime := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			elapsed := time.Since(startTime).Seconds()
			if elapsed == 0 {
				return
			}
			count := s.Metrics.RequestCounter.Value.Load()
			value := float64(count) / elapsed
			fyne.Do(func() {
				s.info.reqPerSecond.SetText(fmt.Sprintf("%f", value))
				s.info.reqPerSecond.Refresh()
			})
			return
		}
	}
}

// setValues set values in stats.
func (s *statistics) setValues() {
	for _, stat := range s.stats {
		stat.value.SetText(strconv.FormatInt(stat.metric.Value.Load(), 10))
	}
}

// resetValues reset values in stats.
func (s *statistics) resetValues() {
	for _, stat := range s.stats {
		stat.value.SetText(zeroValue)
	}
}
