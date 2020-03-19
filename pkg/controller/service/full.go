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

func (s *Service) getUserFull(ctx context.Context, user *models.User) (*models.UserFull, error) {
	if user.InternalUserID != nil {
		internalUser, err := s.internalUsers.GetInternalUser(ctx, *user.InternalUserID)
		if err != nil {
			return nil, errors.Wrap(err, "get internal user")
		}
		return &models.UserFull{
			User:  *user,
			Email: internalUser.Email,
		}, nil
	}

	if user.ExternalUserID != nil {
		externalUser, err := s.externalUsers.GetExternalUser(ctx, *user.ExternalUserID)
		if err != nil {
			return nil, errors.Wrap(err, "get external user")
		}
		return &models.UserFull{
			User:         *user,
			Email:        externalUser.Email,
			ProviderID:   externalUser.ProviderID,
			ProviderName: externalUser.ProviderName,
		}, nil
	}

	return nil, errors.New("user must have an internal or external user ID")
}

func (s *Service) getMembershipFull2(ctx context.Context, membership *models.Membership) (*models.MembershipFull2, error) {
	user, err := s.users.GetUser(ctx, membership.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}
	fullUser, err := s.getUserFull(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "get full user")
	}

	membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(ctx, membership.UserID, membership.ProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "list membership role bindings")
	}

	roles := make([]models.Role, 0)
	for _, membershipRoleBinding := range membershipRoleBindings {
		role, err := s.roles.GetRole(ctx, membershipRoleBinding.RoleID, membership.ProjectID)
		if err != nil {
			return nil, errors.Wrap(err, "get role")
		}
		roles = append(roles, *role)
	}

	return &models.MembershipFull2{
		Membership: *membership,
		User:       *fullUser,
		Roles:      roles,
	}, nil
}
