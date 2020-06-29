package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/authz"
	"github.com/deviceplane/deviceplane/pkg/controller/middleware"
	"github.com/deviceplane/deviceplane/pkg/controller/query"
	"github.com/deviceplane/deviceplane/pkg/controller/scheduling"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/email"
	"github.com/deviceplane/deviceplane/pkg/hash"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/namesgenerator"
	"github.com/deviceplane/deviceplane/pkg/spec"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"gopkg.in/yaml.v2"
)

const (
	sessionCookie = "dp_sess"
)

var (
	errEmailDomainNotAllowed = errors.New("email domain not allowed")
	errEmailAlreadyTaken     = errors.New("email already taken")
	errTokenExpired          = errors.New("token expired")
)

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) health(w http.ResponseWriter, r *http.Request) {
	s.st.Incr("health", nil, 1)
}

func (s *Service) intentional500(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.withSuperUserAuth(w, r, user, func() {
			w.WriteHeader(http.StatusInternalServerError)
		})
	})
}

func (s *Service) registerInternalUser(w http.ResponseWriter, r *http.Request) {
	utils.WithReferrer(w, r, func(referrer *url.URL) {
		var registerRequest struct {
			Email    string `json:"email" validate:"email"`
			Password string `json:"password" validate:"password"`
		}
		if err := read(r, &registerRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		emailDomain, err := utils.GetDomainFromEmail(registerRequest.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		emailDomainAllowed := false
		if len(s.allowedEmailDomains) == 0 {
			emailDomainAllowed = true
		} else {
			for _, allowedEmailDomain := range s.allowedEmailDomains {
				if allowedEmailDomain == emailDomain {
					emailDomainAllowed = true
					break
				}
			}
		}

		if !emailDomainAllowed {
			http.Error(w, errEmailDomainNotAllowed.Error(), http.StatusBadRequest)
			return
		}

		if _, err := s.internalUsers.LookupInternalUser(r.Context(), registerRequest.Email); err == nil {
			http.Error(w, errEmailAlreadyTaken.Error(), http.StatusBadRequest)
			return
		} else if err != store.ErrUserNotFound {
			log.WithError(err).Error("lookup internal user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		internalUser, err := s.internalUsers.CreateInternalUser(r.Context(), registerRequest.Email, hash.Hash(registerRequest.Password))
		if err != nil {
			log.WithError(err).Error("create internal user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if s.email == nil {
			// If email provider is nil then skip the registration workflow
			user, err := s.users.InitializeUser(r.Context(), &internalUser.ID, nil)
			if err != nil {
				log.WithError(err).Error("mark internal user registration completed")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			utils.Respond(w, user)
			return
		}

		registrationTokenValue := ksuid.New().String()

		if _, err := s.registrationTokens.CreateRegistrationToken(r.Context(), internalUser.ID, hash.Hash(registrationTokenValue)); err != nil {
			log.WithError(err).Error("create registration token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.email.Send(email.Request{
			FromName:    s.emailFromName,
			FromAddress: s.emailFromAddress,
			ToName:      internalUser.Email,
			ToAddress:   internalUser.Email,
			Subject:     "Deviceplane Email Confirmation",
			Content: email.Content{
				Title:       "Email Confirmation",
				Body:        "Thank you for using Deviceplane! Please click the button below to confirm your email.",
				ActionTitle: "Confirm Email",
				ActionLink:  fmt.Sprintf("%s://%s/confirm/%s", referrer.Scheme, referrer.Host, registrationTokenValue),
			},
		}); err != nil {
			log.WithError(err).Error("send registration email")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
	})
}

func (s *Service) registerExternalUser(w http.ResponseWriter, r *http.Request) {
	s.withValidatedSsoJWT(w, r, func(ssoJWT models.SsoJWT) {
		emailDomain, err := utils.GetDomainFromEmail(ssoJWT.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		emailDomainAllowed := false
		if len(s.allowedEmailDomains) == 0 {
			emailDomainAllowed = true
		} else {
			for _, allowedEmailDomain := range s.allowedEmailDomains {
				if allowedEmailDomain == emailDomain {
					emailDomainAllowed = true
					break
				}
			}
		}

		if !emailDomainAllowed {
			http.Error(w, errEmailDomainNotAllowed.Error(), http.StatusBadRequest)
			return
		}

		if _, err := s.externalUsers.GetExternalUserByProviderID(r.Context(), ssoJWT.Provider, ssoJWT.Subject); err == nil {
			http.Error(w, "user already registered", http.StatusBadRequest)
			return
		} else if err != store.ErrUserNotFound {
			log.WithError(err).Error("get external user by provider")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		externalUser, err := s.externalUsers.CreateExternalUser(r.Context(), ssoJWT.Provider, ssoJWT.Subject, ssoJWT.Email, ssoJWT.Claims)
		if err != nil {
			log.WithError(err).Error("create external user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := s.users.InitializeUser(r.Context(), nil, &externalUser.ID)
		if err != nil {
			log.WithError(err).Error("initialize external user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user, err = s.users.UpdateUserName(r.Context(), user.ID, ssoJWT.Name)
		if err != nil {
			log.WithError(err).Error("update user name")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.newSession(w, r, user.ID)
	})
}

func (s *Service) confirmRegistration(w http.ResponseWriter, r *http.Request) {
	var confirmRegistrationRequest struct {
		RegistrationTokenValue string `json:"registrationTokenValue" validate:"required,min=1,max=1000"`
	}
	if err := read(r, &confirmRegistrationRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	registrationToken, err := s.registrationTokens.ValidateRegistrationToken(r.Context(),
		hash.Hash(confirmRegistrationRequest.RegistrationTokenValue))
	if err != nil {
		log.WithError(err).Error("validate registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := s.users.InitializeUser(r.Context(), &registrationToken.InternalUserID, nil)
	if err != nil {
		log.WithError(err).Error("mark internal user registration completed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.newSession(w, r, user.ID)
}

func (s *Service) recoverPassword(w http.ResponseWriter, r *http.Request) {
	utils.WithReferrer(w, r, func(referrer *url.URL) {
		var recoverPasswordRequest struct {
			Email string `json:"email" validate:"email"`
		}
		if err := read(r, &recoverPasswordRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		internalUser, err := s.internalUsers.LookupInternalUser(r.Context(), recoverPasswordRequest.Email)
		if err == store.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("lookup user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		passwordRecoveryTokenValue := ksuid.New().String()

		if _, err := s.passwordRecoveryTokens.CreatePasswordRecoveryToken(r.Context(),
			internalUser.ID, hash.Hash(passwordRecoveryTokenValue)); err != nil {
			log.WithError(err).Error("create password recovery token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.email.Send(email.Request{
			FromName:    s.emailFromName,
			FromAddress: s.emailFromAddress,
			ToName:      internalUser.Email,
			ToAddress:   internalUser.Email,
			Subject:     "Deviceplane Password Reset",
			Content: email.Content{
				Title:       "Password Reset",
				Body:        "Please click the button below to reset your password.",
				ActionTitle: "Reset Password",
				ActionLink:  fmt.Sprintf("%s://%s/recover/%s", referrer.Scheme, referrer.Host, passwordRecoveryTokenValue),
			},
		}); err != nil {
			log.WithError(err).Error("send recovery email")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) getPasswordRecoveryToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	passwordRecoveryTokenValue := vars["passwordrecoverytokenvalue"]

	passwordRecoveryToken, err := s.passwordRecoveryTokens.ValidatePasswordRecoveryToken(r.Context(), hash.Hash(passwordRecoveryTokenValue))
	if err == store.ErrPasswordRecoveryTokenNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.WithError(err).Error("validate password recovery token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, passwordRecoveryToken)
}

func (s *Service) changeInternalUserPassword(w http.ResponseWriter, r *http.Request) {
	var changePasswordRequest struct {
		PasswordRecoveryTokenValue string `json:"passwordRecoveryTokenValue" validate:"required,min=1,max=1000"`
		Password                   string `json:"password" validate:"password"`
	}
	if err := read(r, &changePasswordRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordRecoveryToken, err := s.passwordRecoveryTokens.ValidatePasswordRecoveryToken(r.Context(),
		hash.Hash(changePasswordRequest.PasswordRecoveryTokenValue))
	if err == store.ErrPasswordRecoveryTokenNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.WithError(err).Error("validate password recovery token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if time.Now().After(passwordRecoveryToken.ExpiresAt) {
		http.Error(w, errTokenExpired.Error(), http.StatusForbidden)
		return
	}

	if _, err := s.internalUsers.UpdateInternalUserPasswordHash(r.Context(),
		passwordRecoveryToken.InternalUserID, hash.Hash(changePasswordRequest.Password)); err != nil {
		log.WithError(err).Error("update password hash")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func (s *Service) loginInternalUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email" validate:"email"`
		Password string `json:"password" validate:"password"`
	}
	if err := read(r, &loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	internalUser, err := s.internalUsers.ValidateInternalUserWithEmail(r.Context(), loginRequest.Email, hash.Hash(loginRequest.Password))
	if err == store.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("validate internal user with email")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := s.users.GetUserByInternalID(r.Context(), internalUser.ID)
	if err == store.ErrUserNotFound {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		log.WithError(err).Error("get user by internal user id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.newSession(w, r, user.ID)
}

func (s *Service) loginExternalUser(w http.ResponseWriter, r *http.Request) {
	s.withValidatedSsoJWT(w, r, func(ssoJWT models.SsoJWT) {
		externalUser, err := s.externalUsers.GetExternalUserByProviderID(r.Context(), ssoJWT.Provider, ssoJWT.Subject)
		if err == store.ErrUserNotFound {
			http.Error(w, "external user not found", http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get external user by provider ID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := s.users.GetUserByExternalID(r.Context(), externalUser.ID)
		if err == store.ErrUserNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get user by external ID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.newSession(w, r, user.ID)
	})
}

func (s *Service) newSession(w http.ResponseWriter, r *http.Request, userID string) {
	utils.WithReferrer(w, r, func(referrer *url.URL) {
		sessionValue := ksuid.New().String()

		if _, err := s.sessions.CreateSession(r.Context(), userID, hash.Hash(sessionValue)); err != nil {
			log.WithError(err).Error("create session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var secure bool
		switch referrer.Scheme {
		case "http":
			secure = false
		case "https":
			secure = true
		}

		cookie := &http.Cookie{
			Name:  sessionCookie,
			Value: sessionValue,

			Expires: time.Now().AddDate(0, 1, 0),

			Domain:   r.Host,
			Secure:   secure,
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)
	})
}

func (s *Service) logout(w http.ResponseWriter, r *http.Request) {
	sessionValue, err := r.Cookie(sessionCookie)

	switch err {
	case nil:
		session, err := s.sessions.ValidateSession(r.Context(), hash.Hash(sessionValue.Value))
		if err == store.ErrSessionNotFound {
			w.WriteHeader(http.StatusForbidden)
			return
		} else if err != nil {
			log.WithError(err).Error("validate session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.sessions.DeleteSession(r.Context(), session.ID); err != nil {
			log.WithError(err).Error("delete session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.ErrNoCookie:
		return
	default:
		log.WithError(err).Error("get session cookie")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) getMe(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		if user != nil {
			utils.Respond(w, user)
			return
		}

		if serviceAccount != nil {
			utils.Respond(w, serviceAccount)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		return
	})
}

func (s *Service) updateMe(w http.ResponseWriter, r *http.Request) {
	s.withUserAuth(w, r, func(user *models.User) {
		var updateUserRequest struct {
			Password        *string `json:"password"`
			CurrentPassword *string `json:"currentPassword"`
			Name            *string `json:"name"`
		}
		if err := read(r, &updateUserRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if updateUserRequest.Password != nil {
			if user.InternalUserID == nil {
				http.Error(w, "cannot update password for externally authenticated users", http.StatusBadRequest)
				return
			}

			if updateUserRequest.CurrentPassword == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if _, err := s.internalUsers.ValidateInternalUser(r.Context(), *user.InternalUserID, hash.Hash(*updateUserRequest.CurrentPassword)); err == store.ErrUserNotFound {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err != nil {
				log.WithError(err).Error("validate user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if _, err := s.internalUsers.UpdateInternalUserPasswordHash(r.Context(), *user.InternalUserID, hash.Hash(*updateUserRequest.Password)); err != nil {
				log.WithError(err).Error("update password hash")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if updateUserRequest.Name != nil {
			if _, err := s.users.UpdateUserName(r.Context(), user.ID, *updateUserRequest.Name); err != nil {
				log.WithError(err).Error("update user name")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		user, err := s.users.GetUser(r.Context(), user.ID)
		if err != nil {
			log.WithError(err).Error("get user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, user)
	})
}

func (s *Service) listMembershipsByUser(w http.ResponseWriter, r *http.Request) {
	s.withUserAuth(w, r, func(user *models.User) {
		memberships, err := s.memberships.ListMembershipsByUser(r.Context(), user.ID)
		if err != nil {
			log.WithError(err).Error("list memberships by user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ret interface{} = memberships
		if _, ok := r.URL.Query()["full"]; ok {
			membershipsFull := make([]models.MembershipFull1, 0)

			for _, membership := range memberships {
				user, err := s.users.GetUser(r.Context(), membership.UserID)
				if err != nil {
					log.WithError(err).Error("get user")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				project, err := s.projects.GetProject(r.Context(), membership.ProjectID)
				if err != nil {
					log.WithError(err).Error("get project")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				projectDeviceCounts, err := s.projectDeviceCounts.GetProjectDeviceCounts(r.Context(), membership.ProjectID)
				if err != nil {
					log.WithError(err).Error("get project device counts")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				projectApplicationCounts, err := s.projectApplicationCounts.GetProjectApplicationCounts(r.Context(), membership.ProjectID)
				if err != nil {
					log.WithError(err).Error("get project application counts")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				membershipsFull = append(membershipsFull, models.MembershipFull1{
					Membership: membership,
					User:       *user,
					Project: models.ProjectFull{
						Project:           *project,
						DeviceCounts:      *projectDeviceCounts,
						ApplicationCounts: *projectApplicationCounts,
					},
				})
			}

			ret = membershipsFull
		}

		utils.Respond(w, ret)
	})
}

func (s *Service) createUserAccessKey(w http.ResponseWriter, r *http.Request) {
	s.withUserAuth(w, r, func(user *models.User) {
		var createUserAccessKeyRequest struct {
			Description string `json:"description" validate:"description"`
		}
		if err := read(r, &createUserAccessKeyRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userAccessKeyValue := "u" + ksuid.New().String()

		userAccessKey, err := s.userAccessKeys.CreateUserAccessKey(r.Context(),
			user.ID, hash.Hash(userAccessKeyValue), createUserAccessKeyRequest.Description)
		if err != nil {
			log.WithError(err).Error("create user access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, models.UserAccessKeyWithValue{
			UserAccessKey: *userAccessKey,
			Value:         userAccessKeyValue,
		})
	})
}

func (s *Service) getUserAccessKey(w http.ResponseWriter, r *http.Request) {
	s.withUserAuth(w, r, func(user *models.User) {
		vars := mux.Vars(r)
		userAccessKeyID := vars["useraccesskey"]

		userAccessKey, err := s.userAccessKeys.GetUserAccessKey(r.Context(), userAccessKeyID)
		if err == store.ErrUserAccessKeyNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get user access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if userAccessKey.UserID != user.ID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		utils.Respond(w, userAccessKey)
	})
}

func (s *Service) listUserAccessKeys(w http.ResponseWriter, r *http.Request) {
	s.withUserAuth(w, r, func(user *models.User) {
		userAccessKeys, err := s.userAccessKeys.ListUserAccessKeys(r.Context(), user.ID)
		if err != nil {
			log.WithError(err).Error("list users")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, userAccessKeys)
	})
}

func (s *Service) deleteUserAccessKey(w http.ResponseWriter, r *http.Request) {
	s.withUserAuth(w, r, func(user *models.User) {
		vars := mux.Vars(r)
		userAccessKeyID := vars["useraccesskey"]

		userAccessKey, err := s.userAccessKeys.GetUserAccessKey(r.Context(), userAccessKeyID)
		if err == store.ErrUserAccessKeyNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get user access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if userAccessKey.UserID != user.ID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := s.userAccessKeys.DeleteUserAccessKey(r.Context(), userAccessKeyID); err != nil {
			log.WithError(err).Error("delete user access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) createProject(w http.ResponseWriter, r *http.Request) {
	s.withUserAuth(w, r, func(user *models.User) {
		var createProjectRequest struct {
			Name string `json:"name" validate:"name"`
		}
		if err := read(r, &createProjectRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := s.projects.LookupProject(r.Context(), createProjectRequest.Name); err == nil {
			http.Error(w, store.ErrProjectNameAlreadyInUse.Error(), http.StatusBadRequest)
			return
		} else if err != nil && err != store.ErrProjectNotFound {
			log.WithError(err).Error("lookup project")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		project, err := s.projects.CreateProject(r.Context(), createProjectRequest.Name)
		if err != nil {
			log.WithError(err).Error("create project")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = s.memberships.CreateMembership(r.Context(), user.ID, project.ID); err != nil {
			log.WithError(err).Error("create membership")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		adminAllRoleBytes, err := yaml.Marshal(authz.AdminAllRole)
		if err != nil {
			log.WithError(err).Error("marshal admin role config")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		adminRole, err := s.roles.CreateRole(r.Context(), project.ID, "admin-all", "", string(adminAllRoleBytes))
		if err != nil {
			log.WithError(err).Error("create admin role")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeAllRoleBytes, err := yaml.Marshal(authz.WriteAllRole)
		if err != nil {
			log.WithError(err).Error("marshal write role config")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = s.roles.CreateRole(r.Context(), project.ID, "write-all", "", string(writeAllRoleBytes)); err != nil {
			log.WithError(err).Error("create write role")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		readAllRoleBytes, err := yaml.Marshal(authz.ReadAllRole)
		if err != nil {
			log.WithError(err).Error("marshal read role config")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = s.roles.CreateRole(r.Context(), project.ID, "read-all", "", string(readAllRoleBytes)); err != nil {
			log.WithError(err).Error("create read role")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = s.membershipRoleBindings.CreateMembershipRoleBinding(r.Context(),
			user.ID, adminRole.ID, project.ID,
		); err != nil {
			log.WithError(err).Error("create membership role binding")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create default device registration token.
		// It is named "default" and has an unlimited device registration cap.
		_, err = s.deviceRegistrationTokens.CreateDeviceRegistrationToken(
			r.Context(),
			project.ID,
			"default",
			"",
			nil,
		)
		if err != nil {
			log.WithError(err).Error("create default registration token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, project)
	})
}

func (s *Service) getProject(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceProjects, authz.ActionGetProject,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				project, err := s.projects.GetProject(r.Context(), project.ID)
				if err != nil {
					log.WithError(err).Error("get project")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, project)
			},
		)
	})
}

func (s *Service) updateProject(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceProjects, authz.ActionUpdateProject,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				var updateProjectRequest struct {
					Name          string `json:"name" validate:"name"`
					DatadogAPIKey string `json:"datadogApiKey"`
				}
				if err := read(r, &updateProjectRequest); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if p, err := s.projects.LookupProject(r.Context(),
					updateProjectRequest.Name); err == nil && p.ID != project.ID {
					http.Error(w, store.ErrProjectNameAlreadyInUse.Error(), http.StatusBadRequest)
					return
				} else if err != nil && err != store.ErrProjectNotFound {
					log.WithError(err).Error("lookup project")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				p, err := s.projects.UpdateProject(r.Context(), project.ID, updateProjectRequest.Name, updateProjectRequest.DatadogAPIKey)
				if err != nil {
					log.WithError(err).Error("update project")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, p)
			},
		)
	})
}

func (s *Service) deleteProject(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceProjects, authz.ActionDeleteProject,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				if err := s.projects.DeleteProject(r.Context(), project.ID); err != nil {
					log.WithError(err).Error("delete project")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			},
		)
	})
}

func (s *Service) createRole(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceRoles, authz.ActionCreateRole,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				var createRoleRequest struct {
					Name        string `json:"name" validate:"name"`
					Description string `json:"description" validate:"description"`
					Config      string `json:"config" validate:"config"`
				}
				if err := read(r, &createRoleRequest); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if _, err := s.roles.LookupRole(r.Context(), createRoleRequest.Name, project.ID); err == nil {
					http.Error(w, store.ErrRoleNameAlreadyInUse.Error(), http.StatusBadRequest)
					return
				} else if err != nil && err != store.ErrRoleNotFound {
					log.WithError(err).Error("lookup role")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				var roleConfig authz.Config
				if err := yaml.UnmarshalStrict([]byte(createRoleRequest.Config), &roleConfig); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				role, err := s.roles.CreateRole(r.Context(), project.ID, createRoleRequest.Name,
					createRoleRequest.Description, createRoleRequest.Config)
				if err != nil {
					log.WithError(err).Error("create role")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, role)
			},
		)
	})
}

func (s *Service) listRoles(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceRoles, authz.ActionListRoles,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				roles, err := s.roles.ListRoles(r.Context(), project.ID)
				if err != nil {
					log.WithError(err).Error("list roles")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, roles)
			},
		)
	})
}

func (s *Service) getRole(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceRoles, authz.ActionGetRole,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withRole(w, r, project, func(role *models.Role) {
					utils.Respond(w, role)
				})
			},
		)
	})
}

func (s *Service) updateRole(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceRoles, authz.ActionUpdateRole,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withRole(w, r, project, func(role *models.Role) {
					var updateRoleRequest struct {
						Name        string `json:"name" validate:"name"`
						Description string `json:"description" validate:"description"`
						Config      string `json:"config" validate:"config"`
					}
					if err := read(r, &updateRoleRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					var roleConfig authz.Config
					if err := yaml.UnmarshalStrict([]byte(updateRoleRequest.Config), &roleConfig); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					if role, err := s.roles.LookupRole(r.Context(),
						updateRoleRequest.Name, project.ID); err == nil && role.ID != role.ID {
						http.Error(w, store.ErrRoleNameAlreadyInUse.Error(), http.StatusBadRequest)
						return
					} else if err != nil && err != store.ErrRoleNotFound {
						log.WithError(err).Error("lookup role")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					role, err := s.roles.UpdateRole(r.Context(), role.ID, project.ID, updateRoleRequest.Name,
						updateRoleRequest.Description, updateRoleRequest.Config)
					if err != nil {
						log.WithError(err).Error("update role")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, role)
				})
			},
		)
	})
}

func (s *Service) deleteRole(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceRoles, authz.ActionDeleteRole,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withRole(w, r, project, func(role *models.Role) {
					if err := s.roles.DeleteRole(r.Context(), role.ID, project.ID); err != nil {
						log.WithError(err).Error("delete role")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) createServiceAccount(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccounts, authz.ActionCreateServiceAccount,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				var createServiceAccountRequest struct {
					Name        string `json:"name" validate:"name"`
					Description string `json:"description" validate:"description"`
				}
				if err := read(r, &createServiceAccountRequest); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if _, err := s.serviceAccounts.LookupServiceAccount(r.Context(), createServiceAccountRequest.Name, project.ID); err == nil {
					http.Error(w, store.ErrServiceAccountNameAlreadyInUse.Error(), http.StatusBadRequest)
					return
				} else if err != nil && err != store.ErrServiceAccountNotFound {
					log.WithError(err).Error("lookup service account")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				serviceAccount, err := s.serviceAccounts.CreateServiceAccount(r.Context(), project.ID, createServiceAccountRequest.Name,
					createServiceAccountRequest.Description)
				if err != nil {
					log.WithError(err).Error("create service account")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, serviceAccount)
			},
		)
	})
}

func (s *Service) getServiceAccount(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccounts, authz.ActionGetServiceAccount,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withServiceAccount(w, r, project, func(serviceAccount *models.ServiceAccount) {
					var ret interface{} = serviceAccount
					if _, ok := r.URL.Query()["full"]; ok {
						serviceAccountRoleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(), serviceAccount.ID, project.ID)
						if err != nil {
							log.WithError(err).Error("list service account role bindings")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						roles := make([]models.Role, 0)
						for _, serviceAccountRoleBinding := range serviceAccountRoleBindings {
							role, err := s.roles.GetRole(r.Context(), serviceAccountRoleBinding.RoleID, project.ID)
							if err != nil {
								log.WithError(err).Error("get role")
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
							roles = append(roles, *role)
						}

						ret = models.ServiceAccountFull{
							ServiceAccount: *serviceAccount,
							Roles:          roles,
						}
					}

					utils.Respond(w, ret)
				})
			},
		)
	})
}

func (s *Service) listServiceAccounts(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccounts, authz.ActionListServiceAccounts,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				serviceAccounts, err := s.serviceAccounts.ListServiceAccounts(r.Context(), project.ID)
				if err != nil {
					log.WithError(err).Error("list service accounts")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				var ret interface{} = serviceAccounts
				if _, ok := r.URL.Query()["full"]; ok {
					serviceAccountsFull := make([]models.ServiceAccountFull, 0)

					for _, serviceAccount := range serviceAccounts {
						serviceAccountRoleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(), serviceAccount.ID, project.ID)
						if err != nil {
							log.WithError(err).Error("list service account role bindings")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						roles := make([]models.Role, 0)
						for _, serviceAccountRoleBinding := range serviceAccountRoleBindings {
							role, err := s.roles.GetRole(r.Context(), serviceAccountRoleBinding.RoleID, project.ID)
							if err != nil {
								log.WithError(err).Error("get role")
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
							roles = append(roles, *role)
						}

						serviceAccountsFull = append(serviceAccountsFull, models.ServiceAccountFull{
							ServiceAccount: serviceAccount,
							Roles:          roles,
						})
					}

					ret = serviceAccountsFull
				}

				utils.Respond(w, ret)
			},
		)
	})
}

func (s *Service) updateServiceAccount(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccounts, authz.ActionUpdateServiceAccount,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withServiceAccount(w, r, project, func(serviceAccount *models.ServiceAccount) {
					var updateServiceAccountRequest struct {
						Name        string `json:"name" validate:"name"`
						Description string `json:"description" validate:"description"`
					}
					if err := read(r, &updateServiceAccountRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					if sa, err := s.serviceAccounts.LookupServiceAccount(r.Context(),
						updateServiceAccountRequest.Name, project.ID); err == nil && sa.ID != serviceAccount.ID {
						http.Error(w, store.ErrServiceAccountNameAlreadyInUse.Error(), http.StatusBadRequest)
						return
					} else if err != nil && err != store.ErrServiceAccountNotFound {
						log.WithError(err).Error("lookup service account")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					sa, err := s.serviceAccounts.UpdateServiceAccount(r.Context(), serviceAccount.ID, project.ID,
						updateServiceAccountRequest.Name, updateServiceAccountRequest.Description)
					if err != nil {
						log.WithError(err).Error("update service account")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, sa)
				})
			},
		)
	})
}

func (s *Service) deleteServiceAccount(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccounts, authz.ActionDeleteServiceAccount,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withServiceAccount(w, r, project, func(serviceAccount *models.ServiceAccount) {
					if err := s.serviceAccounts.DeleteServiceAccount(r.Context(), serviceAccount.ID, project.ID); err != nil {
						log.WithError(err).Error("delete service account")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) createServiceAccountAccessKey(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccountAccessKeys, authz.ActionCreateServiceAccountAccessKey,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				vars := mux.Vars(r)
				serviceAccountID := vars["serviceaccount"]

				var createServiceAccountAccessKeyRequest struct {
					Description string `json:"description" validate:"description"`
				}
				if err := read(r, &createServiceAccountAccessKeyRequest); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				serviceAccountAccessKeyValue := "s" + ksuid.New().String()

				serviceAccount, err := s.serviceAccountAccessKeys.CreateServiceAccountAccessKey(r.Context(),
					project.ID, serviceAccountID, hash.Hash(serviceAccountAccessKeyValue), createServiceAccountAccessKeyRequest.Description)
				if err != nil {
					log.WithError(err).Error("create service account access key")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, models.ServiceAccountAccessKeyWithValue{
					ServiceAccountAccessKey: *serviceAccount,
					Value:                   serviceAccountAccessKeyValue,
				})
			},
		)
	})
}

func (s *Service) getServiceAccountAccessKey(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccountAccessKeys, authz.ActionGetServiceAccountAccessKey,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				vars := mux.Vars(r)
				serviceAccountAccessKeyID := vars["serviceaccountaccesskey"]

				serviceAccountAccessKey, err := s.serviceAccountAccessKeys.GetServiceAccountAccessKey(r.Context(), serviceAccountAccessKeyID, project.ID)
				if err == store.ErrServiceAccountAccessKeyNotFound {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				} else if err != nil {
					log.WithError(err).Error("get service account access key")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, serviceAccountAccessKey)
			},
		)
	})
}

func (s *Service) listServiceAccountAccessKeys(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccountAccessKeys, authz.ActionListServiceAccountAccessKeys,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				vars := mux.Vars(r)
				serviceAccountID := vars["serviceaccount"]

				serviceAccountAccessKeys, err := s.serviceAccountAccessKeys.ListServiceAccountAccessKeys(r.Context(), project.ID, serviceAccountID)
				if err != nil {
					log.WithError(err).Error("list service accounts")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, serviceAccountAccessKeys)
			},
		)
	})
}

func (s *Service) deleteServiceAccountAccessKey(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccountAccessKeys, authz.ActionDeleteServiceAccountAccessKey,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				vars := mux.Vars(r)
				serviceAccountAccessKeyID := vars["serviceaccountaccesskey"]

				if err := s.serviceAccountAccessKeys.DeleteServiceAccountAccessKey(r.Context(), serviceAccountAccessKeyID, project.ID); err != nil {
					log.WithError(err).Error("delete service account access key")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			},
		)
	})
}

func (s *Service) createServiceAccountRoleBinding(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccountRoleBindings, authz.ActionCreateServiceAccountRoleBinding,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withRole(w, r, project, func(role *models.Role) {
					vars := mux.Vars(r)
					serviceAccountID := vars["serviceaccount"]

					serviceAccountRoleBinding, err := s.serviceAccountRoleBindings.CreateServiceAccountRoleBinding(r.Context(), serviceAccountID, role.ID, project.ID)
					if err != nil {
						log.WithError(err).Error("create service account role binding")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, serviceAccountRoleBinding)
				})
			},
		)
	})
}

func (s *Service) getServiceAccountRoleBinding(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccountRoleBindings, authz.ActionGetServiceAccountRoleBinding,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withRole(w, r, project, func(role *models.Role) {
					vars := mux.Vars(r)
					serviceAccountID := vars["serviceaccount"]

					serviceAccountRoleBinding, err := s.serviceAccountRoleBindings.GetServiceAccountRoleBinding(r.Context(), serviceAccountID, role.ID, project.ID)
					if err != nil {
						log.WithError(err).Error("get service account role binding")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, serviceAccountRoleBinding)
				})
			},
		)
	})
}

func (s *Service) listServiceAccountRoleBindings(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccountRoleBindings, authz.ActionListServiceAccountRoleBinding,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withRole(w, r, project, func(role *models.Role) {
					vars := mux.Vars(r)
					serviceAccountID := vars["serviceaccount"]

					serviceAccountRoleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(), project.ID, serviceAccountID)
					if err != nil {
						log.WithError(err).Error("list service account role bindings")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, serviceAccountRoleBindings)
				})
			},
		)
	})
}

func (s *Service) deleteServiceAccountRoleBinding(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceServiceAccountRoleBindings, authz.ActionDeleteServiceAccountRoleBinding,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withRole(w, r, project, func(role *models.Role) {
					vars := mux.Vars(r)
					serviceAccountID := vars["serviceaccount"]

					if err := s.serviceAccountRoleBindings.DeleteServiceAccountRoleBinding(r.Context(), serviceAccountID, role.ID, project.ID); err != nil {
						log.WithError(err).Error("delete service account role binding")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) createMembership(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceMemberships, authz.ActionCreateMembership,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				var createMembershipRequest struct {
					Email  *string `json:"email" validate:"omitempty,email"`
					UserID *string `json:"userId" validate:"omitempty,id"`
				}
				if err := read(r, &createMembershipRequest); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				var user *models.User
				if createMembershipRequest.UserID != nil {
					var err error
					user, err = s.users.GetUser(r.Context(), *createMembershipRequest.UserID)
					if err == store.ErrUserNotFound {
						http.Error(w, err.Error(), http.StatusNotFound)
						return
					} else if err != nil {
						log.WithError(err).Error("get user")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				} else if createMembershipRequest.Email != nil {
					internalUser, err := s.internalUsers.LookupInternalUser(r.Context(), *createMembershipRequest.Email)
					if err == store.ErrUserNotFound {
						http.Error(w, err.Error(), http.StatusNotFound)
						return
					} else if err != nil {
						log.WithError(err).Error("lookup internal user")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					user, err = s.users.GetUserByInternalID(r.Context(), internalUser.ID)
					if err == store.ErrUserNotFound {
						http.Error(w, err.Error(), http.StatusNotFound)
						return
					} else if err != nil {
						log.WithError(err).Error("get user by internal user id")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				} else {
					http.Error(w, "either email or user ID must exist", http.StatusBadRequest)
					return
				}

				membership, err := s.memberships.CreateMembership(r.Context(), user.ID, project.ID)
				if err != nil {
					log.WithError(err).Error("create membership")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, membership)
			},
		)
	})
}

func (s *Service) getMembership(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceMemberships, authz.ActionGetMembership,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				vars := mux.Vars(r)
				userID := vars["user"]

				membership, err := s.memberships.GetMembership(r.Context(), userID, project.ID)
				if err == store.ErrMembershipNotFound {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				} else if err != nil {
					log.WithError(err).Error("get membership")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				var ret interface{} = membership
				if _, ok := r.URL.Query()["full"]; ok {
					membershipFull2, err := s.getMembershipFull2(r.Context(), membership)
					if err != nil {
						log.WithError(err).Error("get full membership")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					ret = *membershipFull2
				}

				utils.Respond(w, ret)
			},
		)
	})
}

func (s *Service) listMembershipsByProject(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceMemberships, authz.ActionListMembershipsByProject,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				memberships, err := s.memberships.ListMembershipsByProject(r.Context(), project.ID)
				if err != nil {
					log.WithError(err).Error("list memberships by project")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				var ret interface{} = memberships
				if _, ok := r.URL.Query()["full"]; ok {
					membershipsFull := make([]models.MembershipFull2, 0)

					for _, membership := range memberships {
						membershipFull2, err := s.getMembershipFull2(r.Context(), &membership)
						if err != nil {
							log.WithError(err).Error("get full membership")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						membershipsFull = append(membershipsFull, *membershipFull2)
					}

					ret = membershipsFull
				}

				utils.Respond(w, ret)
			},
		)
	})
}

func (s *Service) deleteMembership(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceMemberships, authz.ActionDeleteMembership,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				vars := mux.Vars(r)
				userID := vars["user"]

				if err := s.memberships.DeleteMembership(r.Context(), userID, project.ID); err != nil {
					log.WithError(err).Error("delete membership")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			},
		)
	})
}

func (s *Service) createMembershipRoleBinding(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceMembershipRoleBindings, authz.ActionCreateMembershipRoleBinding,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withRole(w, r, project, func(role *models.Role) {
					vars := mux.Vars(r)
					userID := vars["user"]

					membershipRoleBinding, err := s.membershipRoleBindings.CreateMembershipRoleBinding(r.Context(), userID, role.ID, project.ID)
					if err != nil {
						log.WithError(err).Error("create membership role binding")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, membershipRoleBinding)
				})
			},
		)
	})
}

func (s *Service) getMembershipRoleBinding(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceMembershipRoleBindings, authz.ActionGetMembershipRoleBinding,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withRole(w, r, project, func(role *models.Role) {
					vars := mux.Vars(r)
					userID := vars["user"]

					membershipRoleBinding, err := s.membershipRoleBindings.GetMembershipRoleBinding(r.Context(), userID, role.ID, project.ID)
					if err == store.ErrMembershipRoleBindingNotFound {
						http.Error(w, err.Error(), http.StatusNotFound)
						return
					} else if err != nil {
						log.WithError(err).Error("get membership role binding")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, membershipRoleBinding)
				})
			},
		)
	})
}

