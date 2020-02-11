package service

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/deviceplane/deviceplane/pkg/agent/service/client"
	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/controller/authz"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
)

func (s *Service) initiateDeviceConnection(w http.ResponseWriter, r *http.Request) {
	s.withDeviceAuth(w, r, func(project *models.Project, device *models.Device) {
		s.withHijackedWebSocketConnection(w, r, func(clientConn net.Conn) {
			s.connman.Set(project.ID+device.ID, clientConn)
		})
	})
}

var currentSSHCount int64

const currentSSHCountName = "internal.current_ssh_connection_count"

func (s *Service) initiateSSH(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionSSH,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					s.withHijackedWebSocketConnection(w, r, func(clientConn net.Conn) {
						s.withDeviceConnection(w, r, project, device, func(deviceConn net.Conn) {
							err := client.InitiateSSH(r.Context(), clientConn)
							if err != nil {
								http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
								return
							}

							sshCount := atomic.AddInt64(&currentSSHCount, 1)
							s.st.Gauge(currentSSHCountName, float64(sshCount), utils.InternalTags(project.ID), 1)
							defer func() {
								sshCount := atomic.AddInt64(&currentSSHCount, -1)
								s.st.Gauge(currentSSHCountName, float64(sshCount), utils.InternalTags(project.ID), 1)
							}()

							go io.Copy(deviceConn, clientConn)
							io.Copy(clientConn, deviceConn)
						})
					})
				})
			},
		)
	})
}

func (s *Service) initiateReboot(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionReboot,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					s.withDeviceConnection(w, r, project, device, func(deviceConn net.Conn) {
						resp, err := client.InitiateReboot(r.Context(), deviceConn)
						if err != nil {
							http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
							return
						}

						utils.ProxyResponseFromDevice(w, resp)
					})
				})
			},
		)
	})
}

func (s *Service) deviceDebug(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionGetMetrics,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					s.withDeviceConnection(w, r, project, device, func(deviceConn net.Conn) {
						path := r.URL.EscapedPath()
						dIndex := strings.Index(path, "/debug/")
						if dIndex == -1 {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						debugPath := path[dIndex:]

						req, err := http.NewRequestWithContext(
							r.Context(),
							r.Method,
							debugPath,
							r.Body,
						)
						if err != nil {
							http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
							return
						}

						if err := req.Write(deviceConn); err != nil {
							http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
							return
						}

						resp, err := http.ReadResponse(bufio.NewReader(deviceConn), req)
						if err != nil {
							http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
							return
						}

						utils.ProxyResponseFromDevice(w, resp)
					})
				})
			},
		)
	})
}

func (s *Service) imagePullProgress(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionGetImagePullProgress,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					s.withDeviceConnection(w, r, project, device, func(deviceConn net.Conn) {
						vars := mux.Vars(r)
						applicationID := vars["application"]
						service := vars["service"]

						resp, err := client.GetImagePullProgress(r.Context(), deviceConn, applicationID, service)
						if err != nil {
							http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
							return
						}

						utils.ProxyResponseFromDevice(w, resp)
					})
				})
			},
		)
	})
}

func (s *Service) hostMetrics(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionGetMetrics,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					s.withDeviceConnection(w, r, project, device, func(deviceConn net.Conn) {
						resp, err := client.GetDeviceMetrics(r.Context(), deviceConn)
						if err != nil {
							http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
							return
						}

						utils.ProxyResponseFromDevice(w, resp)
					})
				})
			},
		)
	})
}

func (s *Service) agentMetrics(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionGetMetrics,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					s.withDeviceConnection(w, r, project, device, func(deviceConn net.Conn) {
						resp, err := client.GetAgentMetrics(r.Context(), deviceConn)
						if err != nil {
							http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
							return
						}

						utils.ProxyResponseFromDevice(w, resp)
					})
				})
			},
		)
	})
}

func (s *Service) serviceMetrics(w http.ResponseWriter, r *http.Request) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		s.validateAuthorization(
			authz.ResourceDevices, authz.ActionGetServiceMetrics,
			w, r,
			user, serviceAccount,
			func(project *models.Project) {
				s.withDevice(w, r, project, func(device *models.Device) {
					s.withApplication(w, r, project, func(application *models.Application) {
						s.withDeviceConnection(w, r, project, device, func(deviceConn net.Conn) {
							vars := mux.Vars(r)
							service := vars["service"]

							serviceMetricEndpointConfig, exists := application.MetricEndpointConfigs[service]
							if !exists {
								serviceMetricEndpointConfig.Port = models.DefaultMetricPort
								serviceMetricEndpointConfig.Path = models.DefaultMetricPath
							}

							resp, err := client.GetServiceMetrics(
								r.Context(), deviceConn, application.ID, service,
								serviceMetricEndpointConfig.Path, serviceMetricEndpointConfig.Port,
							)
							if err != nil {
								http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
								return
							}

							utils.ProxyResponseFromDevice(w, resp)
						})
					})
				})
			},
		)
	})
}
