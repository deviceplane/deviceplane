package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deviceplane/deviceplane/pkg/spec"
	"gopkg.in/yaml.v2"

	"github.com/segmentio/ksuid"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/scheduler"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/email"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/gorilla/mux"
)

const (
	sessionCookie = "dp_sess"
)

type Service struct {
	users                    store.Users
	registrationTokens       store.RegistrationTokens
	accessKeys               store.AccessKeys
	sessions                 store.Sessions
	projects                 store.Projects
	memberships              store.Memberships
	devices                  store.Devices
	deviceLabels             store.DeviceLabels
	deviceRegistrationTokens store.DeviceRegistrationTokens
	deviceAccessKeys         store.DeviceAccessKeys
	applications             store.Applications
	releases                 store.Releases
	email                    email.Interface
	router                   *mux.Router
	cookieDomain             string
}

func NewService(
	users store.Users,
	registrationTokens store.RegistrationTokens,
	sessions store.Sessions,
	accessKeys store.AccessKeys,
	projects store.Projects,
	memberships store.Memberships,
	devices store.Devices,
	deviceLabels store.DeviceLabels,
	deviceRegistrationTokens store.DeviceRegistrationTokens,
	deviceAccessKeys store.DeviceAccessKeys,
	applications store.Applications,
	releases store.Releases,
	email email.Interface,
	cookieDomain string,
) *Service {
	s := &Service{
		users:                    users,
		registrationTokens:       registrationTokens,
		sessions:                 sessions,
		accessKeys:               accessKeys,
		projects:                 projects,
		memberships:              memberships,
		devices:                  devices,
		deviceLabels:             deviceLabels,
		deviceRegistrationTokens: deviceRegistrationTokens,
		deviceAccessKeys:         deviceAccessKeys,
		applications:             applications,
		releases:                 releases,
		email:                    email,
		cookieDomain:             cookieDomain,
		router:                   mux.NewRouter(),
	}

	s.router.HandleFunc("/health", s.health).Methods("GET")

	s.router.HandleFunc("/register", s.register).Methods("POST")
	s.router.HandleFunc("/completeregistration", s.confirmRegistration).Methods("POST")
	s.router.HandleFunc("/login", s.login).Methods("POST")
	s.router.HandleFunc("/logout", s.logout).Methods("POST")
	s.router.HandleFunc("/me", s.withUserAuth(s.me)).Methods("GET")

	s.router.HandleFunc("/users/{user}/memberships", s.withUserAuth(s.listMembershipsByUser)).Methods("GET")

	s.router.HandleFunc("/projects", s.withUserAuth(s.createProject)).Methods("POST")
	s.router.HandleFunc("/projects/{project}", s.validateMembershipLevel("write", s.getProject)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/memberships", s.validateMembershipLevel("admin", s.createMembership)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/memberships", s.validateMembershipLevel("read", s.listMembershipsByProject)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/applications", s.validateMembershipLevel("write", s.createApplication)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/applications/{id}", s.validateMembershipLevel("read", s.getApplication)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications", s.validateMembershipLevel("read", s.listApplications)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/applications/{application}/releases", s.validateMembershipLevel("write", s.createRelease)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/applications/{application}/releases/latest", s.validateMembershipLevel("read", s.getLatestRelease)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications/{application}/releases/{id}", s.validateMembershipLevel("read", s.getRelease)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications/{application}/releases", s.validateMembershipLevel("read", s.listReleases)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/devices/{id}", s.validateMembershipLevel("read", s.getDevice)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/devices", s.validateMembershipLevel("read", s.listDevices)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/devices/{device}/labels", s.validateMembershipLevel("write", s.setDeviceLabel)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/devices/{device}/labels/{key}", s.validateMembershipLevel("read", s.getDeviceLabel)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/devices/{device}/labels", s.validateMembershipLevel("read", s.listDeviceLabels)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/devices/{device}/labels/{key}", s.validateMembershipLevel("write", s.deleteDeviceLabel)).Methods("DELETE")

	s.router.HandleFunc("/projects/{project}/deviceregistrationtokens", s.validateMembershipLevel("write", s.createDeviceRegistrationToken)).Methods("POST")

	s.router.HandleFunc("/projects/{project}/devices/register", s.registerDevice).Methods("POST")
	s.router.HandleFunc("/projects/{project}/devices/{device}/bundle", s.withDeviceAuth(s.getBundle)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/devices/{device}/info", s.withDeviceAuth(s.setDeviceInfo)).Methods("POST")

	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Service) withUserAuth(handler func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID string

		sessionValue, err := r.Cookie(sessionCookie)

		switch err {
		case nil:
			session, err := s.sessions.ValidateSession(r.Context(), hash(sessionValue.Value))
			if err == store.ErrSessionNotFound {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err != nil {
				log.WithError(err).Error("validate session")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userID = session.UserID
		case http.ErrNoCookie:
			accessKeyValue, _, _ := r.BasicAuth()
			if accessKeyValue == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			accessKey, err := s.accessKeys.ValidateAccessKey(r.Context(), hash(accessKeyValue))
			if err == store.ErrAccessKeyNotFound {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err != nil {
				log.WithError(err).Error("validate access key")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userID = accessKey.UserID
		default:
			log.WithError(err).Error("get session cookie")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if userID == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		user, err := s.users.GetUser(r.Context(), userID)
		if err == store.ErrUserNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !user.RegistrationCompleted {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(w, r, userID)
	}
}

func (s *Service) validateMembershipLevel(requiredLevel string, handler func(http.ResponseWriter, *http.Request, string, string)) func(http.ResponseWriter, *http.Request) {
	return s.withUserAuth(func(w http.ResponseWriter, r *http.Request, userID string) {
		vars := mux.Vars(r)
		projectID := vars["project"]
		if projectID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		membership, err := s.memberships.GetMembership(r.Context(), userID, projectID)
		if err != nil {
			log.WithError(err).Error("get membership")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if membership == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		accessGranted := false
		switch requiredLevel {
		case "admin":
			if membership.Level == "admin" {
				accessGranted = true
			}
		case "write":
			if membership.Level == "admin" || membership.Level == "write" {
				accessGranted = true
			}
		case "read":
			if membership.Level == "admin" || membership.Level == "write" || membership.Level == "read" {
				accessGranted = true
			}
		}

		if !accessGranted {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(w, r, projectID, userID)
	})
}

func (s *Service) register(w http.ResponseWriter, r *http.Request) {
	var registerRequest struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.users.CreateUser(r.Context(), registerRequest.Email, hash(registerRequest.Password),
		registerRequest.FirstName, registerRequest.LastName)
	if err != nil {
		log.WithError(err).Error("create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	registrationTokenValue := ksuid.New().String()

	if _, err := s.registrationTokens.CreateRegistrationToken(r.Context(), user.ID, hash(registrationTokenValue)); err != nil {
		log.WithError(err).Error("create registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	name := user.FirstName + " " + user.LastName

	if err := s.email.Send(email.Request{
		FromName:         "Device Plane",
		FromAddress:      "noreply@deviceplane.io",
		ToName:           name,
		ToAddress:        user.Email,
		Subject:          "Device Plane Registration Confirmation",
		PlainTextContent: "Please go to the following URL to complete registration. https://app.deviceplane.io/confirm/" + registrationTokenValue,
		HTMLContent:      "Please go to the following URL to complete registration. https://app.deviceplane.io/confirm/" + registrationTokenValue,
	}); err != nil {
		log.WithError(err).Error("send registration email")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) confirmRegistration(w http.ResponseWriter, r *http.Request) {
	var confirmRegistrationRequest struct {
		RegistrationTokenValue string `json:"registrationTokenValue"`
	}
	if err := json.NewDecoder(r.Body).Decode(&confirmRegistrationRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	registrationToken, err := s.registrationTokens.ValidateRegistrationToken(r.Context(),
		hash(confirmRegistrationRequest.RegistrationTokenValue))
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

	s.newSession(w, r, registrationToken.UserID)
}

func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.users.ValidateUser(r.Context(), loginRequest.Email, hash(loginRequest.Password))
	if err == store.ErrUserNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("validate user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !user.RegistrationCompleted {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	s.newSession(w, r, user.ID)
}

func (s *Service) newSession(w http.ResponseWriter, r *http.Request, userID string) {
	sessionValue := ksuid.New().String()

	if _, err := s.sessions.CreateSession(r.Context(), userID, hash(sessionValue)); err != nil {
		log.WithError(err).Error("create session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  sessionCookie,
		Value: sessionValue,

		Domain:  s.cookieDomain,
		Expires: time.Now().AddDate(0, 1, 0),

		Secure:   true,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

func (s *Service) logout(w http.ResponseWriter, r *http.Request) {
	sessionValue, err := r.Cookie(sessionCookie)

	switch err {
	case nil:
		session, err := s.sessions.ValidateSession(r.Context(), hash(sessionValue.Value))
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
		w.WriteHeader(http.StatusOK)
		return
	default:
		log.WithError(err).Error("get session cookie")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) me(w http.ResponseWriter, r *http.Request, authenticatedUserID string) {
	user, err := s.users.GetUser(r.Context(), authenticatedUserID)
	if err != nil {
		log.WithError(err).Error("get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) listMembershipsByUser(w http.ResponseWriter, r *http.Request, authenticatedUserID string) {
	vars := mux.Vars(r)
	userID := vars["user"]

	if userID != authenticatedUserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	memberships, err := s.memberships.ListMembershipsByUser(r.Context(), userID)
	if err != nil {
		log.WithError(err).Error("list memberships by user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memberships)
}

func (s *Service) createProject(w http.ResponseWriter, r *http.Request, userID string) {
	var createProjectRequest struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createProjectRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := s.projects.CreateProject(r.Context(), createProjectRequest.Name)
	if err != nil {
		log.WithError(err).Error("create project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = s.memberships.CreateMembership(r.Context(), userID, project.ID, "admin")
	if err != nil {
		log.WithError(err).Error("create membership")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func (s *Service) getProject(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	project, err := s.projects.GetProject(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("get project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func (s *Service) createMembership(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	var createMembershipRequest struct {
		UserID string `json:"userId"`
		Level  string `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createMembershipRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	membership, err := s.memberships.CreateMembership(r.Context(),
		createMembershipRequest.UserID, projectID, createMembershipRequest.Level)
	if err != nil {
		log.WithError(err).Error("create membership")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(membership)
}

func (s *Service) listMembershipsByProject(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	memberships, err := s.memberships.ListMembershipsByProject(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list memberships by project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memberships)
}

func (s *Service) createApplication(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	var createApplicationRequest struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createApplicationRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	application, err := s.applications.CreateApplication(r.Context(), projectID, createApplicationRequest.Name)
	if err != nil {
		log.WithError(err).Error("create application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(application)
}

func (s *Service) getApplication(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	id := vars["id"]

	application, err := s.applications.GetApplication(r.Context(), id, projectID)
	if err != nil {
		log.WithError(err).Error("get application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(application)
}

func (s *Service) listApplications(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	applications, err := s.applications.ListApplications(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list applications")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(applications)
}

// TOOD: this has a vulnerability!
func (s *Service) createRelease(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	applicationID := vars["application"]

	var createReleaseRequest struct {
		Config string `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createReleaseRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	release, err := s.releases.CreateRelease(r.Context(), projectID, applicationID, createReleaseRequest.Config)
	if err != nil {
		log.WithError(err).Error("create release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(release)
}

func (s *Service) getRelease(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	id := vars["id"]

	release, err := s.releases.GetRelease(r.Context(), id, projectID)
	if err != nil {
		log.WithError(err).Error("get release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(release)
}

func (s *Service) getLatestRelease(w http.ResponseWriter, r *http.Request, projectID string, userID string) {
	vars := mux.Vars(r)
	applicationID := vars["application"]

	release, err := s.releases.GetLatestRelease(r.Context(), projectID, applicationID)
	if err != nil {
		log.WithError(err).Error("get latest release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(release)
}

func (s *Service) listReleases(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	applicationID := vars["application"]

	releases, err := s.releases.ListReleases(r.Context(), projectID, applicationID)
	if err != nil {
		log.WithError(err).Error("list releases")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(releases)
}

func (s *Service) listDevices(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	devices, err := s.devices.ListDevices(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list devices")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(devices)
}

func (s *Service) getDevice(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	id := vars["id"]

	device, err := s.devices.GetDevice(r.Context(), id, projectID)
	if err != nil {
		log.WithError(err).Error("get device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(device)
}

func (s *Service) setDeviceLabel(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	deviceID := vars["device"]

	var setDeviceLabelRequest struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&setDeviceLabelRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deviceLabel, err := s.deviceLabels.SetDeviceLabel(r.Context(), setDeviceLabelRequest.Key,
		deviceID, projectID, setDeviceLabelRequest.Value)
	if err != nil {
		log.WithError(err).Error("set device label")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceLabel)
}

func (s *Service) getDeviceLabel(w http.ResponseWriter, r *http.Request, projectID string, userID string) {
	vars := mux.Vars(r)
	deviceID := vars["device"]
	key := vars["key"]

	deviceLabel, err := s.deviceLabels.GetDeviceLabel(r.Context(), key, deviceID, projectID)
	if err != nil {
		log.WithError(err).Error("get device label")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceLabel)
}

func (s *Service) listDeviceLabels(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	deviceID := vars["device"]

	deviceLabels, err := s.deviceLabels.ListDeviceLabels(r.Context(), deviceID, projectID)
	if err != nil {
		log.WithError(err).Error("list device labels")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceLabels)
}

func (s *Service) deleteDeviceLabel(w http.ResponseWriter, r *http.Request, projectID string, userID string) {
	vars := mux.Vars(r)
	deviceID := vars["device"]
	key := vars["key"]

	if err := s.deviceLabels.DeleteDeviceLabel(r.Context(), key, deviceID, projectID); err != nil {
		log.WithError(err).Error("delete device label")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) createDeviceRegistrationToken(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	deviceRegistrationToken, err := s.deviceRegistrationTokens.CreateDeviceRegistrationToken(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("create device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceRegistrationToken)
}

func (s *Service) withDeviceAuth(handler func(http.ResponseWriter, *http.Request, string, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectID := vars["project"]

		deviceAccessKeyValue, _, _ := r.BasicAuth()
		if deviceAccessKeyValue == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		deviceAccessKey, err := s.deviceAccessKeys.ValidateDeviceAccessKey(r.Context(), projectID, hash(deviceAccessKeyValue))
		if err == store.ErrDeviceAccessKeyNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			log.WithError(err).Error("validate device access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler(w, r, projectID, deviceAccessKey.DeviceID)
	}
}

// TODO: verify project ID
func (s *Service) registerDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]

	var registerDeviceRequest models.RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&registerDeviceRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deviceRegistrationToken, err := s.deviceRegistrationTokens.GetDeviceRegistrationToken(r.Context(), registerDeviceRequest.DeviceRegistrationTokenID, projectID)
	if err != nil {
		log.WithError(err).Error("get device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if deviceRegistrationToken.DeviceAccessKeyID != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	device, err := s.devices.CreateDevice(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("create device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deviceAccessKeyValue := ksuid.New().String()

	deviceAccessKey, err := s.deviceAccessKeys.CreateDeviceAccessKey(r.Context(), projectID, device.ID, hash(deviceAccessKeyValue))
	if err != nil {
		log.WithError(err).Error("create device access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = s.deviceRegistrationTokens.BindDeviceRegistrationToken(r.Context(), deviceRegistrationToken.ID, projectID, deviceAccessKey.ID); err != nil {
		log.WithError(err).Error("bind device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.RegisterDeviceResponse{
		DeviceID:             device.ID,
		DeviceAccessKeyValue: deviceAccessKeyValue,
	})
}

func (s *Service) getBundle(w http.ResponseWriter, r *http.Request, projectID, deviceID string) {
	var bundle models.Bundle

	applications, err := s.applications.ListApplications(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list applications")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deviceLabels, err := s.deviceLabels.ListDeviceLabels(r.Context(), deviceID, projectID)
	if err != nil {
		log.WithError(err).Error("list device labels")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i, application := range applications {
		release, err := s.releases.GetLatestRelease(r.Context(), projectID, application.ID)
		if err != nil && err != store.ErrReleaseNotFound {
			log.WithError(err).Error("get latest release")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var applicationSpec spec.Application
		if err = yaml.Unmarshal([]byte(release.Config), &applicationSpec); err != nil {
			log.WithError(err).Error("invalid application spec")
			continue
		}

		transformedApplicationSpec, err := scheduler.TransformSpec(applicationSpec, deviceLabels)
		if err != nil {
			log.WithError(err).Error("transform application spec")
			continue
		}

		transformedApplicationSpecBytes, err := yaml.Marshal(transformedApplicationSpec)
		if err != nil {
			log.WithError(err).Error("marshal transformed application spec")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		release.Config = string(transformedApplicationSpecBytes)

		bundle.Applications = append(bundle.Applications, models.ApplicationAndLatestRelease{
			Application:   applications[i],
			LatestRelease: release,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bundle)
}

func (s *Service) setDeviceInfo(w http.ResponseWriter, r *http.Request, projectID, deviceID string) {
	var setDeviceInfoRequest models.SetDeviceInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&setDeviceInfoRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := s.devices.SetDeviceInfo(r.Context(), deviceID, projectID, setDeviceInfoRequest.Info); err != nil {
		log.WithError(err).Error("set device info")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func hash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}