func (s *Service) listMembershipRoleBindings(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceMembershipRoleBindings, authz.ActionListMembershipRoleBindings,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				vars := mux.Vars(r)
				userID := vars["user"]

				membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), userID, project.ID)
				if err != nil {
					log.WithError(err).Error("list membership role bindings")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, membershipRoleBindings)
			},
		)
	})
}

func (s *Service) deleteMembershipRoleBinding(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceMembershipRoleBindings, authz.ActionDeleteMembershipRoleBinding,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withRole(w, r, project, func(role *models.Role) {
					vars := mux.Vars(r)
					userID := vars["user"]

					if err := s.membershipRoleBindings.DeleteMembershipRoleBinding(r.Context(), userID, role.ID, project.ID); err != nil {
						log.WithError(err).Error("delete membership role binding")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) createConnection(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceConnections, authz.ActionCreateConnection,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				var createConnectionRequest struct {
					Name     string          `json:"name" validate:"name"`
					Protocol models.Protocol `json:"protocol" validate:"protocol"`
					Port     uint            `json:"port" validate:"port"`
				}
				if err := read(r, &createConnectionRequest); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if _, err := s.connections.LookupConnection(r.Context(), createConnectionRequest.Name, project.ID); err == nil {
					http.Error(w, store.ErrConnectionNameAlreadyInUse.Error(), http.StatusBadRequest)
					return
				} else if err != nil && err != store.ErrConnectionNotFound {
					log.WithError(err).Error("lookup connection")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				connection, err := s.connections.CreateConnection(
					r.Context(),
					project.ID,
					createConnectionRequest.Name,
					createConnectionRequest.Protocol,
					createConnectionRequest.Port,
				)
				if err != nil {
					log.WithError(err).Error("create connection")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, connection)
			},
		)
	})
}

