package metrics

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/node_exporter/collector"

	kingpin "gopkg.in/alecthomas/kingpin.v2" // This one specifically...
)

// This should be more than enough concurrent requests. We're only polling this
// once per device (non-concurrently), so other requests come from customer API
// calls. If we exceed this count, we're probably doing something wrong.
const MaxRequestsInFlight = 3

var (
	ErrInvalidCollector = errors.New("invalid collector")

	DefaultNodeCollectorConfig = NodeCollectorConfig{
		ProcFSPath: "/proc",
		SysFSPath:  "/sys",
		RootFSPath: "/",
		Collectors: []string{
			"cpu",
			"diskstats",
			"filesystem",
			"loadavg",
			"meminfo",
			"textfile",
			"time",
			"netdev",
		},
	}
)

type NodeCollectorConfig struct {
	ProcFSPath string
	SysFSPath  string
	RootFSPath string
	Collectors []string
}

func HostMetricsHandler(ncConfig *NodeCollectorConfig) (*http.Handler, error) {
	if ncConfig == nil {
		ncConfig = &DefaultNodeCollectorConfig
	}

	nc, err := NewNodeCollector(ncConfig)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	r := prometheus.NewRegistry()
	if err := r.Register(nc); err != nil {
		return nil, fmt.Errorf("couldn't register node collector: %s", err)
	}

	metricsRegistry := prometheus.NewRegistry()
	handler := promhttp.InstrumentMetricHandler(
		metricsRegistry,
		promhttp.HandlerFor(
			prometheus.Gatherers{metricsRegistry, r},
			promhttp.HandlerOpts{
				ErrorLog:            log.NewErrorLogger(),
				ErrorHandling:       promhttp.ContinueOnError,
				MaxRequestsInFlight: MaxRequestsInFlight,
				Registry:            metricsRegistry,
			},
		),
	)
	return &handler, nil
}

var collectorCreators = map[string]func() (collector.Collector, error){
	"cpu":         collector.NewCPUCollector,
	"diskstats":   collector.NewDiskstatsCollector,
	"filesystem":  collector.NewFilesystemCollector,
	"loadavg":     collector.NewLoadavgCollector,
	"meminfo":     collector.NewMeminfoCollector,
	"textfile":    collector.NewTextFileCollector,
	"time":        collector.NewTimeCollector,
	"runit":       collector.NewRunitCollector,
	"supervisord": collector.NewSupervisordCollector,
	"netdev":      collector.NewNetDevCollector,
	"ntp":         collector.NewNtpCollector,
}

func NewNodeCollector(config *NodeCollectorConfig) (*collector.NodeCollector, error) {
	// We need to do this because node_exporter collectors directly read CLI
	// arguments for config
	_, err := kingpin.CommandLine.Parse([]string{
		"--path.procfs", config.ProcFSPath,
		"--path.sysfs", config.SysFSPath,
		"--path.rootfs", config.RootFSPath,
	})
	if err != nil {
		return nil, err
	}

	collectors := make(map[string]collector.Collector)

	for _, name := range config.Collectors {
		cCreator, exists := collectorCreators[name]
		if !exists {
			return nil, ErrInvalidCollector
		}

		c, err := cCreator()
		if err != nil {
			return nil, err
		}

		collectors[name] = c
	}
	return &collector.NodeCollector{Collectors: collectors}, nil
}
