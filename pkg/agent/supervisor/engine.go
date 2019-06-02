package supervisor

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/engine"
	canonical_image "github.com/deviceplane/deviceplane/pkg/image"
	"github.com/deviceplane/deviceplane/pkg/spec"
)

func containerCreate(ctx context.Context, eng engine.Engine, name string, service spec.Service) string {
	var id string

	retry(ctx, func(ctx context.Context) error {
		var err error
		id, err = eng.CreateContainer(ctx, name, service)
		if err != nil {
			log.WithError(err).Error("create container")
			return err
		}
		return nil
	}, 2*time.Minute)

	return id
}

func containerStart(ctx context.Context, eng engine.Engine, id string) {
	retry(ctx, func(ctx context.Context) error {
		if err := eng.StartContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("start container")
			return err
		}
		return nil
	}, 2*time.Minute)
}

func containerList(ctx context.Context, eng engine.Engine, keyFilters map[string]struct{}, keyAndValueFilters map[string]string, all bool) []engine.Instance {
	var instances []engine.Instance

	retry(ctx, func(ctx context.Context) error {
		var err error
		instances, err = eng.ListContainers(ctx, keyFilters, keyAndValueFilters, all)
		if err != nil {
			log.WithError(err).Error("list containers")
			return err
		}
		return nil
	}, 2*time.Minute)

	return instances
}

func containerStop(ctx context.Context, eng engine.Engine, id string) {
	retry(ctx, func(ctx context.Context) error {
		if err := eng.StopContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("stop container")
			return err
		}
		return nil
	}, 2*time.Minute)
}

func containerRemove(ctx context.Context, eng engine.Engine, id string) {
	retry(ctx, func(ctx context.Context) error {
		if err := eng.RemoveContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("remove container")
			return err
		}
		return nil
	}, 2*time.Minute)
}

func imagePull(ctx context.Context, eng engine.Engine, image string) {
	retry(ctx, func(ctx context.Context) error {
		if err := eng.PullImage(ctx, canonical_image.ToCanonical(image)); err != nil {
			log.WithError(err).Error("pull image")
			return err
		}
		return nil
	}, 30*time.Minute)
}