func (s *Service) getConnection(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceConnections, authz.ActionGetConnection,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withConnection(w, r, project, func(connection *models.Connection) {
					utils.Respond(w, connection)
				})
			},
		)
	})
}

func (s *Service) listConnections(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceConnections, authz.ActionListConnections,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				connections, err := s.connections.ListConnections(r.Context(), project.ID)
				if err != nil {
					log.WithError(err).Error("list connections")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, connections)
			},
		)
	})
}

func (s *Service) updateConnection(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceConnections, authz.ActionUpdateConnection,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withConnection(w, r, project, func(connection *models.Connection) {
					var updateConnectionRequest struct {
						Name     string          `json:"name" validate:"name"`
						Protocol models.Protocol `json:"protocol" validate:"protocol"`
						Port     uint            `json:"port" validate:"port"`
					}
					if err := read(r, &updateConnectionRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					if c, err := s.connections.LookupConnection(r.Context(),
						updateConnectionRequest.Name, project.ID); err == nil && c.ID != connection.ID {
						http.Error(w, store.ErrConnectionNameAlreadyInUse.Error(), http.StatusBadRequest)
						return
					} else if err != nil && err != store.ErrConnectionNotFound {
						log.WithError(err).Error("lookup connection")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					c, err := s.connections.UpdateConnection(r.Context(), connection.ID, project.ID,
						updateConnectionRequest.Name, updateConnectionRequest.Protocol,
						updateConnectionRequest.Port,
					)
					if err != nil {
						log.WithError(err).Error("update connection")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, c)
				})
			},
		)
	})
}

