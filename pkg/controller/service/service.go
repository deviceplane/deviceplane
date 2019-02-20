package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/gorilla/mux"
)

type Service struct {
	users        store.Users
	accessKeys   store.AccessKeys
	projects     store.Projects
	memberships  store.Memberships
	devices      store.Devices
	applications store.Applications
	releases     store.Releases
	router       *mux.Router
}

func NewService(
	users store.Users,
	accessKeys store.AccessKeys,
	projects store.Projects,
	memberships store.Memberships,
	devices store.Devices,
	applications store.Applications,
	releases store.Releases,
) *Service {
	s := &Service{
		users:        users,
		accessKeys:   accessKeys,
		projects:     projects,
		memberships:  memberships,
		devices:      devices,
		applications: applications,
		releases:     releases,
		router:       mux.NewRouter(),
	}

	s.router.HandleFunc("/health", s.health).Methods("GET")

	s.router.HandleFunc("/users/{user}/memberships", s.withUserAuth(s.listMembershipsByUser)).Methods("GET")

	s.router.HandleFunc("/projects", s.withUserAuth(s.createProject)).Methods("POST")
	s.router.HandleFunc("/projects/{project}", s.validateMembershipLevel("write", s.getProject)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/memberships", s.validateMembershipLevel("admin", s.createMembership)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/memberships", s.validateMembershipLevel("read", s.listMembershipsByProject)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/applications", s.validateMembershipLevel("write", s.createApplication)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/applications/{id}", s.validateMembershipLevel("read", s.getApplication)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications", s.validateMembershipLevel("read", s.listApplications)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/applications/{application}/releases", s.validateMembershipLevel("write", s.createRelease)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/applications/{application}/releases/{id}", s.validateMembershipLevel("read", s.getRelease)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications/{application}/releases/latest", s.validateMembershipLevel("read", s.getLatestRelease)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications/{application}/releases", s.validateMembershipLevel("read", s.listReleases)).Methods("GET")

	s.router.HandleFunc("/{project}/bundle", s.getBundle).Methods("GET")

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
		accessKeyValue, _, _ := r.BasicAuth()
		if accessKeyValue == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		accessKey, err := s.accessKeys.ValidateAccessKey(r.Context(), hash(accessKeyValue))
		if err != nil {
			log.WithError(err).Error("validate access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if accessKey == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(w, r, accessKey.UserID)
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

	application, err := s.applications.GetApplication(r.Context(), id)
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

	release, err := s.releases.CreateRelease(r.Context(), applicationID, createReleaseRequest.Config)
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

	release, err := s.releases.GetRelease(r.Context(), id)
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

	release, err := s.releases.GetLatestRelease(r.Context(), applicationID)
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

	releases, err := s.releases.ListReleases(r.Context(), applicationID)
	if err != nil {
		log.WithError(err).Error("list releases")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(releases)
}

func (s *Service) getBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]

	var bundle models.Bundle

	applications, err := s.applications.ListApplications(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list applications")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i, application := range applications {
		release, err := s.releases.GetLatestRelease(r.Context(), application.ID)
		if err != nil {
			log.WithError(err).Error("get latest release")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bundle.Applications = append(bundle.Applications, models.ApplicationAndLatestRelease{
			Application:   applications[i],
			LatestRelease: *release,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bundle)
}

func hash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}
