package spec

import (
	"fmt"
	"sort"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/hash"
	"github.com/deviceplane/deviceplane/pkg/models"
)

func WithStandardLabels(s models.Service, applicationID, serviceName string) models.Service {
	// Calculate hash before adding standard labels
	hash := Hash(s, serviceName)

	// TODO
	if s.Labels == nil {
		s.Labels = make(map[string]string)
	}
	s.Labels[models.ApplicationLabel] = applicationID
	s.Labels[models.ServiceLabel] = serviceName
	s.Labels[models.HashLabel] = hash

	return s
}

func Hash(s models.Service, name string) string {
	return applyHash(s, name, hash.Hash)
}

func ShortHash(s models.Service, name string) string {
	return applyHash(s, name, hash.ShortHash)
}

func applyHash(s models.Service, name string, hash func(string) string) string {
	mapToSlice := func(m map[string]string) []string {
		var s []string
		for k, v := range m {
			s = append(s, fmt.Sprintf("%s::%s", k, v))
		}
		sort.Strings(s)
		return s
	}

	var parts []string
	parts = append(parts, name)
	parts = append(parts, s.CapAdd...)
	parts = append(parts, s.CapDrop...)
	parts = append(parts, s.Command...)
	parts = append(parts, s.CPUSet)
	parts = append(parts, fmt.Sprint(s.CPUShares))
	parts = append(parts, fmt.Sprint(s.CPUQuota))
	parts = append(parts, s.DNS...)
	parts = append(parts, s.DNSOpts...)
	parts = append(parts, s.DNSSearch...)
	parts = append(parts, s.DomainName)
	parts = append(parts, s.Entrypoint...)
	parts = append(parts, s.Environment...)
	parts = append(parts, s.ExtraHosts...)
	parts = append(parts, s.GroupAdd...)
	parts = append(parts, s.Image)
	parts = append(parts, s.Hostname)
	parts = append(parts, s.Ipc)
	parts = append(parts, mapToSlice(s.Labels)...)
	parts = append(parts, fmt.Sprint(s.MemLimit))
	parts = append(parts, fmt.Sprint(s.MemReservation, 10))
	parts = append(parts, fmt.Sprint(s.MemSwapLimit, 10))
	parts = append(parts, fmt.Sprint(s.NetworkMode, 10))
	parts = append(parts, fmt.Sprint(s.OomKillDisable))
	parts = append(parts, fmt.Sprint(s.OomScoreAdj))
	parts = append(parts, s.Pid)
	parts = append(parts, s.Ports...)
	parts = append(parts, fmt.Sprint(s.Privileged))
	parts = append(parts, fmt.Sprint(s.ReadOnly))
	parts = append(parts, s.SecurityOpt...)
	parts = append(parts, fmt.Sprint(s.ShmSize))
	parts = append(parts, s.StopSignal)
	parts = append(parts, s.User)
	parts = append(parts, s.Uts)
	parts = append(parts, s.Volumes.HashString())
	parts = append(parts, s.WorkingDir)

	return hash(strings.Join(parts, ":"))
}