func (s *Service) deleteConnection(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceConnections, authz.ActionDeleteConnection,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withConnection(w, r, project, func(connection *models.Connection) {
					if err := s.connections.DeleteConnection(r.Context(), connection.ID, project.ID); err != nil {
						log.WithError(err).Error("delete connection")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) createApplication(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceApplications, authz.ActionCreateApplication,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				var createApplicationRequest struct {
					Name        string `json:"name" validate:"name"`
					Description string `json:"description" validate:"description"`
				}
				if err := read(r, &createApplicationRequest); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if _, err := s.applications.LookupApplication(r.Context(), createApplicationRequest.Name, project.ID); err == nil {
					http.Error(w, store.ErrApplicationNameAlreadyInUse.Error(), http.StatusBadRequest)
					return
				} else if err != nil && err != store.ErrApplicationNotFound {
					log.WithError(err).Error("lookup application")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				application, err := s.applications.CreateApplication(
					r.Context(),
					project.ID,
					createApplicationRequest.Name,
					createApplicationRequest.Description)
				if err != nil {
					log.WithError(err).Error("create application")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, application)
			},
		)
	})
}

func (s *Service) getApplication(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceApplications, authz.ActionGetApplication,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withApplication(w, r, project, func(application *models.Application) {
					var ret interface{} = application
					if _, ok := r.URL.Query()["full"]; ok {
						latestRelease, err := s.releases.GetLatestRelease(r.Context(), project.ID, application.ID)
						if err != nil && err != store.ErrReleaseNotFound {
							log.WithError(err).Error("get latest release")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						applicationDeviceCounts, err := s.applicationDeviceCounts.GetApplicationDeviceCounts(r.Context(), project.ID, application.ID)
						if err != nil {
							log.WithError(err).Error("get application device counts")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						serviceStateCounts, err := s.deviceServiceStates.ListApplicationServiceStateCounts(
							r.Context(), project.ID, application.ID)
						if err != nil {
							log.WithError(err).Error("get device service states")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						ret = models.ApplicationFull1{
							Application:        *application,
							LatestRelease:      latestRelease,
							DeviceCounts:       *applicationDeviceCounts,
							ServiceStateCounts: serviceStateCounts,
						}
					}

					utils.Respond(w, ret)
				})
			},
		)
	})
}

