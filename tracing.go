package gotracing

import (
	"runtime"
	"strings"
	"time"

	logger "github.com/Hnampk/go-tracing/fabric-flogging"
	"github.com/prometheus/client_golang/prometheus"
)

type Tracer struct {
	flogging *logger.FabricLogger
	metrics  *MetricsManager
}

type MetricsManager struct {
	countStartMetrics *prometheus.GaugeVec
	countEndMetrics   *prometheus.GaugeVec
	durationMetrics   *prometheus.HistogramVec
}

var (
	ProcessTimeBuckets = []float64{0.5, 0.8, 1, 1.2, 1.5, 2, 2.5, 10, 20, 60}
)

// MustGetTracer creates a tracer with the specified name. If an invalid name
// is provided, the operation will panic.
func MustGetTracer(moduleName string) *Tracer {
	countStartMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "api",
			Subsystem: moduleName,
			Name:      "count_start",
			Help:      "Number of function called",
		}, []string{"func"},
	)
	countEndMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "api",
			Subsystem: moduleName,
			Name:      "count_end",
			Help:      "Number of function done",
		}, []string{"func"},
	)
	durationMetrics := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "api",
			Subsystem: moduleName,
			Name:      "duration",
			Help:      "Amount of time spent to process a transaction",
			Buckets:   ProcessTimeBuckets,
		}, []string{"func"},
	)
	prometheus.MustRegister(
		countStartMetrics,
		countEndMetrics,
		durationMetrics,
	)

	return &Tracer{
		flogging: logger.MustGetLogger(moduleName),
		metrics: &MetricsManager{
			countStartMetrics: countStartMetrics,
			countEndMetrics:   countEndMetrics,
			durationMetrics:   durationMetrics,
		},
	}
}

// StartFunction must be called at the begin of function
func (t *Tracer) StartFunction(traceNo string) (startTime time.Time) {
	startTime, millis := currentMillisWithTime()
	t.flogging.GetRootLogger().Infof("[%s] StartFunction at %d", traceNo, millis)
	t.metrics.countStartMetrics.WithLabelValues(getCallerFuncName()).Add(1)
	return
}

// EndFunction must be called at the end of function
func (t *Tracer) EndFunction(traceNo string) {
	t.flogging.GetRootLogger().Infof("[%s] EndFunction at %d", traceNo, currentMillis())
	t.metrics.countEndMetrics.WithLabelValues(getCallerFuncName()).Add(1)
}

// EndFunctionWithDurationSince same as EndFunction(), but with duration metrics
//	for example, put this line at the beginning of function:
//	defer mylogger.EndFunctionWithDurationSince(traceNo, time.Now())
// 	or
//	mylogger.EndFunctionWithDurationSince(traceNo, startTime)
func (t *Tracer) EndFunctionWithDurationSince(traceNo string, startTime time.Time) {
	duration := time.Since(startTime)
	t.flogging.GetRootLogger().Infof("[%s] EndFunction at %d, duration=%dms", traceNo, currentMillis(), duration.Milliseconds())

	t.metrics.countEndMetrics.WithLabelValues(getCallerFuncName()).Add(1)
	t.metrics.durationMetrics.WithLabelValues(getCallerFuncName()).Observe(float64(duration.Milliseconds()))
}

func getCallerFuncName() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()[strings.LastIndex(f.Name(), ".")+1:]
}

func currentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func currentMillisWithTime() (time.Time, int64) {
	now := time.Now()
	nowMillis := now.UnixNano() / int64(time.Millisecond)
	return now, nowMillis
}
