package translation

import (
	"strings"

	"github.com/deviceplane/deviceplane/pkg/models"
)

type StatsCache struct {
	counterCache map[string]float64
}

func NewStatsCache() *StatsCache {
	return &StatsCache{
		counterCache: make(map[string]float64),
	}
}

func squish(prefix, metric string, tags []string) string {
	return strings.Join(
		append([]string{prefix, metric}, tags...),
		"/",
	)
}

func GetMetricsPrefix(project *models.Project, device *models.Device, endpoint string) string {
	return strings.Join([]string{project.Name, device.Name, endpoint}, "/")
}

func (s *StatsCache) UpdateCount(prefix, metric string, tags []string, newCount float64) (delta float64, ok bool) {
	key := squish(prefix, metric, tags)
	currentCount, ok := s.counterCache[key]
	if ok {
		delta = newCount - currentCount
	}
	s.counterCache[key] = currentCount

	return delta, ok
}