func (s *Service) listApplications(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceApplications, authz.ActionListApplications,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				applications, err := s.applications.ListApplications(r.Context(), project.ID)
				if err != nil {
					log.WithError(err).Error("list applications")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				var ret interface{} = applications
				if _, ok := r.URL.Query()["full"]; ok {
					applicationsFull := make([]models.ApplicationFull1, 0)

					for _, application := range applications {
						latestRelease, err := s.releases.GetLatestRelease(r.Context(), project.ID, application.ID)
						if err != nil && err != store.ErrReleaseNotFound {
							log.WithError(err).Error("get latest release")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						applicationDeviceCounts, err := s.applicationDeviceCounts.GetApplicationDeviceCounts(r.Context(), project.ID, application.ID)
						if err != nil {
							log.WithError(err).Error("get application device counts")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						serviceStateCounts, err := s.deviceServiceStates.ListApplicationServiceStateCounts(
							r.Context(), project.ID, application.ID)
						if err != nil {
							log.WithError(err).Error("get device service states")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						applicationsFull = append(applicationsFull, models.ApplicationFull1{
							Application:        application,
							LatestRelease:      latestRelease,
							DeviceCounts:       *applicationDeviceCounts,
							ServiceStateCounts: serviceStateCounts,
						})
					}

					ret = applicationsFull
				}

				utils.Respond(w, ret)
			},
		)
	})
}

func (s *Service) updateApplication(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceApplications, authz.ActionUpdateApplication,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withApplication(w, r, project, func(application *models.Application) {
					var updateApplicationRequest struct {
						Name                  *string                                 `json:"name" validate:"name,omitempty"`
						Description           *string                                 `json:"description" validate:"description,omitempty"`
						SchedulingRule        *models.SchedulingRule                  `json:"schedulingRule"`
						MetricEndpointConfigs *map[string]models.MetricEndpointConfig `json:"metricEndpointConfigs"`
					}
					if err := read(r, &updateApplicationRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					var app *models.Application
					var err error
					if updateApplicationRequest.Name != nil {
						if app, err = s.applications.LookupApplication(r.Context(),
							*updateApplicationRequest.Name, project.ID); err == nil && app.ID != application.ID {
							http.Error(w, store.ErrApplicationNameAlreadyInUse.Error(), http.StatusBadRequest)
							return
						} else if err != nil && err != store.ErrApplicationNotFound {
							log.WithError(err).Error("lookup application")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						if app, err = s.applications.UpdateApplicationName(r.Context(), application.ID, project.ID, *updateApplicationRequest.Name); err != nil {
							log.WithError(err).Error("update application name")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
					}
					if updateApplicationRequest.Description != nil {
						if app, err = s.applications.UpdateApplicationDescription(r.Context(), application.ID, project.ID, *updateApplicationRequest.Description); err != nil {
							log.WithError(err).Error("update application description")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
					}
					if updateApplicationRequest.SchedulingRule != nil {
						validationErr, err := scheduling.ValidateSchedulingRule(
							*updateApplicationRequest.SchedulingRule,
							func(releaseID string) (bool, error) {
								r, err := utils.GetReleaseByIdentifier(s.releases, r.Context(), project.ID, application.ID, releaseID)
								if err != nil {
									return false, err
								}
								return r != nil, nil
							},
						)
						if validationErr != nil {
							http.Error(w, validationErr.Error(), http.StatusBadRequest)
							return
						}
						if err != nil {
							log.WithError(err).Error("check application release")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						if app, err = s.applications.UpdateApplicationSchedulingRule(r.Context(), application.ID, project.ID, *updateApplicationRequest.SchedulingRule); err != nil {
							log.WithError(err).Error("update application scheduling rule")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
					}
					if updateApplicationRequest.MetricEndpointConfigs != nil {
						if app, err = s.applications.UpdateApplicationMetricEndpointConfigs(r.Context(), application.ID, project.ID, *updateApplicationRequest.MetricEndpointConfigs); err != nil {
							log.WithError(err).Error("update application service metrics config")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
					}

					utils.Respond(w, app)
				})
			},
		)
	})
}

func (s *Service) deleteApplication(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceApplications, authz.ActionDeleteApplication,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withApplication(w, r, project, func(application *models.Application) {
					if err := s.applications.DeleteApplication(r.Context(), application.ID, project.ID); err != nil {
						log.WithError(err).Error("delete application")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

// TODO: this has a vulnerability!
func (s *Service) createRelease(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceReleases, authz.ActionCreateRelease,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withApplication(w, r, project, func(application *models.Application) {
					var createReleaseRequest models.CreateReleaseRequest
					if err := read(r, &createReleaseRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					if err := spec.Validate([]byte(createReleaseRequest.RawConfig)); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					var applicationConfig map[string]models.Service
					if err := yaml.UnmarshalStrict([]byte(createReleaseRequest.RawConfig), &applicationConfig); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					jsonApplicationConfig, err := json.Marshal(applicationConfig)
					if err != nil {
						log.WithError(err).Error("marshal json application config")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					var userID string
					if user != nil {
						userID = user.ID
					}
					var serviceAccountID string
					if serviceAccount != nil {
						serviceAccountID = serviceAccount.ID
					}
					release, err := s.releases.CreateRelease(
						r.Context(),
						project.ID,
						application.ID,
						createReleaseRequest.RawConfig,
						string(jsonApplicationConfig),
						userID,
						serviceAccountID,
					)
					if err != nil {
						log.WithError(err).Error("create release")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, release)
				})
			},
		)
	})
}

func (s *Service) getRelease(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceReleases, authz.ActionGetRelease,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withApplication(w, r, project, func(application *models.Application) {
					s.withRelease(w, r, project, application, func(release *models.Release) {
						var ret interface{} = release
						var err error
						if _, ok := r.URL.Query()["full"]; ok {
							ret, err = s.getReleaseFull(r.Context(), *release)
							if err != nil {
								log.WithError(err).Error("get release full")
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
						}

						utils.Respond(w, ret)
					})
				})
			},
		)
	})
}

func (s *Service) getLatestRelease(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceReleases, authz.ActionGetLatestRelease,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withApplication(w, r, project, func(application *models.Application) {
					release, err := s.releases.GetLatestRelease(r.Context(), project.ID, application.ID)
					if err == store.ErrReleaseNotFound {
						http.Error(w, err.Error(), http.StatusNotFound)
						return
					} else if err != nil {
						log.WithError(err).Error("get latest release")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					var ret interface{} = release
					if _, ok := r.URL.Query()["full"]; ok {
						ret, err = s.getReleaseFull(r.Context(), *release)
						if err != nil {
							log.WithError(err).Error("get release full")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
					}

					utils.Respond(w, ret)
				})
			},
		)
	})
}

