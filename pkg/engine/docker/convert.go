package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"

	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/spec"
	"github.com/docker/docker/api/types/container"
)

func convert(s spec.Service) (*container.Config, *container.HostConfig) {
	return &container.Config{
			Cmd:        strslice.StrSlice(s.Command),
			Domainname: s.DomainName,
			Entrypoint: strslice.StrSlice(s.Entrypoint),
			Env:        s.Environment,
			Hostname:   s.Hostname,
			Image:      s.Image,
			Labels:     s.Labels,
			StopSignal: s.StopSignal,
			User:       s.User,
			WorkingDir: s.WorkingDir,
		}, &container.HostConfig{
			CapAdd:         strslice.StrSlice(s.CapAdd),
			CapDrop:        strslice.StrSlice(s.CapDrop),
			DNS:            s.DNS,
			DNSOptions:     s.DNSOpts,
			DNSSearch:      s.DNSSearch,
			ExtraHosts:     s.ExtraHosts,
			GroupAdd:       s.GroupAdd,
			IpcMode:        container.IpcMode(s.Ipc),
			NetworkMode:    container.NetworkMode(s.NetworkMode),
			OomScoreAdj:    int(s.OomScoreAdj),
			PidMode:        container.PidMode(s.Pid),
			Privileged:     s.Privileged,
			ReadonlyRootfs: s.ReadOnly,
			Resources: container.Resources{
				CpusetCpus:        s.CPUSet,
				CPUShares:         int64(s.CPUShares),
				CPUQuota:          int64(s.CPUQuota),
				Memory:            int64(s.MemLimit),
				MemoryReservation: int64(s.MemReservation),
				MemorySwap:        int64(s.MemSwapLimit),
				OomKillDisable:    &s.OomKillDisable, // TODO: this might have the wrong default value
			},
			ShmSize:     int64(s.ShmSize),
			SecurityOpt: s.SecurityOpt,
			UTSMode:     container.UTSMode(s.Uts),
		}
}

func convertToInstance(c types.Container) engine.Instance {
	return engine.Instance{
		ID:     c.ID,
		Labels: c.Labels,
		// TODO
		Running: c.State == "running",
	}
}
