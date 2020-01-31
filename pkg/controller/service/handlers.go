package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) health(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.st.Incr("health", nil, 1)
	})
}

func (s *Service) intentional500(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
}

func (s *Service) register(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.WithReferrer(w, r, func(referrer *url.URL) {
			var registerRequest struct {
				Email     string `json:"email" validate:"email"`
				Password  string `json:"password" validate:"password"`
				FirstName string `json:"firstName" validate:"required,min=1,max=100"`
				LastName  string `json:"lastName" validate:"required,min=1,max=100"`
				Company   string `json:"company" validate:"max=100"`
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

			if _, err := s.users.LookupUser(r.Context(), registerRequest.Email); err == nil {
				http.Error(w, errEmailAlreadyTaken.Error(), http.StatusBadRequest)
				return
			} else if err != store.ErrUserNotFound {
				log.WithError(err).Error("lookup user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			user, err := s.users.CreateUser(r.Context(), registerRequest.Email, hash.Hash(registerRequest.Password),
				registerRequest.FirstName, registerRequest.LastName, registerRequest.Company)
			if err != nil {
				log.WithError(err).Error("create user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if s.email == nil {
				// If email provider is nil then skip the registration workflow
				if _, err := s.users.MarkRegistrationCompleted(r.Context(), user.ID); err != nil {
					log.WithError(err).Error("mark registration completed")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {
				registrationTokenValue := ksuid.New().String()

				if _, err := s.registrationTokens.CreateRegistrationToken(r.Context(), user.ID, hash.Hash(registrationTokenValue)); err != nil {
					log.WithError(err).Error("create registration token")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				content := fmt.Sprintf(
					"Please go to the following URL to complete registration. %s://%s/confirm/%s",
					referrer.Scheme,
					referrer.Host,
					registrationTokenValue,
				)

				if err := s.email.Send(email.Request{
					FromName:    s.emailFromName,
					FromAddress: s.emailFromAddress,
					ToName:      user.FirstName + " " + user.LastName,
					ToAddress:   user.Email,
					Subject:     "Deviceplane Registration Confirmation",
					Body:        content,
				}); err != nil {
					log.WithError(err).Error("send registration email")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			utils.Respond(w, user)
		})
	})
}

func (s *Service) confirmRegistration(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.WithReferrer(w, r, func(referrer *url.URL) {
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

			if _, err := s.users.MarkRegistrationCompleted(r.Context(), registrationToken.UserID); err != nil {
				log.WithError(err).Error("mark registration completed")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			user, err := s.users.GetUser(r.Context(), registrationToken.UserID)
			if err != nil {
				log.WithError(err).Error("mark registration completed")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			with.User = user
			s.newSession(with)(w, r)
		})
	})
}

func (s *Service) recoverPassword(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.WithReferrer(w, r, func(referrer *url.URL) {
			var recoverPasswordRequest struct {
				Email string `json:"email" validate:"email"`
			}
			if err := read(r, &recoverPasswordRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			user, err := s.users.LookupUser(r.Context(), recoverPasswordRequest.Email)
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
				user.ID, hash.Hash(passwordRecoveryTokenValue)); err != nil {
				log.WithError(err).Error("create password recovery token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			content := fmt.Sprintf(
				"Please go to the following URL to recover your password. %s://%s/recover/%s",
				referrer.Scheme,
				referrer.Host,
				passwordRecoveryTokenValue,
			)

			if err := s.email.Send(email.Request{
				FromName:    s.emailFromName,
				FromAddress: s.emailFromAddress,
				ToName:      user.FirstName + " " + user.LastName,
				ToAddress:   user.Email,
				Subject:     "Deviceplane Password Recovery",
				Body:        content,
			}); err != nil {
				log.WithError(err).Error("send recovery email")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})
	})
}

func (s *Service) getPasswordRecoveryToken(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func (s *Service) changePassword(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		if _, err := s.users.UpdatePasswordHash(r.Context(),
			passwordRecoveryToken.UserID, hash.Hash(changePasswordRequest.Password)); err != nil {
			log.WithError(err).Error("update password hash")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := s.users.GetUser(r.Context(), passwordRecoveryToken.UserID)
		if err != nil {
			log.WithError(err).Error("get user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, user)
	})
}

func (s *Service) login(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var loginRequest struct {
			Email    string `json:"email" validate:"email"`
			Password string `json:"password" validate:"password"`
		}
		if err := read(r, &loginRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := s.users.ValidateUserWithEmail(r.Context(), loginRequest.Email, hash.Hash(loginRequest.Password))
		if err == store.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("validate user with email")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !user.RegistrationCompleted {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		with.User = user
		s.newSession(with)(w, r)
	})
}

func (s *Service) newSession(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.WithReferrer(w, r, func(referrer *url.URL) {
			sessionValue := ksuid.New().String()

			if _, err := s.sessions.CreateSession(r.Context(), with.User.ID, hash.Hash(sessionValue)); err != nil {
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
	})
}

func (s *Service) logout(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func (s *Service) getMe(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		if with.User == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, with.User)
	})
}

func (s *Service) updateMe(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if with.User == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO: validation
		var updateUserRequest struct {
			Password        *string `json:"password"`
			CurrentPassword *string `json:"currentPassword"`
			FirstName       *string `json:"firstName"`
			LastName        *string `json:"lastName"`
			Company         *string `json:"company"`
		}
		if err := read(r, &updateUserRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if updateUserRequest.Password != nil {
			if updateUserRequest.CurrentPassword == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if _, err := s.users.ValidateUser(r.Context(), with.User.ID, hash.Hash(*updateUserRequest.CurrentPassword)); err == store.ErrUserNotFound {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err != nil {
				log.WithError(err).Error("validate user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if _, err := s.users.UpdatePasswordHash(r.Context(), with.User.ID, hash.Hash(*updateUserRequest.Password)); err != nil {
				log.WithError(err).Error("update password hash")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if updateUserRequest.FirstName != nil {
			if _, err := s.users.UpdateFirstName(r.Context(), with.User.ID, *updateUserRequest.FirstName); err != nil {
				log.WithError(err).Error("update first name")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if updateUserRequest.LastName != nil {
			if _, err := s.users.UpdateLastName(r.Context(), with.User.ID, *updateUserRequest.LastName); err != nil {
				log.WithError(err).Error("update last name")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if updateUserRequest.Company != nil {
			if _, err := s.users.UpdateCompany(r.Context(), with.User.ID, *updateUserRequest.Company); err != nil {
				log.WithError(err).Error("update company")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		user, err := s.users.GetUser(r.Context(), with.User.ID)
		if err != nil {
			log.WithError(err).Error("get user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, user)
	})
}

func (s *Service) listMembershipsByUser(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		if with.User.ID == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		memberships, err := s.memberships.ListMembershipsByUser(r.Context(), with.User.ID)
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

func (s *Service) createUserAccessKey(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		if with.User.ID == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var createUserAccessKeyRequest struct {
			Description string `json:"description" validate:"description"`
		}
		if err := read(r, &createUserAccessKeyRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userAccessKeyValue := "u" + ksuid.New().String()

		user, err := s.userAccessKeys.CreateUserAccessKey(r.Context(),
			with.User.ID, hash.Hash(userAccessKeyValue), createUserAccessKeyRequest.Description)
		if err != nil {
			log.WithError(err).Error("create user access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, models.UserAccessKeyWithValue{
			UserAccessKey: *user,
			Value:         userAccessKeyValue,
		})
	})
}

// TODO: verify that the user owns this access key
func (s *Service) getUserAccessKey(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		if with.User.ID == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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

		utils.Respond(w, userAccessKey)
	})
}

func (s *Service) listUserAccessKeys(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		if with.User.ID == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userAccessKeys, err := s.userAccessKeys.ListUserAccessKeys(r.Context(), with.User.ID)
		if err != nil {
			log.WithError(err).Error("list users")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, userAccessKeys)
	})
}

// TODO: verify that the user owns this access key
func (s *Service) deleteUserAccessKey(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userAccessKeyID := vars["useraccesskey"]

		if err := s.userAccessKeys.DeleteUserAccessKey(r.Context(), userAccessKeyID); err != nil {
			log.WithError(err).Error("delete user access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) createProject(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		if _, err = s.memberships.CreateMembership(r.Context(), with.User.ID, with.Project.ID); err != nil {
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

		adminRole, err := s.roles.CreateRole(r.Context(), with.Project.ID, "admin-all", "", string(adminAllRoleBytes))
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

		if _, err = s.roles.CreateRole(r.Context(), with.Project.ID, "write-all", "", string(writeAllRoleBytes)); err != nil {
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

		if _, err = s.roles.CreateRole(r.Context(), with.Project.ID, "read-all", "", string(readAllRoleBytes)); err != nil {
			log.WithError(err).Error("create read role")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = s.membershipRoleBindings.CreateMembershipRoleBinding(r.Context(),
			with.User.ID, adminRole.ID, with.Project.ID,
		); err != nil {
			log.WithError(err).Error("create membership role binding")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create default device registration token.
		// It is named "default" and has an unlimited device registration cap.
		_, err = s.deviceRegistrationTokens.CreateDeviceRegistrationToken(
			r.Context(),
			with.Project.ID,
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

func (s *Service) getProject(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project, err := s.projects.GetProject(r.Context(), with.Project.ID)
		if err != nil {
			log.WithError(err).Error("get project")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, project)
	})
}

func (s *Service) updateProject(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var updateProjectRequest struct {
			Name          string `json:"name" validate:"name"`
			DatadogApiKey string `json:"datadogApiKey"`
		}
		if err := read(r, &updateProjectRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		project, err := s.projects.UpdateProject(r.Context(), with.Project.ID, updateProjectRequest.Name, updateProjectRequest.DatadogApiKey)
		if err != nil {
			log.WithError(err).Error("update project")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, project)
	})
}

func (s *Service) deleteProject(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.projects.DeleteProject(r.Context(), with.Project.ID); err != nil {
			log.WithError(err).Error("delete project")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) createRole(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var createRoleRequest struct {
			Name        string `json:"name" validate:"name"`
			Description string `json:"description" validate:"description"`
			Config      string `json:"config" validate:"config"`
		}
		if err := read(r, &createRoleRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := s.roles.LookupRole(r.Context(), createRoleRequest.Name, with.Project.ID); err == nil {
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

		role, err := s.roles.CreateRole(r.Context(), with.Project.ID, createRoleRequest.Name,
			createRoleRequest.Description, createRoleRequest.Config)
		if err != nil {
			log.WithError(err).Error("create role")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, role)
	})
}

func (s *Service) getRole(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, err := s.roles.GetRole(r.Context(), with.Role.ID, with.Project.ID)
		if err == store.ErrRoleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get role")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, role)
	})
}

func (s *Service) listRoles(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roles, err := s.roles.ListRoles(r.Context(), with.Project.ID)
		if err != nil {
			log.WithError(err).Error("list roles")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, roles)
	})
}

func (s *Service) updateRole(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roleID := vars["role"]

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
			updateRoleRequest.Name, with.Project.ID); err == nil && role.ID != roleID {
			http.Error(w, store.ErrRoleNameAlreadyInUse.Error(), http.StatusBadRequest)
			return
		} else if err != nil && err != store.ErrRoleNotFound {
			log.WithError(err).Error("lookup role")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		role, err := s.roles.UpdateRole(r.Context(), with.Role.ID, with.Project.ID, updateRoleRequest.Name,
			updateRoleRequest.Description, updateRoleRequest.Config)
		if err != nil {
			log.WithError(err).Error("update role")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, role)
	})
}

func (s *Service) deleteRole(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.roles.DeleteRole(r.Context(), with.Role.ID, with.Project.ID); err != nil {
			log.WithError(err).Error("delete role")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) createServiceAccount(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var createServiceAccountRequest struct {
			Name        string `json:"name" validate:"name"`
			Description string `json:"description" validate:"description"`
		}
		if err := read(r, &createServiceAccountRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := s.serviceAccounts.LookupServiceAccount(r.Context(), createServiceAccountRequest.Name, with.Project.ID); err == nil {
			http.Error(w, store.ErrServiceAccountNameAlreadyInUse.Error(), http.StatusBadRequest)
			return
		} else if err != nil && err != store.ErrServiceAccountNotFound {
			log.WithError(err).Error("lookup service account")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		serviceAccount, err := s.serviceAccounts.CreateServiceAccount(r.Context(), with.Project.ID, createServiceAccountRequest.Name,
			createServiceAccountRequest.Description)
		if err != nil {
			log.WithError(err).Error("create service account")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, serviceAccount)
	})
}

func (s *Service) getServiceAccount(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ret interface{} = with.ServiceAccount
		if _, ok := r.URL.Query()["full"]; ok {
			serviceAccountRoleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(), with.ServiceAccount.ID, with.Project.ID)
			if err != nil {
				log.WithError(err).Error("list service account role bindings")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			roles := make([]models.Role, 0)
			for _, serviceAccountRoleBinding := range serviceAccountRoleBindings {
				role, err := s.roles.GetRole(r.Context(), serviceAccountRoleBinding.RoleID, with.Project.ID)
				if err != nil {
					log.WithError(err).Error("get role")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				roles = append(roles, *role)
			}

			ret = models.ServiceAccountFull{
				ServiceAccount: *with.ServiceAccount,
				Roles:          roles,
			}
		}

		utils.Respond(w, ret)
	})
}

func (s *Service) listServiceAccounts(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceAccounts, err := s.serviceAccounts.ListServiceAccounts(r.Context(), with.Project.ID)
		if err != nil {
			log.WithError(err).Error("list service accounts")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ret interface{} = serviceAccounts
		if _, ok := r.URL.Query()["full"]; ok {
			serviceAccountsFull := make([]models.ServiceAccountFull, 0)

			for _, serviceAccount := range serviceAccounts {
				serviceAccountRoleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(), serviceAccount.ID, with.Project.ID)
				if err != nil {
					log.WithError(err).Error("list service account role bindings")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				roles := make([]models.Role, 0)
				for _, serviceAccountRoleBinding := range serviceAccountRoleBindings {
					role, err := s.roles.GetRole(r.Context(), serviceAccountRoleBinding.RoleID, with.Project.ID)
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
	})
}

func (s *Service) updateServiceAccount(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var updateServiceAccountRequest struct {
			Name        string `json:"name" validate:"name"`
			Description string `json:"description" validate:"description"`
		}
		if err := read(r, &updateServiceAccountRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		serviceAccount, err := s.serviceAccounts.UpdateServiceAccount(r.Context(), with.ServiceAccount.ID, with.Project.ID,
			updateServiceAccountRequest.Name, updateServiceAccountRequest.Description)
		if err != nil {
			log.WithError(err).Error("update service account")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, serviceAccount)
	})
}

func (s *Service) deleteServiceAccount(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.serviceAccounts.DeleteServiceAccount(r.Context(), with.ServiceAccount.ID, with.Project.ID); err != nil {
			log.WithError(err).Error("delete service account")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) createServiceAccountAccessKey(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			with.Project.ID, serviceAccountID, hash.Hash(serviceAccountAccessKeyValue), createServiceAccountAccessKeyRequest.Description)
		if err != nil {
			log.WithError(err).Error("create service account access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, models.ServiceAccountAccessKeyWithValue{
			ServiceAccountAccessKey: *serviceAccount,
			Value:                   serviceAccountAccessKeyValue,
		})
	})
}

func (s *Service) getServiceAccountAccessKey(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceAccountAccessKeyID := vars["serviceaccountaccesskey"]

		serviceAccountAccessKey, err := s.serviceAccountAccessKeys.GetServiceAccountAccessKey(r.Context(), serviceAccountAccessKeyID, with.Project.ID)
		if err == store.ErrServiceAccountAccessKeyNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get service account access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, serviceAccountAccessKey)
	})
}

func (s *Service) listServiceAccountAccessKeys(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceAccountID := vars["serviceaccount"]

		serviceAccountAccessKeys, err := s.serviceAccountAccessKeys.ListServiceAccountAccessKeys(r.Context(), with.Project.ID, serviceAccountID)
		if err != nil {
			log.WithError(err).Error("list service accounts")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, serviceAccountAccessKeys)
	})
}

func (s *Service) deleteServiceAccountAccessKey(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceAccountAccessKeyID := vars["serviceaccountaccesskey"]

		if err := s.serviceAccountAccessKeys.DeleteServiceAccountAccessKey(r.Context(), serviceAccountAccessKeyID, with.Project.ID); err != nil {
			log.WithError(err).Error("delete service account access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) createServiceAccountRoleBinding(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceAccountID := vars["serviceaccount"]

		serviceAccountRoleBinding, err := s.serviceAccountRoleBindings.CreateServiceAccountRoleBinding(r.Context(), serviceAccountID, with.Role.ID, with.Project.ID)
		if err != nil {
			log.WithError(err).Error("create service account role binding")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, serviceAccountRoleBinding)
	})
}

func (s *Service) getServiceAccountRoleBinding(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceAccountID := vars["serviceaccount"]

		serviceAccountRoleBinding, err := s.serviceAccountRoleBindings.GetServiceAccountRoleBinding(r.Context(), serviceAccountID, with.Role.ID, with.Project.ID)
		if err != nil {
			log.WithError(err).Error("get service account role binding")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, serviceAccountRoleBinding)
	})
}

func (s *Service) listServiceAccountRoleBindings(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceAccountID := vars["serviceaccount"]

		serviceAccountRoleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(), with.Project.ID, serviceAccountID)
		if err != nil {
			log.WithError(err).Error("list service account role bindings")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, serviceAccountRoleBindings)
	})
}

func (s *Service) deleteServiceAccountRoleBinding(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceAccountID := vars["serviceaccount"]

		if err := s.serviceAccountRoleBindings.DeleteServiceAccountRoleBinding(r.Context(), serviceAccountID, with.Role.ID, with.Project.ID); err != nil {
			log.WithError(err).Error("delete service account role binding")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) createMembership(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var createMembershipRequest struct {
			Email string `json:"email" validate:"email"`
		}
		if err := read(r, &createMembershipRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := s.users.LookupUser(r.Context(), createMembershipRequest.Email)
		if err == store.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("lookup user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		membership, err := s.memberships.CreateMembership(r.Context(), user.ID, with.Project.ID)
		if err != nil {
			log.WithError(err).Error("create membership")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, membership)
	})
}

func (s *Service) getMembership(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user"]

		membership, err := s.memberships.GetMembership(r.Context(), userID, with.Project.ID)
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
			user, err := s.users.GetUser(r.Context(), membership.UserID)
			if err != nil {
				log.WithError(err).Error("get user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), userID, with.Project.ID)
			if err != nil {
				log.WithError(err).Error("list membership role bindings")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			roles := make([]models.Role, 0)
			for _, membershipRoleBinding := range membershipRoleBindings {
				role, err := s.roles.GetRole(r.Context(), membershipRoleBinding.RoleID, with.Project.ID)
				if err != nil {
					log.WithError(err).Error("get role")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				roles = append(roles, *role)
			}

			ret = models.MembershipFull2{
				Membership: *membership,
				User:       *user,
				Roles:      roles,
			}
		}

		utils.Respond(w, ret)
	})
}

func (s *Service) listMembershipsByProject(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		memberships, err := s.memberships.ListMembershipsByProject(r.Context(), with.Project.ID)
		if err != nil {
			log.WithError(err).Error("list memberships by project")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ret interface{} = memberships
		if _, ok := r.URL.Query()["full"]; ok {
			membershipsFull := make([]models.MembershipFull2, 0)

			for _, membership := range memberships {
				user, err := s.users.GetUser(r.Context(), membership.UserID)
				if err != nil {
					log.WithError(err).Error("get user")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), membership.UserID, with.Project.ID)
				if err != nil {
					log.WithError(err).Error("list membership role bindings")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				roles := make([]models.Role, 0)
				for _, membershipRoleBinding := range membershipRoleBindings {
					role, err := s.roles.GetRole(r.Context(), membershipRoleBinding.RoleID, with.Project.ID)
					if err != nil {
						log.WithError(err).Error("get role")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					roles = append(roles, *role)
				}

				membershipsFull = append(membershipsFull, models.MembershipFull2{
					Membership: membership,
					User:       *user,
					Roles:      roles,
				})
			}

			ret = membershipsFull
		}

		utils.Respond(w, ret)
	})
}

func (s *Service) deleteMembership(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user"]

		if err := s.memberships.DeleteMembership(r.Context(), userID, with.Project.ID); err != nil {
			log.WithError(err).Error("delete membership")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) createMembershipRoleBinding(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user"]

		membershipRoleBinding, err := s.membershipRoleBindings.CreateMembershipRoleBinding(r.Context(), userID, with.Role.ID, with.Project.ID)
		if err != nil {
			log.WithError(err).Error("create membership role binding")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, membershipRoleBinding)
	})
}

func (s *Service) getMembershipRoleBinding(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user"]

		membershipRoleBinding, err := s.membershipRoleBindings.GetMembershipRoleBinding(r.Context(), userID, with.Role.ID, with.Project.ID)
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
}

func (s *Service) listMembershipRoleBindings(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user"]

		membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), userID, with.Project.ID)
		if err != nil {
			log.WithError(err).Error("list membership role bindings")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, membershipRoleBindings)
	})
}

