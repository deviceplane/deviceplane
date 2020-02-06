package translation

import (
	"strings"
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

func (s *StatsCache) UpdateCount(prefix, metric string, tags []string, newCount float64) (delta float64, ok bool) {
	key := squish(prefix, metric, tags)
	currentCount, ok := s.counterCache[key]
	if ok {
		delta = newCount - currentCount
	}
	s.counterCache[key] = currentCount

	return delta, ok
}
