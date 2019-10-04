package spec

import (
	"fmt"
	"sort"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/hash"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/yamltypes"
)

type Service struct {
	CapAdd         []string                  `yaml:"cap_add,omitempty"`
	CapDrop        []string                  `yaml:"cap_drop,omitempty"`
	Command        yamltypes.Command         `yaml:"command,flow,omitempty"`
	CPUSet         string                    `yaml:"cpuset,omitempty"`
	CPUShares      yamltypes.StringorInt     `yaml:"cpu_shares,omitempty"`
	CPUQuota       yamltypes.StringorInt     `yaml:"cpu_quota,omitempty"`
	DNS            yamltypes.Stringorslice   `yaml:"dns,omitempty"`
	DNSOpts        []string                  `yaml:"dns_opt,omitempty"`
	DNSSearch      yamltypes.Stringorslice   `yaml:"dns_search,omitempty"`
	DomainName     string                    `yaml:"domainname,omitempty"`
	Entrypoint     yamltypes.Command         `yaml:"entrypoint,flow,omitempty"`
	Environment    yamltypes.MaporEqualSlice `yaml:"environment,omitempty"`
	ExtraHosts     []string                  `yaml:"extra_hosts,omitempty"`
	GroupAdd       []string                  `yaml:"group_add,omitempty"`
	Image          string                    `yaml:"image,omitempty"`
	Hostname       string                    `yaml:"hostname,omitempty"`
	Ipc            string                    `yaml:"ipc,omitempty"`
	Labels         yamltypes.SliceorMap      `yaml:"labels,omitempty"`
	MemLimit       yamltypes.MemStringorInt  `yaml:"mem_limit,omitempty"`
	MemReservation yamltypes.MemStringorInt  `yaml:"mem_reservation,omitempty"`
	MemSwapLimit   yamltypes.MemStringorInt  `yaml:"memswap_limit,omitempty"`
	NetworkMode    string                    `yaml:"network_mode,omitempty"`
	OomKillDisable bool                      `yaml:"oom_kill_disable,omitempty"`
	OomScoreAdj    yamltypes.StringorInt     `yaml:"oom_score_adj,omitempty"`
	Pid            string                    `yaml:"pid,omitempty"`
	Ports          []string                  `yaml:"ports,omitempty"`
	Privileged     bool                      `yaml:"privileged,omitempty"`
	ReadOnly       bool                      `yaml:"read_only,omitempty"`
	Restart        string                    `yaml:"restart,omitempty"`
	SecurityOpt    []string                  `yaml:"security_opt,omitempty"`
	ShmSize        yamltypes.MemStringorInt  `yaml:"shm_size,omitempty"`
	StopSignal     string                    `yaml:"stop_signal,omitempty"`
	User           string                    `yaml:"user,omitempty"`
	Uts            string                    `yaml:"uts,omitempty"`
	Volumes        *yamltypes.Volumes        `yaml:"volumes,omitempty"`
	WorkingDir     string                    `yaml:"working_dir,omitempty"`
}

func (s Service) WithStandardLabels(applicationID, serviceName string) Service {
	// Calculate hash before adding standard labels
	hash := s.Hash(serviceName)

	// TODO
	if s.Labels == nil {
		s.Labels = make(map[string]string)
	}
	s.Labels[models.ApplicationLabel] = applicationID
	s.Labels[models.ServiceLabel] = serviceName
	s.Labels[models.HashLabel] = hash

	return s
}

func (s Service) Hash(name string) string {
	return s.hash(name, hash.Hash)
}

func (s Service) ShortHash(name string) string {
	return s.hash(name, hash.ShortHash)
}

func (s Service) hash(name string, hash func(string) string) string {
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
