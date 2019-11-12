package models

import "github.com/deviceplane/deviceplane/pkg/yamltypes"

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