func (s *Service) listReleases(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceReleases, authz.ActionListReleases,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withApplication(w, r, project, func(application *models.Application) {
					releases, err := s.releases.ListReleases(r.Context(), project.ID, application.ID)
					if err != nil {
						log.WithError(err).Error("list releases")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					var ret interface{} = releases
					if _, ok := r.URL.Query()["full"]; ok {
						releasesFull := make([]models.ReleaseFull, 0)
						for _, release := range releases {
							releaseFull, err := s.getReleaseFull(r.Context(), release)
							if err != nil {
								log.WithError(err).Error("get release full")
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
							releasesFull = append(releasesFull, *releaseFull)
						}
						ret = releasesFull
					}

					utils.Respond(w, ret)
				})
			},
		)
	})
}

func (s *Service) listDevices(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionListDevices,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				searchQuery := r.URL.Query().Get("search")

				devices, err := s.devices.ListDevices(r.Context(), project.ID, searchQuery)
				if err != nil {
					log.WithError(err).Error("list devices")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Header().Set("Total-Device-Count", strconv.Itoa(len(devices)))

				filters, err := query.FiltersFromQuery(r.URL.Query())
				if err != nil {
					http.Error(w, errors.Wrap(err, "get filters from query").Error(), http.StatusBadRequest)
					return
				}

				if len(filters) != 0 {
					appStatuses, err := s.deviceApplicationStatuses.ListAllDeviceApplicationStatuses(r.Context(), project.ID)
					if err != nil {
						http.Error(w, errors.Wrap(err, "get filter dependencies").Error(), http.StatusBadRequest)
						return
					}
					appStatusMap, err := utils.DeviceApplicationStatusesListToMap(appStatuses)
					if err != nil {
						http.Error(w, errors.Wrap(err, "get filter dependencies").Error(), http.StatusBadRequest)
						return
					}

					serviceStates, err := s.deviceServiceStates.ListAllDeviceServiceStates(r.Context(), project.ID)
					if err != nil {
						http.Error(w, errors.Wrap(err, "get filter dependencies").Error(), http.StatusBadRequest)
						return
					}
					serviceStateMap, err := utils.DeviceServiceStatesListToMap(serviceStates)
					if err != nil {
						http.Error(w, errors.Wrap(err, "get filter dependencies").Error(), http.StatusBadRequest)
						return
					}

					devices, _, err = query.QueryDevices(
						query.QueryDependencies{
							DeviceApplicationStatuses: appStatusMap,
							DeviceServiceStates:       serviceStateMap,
							Releases:                  s.releases,
							Context:                   r.Context(),
						},
						devices,
						filters,
					)
					if err != nil {
						http.Error(w, errors.Wrap(err, "filter devices").Error(), http.StatusBadRequest)
						return
					}
				}

				ds := make([]interface{}, len(devices))
				for i := range devices {
					ds[i] = devices[i]
				}

				middleware.SortAndPaginateAndRespond(*r, w, ds)
			},
		)
	})
}

func (s *Service) previewScheduledDevices(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionPreviewApplicationScheduling,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				searchQuery := r.URL.Query().Get("search")

				devices, err := s.devices.ListDevices(r.Context(), project.ID, searchQuery)
				if err != nil {
					log.WithError(err).Error("list devices")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				filters, err := query.FiltersFromQuery(r.URL.Query())
				if err != nil {
					http.Error(w, errors.Wrap(err, "get filters from query").Error(), http.StatusBadRequest)
					return
				}

				if len(filters) != 0 {
					devices, _, err = query.QueryDevices(query.QueryDependencies{}, devices, filters)
					if err != nil {
						http.Error(w, errors.Wrap(err, "filter devices").Error(), http.StatusBadRequest)
						return
					}
				}

				schedulingRule, err := scheduling.SchedulingRuleFromQuery(r.URL.Query())
				if schedulingRule == nil && err == nil {
					err = scheduling.ErrNonexistentSchedulingRule
				}
				if err != nil {
					http.Error(w, errors.Wrap(err, "get scheduling rule from query").Error(), http.StatusBadRequest)
					return
				}

				scheduledDevices, err := scheduling.GetScheduledDevices(devices, *schedulingRule)
				if err != nil {
					http.Error(w, errors.Wrap(err, "preview scheduling rule").Error(), http.StatusBadRequest)
					return
				}

				utils.Respond(w, scheduledDevices)
			},
		)
	})
}

func (s *Service) getDevice(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionGetDevice,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					var ret interface{} = device
					if _, ok := r.URL.Query()["full"]; ok {
						applications, err := s.applications.ListApplications(r.Context(), project.ID)
						if err != nil {
							log.WithError(err).Error("list applications")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						allApplicationStatusInfo := make([]models.DeviceApplicationStatusInfo, 0)
						for _, application := range applications {
							applicationStatusInfo := models.DeviceApplicationStatusInfo{
								Application:     application,
								ServiceStatuses: []models.DeviceServiceStatusFull{},
								ServiceStates:   []models.DeviceServiceState{},
							}

							deviceApplicationStatus, err := s.deviceApplicationStatuses.GetDeviceApplicationStatus(
								r.Context(), project.ID, device.ID, application.ID)
							if err == nil {
								currentRelease, err := s.releases.GetRelease(r.Context(),
									deviceApplicationStatus.CurrentReleaseID,
									deviceApplicationStatus.ProjectID,
									deviceApplicationStatus.ApplicationID,
								)
								if err != nil {
									log.WithError(err).Error("get release")
									w.WriteHeader(http.StatusInternalServerError)
									return
								}

								applicationStatusInfo.ApplicationStatus = &models.DeviceApplicationStatusFull{
									DeviceApplicationStatus: *deviceApplicationStatus,
									CurrentRelease:          *currentRelease,
								}
							} else if err != store.ErrDeviceApplicationStatusNotFound {
								log.WithError(err).Error("get device application status")
								w.WriteHeader(http.StatusInternalServerError)
								return
							}

							deviceServiceStatuses, err := s.deviceServiceStatuses.GetDeviceServiceStatuses(
								r.Context(), project.ID, device.ID, application.ID)
							if err == nil {
								for _, deviceServiceStatus := range deviceServiceStatuses {
									currentRelease, err := s.releases.GetRelease(r.Context(),
										deviceServiceStatus.CurrentReleaseID,
										deviceServiceStatus.ProjectID,
										deviceServiceStatus.ApplicationID,
									)
									if err != nil {
										log.WithError(err).Error("get release")
										w.WriteHeader(http.StatusInternalServerError)
										return
									}

									applicationStatusInfo.ServiceStatuses = append(applicationStatusInfo.ServiceStatuses, models.DeviceServiceStatusFull{
										DeviceServiceStatus: deviceServiceStatus,
										CurrentRelease:      *currentRelease,
									})
								}
							} else {
								log.WithError(err).Error("get device service statuses")
								w.WriteHeader(http.StatusInternalServerError)
								return
							}

							deviceServiceStates, err := s.deviceServiceStates.GetDeviceServiceStates(
								r.Context(), project.ID, device.ID, application.ID)
							if err == nil {
								applicationStatusInfo.ServiceStates = deviceServiceStates
							} else {
								log.WithError(err).Error("get device service states")
								w.WriteHeader(http.StatusInternalServerError)
								return
							}

							allApplicationStatusInfo = append(allApplicationStatusInfo, applicationStatusInfo)
						}

						ret = models.DeviceFull{
							Device:                *device,
							ApplicationStatusInfo: allApplicationStatusInfo,
						}
					}

					utils.Respond(w, ret)
				})
			},
		)
	})
}

func (s *Service) updateDevice(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionUpdateDevice,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					var updateDeviceRequest struct {
						Name string `json:"name" validate:"name"`
					}
					if err := read(r, &updateDeviceRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					if d, err := s.devices.LookupDevice(r.Context(),
						updateDeviceRequest.Name, project.ID); err == nil && d.ID != device.ID {
						http.Error(w, store.ErrDeviceNameAlreadyInUse.Error(), http.StatusBadRequest)
						return
					} else if err != nil && err != store.ErrDeviceNotFound {
						log.WithError(err).Error("lookup device")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					d, err := s.devices.UpdateDeviceName(r.Context(), device.ID, project.ID, updateDeviceRequest.Name)
					if err != nil {
						log.WithError(err).Error("update device name")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, d)
				})
			},
		)
	})
}

func (s *Service) deleteDevice(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionDeleteDevice,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					if err := s.devices.DeleteDevice(r.Context(), device.ID, project.ID); err != nil {
						log.WithError(err).Error("delete device")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) setDeviceEnvironmentVariable(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceEnvironmentVariables, authz.ActionSetDeviceEnvironmentVariable,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					var setDeviceEnvironmentVariableRequest struct {
						Key   string `json:"key" validate:"environmentvariablekey"`
						Value string `json:"value" validate:"environmentvariablevalue"`
					}
					if err := read(r, &setDeviceEnvironmentVariableRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					deviceEnvironmentVariable, err := s.devices.SetDeviceEnvironmentVariable(
						r.Context(),
						device.ID,
						project.ID,
						setDeviceEnvironmentVariableRequest.Key,
						setDeviceEnvironmentVariableRequest.Value,
					)
					if err != nil {
						log.WithError(err).Error("set device environment variable")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, deviceEnvironmentVariable)
				})
			},
		)
	})
}

