package spec

import (
	"strings"

	"github.com/deviceplane/deviceplane/pkg/hash"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/yamltypes"
)

type Service struct {
	CapAdd         []string                  `yaml:"cap_add,omitempty"`
	CapDrop        []string                  `yaml:"cap_drop,omitempty"`
	CPUSet         string                    `yaml:"cpuset,omitempty"`
	CPUShares      yamltypes.StringorInt     `yaml:"cpu_shares,omitempty"`
	CPUQuota       yamltypes.StringorInt     `yaml:"cpu_quota,omitempty"`
	Command        yamltypes.Command         `yaml:"command,flow,omitempty"`
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
	Privileged     bool                      `yaml:"privileged,omitempty"`
	SecurityOpt    []string                  `yaml:"security_opt,omitempty"`
	ShmSize        yamltypes.MemStringorInt  `yaml:"shm_size,omitempty"`
	StopSignal     string                    `yaml:"stop_signal,omitempty"`
	Uts            string                    `yaml:"uts,omitempty"`
	ReadOnly       bool                      `yaml:"read_only,omitempty"`
	User           string                    `yaml:"user,omitempty"`
	WorkingDir     string                    `yaml:"working_dir,omitempty"`
}

func (s Service) WithStandardLabels(applicationID, serviceName string) Service {
	// TODO
	if s.Labels == nil {
		s.Labels = make(map[string]string)
	}
	s.Labels[models.ApplicationLabel] = applicationID
	s.Labels[models.ServiceLabel] = serviceName
	s.Labels[models.HashLabel] = s.Hash(serviceName)
	return s
}

func (s Service) Hash(name string) string {
	return s.hash(name, hash.Hash)
}

func (s Service) ShortHash(name string) string {
	return s.hash(name, hash.ShortHash)
}

func (s Service) hash(name string, hash func(string) string) string {
	var parts []string

	parts = append(parts, name)
	parts = append(parts, s.Image)
	parts = append(parts, s.Entrypoint...)
	parts = append(parts, s.Command...)
	// TODO: labels

	return hash(strings.Join(parts, ""))
}
