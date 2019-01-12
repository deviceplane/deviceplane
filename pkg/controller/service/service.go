package service

import (
	"encoding/json"
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/gorilla/mux"
)

type Service struct {
	users        store.Users
	projects     store.Projects
	devices      store.Devices
	applications store.Applications
	releases     store.Releases
	router       *mux.Router
}

func NewService(
	users store.Users,
	projects store.Projects,
	devices store.Devices,
	applications store.Applications,
	releases store.Releases,
) *Service {
	s := &Service{
		users:        users,
		projects:     projects,
		devices:      devices,
		applications: applications,
		releases:     releases,
		router:       mux.NewRouter(),
	}

	s.router.HandleFunc("/users", s.createProject).Methods("POST")

	s.router.HandleFunc("/projects", s.createProject).Methods("POST")
	s.router.HandleFunc("/projects/{id}", s.getProject).Methods("GET")

	s.router.HandleFunc("/{project}/applications", s.createApplication).Methods("POST")
	s.router.HandleFunc("/{project}/applications/{id}", s.createProject).Methods("GET")

	s.router.HandleFunc("/{project}/applications/{application}/releases", s.createRelease).Methods("POST")

	s.router.HandleFunc("/{project}/bundle", s.getBundle).Methods("GET")

	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := s.users.CreateUser(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) createProject(w http.ResponseWriter, r *http.Request) {
	project, err := s.projects.CreateProject(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func (s *Service) getProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	project, err := s.projects.GetProject(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func (s *Service) createApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]

	application, err := s.applications.CreateApplication(r.Context(), projectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(application)
}

func (s *Service) createRelease(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationID := vars["application"]

	var createReleaseRequest models.CreateRelease
	if err := json.NewDecoder(r.Body).Decode(&createReleaseRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	release, err := s.releases.CreateRelease(r.Context(), applicationID, createReleaseRequest.Config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(release)
}

func (s *Service) getBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]

	var bundle models.Bundle

	applications, err := s.applications.ListApplications(r.Context(), projectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i, application := range applications {
		release, err := s.releases.GetLatestRelease(r.Context(), application.ID)
		if err != nil {
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
