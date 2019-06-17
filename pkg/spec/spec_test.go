package spec

import (
	"testing"

	"github.com/deviceplane/deviceplane/pkg/yamltypes"
	"github.com/stretchr/testify/require"
)

func fullService() Service {
	return Service{
		CapAdd:      []string{"x", "y", "z"},
		CapDrop:     []string{"x", "y", "z"},
		CPUSet:      "x",
		CPUShares:   yamltypes.StringorInt(1),
		CPUQuota:    yamltypes.StringorInt(1),
		Command:     yamltypes.Command([]string{"x", "y", "z"}),
		DNS:         yamltypes.Stringorslice([]string{"x", "y", "z"}),
		DNSOpts:     []string{"x", "y", "z"},
		DNSSearch:   yamltypes.Stringorslice([]string{"x", "y", "z"}),
		DomainName:  "x",
		Entrypoint:  yamltypes.Command([]string{"x", "y", "z"}),
		Environment: yamltypes.MaporEqualSlice([]string{"x", "y", "z"}),
		ExtraHosts:  []string{"x", "y", "z"},
		GroupAdd:    []string{"x", "y", "z"},
		Image:       "x",
		Hostname:    "x",
		Ipc:         "x",
		Labels: yamltypes.SliceorMap(map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		}),
		MemLimit:       yamltypes.MemStringorInt(1),
		MemReservation: yamltypes.MemStringorInt(1),
		MemSwapLimit:   yamltypes.MemStringorInt(1),
		NetworkMode:    "x",
		OomKillDisable: true,
		OomScoreAdj:    yamltypes.StringorInt(1),
		Pid:            "x",
		Privileged:     true,
		SecurityOpt:    []string{"x", "y", "z"},
		ShmSize:        yamltypes.MemStringorInt(1),
		StopSignal:     "x",
		Uts:            "x",
		ReadOnly:       true,
		User:           "x",
		WorkingDir:     "x",
	}
}

func TestHash(t *testing.T) {
	s := fullService()
	require.Equal(t, s.Hash(""), s.Hash(""))
	require.Equal(t, s.Hash("s"), s.Hash("s"))

	require.NotEqual(t, s.Hash("s1"), s.Hash("s2"))

	for _, f := range []func(Service) Service{
		func(s Service) Service {
			s.Image = "xx"
			return s
		},
		func(s Service) Service {
			s.Command = yamltypes.Command([]string{"xx", "yy", "zz"})
			return s
		},
		func(s Service) Service {
			s.MemLimit = yamltypes.MemStringorInt(2)
			return s
		},
		func(s Service) Service {
			s.ReadOnly = false
			return s
		},
		func(s Service) Service {
			s.Labels = yamltypes.SliceorMap(map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
				"k4": "v4",
			})
			return s
		},
		func(s Service) Service {
			s.Labels = yamltypes.SliceorMap(map[string]string{
				"k1": "v1",
				"k2": "v2",
			})
			return s
		},
		func(s Service) Service {
			s.Labels = yamltypes.SliceorMap(map[string]string{
				"k1": "vv1",
				"k2": "vv2",
				"k3": "vv3",
			})
			return s
		},
	} {
		require.NotEqual(t, s.Hash(""), f(s).Hash(""))
	}
}
