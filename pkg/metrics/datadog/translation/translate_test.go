package translation

import (
	"fmt"
	"strings"
	"testing"

	prometheus "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func TestParsing(t *testing.T) {
	test := func(metricsString string, statsCache *StatsCache, metricTypesToCount []prometheus.MetricType) {
		fmt.Println("Run...")

		parsed := make(map[string]bool, 0)
		parsable := make(map[string]bool, 0)

		reader := strings.NewReader(exampleMetrics)
		metrics, err := ConvertOpenMetricsToDataDog(reader, statsCache, "")
		if err != nil {
			t.Error("Could not parse metrics")
		}

		for _, metric := range metrics {
			parsed[metric.Metric] = true
		}

		parser := expfmt.TextParser{}
		promMetrics, err := parser.TextToMetricFamilies(strings.NewReader(exampleMetrics))
		if err != nil {
			t.Error("For some reason, prometheus's parser could not parse metrics")
		}

		var parsableMetricCount int
		for _, promMetric := range promMetrics {
			for _, countableMetricType := range metricTypesToCount {
				if *promMetric.Type == countableMetricType {
					parsable[promMetric.GetName()] = true
					parsableMetricCount += len(promMetric.Metric)
					break
				}
			}
		}

		if len(metrics) != parsableMetricCount {
			t.Errorf("metric count (%d) is not equal to total parsable metric count (%d)",
				len(metrics),
				parsableMetricCount,
			)
		}

		for k := range parsed {
			if !parsable[k] {
				fmt.Println("Parsed but not parsable?", k)
			}
		}
		for k := range parsable {
			if !parsed[k] {
				fmt.Println("Parsable but not parsed...", k)
			}
		}
		fmt.Println()
	}

	cache := NewStatsCache()

	// statsCache's counter cache should be unpopulated
	test(exampleMetrics, cache, []prometheus.MetricType{
		prometheus.MetricType_GAUGE,
	})

	// statsCache's counters should now be populated (but with delta 0)
	test(exampleMetrics, cache, []prometheus.MetricType{
		prometheus.MetricType_GAUGE,
		prometheus.MetricType_COUNTER,
	})
}

func TestPrintParsing(t *testing.T) {
	parser := expfmt.TextParser{}
	promMetrics, err := parser.TextToMetricFamilies(strings.NewReader(exampleMetrics))
	if err != nil {
		t.Error("For some reason, prometheus's parser could not parse metrics")
	}

	for _, m := range promMetrics {
		fmt.Println("---", m.GetName())
		for _, j := range m.GetMetric() {
			fmt.Println("|>", j.String())
		}
	}
}

const exampleMetrics = `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 8.263e-06
go_gc_duration_seconds{quantile="0.25"} 1.6897e-05
go_gc_duration_seconds{quantile="0.5"} 2.2286e-05
go_gc_duration_seconds{quantile="0.75"} 3.2843e-05
go_gc_duration_seconds{quantile="1"} 0.009978807
go_gc_duration_seconds_sum 0.01952826
go_gc_duration_seconds_count 130
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 24
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.12"} 0
go_info{version="go1.13"} 1
go_info{version="go1.14"} 0
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.037896e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 8.325016e+07
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.462344e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 997792
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 2.682178892473599e-06
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.377728e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 1.037896e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 6.316032e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 3.03104e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 6081
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 6.1693952e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.619136e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.573080119690432e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 1.003873e+06
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 6944
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 55760
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 81920
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.194304e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.239216e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 917504
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 917504
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.2286456e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 18
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 15
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
`
