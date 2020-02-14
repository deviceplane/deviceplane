package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

var _ engine.Engine = &Engine{}

type Engine struct {
	client *client.Client
}

func NewEngine() (*Engine, error) {
	client, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &Engine{
		client: client,
	}, nil
}

func (e *Engine) CreateContainer(ctx context.Context, name string, s models.Service) (string, error) {
	config, hostConfig, err := convert(s)
	if err != nil {
		return "", err
	}

	resp, err := e.client.ContainerCreate(ctx, config, hostConfig, nil, name)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (e *Engine) InspectContainer(ctx context.Context, id string) (*engine.InspectResponse, error) {
	container, err := e.client.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}

	var exitCode *int
	var containerErr string
	if container.State != nil {
		exitCode = &container.State.ExitCode
		containerErr = container.State.Error
	}
	return &engine.InspectResponse{
		PID:      container.State.Pid,
		ExitCode: exitCode,
		Error:    containerErr,
	}, nil
}

func (e *Engine) StartContainer(ctx context.Context, id string) error {
	if err := e.client.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		// TODO
		if strings.Contains(err.Error(), "No such container") {
			return engine.ErrInstanceNotFound
		}
		return err
	}
	return nil
}

func (e *Engine) ListContainers(ctx context.Context, keyFilters map[string]struct{}, keyAndValueFilters map[string]string, all bool) ([]engine.Instance, error) {
	args := filters.NewArgs()
	for k := range keyFilters {
		args.Add("label", k)
	}
	for k, v := range keyAndValueFilters {
		args.Add("label", fmt.Sprintf("%s=%s", k, v))
	}

	containers, err := e.client.ContainerList(ctx, types.ContainerListOptions{
		Filters: args,
		All:     all,
	})
	if err != nil {
		return nil, err
	}

	var instances []engine.Instance
	for _, container := range containers {
		instances = append(instances, convertToInstance(container))
	}

	return instances, nil
}

func (e *Engine) StopContainer(ctx context.Context, id string) error {
	if err := e.client.ContainerStop(ctx, id, nil); err != nil {
		// TODO
		if strings.Contains(err.Error(), "No such container") {
			return engine.ErrInstanceNotFound
		}
		return engine.ErrInstanceNotFound
	}
	return nil
}

func (e *Engine) RemoveContainer(ctx context.Context, id string) error {
	if err := e.client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{}); err != nil {
		// TODO
		if strings.Contains(err.Error(), "No such container") {
			return engine.ErrInstanceNotFound
		}
		return engine.ErrInstanceNotFound
	}
	return nil
}

func (e *Engine) PullImage(ctx context.Context, image, registryAuth string, w io.Writer) error {
	processedRegistryAuth := ""
	if registryAuth != "" {
		var err error
		processedRegistryAuth, err = getProcessedRegistryAuth(registryAuth)
		if err != nil {
			return err
		}
	}

	out, err := e.client.ImagePull(ctx, image, types.ImagePullOptions{
		RegistryAuth: processedRegistryAuth,
	})
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(w, out)
	return err
}

func getProcessedRegistryAuth(registryAuth string) (string, error) {
	decodedRegistryAuth, err := base64.StdEncoding.DecodeString(registryAuth)
	if err != nil {
		return "", errors.Wrap(err, "invalid registry auth")
	}

	registryAuthParts := strings.SplitN(string(decodedRegistryAuth), ":", 2)
	if len(registryAuthParts) != 2 {
		return "", errors.New("invalid registry auth")
	}

	processedRegistryAuthBytes, err := json.Marshal(types.AuthConfig{
		Username: registryAuthParts[0],
		Password: registryAuthParts[1],
	})
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(processedRegistryAuthBytes), nil
}
