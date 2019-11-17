package utils

import (
	"context"
	"io"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/engine"
	canonical_image "github.com/deviceplane/deviceplane/pkg/image"
	"github.com/deviceplane/deviceplane/pkg/models"
)

func ContainerCreate(ctx context.Context, eng engine.Engine, name string, service models.Service) (string, error) {
	var id string

	err := Retry(ctx, func(ctx context.Context) error {
		var err error
		id, err = eng.CreateContainer(ctx, name, service)
		if err != nil {
			log.WithError(err).Error("create container")
			return err
		}
		return nil
	}, 2*time.Minute)

	return id, err
}

func ContainerStart(ctx context.Context, eng engine.Engine, id string) error {
	return Retry(ctx, func(ctx context.Context) error {
		if err := eng.StartContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("start container")
			return err
		}
		return nil
	}, 2*time.Minute)
}

func ContainerList(ctx context.Context, eng engine.Engine, keyFilters map[string]struct{}, keyAndValueFilters map[string]string, all bool) ([]engine.Instance, error) {
	var instances []engine.Instance

	err := Retry(ctx, func(ctx context.Context) error {
		var err error
		instances, err = eng.ListContainers(ctx, keyFilters, keyAndValueFilters, all)
		if err != nil {
			log.WithError(err).Error("list containers")
			return err
		}
		return nil
	}, 2*time.Minute)

	return instances, err
}

func ContainerStop(ctx context.Context, eng engine.Engine, id string) error {
	return Retry(ctx, func(ctx context.Context) error {
		if err := eng.StopContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("stop container")
			return err
		}
		return nil
	}, 2*time.Minute)
}

func ContainerRemove(ctx context.Context, eng engine.Engine, id string) error {
	return Retry(ctx, func(ctx context.Context) error {
		if err := eng.RemoveContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("remove container")
			return err
		}
		return nil
	}, 2*time.Minute)
}

func ImagePull(ctx context.Context, eng engine.Engine, image string, w io.Writer) error {
	return Retry(ctx, func(ctx context.Context) error {
		if err := eng.PullImage(ctx, canonical_image.ToCanonical(image), w); err != nil {
			log.WithError(err).Error("pull image")
			return err
		}
		return nil
	}, time.Hour)
}