func (s *Service) deleteDeviceEnvironmentVariable(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceEnvironmentVariables, authz.ActionDeleteDeviceEnvironmentVariable,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					vars := mux.Vars(r)
					key := vars["key"]

					if err := s.devices.DeleteDeviceEnvironmentVariable(r.Context(), device.ID, project.ID, key); err != nil {
						log.WithError(err).Error("delete device environment variable")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) listAllDeviceLabelKeys(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceLabels, authz.ActionListAllDeviceLabels,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				deviceLabels, err := s.devices.ListAllDeviceLabelKeys(
					r.Context(),
					project.ID,
				)
				if err != nil {
					log.WithError(err).Error("list device labels")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, deviceLabels)
			},
		)
	})
}

func (s *Service) setDeviceLabel(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceLabels, authz.ActionSetDeviceLabel,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					var setDeviceLabelRequest struct {
						Key   string `json:"key" validate:"labelkey"`
						Value string `json:"value" validate:"labelvalue"`
					}
					if err := read(r, &setDeviceLabelRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					deviceLabel, err := s.devices.SetDeviceLabel(
						r.Context(),
						device.ID,
						project.ID,
						setDeviceLabelRequest.Key,
						setDeviceLabelRequest.Value,
					)
					if err != nil {
						log.WithError(err).Error("set device label")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, deviceLabel)
				})
			},
		)
	})
}

func (s *Service) deleteDeviceLabel(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceLabels, authz.ActionDeleteDeviceLabel,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					vars := mux.Vars(r)
					key := vars["key"]

					if err := s.devices.DeleteDeviceLabel(r.Context(), device.ID, project.ID, key); err != nil {
						log.WithError(err).Error("delete device label")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) createDeviceRegistrationToken(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceRegistrationTokens, authz.ActionCreateDeviceRegistrationToken,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				var createDeviceRegistrationTokenRequest struct {
					Name             string `json:"name" validate:"name"`
					Description      string `json:"description" validate:"description"`
					MaxRegistrations *int   `json:"maxRegistrations"`
				}
				if err := read(r, &createDeviceRegistrationTokenRequest); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if _, err := s.deviceRegistrationTokens.LookupDeviceRegistrationToken(
					r.Context(),
					createDeviceRegistrationTokenRequest.Name,
					project.ID,
				); err == nil {
					http.Error(w, store.ErrDeviceRegistrationTokenNameAlreadyInUse.Error(), http.StatusBadRequest)
					return
				} else if err != nil && err != store.ErrDeviceRegistrationTokenNotFound {
					log.WithError(err).Error("lookup device registration token")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				deviceRegistrationToken, err := s.deviceRegistrationTokens.CreateDeviceRegistrationToken(
					r.Context(),
					project.ID,
					createDeviceRegistrationTokenRequest.Name,
					createDeviceRegistrationTokenRequest.Description,
					createDeviceRegistrationTokenRequest.MaxRegistrations)
				if err != nil {
					log.WithError(err).Error("create device registration token")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, deviceRegistrationToken)
			},
		)
	})
}

func (s *Service) getDeviceRegistrationToken(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceRegistrationTokens, authz.ActionUpdateDeviceRegistrationToken,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDeviceRegistrationToken(w, r, project, func(deviceRegistrationToken *models.DeviceRegistrationToken) {
					var ret interface{} = deviceRegistrationToken
					if _, ok := r.URL.Query()["full"]; ok {
						devicesRegistered, err := s.devicesRegisteredWithToken.GetDevicesRegisteredWithTokenCount(r.Context(), deviceRegistrationToken.ID, project.ID)
						if err != nil {
							log.WithError(err).Error("get registered device count")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						ret = models.DeviceRegistrationTokenFull{
							DeviceRegistrationToken: *deviceRegistrationToken,
							DeviceCounts:            *devicesRegistered,
						}
					}

					utils.Respond(w, ret)
				})
			},
		)
	})
}

func (s *Service) updateDeviceRegistrationToken(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceLabels, authz.ActionDeleteDeviceLabel,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDeviceRegistrationToken(w, r, project, func(deviceRegistrationToken *models.DeviceRegistrationToken) {
					var updateDeviceRegistrationTokenRequest struct {
						Name             string `json:"name" validate:"name"`
						Description      string `json:"description" validate:"description"`
						MaxRegistrations *int   `json:"maxRegistrations"`
					}
					if err := read(r, &updateDeviceRegistrationTokenRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					token, err := s.deviceRegistrationTokens.UpdateDeviceRegistrationToken(
						r.Context(),
						deviceRegistrationToken.ID,
						project.ID,
						updateDeviceRegistrationTokenRequest.Name,
						updateDeviceRegistrationTokenRequest.Description,
						updateDeviceRegistrationTokenRequest.MaxRegistrations,
					)
					if err != nil {
						log.WithError(err).Error("update device registration token")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, token)
				})
			},
		)
	})
}

func (s *Service) deleteDeviceRegistrationToken(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceRegistrationTokens, authz.ActionDeleteDeviceRegistrationToken,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDeviceRegistrationToken(w, r, project, func(deviceRegistrationToken *models.DeviceRegistrationToken) {
					if err := s.deviceRegistrationTokens.DeleteDeviceRegistrationToken(r.Context(), deviceRegistrationToken.ID, project.ID); err != nil {
						log.WithError(err).Error("delete device registration token")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) listDeviceRegistrationTokens(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceRegistrationTokens, authz.ActionListDeviceRegistrationTokens,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				deviceRegistrationTokens, err := s.deviceRegistrationTokens.ListDeviceRegistrationTokens(r.Context(), project.ID)
				if err != nil {
					log.WithError(err).Error("list device registration tokens")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				var ret interface{} = deviceRegistrationTokens
				if _, ok := r.URL.Query()["full"]; ok {
					deviceRegistrationTokensFull := make([]models.DeviceRegistrationTokenFull, 0)

					for _, deviceRegistrationToken := range deviceRegistrationTokens {
						deviceRegistrationTokenCounts, err := s.devicesRegisteredWithToken.GetDevicesRegisteredWithTokenCount(r.Context(), deviceRegistrationToken.ID, project.ID)
						if err != nil {
							log.WithError(err).Error("get count of devices registered with device registration token")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						deviceRegistrationTokensFull = append(deviceRegistrationTokensFull, models.DeviceRegistrationTokenFull{
							DeviceRegistrationToken: deviceRegistrationToken,
							DeviceCounts:            *deviceRegistrationTokenCounts,
						})
					}

					ret = deviceRegistrationTokensFull
				}

				utils.Respond(w, ret)
			},
		)
	})
}

func (s *Service) setDeviceRegistrationTokenEnvironmentVariable(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceRegistrationTokenEnvironmentVariables, authz.ActionSetDeviceRegistrationTokenEnvironmentVariable,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDeviceRegistrationToken(w, r, project, func(deviceRegistrationToken *models.DeviceRegistrationToken) {
					var setDeviceRegistrationTokenEnvironmentVariableRequest struct {
						Key   string `json:"key" validate:"environmentvariablekey"`
						Value string `json:"value" validate:"environmentvariablevalue"`
					}
					if err := read(r, &setDeviceRegistrationTokenEnvironmentVariableRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					deviceRegistrationTokenEnvironmentVariable, err := s.deviceRegistrationTokens.SetDeviceRegistrationTokenEnvironmentVariable(
						r.Context(),
						deviceRegistrationToken.ID,
						project.ID,
						setDeviceRegistrationTokenEnvironmentVariableRequest.Key,
						setDeviceRegistrationTokenEnvironmentVariableRequest.Value,
					)
					if err != nil {
						log.WithError(err).Error("set device registration token environment variable")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, deviceRegistrationTokenEnvironmentVariable)
				})
			},
		)
	})
}

func (s *Service) deleteDeviceRegistrationTokenEnvironmentVariable(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceRegistrationTokenEnvironmentVariables, authz.ActionDeleteDeviceRegistrationTokenEnvironmentVariable,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDeviceRegistrationToken(w, r, project, func(deviceRegistrationToken *models.DeviceRegistrationToken) {
					vars := mux.Vars(r)
					key := vars["key"]

					if err := s.deviceRegistrationTokens.DeleteDeviceRegistrationTokenEnvironmentVariable(r.Context(),
						deviceRegistrationToken.ID, project.ID, key,
					); err != nil {
						log.WithError(err).Error("delete device registration token environment variable")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) setDeviceRegistrationTokenLabel(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceRegistrationTokenLabels, authz.ActionSetDeviceRegistrationTokenLabel,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDeviceRegistrationToken(w, r, project, func(deviceRegistrationToken *models.DeviceRegistrationToken) {
					var setLabelRequest struct {
						Key   string `json:"key" validate:"labelkey"`
						Value string `json:"value" validate:"labelvalue"`
					}
					if err := read(r, &setLabelRequest); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					label, err := s.deviceRegistrationTokens.SetDeviceRegistrationTokenLabel(
						r.Context(),
						deviceRegistrationToken.ID,
						project.ID,
						setLabelRequest.Key,
						setLabelRequest.Value,
					)
					if err != nil {
						log.WithError(err).Error("set device registration token label")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					utils.Respond(w, label)
				})
			},
		)
	})
}

func (s *Service) deleteDeviceRegistrationTokenLabel(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceRegistrationTokenLabels, authz.ActionDeleteDeviceRegistrationTokenLabel,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDeviceRegistrationToken(w, r, project, func(deviceRegistrationToken *models.DeviceRegistrationToken) {
					vars := mux.Vars(r)
					key := vars["key"]

					if err := s.deviceRegistrationTokens.DeleteDeviceRegistrationTokenLabel(r.Context(), deviceRegistrationToken.ID, project.ID, key); err != nil {
						log.WithError(err).Error("delete device registration token label")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			},
		)
	})
}

