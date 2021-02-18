package skale

import "github.com/figment-networks/indexing-engine/metrics"

var (
	rawRequestDuration = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexers",
		Subsystem: "api",
		Name:      "request_duration",
		Desc:      "Duration how long it takes to take data",
		Tags:      []string{"endpoint", "status"},
	})
)
