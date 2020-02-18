package supervisor

import (
	"context"
	"io"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/engine"
	canonical_image "github.com/deviceplane/deviceplane/pkg/image"
	"github.com/deviceplane/deviceplane/pkg/models"
)

const containerCreateTimeout = time.Minute

func containerCreate(ctx context.Context, eng engine.Engine, name string, service models.Service) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, containerCreateTimeout)
	defer cancel()

	id, err := eng.CreateContainer(ctx, name, service)
	if err != nil {
		log.WithError(err).Error("create container")
		return "", err
	}

	return id, nil
}

const containerStartTimeout = time.Minute

func containerStart(ctx context.Context, eng engine.Engine, id string) error {
	ctx, cancel := context.WithTimeout(ctx, containerStartTimeout)
	defer cancel()

	if err := eng.StartContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
		log.WithError(err).Error("start container")
		return err
	}

	return nil
}

const containerListTimeout = time.Minute

func containerList(ctx context.Context, eng engine.Engine, keyFilters map[string]struct{}, keyAndValueFilters map[string]string, all bool) ([]engine.Instance, error) {
	ctx, cancel := context.WithTimeout(ctx, containerListTimeout)
	defer cancel()

	instances, err := eng.ListContainers(ctx, keyFilters, keyAndValueFilters, all)
	if err != nil {
		log.WithError(err).Error("list containers")
		return nil, err
	}

	return instances, nil
}

const containerStopTimeout = time.Minute

func containerStop(ctx context.Context, eng engine.Engine, id string) error {
	ctx, cancel := context.WithTimeout(ctx, containerStopTimeout)
	defer cancel()

	if err := eng.StopContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
		log.WithError(err).Error("stop container")
		return err
	}

	return nil
}

const containerRemoveTimeout = time.Minute

func containerRemove(ctx context.Context, eng engine.Engine, id string) error {
	ctx, cancel := context.WithTimeout(ctx, containerRemoveTimeout)
	defer cancel()

	if err := eng.RemoveContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
		log.WithError(err).Error("remove container")
		return err
	}

	return nil
}

const imagePullTimeout = 48 * time.Hour

func imagePull(ctx context.Context, eng engine.Engine, image string, getRegistryAuth func() string, w io.Writer) error {
	ctx, cancel := context.WithTimeout(ctx, imagePullTimeout)
	defer cancel()

	if err := eng.PullImage(ctx, canonical_image.ToCanonical(image), getRegistryAuth(), w); err != nil {
		log.WithError(err).Error("pull image")
		return err
	}

	return nil
}
