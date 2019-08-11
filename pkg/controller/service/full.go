package service

import (
	"context"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/pkg/errors"
)

func (s *Service) getReleaseFull(ctx context.Context, release models.Release) (*models.ReleaseFull, error) {
	var err error

	var createdByUser *models.User
	if release.CreatedByUserID != nil {
		createdByUser, err = s.users.GetUser(ctx, *release.CreatedByUserID)
		if err != nil {
			return nil, errors.Wrap(err, "get user")
		}
	}

	var createdByServiceAccount *models.ServiceAccount
	if release.CreatedByServiceAccountID != nil {
		createdByServiceAccount, err = s.serviceAccounts.GetServiceAccount(ctx, *release.CreatedByServiceAccountID, release.ProjectID)
		if err != nil {
			return nil, errors.Wrap(err, "get service account")
		}
	}

	releaseDeviceCounts, err := s.releaseDeviceCounts.GetReleaseDeviceCounts(ctx, release.ProjectID, release.ApplicationID, release.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get release device counts")
	}

	return &models.ReleaseFull{
		Release:                 release,
		CreatedByUser:           createdByUser,
		CreatedByServiceAccount: createdByServiceAccount,
		DeviceCounts:            *releaseDeviceCounts,
	}, nil
}
