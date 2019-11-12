package docker

import (
	"path/filepath"

	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/yamltypes"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

func convert(s models.Service) (*container.Config, *container.HostConfig, error) {
	exposedPorts, portBindings, err := ports(s.Ports)
	if err != nil {
		return nil, nil, err
	}
	return &container.Config{
			Cmd:          strslice.StrSlice(s.Command),
			Domainname:   s.DomainName,
			Entrypoint:   strslice.StrSlice(s.Entrypoint),
			Env:          s.Environment,
			ExposedPorts: exposedPorts,
			Hostname:     s.Hostname,
			Image:        s.Image,
			Labels:       s.Labels,
			StopSignal:   s.StopSignal,
			User:         s.User,
			WorkingDir:   s.WorkingDir,
		}, &container.HostConfig{
			Binds:          volumes(s.Volumes),
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
			PortBindings:   portBindings,
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
			RestartPolicy: container.RestartPolicy{
				Name: s.Restart,
			},
			ShmSize:     int64(s.ShmSize),
			SecurityOpt: s.SecurityOpt,
			UTSMode:     container.UTSMode(s.Uts),
		}, nil
}

func ports(portSpecs []string) (map[nat.Port]struct{}, nat.PortMap, error) {
	ports, binding, err := nat.ParsePortSpecs(portSpecs)
	if err != nil {
		return nil, nil, err
	}

	exposedPorts := map[nat.Port]struct{}{}
	for k, v := range ports {
		exposedPorts[nat.Port(k)] = v
	}

	portBindings := nat.PortMap{}
	for k, bv := range binding {
		dcbs := make([]nat.PortBinding, len(bv))
		for k, v := range bv {
			dcbs[k] = nat.PortBinding{
				HostIP:   v.HostIP,
				HostPort: v.HostPort,
			}
		}
		portBindings[nat.Port(k)] = dcbs
	}

	return exposedPorts, portBindings, nil
}

func volumes(volumes *yamltypes.Volumes) []string {
	if volumes == nil {
		return nil
	}

	var vols []string
	for _, v := range volumes.Volumes {
		if filepath.IsAbs(v.Source) {
			vols = append(vols, v.String())
		}
	}

	return vols
}

func convertToInstance(c types.Container) engine.Instance {
	return engine.Instance{
		ID:     c.ID,
		Labels: c.Labels,
		// TODO
		Running: c.State == "running",
	}
}