func (s *Service) getProjectConfig(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceProjectConfigs, authz.ActionSetProjectConfig,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				vars := mux.Vars(r)
				key := vars["key"]

				var value interface{}
				var err error
				switch key {
				case string(models.ProjectMetricsConfigKey):
					value, err = s.metricConfigs.GetProjectMetricsConfig(r.Context(), project.ID)
				case string(models.DeviceMetricsConfigKey):
					value, err = s.metricConfigs.GetDeviceMetricsConfig(r.Context(), project.ID)
				case string(models.ServiceMetricsConfigKey):
					value, err = s.metricConfigs.GetServiceMetricsConfigs(r.Context(), project.ID)
				default:
					http.Error(w, store.ErrProjectConfigNotFound.Error(), http.StatusBadRequest)
					return
				}

				if err != nil {
					log.WithError(err).Error("get project config with key " + key)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				utils.Respond(w, value)
			},
		)
	})
}

func (s *Service) setProjectConfig(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDeviceLabels, authz.ActionDeleteDeviceLabel,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				vars := mux.Vars(r)
				key := vars["key"]

				var err error
				switch key {
				case string(models.ProjectMetricsConfigKey):
					var value models.ProjectMetricsConfig
					if err := read(r, &value); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					err = s.metricConfigs.SetProjectMetricsConfig(r.Context(), project.ID, value)
				case string(models.DeviceMetricsConfigKey):
					var value models.DeviceMetricsConfig
					if err := read(r, &value); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					err = s.metricConfigs.SetDeviceMetricsConfig(r.Context(), project.ID, value)
				case string(models.ServiceMetricsConfigKey):
					var values []models.ServiceMetricsConfig
					// TODO: use read() here
					if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					err = s.metricConfigs.SetServiceMetricsConfigs(r.Context(), project.ID, values)
				default:
					http.Error(w, store.ErrProjectConfigNotFound.Error(), http.StatusBadRequest)
					return
				}

				if err != nil {
					log.WithError(err).Error("set project config with key " + key)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(200)
			},
		)
	})
}

// TODO: verify project ID
func (s *Service) registerDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]

	var registerDeviceRequest models.RegisterDeviceRequest
	if err := read(r, &registerDeviceRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deviceRegistrationToken, err := s.deviceRegistrationTokens.GetDeviceRegistrationToken(r.Context(), registerDeviceRequest.DeviceRegistrationTokenID, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if deviceRegistrationToken.MaxRegistrations != nil {
		devicesRegisteredCount, err := s.devicesRegisteredWithToken.GetDevicesRegisteredWithTokenCount(r.Context(), registerDeviceRequest.DeviceRegistrationTokenID, projectID)
		if err != nil {
			log.WithError(err).Error("get devices registered with token count")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if devicesRegisteredCount.AllCount >= *deviceRegistrationToken.MaxRegistrations {
			log.WithError(err).Error("device allocation limit reached")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	device, err := s.devices.CreateDevice(r.Context(),
		projectID, namesgenerator.GetRandomName(), deviceRegistrationToken.ID,
		deviceRegistrationToken.Labels, deviceRegistrationToken.EnvironmentVariables,
	)
	if err != nil {
		log.WithError(err).Error("create device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deviceAccessKeyValue := ksuid.New().String()

	_, err = s.deviceAccessKeys.CreateDeviceAccessKey(r.Context(), projectID, device.ID, hash.Hash(deviceAccessKeyValue))
	if err != nil {
		log.WithError(err).Error("create device access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, models.RegisterDeviceResponse{
		DeviceID:             device.ID,
		DeviceAccessKeyValue: deviceAccessKeyValue,
	})
}

func (s *Service) getBundle(w http.ResponseWriter, r *http.Request) {
	s.withDeviceAuth(w, r, func(project *models.Project, device *models.Device) {
		s.st.Incr("get_bundle",
			utils.WithTags(
				[]string{},
				utils.TagItems{Project: project},
			),
			1,
		)

		if err := s.devices.UpdateDeviceLastSeenAt(r.Context(), device.ID, project.ID); err != nil {
			log.WithError(err).Error("update device last seen at")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		applications, err := s.applications.ListApplications(r.Context(), project.ID)
		if err != nil {
			log.WithError(err).Error("list applications")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bundle := models.Bundle{
			DeviceID:             device.ID,
			DeviceName:           device.Name,
			EnvironmentVariables: device.EnvironmentVariables,
			DesiredAgentVersion:  device.DesiredAgentVersion,
		}

		for _, application := range applications {
			scheduled, scheduledDevice, err := scheduling.IsApplicationScheduled(*device, application.SchedulingRule)
			if err != nil {
				log.WithError(err).Error("evaluate application scheduling rule")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !scheduled {
				continue
			}

			release, err := utils.GetReleaseByIdentifier(s.releases, r.Context(), project.ID, application.ID, scheduledDevice.ReleaseID)
			if err == store.ErrReleaseNotFound {
				continue
			}
			if err != nil {
				log.WithError(err).Errorf("get release by ID %s", scheduledDevice.ReleaseID)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			bundle.Applications = append(bundle.Applications, models.FullBundledApplication{
				Application: models.BundledApplication{
					ID:                    application.ID,
					ProjectID:             application.ProjectID,
					Name:                  application.Name,
					MetricEndpointConfigs: application.MetricEndpointConfigs,
				},
				LatestRelease: *release,
			})
		}

		deviceApplicationStatuses, err := s.deviceApplicationStatuses.ListDeviceApplicationStatuses(
			r.Context(), project.ID, device.ID)
		if err == nil {
			bundle.ApplicationStatuses = deviceApplicationStatuses
		} else {
			log.WithError(err).Error("list device application statuses")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		deviceServiceStatuses, err := s.deviceServiceStatuses.ListDeviceServiceStatuses(
			r.Context(), project.ID, device.ID)
		if err == nil {
			bundle.ServiceStatuses = deviceServiceStatuses
		} else {
			log.WithError(err).Error("list device service statuses")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		deviceServiceStates, err := s.deviceServiceStates.ListDeviceServiceStates(
			r.Context(), project.ID, device.ID)
		if err == nil {
			bundle.ServiceStates = deviceServiceStates
		} else {
			log.WithError(err).Error("list device service states")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, bundle)
	})
}

func (s *Service) setDeviceInfo(w http.ResponseWriter, r *http.Request) {
	s.withDeviceAuth(w, r, func(project *models.Project, device *models.Device) {
		var setDeviceInfoRequest models.SetDeviceInfoRequest
		if err := read(r, &setDeviceInfoRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := s.devices.SetDeviceInfo(r.Context(), device.ID, project.ID, setDeviceInfoRequest.DeviceInfo); err != nil {
			log.WithError(err).Error("set device info")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) setDeviceApplicationStatus(w http.ResponseWriter, r *http.Request) {
	s.withDeviceAuth(w, r, func(project *models.Project, device *models.Device) {
		vars := mux.Vars(r)
		applicationID := vars["application"]

		var setDeviceApplicationStatusRequest models.SetDeviceApplicationStatusRequest
		if err := read(r, &setDeviceApplicationStatusRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.deviceApplicationStatuses.SetDeviceApplicationStatus(r.Context(), project.ID, device.ID,
			applicationID, setDeviceApplicationStatusRequest.CurrentReleaseID,
		); err != nil {
			log.WithError(err).Error("set device application status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) deleteDeviceApplicationStatus(w http.ResponseWriter, r *http.Request) {
	s.withDeviceAuth(w, r, func(project *models.Project, device *models.Device) {
		vars := mux.Vars(r)
		applicationID := vars["application"]

		if err := s.deviceApplicationStatuses.DeleteDeviceApplicationStatus(r.Context(),
			project.ID, device.ID, applicationID,
		); err != nil {
			log.WithError(err).Error("delete device application status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) setDeviceServiceStatus(w http.ResponseWriter, r *http.Request) {
	s.withDeviceAuth(w, r, func(project *models.Project, device *models.Device) {
		vars := mux.Vars(r)
		applicationID := vars["application"]
		service := vars["service"]

		var setDeviceServiceStatusRequest models.SetDeviceServiceStatusRequest
		if err := read(r, &setDeviceServiceStatusRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.deviceServiceStatuses.SetDeviceServiceStatus(r.Context(), project.ID, device.ID,
			applicationID, service, setDeviceServiceStatusRequest.CurrentReleaseID,
		); err != nil {
			log.WithError(err).Error("set device service status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) deleteDeviceServiceStatus(w http.ResponseWriter, r *http.Request) {
	s.withDeviceAuth(w, r, func(project *models.Project, device *models.Device) {
		vars := mux.Vars(r)
		applicationID := vars["application"]
		service := vars["service"]

		if err := s.deviceServiceStatuses.DeleteDeviceServiceStatus(r.Context(),
			project.ID, device.ID, applicationID, service,
		); err != nil {
			log.WithError(err).Error("delete device service status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) setDeviceServiceState(w http.ResponseWriter, r *http.Request) {
	s.withDeviceAuth(w, r, func(project *models.Project, device *models.Device) {
		vars := mux.Vars(r)
		applicationID := vars["application"]
		service := vars["service"]

		var setDeviceServiceStateRequest models.SetDeviceServiceStateRequest
		if err := read(r, &setDeviceServiceStateRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.deviceServiceStates.SetDeviceServiceState(r.Context(), project.ID, device.ID,
			applicationID,
			service,
			setDeviceServiceStateRequest.State,
			setDeviceServiceStateRequest.ErrorMessage,
		); err != nil {
			log.WithError(err).Error("set device service state")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) deleteDeviceServiceState(w http.ResponseWriter, r *http.Request) {
	s.withDeviceAuth(w, r, func(project *models.Project, device *models.Device) {
		vars := mux.Vars(r)
		applicationID := vars["application"]
		service := vars["service"]

		if err := s.deviceServiceStates.DeleteDeviceServiceState(r.Context(),
			project.ID, device.ID, applicationID, service,
		); err != nil {
			log.WithError(err).Error("delete device service state")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