func (s *Service) deleteMembershipRoleBinding(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user"]

		if err := s.membershipRoleBindings.DeleteMembershipRoleBinding(r.Context(), userID, with.Role.ID, with.Project.ID); err != nil {
			log.WithError(err).Error("delete membership role binding")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) createApplication(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var createApplicationRequest struct {
			Name        string `json:"name" validate:"name"`
			Description string `json:"description" validate:"description"`
		}
		if err := read(r, &createApplicationRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := s.applications.LookupApplication(r.Context(), createApplicationRequest.Name, with.Project.ID); err == nil {
			http.Error(w, store.ErrApplicationNameAlreadyInUse.Error(), http.StatusBadRequest)
			return
		} else if err != nil && err != store.ErrApplicationNotFound {
			log.WithError(err).Error("lookup application")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		application, err := s.applications.CreateApplication(
			r.Context(),
			with.Project.ID,
			createApplicationRequest.Name,
			createApplicationRequest.Description)
		if err != nil {
			log.WithError(err).Error("create application")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, application)
	})
}

func (s *Service) getApplication(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ret interface{} = with.Application
		if _, ok := r.URL.Query()["full"]; ok {
			latestRelease, err := s.releases.GetLatestRelease(r.Context(), with.Project.ID, with.Application.ID)
			if err != nil && err != store.ErrReleaseNotFound {
				log.WithError(err).Error("get latest release")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			applicationDeviceCounts, err := s.applicationDeviceCounts.GetApplicationDeviceCounts(r.Context(), with.Project.ID, with.Application.ID)
			if err != nil {
				log.WithError(err).Error("get application device counts")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ret = models.ApplicationFull1{
				Application:   *with.Application,
				LatestRelease: latestRelease,
				DeviceCounts:  *applicationDeviceCounts,
			}
		}

		utils.Respond(w, ret)
	})
}

func (s *Service) listApplications(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		applications, err := s.applications.ListApplications(r.Context(), with.Project.ID)
		if err != nil {
			log.WithError(err).Error("list applications")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ret interface{} = applications
		if _, ok := r.URL.Query()["full"]; ok {
			applicationsFull := make([]models.ApplicationFull1, 0)

			for _, application := range applications {
				latestRelease, err := s.releases.GetLatestRelease(r.Context(), with.Project.ID, application.ID)
				if err != nil && err != store.ErrReleaseNotFound {
					log.WithError(err).Error("get latest release")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				applicationDeviceCounts, err := s.applicationDeviceCounts.GetApplicationDeviceCounts(r.Context(), with.Project.ID, application.ID)
				if err != nil {
					log.WithError(err).Error("get application device counts")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				applicationsFull = append(applicationsFull, models.ApplicationFull1{
					Application:   application,
					LatestRelease: latestRelease,
					DeviceCounts:  *applicationDeviceCounts,
				})
			}

			ret = applicationsFull
		}

		utils.Respond(w, ret)
	})
}

func (s *Service) updateApplication(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		var application *models.Application
		var err error
		if updateApplicationRequest.Name != nil {
			if application, err = s.applications.LookupApplication(r.Context(),
				*updateApplicationRequest.Name, with.Project.ID); err == nil && application.ID != with.Application.ID {
				http.Error(w, store.ErrApplicationNameAlreadyInUse.Error(), http.StatusBadRequest)
				return
			} else if err != nil && err != store.ErrApplicationNotFound {
				log.WithError(err).Error("lookup application")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if application, err = s.applications.UpdateApplicationName(r.Context(), with.Application.ID, with.Project.ID, *updateApplicationRequest.Name); err != nil {
				log.WithError(err).Error("update application name")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if updateApplicationRequest.Description != nil {
			if application, err = s.applications.UpdateApplicationDescription(r.Context(), with.Application.ID, with.Project.ID, *updateApplicationRequest.Description); err != nil {
				log.WithError(err).Error("update application description")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if updateApplicationRequest.SchedulingRule != nil {
			validationErr, err := scheduling.ValidateSchedulingRule(
				*updateApplicationRequest.SchedulingRule,
				func(releaseID string) (bool, error) {
					var release *models.Release
					var err error
					if strings.Contains(releaseID, "_") {
						release, err = s.releases.GetRelease(r.Context(), releaseID, with.Project.ID, application.ID)
					} else if releaseID == "latest" { // TODO: models.LatestRelease
						release, err = s.releases.GetLatestRelease(r.Context(), with.Project.ID, application.ID)
					} else {
						id, parseErr := strconv.ParseUint(releaseID, 10, 32)
						if parseErr != nil {
							return false, parseErr
						}
						release, err = s.releases.GetReleaseByNumber(r.Context(), uint32(id), with.Project.ID, application.ID)
					}
					if err == store.ErrReleaseNotFound {
						return false, nil
					} else if err != nil {
						return false, err
					}
					return release != nil, nil
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

			if application, err = s.applications.UpdateApplicationSchedulingRule(r.Context(), with.Application.ID, with.Project.ID, *updateApplicationRequest.SchedulingRule); err != nil {
				log.WithError(err).Error("update application scheduling rule")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if updateApplicationRequest.MetricEndpointConfigs != nil {
			if application, err = s.applications.UpdateApplicationMetricEndpointConfigs(r.Context(), with.Application.ID, with.Project.ID, *updateApplicationRequest.MetricEndpointConfigs); err != nil {
				log.WithError(err).Error("update application service metrics config")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		utils.Respond(w, application)
	})
}

func (s *Service) deleteApplication(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.applications.DeleteApplication(r.Context(), with.Application.ID, with.Project.ID); err != nil {
			log.WithError(err).Error("delete application")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

// TODO: this has a vulnerability!
func (s *Service) createRelease(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		release, err := s.releases.CreateRelease(
			r.Context(),
			with.Project.ID,
			with.Application.ID,
			createReleaseRequest.RawConfig,
			string(jsonApplicationConfig),
			with.User.ID,
			with.ServiceAccount.ID,
		)
		if err != nil {
			log.WithError(err).Error("create release")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, release)
	})
}

func (s *Service) getRelease(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ret interface{} = with.Release
		var err error
		if _, ok := r.URL.Query()["full"]; ok {
			ret, err = s.getReleaseFull(r.Context(), *with.Release)
			if err != nil {
				log.WithError(err).Error("get release full")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		utils.Respond(w, ret)
	})
}

func (s *Service) listReleases(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		releases, err := s.releases.ListReleases(r.Context(), with.Project.ID, with.Application.ID)
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
}

func (s *Service) listDevices(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		searchQuery := r.URL.Query().Get("search")

		devices, err := s.devices.ListDevices(r.Context(), with.Project.ID, searchQuery)
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
			devices, _, err = query.QueryDevices(devices, filters)
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
	})
}

func (s *Service) previewScheduledDevices(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		searchQuery := r.URL.Query().Get("search")

		devices, err := s.devices.ListDevices(r.Context(), with.Project.ID, searchQuery)
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
			devices, _, err = query.QueryDevices(devices, filters)
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
	})
}

func (s *Service) getDevice(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ret interface{} = with.Device
		if _, ok := r.URL.Query()["full"]; ok {
			applications, err := s.applications.ListApplications(r.Context(), with.Project.ID)
			if err != nil {
				log.WithError(err).Error("list applications")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			allApplicationStatusInfo := make([]models.DeviceApplicationStatusInfo, 0)
			for _, application := range applications {
				applicationStatusInfo := models.DeviceApplicationStatusInfo{
					Application: application,
				}

				deviceApplicationStatus, err := s.deviceApplicationStatuses.GetDeviceApplicationStatus(
					r.Context(), with.Project.ID, with.Device.ID, application.ID)
				if err == nil {
					applicationStatusInfo.ApplicationStatus = deviceApplicationStatus
				} else if err != store.ErrDeviceApplicationStatusNotFound {
					log.WithError(err).Error("get device application status")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				deviceServiceStatuses, err := s.deviceServiceStatuses.GetDeviceServiceStatuses(
					r.Context(), with.Project.ID, with.Device.ID, application.ID)
				if err == nil {
					applicationStatusInfo.ServiceStatuses = deviceServiceStatuses
				} else {
					log.WithError(err).Error("get device service statuses")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				allApplicationStatusInfo = append(allApplicationStatusInfo, applicationStatusInfo)
			}

			ret = models.DeviceFull{
				Device:                *with.Device,
				ApplicationStatusInfo: allApplicationStatusInfo,
			}
		}

		utils.Respond(w, ret)
	})
}

func (s *Service) updateDevice(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var updateDeviceRequest struct {
			Name string `json:"name" validate:"name"`
		}
		if err := read(r, &updateDeviceRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		device, err := s.devices.UpdateDeviceName(r.Context(), with.Device.ID, with.Project.ID, updateDeviceRequest.Name)
		if err != nil {
			log.WithError(err).Error("update device name")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, device)
	})
}

func (s *Service) deleteDevice(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.devices.DeleteDevice(r.Context(), with.Device.ID, with.Project.ID); err != nil {
			log.WithError(err).Error("delete device")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) listAllDeviceLabelKeys(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceLabels, err := s.devices.ListAllDeviceLabelKeys(
			r.Context(),
			with.Project.ID,
		)
		if err != nil {
			log.WithError(err).Error("list device labels")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, deviceLabels)
	})
}

func (s *Service) setDeviceLabel(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			with.Device.ID,
			with.Project.ID,
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
}

func (s *Service) deleteDeviceLabel(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		if err := s.devices.DeleteDeviceLabel(r.Context(), with.Device.ID, with.Project.ID, key); err != nil {
			log.WithError(err).Error("delete device label")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) createDeviceRegistrationToken(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			with.Project.ID,
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
			with.Project.ID,
			createDeviceRegistrationTokenRequest.Name,
			createDeviceRegistrationTokenRequest.Description,
			createDeviceRegistrationTokenRequest.MaxRegistrations)
		if err != nil {
			log.WithError(err).Error("create device registration token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, deviceRegistrationToken)
	})
}

func (s *Service) getDeviceRegistrationToken(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ret interface{} = with.DeviceRegistrationToken
		if _, ok := r.URL.Query()["full"]; ok {
			devicesRegistered, err := s.devicesRegisteredWithToken.GetDevicesRegisteredWithTokenCount(r.Context(), with.DeviceRegistrationToken.ID, with.Project.ID)
			if err != nil {
				log.WithError(err).Error("get registered device count")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ret = models.DeviceRegistrationTokenFull{
				DeviceRegistrationToken: *with.DeviceRegistrationToken,
				DeviceCounts:            *devicesRegistered,
			}
		}

		utils.Respond(w, ret)
	})
}

func (s *Service) updateDeviceRegistrationToken(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var updateDeviceRegistrationTokenRequest struct {
			Name             string `json:"name" validate:"name"`
			Description      string `json:"description" validate:"description"`
			MaxRegistrations *int   `json:"maxRegistrations"`
		}
		if err := read(r, &updateDeviceRegistrationTokenRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		deviceRegistrationToken, err := s.deviceRegistrationTokens.UpdateDeviceRegistrationToken(
			r.Context(),
			with.DeviceRegistrationToken.ID,
			with.Project.ID,
			updateDeviceRegistrationTokenRequest.Name,
			updateDeviceRegistrationTokenRequest.Description,
			updateDeviceRegistrationTokenRequest.MaxRegistrations,
		)
		if err != nil {
			log.WithError(err).Error("update device registration token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, deviceRegistrationToken)
	})
}

func (s *Service) deleteDeviceRegistrationToken(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.deviceRegistrationTokens.DeleteDeviceRegistrationToken(r.Context(), with.DeviceRegistrationToken.ID, with.Project.ID); err != nil {
			log.WithError(err).Error("delete device registration token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) listDeviceRegistrationTokens(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceRegistrationTokens, err := s.deviceRegistrationTokens.ListDeviceRegistrationTokens(r.Context(), with.Project.ID)
		if err != nil {
			log.WithError(err).Error("list device registration tokens")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ret interface{} = deviceRegistrationTokens
		if _, ok := r.URL.Query()["full"]; ok {
			deviceRegistrationTokensFull := make([]models.DeviceRegistrationTokenFull, 0)

			for _, deviceRegistrationToken := range deviceRegistrationTokens {
				deviceRegistrationTokenCounts, err := s.devicesRegisteredWithToken.GetDevicesRegisteredWithTokenCount(r.Context(), deviceRegistrationToken.ID, with.Project.ID)
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
	})
}

func (s *Service) setDeviceRegistrationTokenLabel(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			with.DeviceRegistrationToken.ID,
			with.Project.ID,
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
}

func (s *Service) deleteDeviceRegistrationTokenLabel(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		if err := s.deviceRegistrationTokens.DeleteDeviceRegistrationTokenLabel(r.Context(), with.DeviceRegistrationToken.ID, with.Project.ID, key); err != nil {
			log.WithError(err).Error("delete device registration token label")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) getProjectConfig(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		var value interface{}
		var err error
		switch key {
		case string(models.ProjectMetricsConfigKey):
			value, err = s.metricConfigs.GetProjectMetricsConfig(r.Context(), with.Project.ID)
		case string(models.DeviceMetricsConfigKey):
			value, err = s.metricConfigs.GetDeviceMetricsConfig(r.Context(), with.Project.ID)
		case string(models.ServiceMetricsConfigKey):
			value, err = s.metricConfigs.GetServiceMetricsConfigs(r.Context(), with.Project.ID)
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
	})
}

func (s *Service) setProjectConfig(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			err = s.metricConfigs.SetProjectMetricsConfig(r.Context(), with.Project.ID, value)
		case string(models.DeviceMetricsConfigKey):
			var value models.DeviceMetricsConfig
			if err := read(r, &value); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = s.metricConfigs.SetDeviceMetricsConfig(r.Context(), with.Project.ID, value)
		case string(models.ServiceMetricsConfigKey):
			var values []models.ServiceMetricsConfig
			// TODO: use read() here
			if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = s.metricConfigs.SetServiceMetricsConfigs(r.Context(), with.Project.ID, values)
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
	})
}

// TODO: verify project ID
func (s *Service) registerDevice(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var registerDeviceRequest models.RegisterDeviceRequest
		if err := read(r, &registerDeviceRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		deviceRegistrationToken, err := s.deviceRegistrationTokens.GetDeviceRegistrationToken(r.Context(), registerDeviceRequest.DeviceRegistrationTokenID, with.Project.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if deviceRegistrationToken.MaxRegistrations != nil {
			devicesRegisteredCount, err := s.devicesRegisteredWithToken.GetDevicesRegisteredWithTokenCount(r.Context(), registerDeviceRequest.DeviceRegistrationTokenID, with.Project.ID)
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

		device, err := s.devices.CreateDevice(r.Context(), with.Project.ID, namesgenerator.GetRandomName(), deviceRegistrationToken.ID, deviceRegistrationToken.Labels)
		if err != nil {
			log.WithError(err).Error("create device")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		deviceAccessKeyValue := ksuid.New().String()

		_, err = s.deviceAccessKeys.CreateDeviceAccessKey(r.Context(), with.Project.ID, device.ID, hash.Hash(deviceAccessKeyValue))
		if err != nil {
			log.WithError(err).Error("create device access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, models.RegisterDeviceResponse{
			DeviceID:             device.ID,
			DeviceAccessKeyValue: deviceAccessKeyValue,
		})
	})
}

func (s *Service) getBundle(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.st.Incr("get_bundle", []string{
			fmt.Sprintf("project_id:%s", with.Project.ID),
			fmt.Sprintf("project_name:%s", with.Project.Name),
		}, 1)

		if err := s.devices.UpdateDeviceLastSeenAt(r.Context(), with.Device.ID, with.Project.ID); err != nil {
			log.WithError(err).Error("update device last seen at")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		applications, err := s.applications.ListApplications(r.Context(), with.Project.ID)
		if err != nil {
			log.WithError(err).Error("list applications")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bundle := models.Bundle{
			DesiredAgentSpec:    with.Device.DesiredAgentSpec,
			DesiredAgentVersion: with.Device.DesiredAgentVersion,
		}

		for _, application := range applications {
			scheduled, scheduledDevice, err := scheduling.IsApplicationScheduled(*with.Device, application.SchedulingRule)
			if err != nil {
				log.WithError(err).Error("evaluate application scheduling rule")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !scheduled {
				continue
			}

			release, err := utils.GetReleaseByIdentifier(s.releases, r.Context(), with.Project.ID, application.ID, scheduledDevice.ReleaseID)
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
			r.Context(), with.Project.ID, with.Device.ID)
		if err == nil {
			bundle.ApplicationStatuses = deviceApplicationStatuses
		} else {
			log.WithError(err).Error("list device application statuses")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		deviceServiceStatuses, err := s.deviceServiceStatuses.ListDeviceServiceStatuses(
			r.Context(), with.Project.ID, with.Device.ID)
		if err == nil {
			bundle.ServiceStatuses = deviceServiceStatuses
		} else {
			log.WithError(err).Error("list device service statuses")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Respond(w, bundle)
	})
}

func (s *Service) setDeviceInfo(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var setDeviceInfoRequest models.SetDeviceInfoRequest
		if err := read(r, &setDeviceInfoRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := s.devices.SetDeviceInfo(r.Context(), with.Device.ID, with.Project.ID, setDeviceInfoRequest.DeviceInfo); err != nil {
			log.WithError(err).Error("set device info")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) setDeviceApplicationStatus(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		applicationID := vars["application"]

		var setDeviceApplicationStatusRequest models.SetDeviceApplicationStatusRequest
		if err := read(r, &setDeviceApplicationStatusRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.deviceApplicationStatuses.SetDeviceApplicationStatus(r.Context(), with.Project.ID, with.Device.ID,
			applicationID, setDeviceApplicationStatusRequest.CurrentReleaseID,
		); err != nil {
			log.WithError(err).Error("set device application status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) deleteDeviceApplicationStatus(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		applicationID := vars["application"]

		if err := s.deviceApplicationStatuses.DeleteDeviceApplicationStatus(r.Context(),
			with.Project.ID, with.Device.ID, applicationID,
		); err != nil {
			log.WithError(err).Error("delete device application status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) setDeviceServiceStatus(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		applicationID := vars["application"]
		service := vars["service"]

		var setDeviceServiceStatusRequest models.SetDeviceServiceStatusRequest
		if err := read(r, &setDeviceServiceStatusRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.deviceServiceStatuses.SetDeviceServiceStatus(r.Context(), with.Project.ID, with.Device.ID,
			applicationID, service, setDeviceServiceStatusRequest.CurrentReleaseID,
		); err != nil {
			log.WithError(err).Error("set device service status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) deleteDeviceServiceStatus(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		applicationID := vars["application"]
		service := vars["service"]

		if err := s.deviceServiceStatuses.DeleteDeviceServiceStatus(r.Context(),
			with.Project.ID, with.Device.ID, applicationID, service,
		); err != nil {
			log.WithError(err).Error("delete device service status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) submitStats(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		applicationID := vars["application"]
		service := vars["service"]

		if err := s.deviceServiceStatuses.DeleteDeviceServiceStatus(r.Context(),
			with.Project.ID, with.Device.ID, applicationID, service,
		); err != nil {
			log.WithError(err).Error("delete device service status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
