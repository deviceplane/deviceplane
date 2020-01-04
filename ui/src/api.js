import axios from 'axios';

import segment from './lib/segment';
import config from './config';

axios.defaults.withCredentials = true;

const url = path => `${config.endpoint}/${path}`;
const get = (path, ...rest) => axios.get(url(path), ...rest);
const post = (path, ...rest) => axios.post(url(path), ...rest);
const del = (path, ...rest) => axios.delete(url(path), ...rest);
const put = (path, ...rest) => axios.put(url(path), ...rest);
const patch = (path, ...rest) => axios.patch(url(path), ...rest);

const api = {
  login: ({ email, password }) => post('login', { email, password }),

  logout: () => post('logout'),

  signup: ({ email, password, firstName, lastName, company }) =>
    post(`register`, {
      email,
      password,
      firstName,
      lastName,
      company,
    }),

  completeRegistration: ({ registrationTokenValue }) =>
    post('completeregistration', { registrationTokenValue }),

  resetPassword: ({ email }) => post('recoverpassword', { email }),

  verifyPasswordResetToken: ({ token }) =>
    get(`passwordrecoverytokens/${token}`),

  updatePassword: ({ token, password }) =>
    post('changepassword', {
      passwordRecoveryTokenValue: token,
      password,
    }),

  user: () => get('me'),

  updateUser: data => patch('me', data),

  project: ({ projectId }) => get(`projects/${projectId}`),

  projects: () =>
    get(`memberships?full`).then(({ data }) =>
      data.map(({ project }) => project)
    ),

  createProject: data =>
    post(`projects`, data).then(response => {
      segment.track('Project Created');

      return response;
    }),

  updateProject: ({ projectId, data }) => put(`projects/${projectId}`, data),

  deleteProject: ({ projectId }) => del(`projects/${projectId}`),

  devices: ({ projectId, queryString = '' }) =>
    get(`projects/${projectId}/devices${queryString}`),

  device: ({ projectId, deviceId }) =>
    get(`projects/${projectId}/devices/${deviceId}?full`),

  updateDevice: ({ projectId, deviceId, data: { name } }) =>
    patch(`projects/${projectId}/devices/${deviceId}`, { name }),

  deleteDevice: ({ projectId, deviceId }) =>
    del(`projects/${projectId}/devices/${deviceId}`),

  reboot: ({ projectId, deviceId }) =>
    post(`projects/${projectId}/devices/${deviceId}/reboot`, {}),

  addDeviceLabel: ({ projectId, deviceId, data }) =>
    put(`projects/${projectId}/devices/${deviceId}/labels`, data),

  removeDeviceLabel: ({ projectId, deviceId, labelId }) =>
    del(`projects/${projectId}/devices/${deviceId}/labels/${labelId}`),

  defaultRegistrationToken: ({ projectId }) =>
    get(`projects/${projectId}/deviceregistrationtokens/default`),

  registrationToken: ({ projectId, tokenId }) =>
    get(`projects/${projectId}/deviceregistrationtokens/${tokenId}?full`),

  registrationTokens: ({ projectId }) =>
    get(`projects/${projectId}/deviceregistrationtokens?full`),

  createRegistrationToken: ({
    projectId,
    data: { name, description, maxRegistrations },
  }) =>
    post(`projects/${projectId}/deviceregistrationtokens`, {
      name,
      description,
      maxRegistrations: Number.parseInt(maxRegistrations),
    }),

  updateRegistrationToken: ({
    projectId,
    tokenId,
    data: { name, description, maxRegistrations, settings },
  }) =>
    put(`projects/${projectId}/deviceregistrationtokens/${tokenId}`, {
      name,
      description,
      maxRegistrations: Number.parseInt(maxRegistrations),
      settings,
    }),

  deleteRegistrationToken: ({ projectId, tokenId }) =>
    del(`projects/${projectId}/deviceregistrationtokens/${tokenId}`),

  addRegistrationTokenLabel: ({ projectId, tokenId, data }) =>
    put(
      `projects/${projectId}/deviceregistrationtokens/${tokenId}/labels`,
      data
    ),

  removeRegistrationTokenLabel: ({ projectId, tokenId, labelId }) =>
    del(
      `projects/${projectId}/deviceregistrationtokens/${tokenId}/labels/${labelId}`
    ),

  applications: ({ projectId }) =>
    get(`projects/${projectId}/applications?full`),

  application: ({ projectId, applicationId }) =>
    get(`projects/${projectId}/applications/${applicationId}?full`),

  createApplication: ({ projectId, data: { name, description } }) =>
    post(`projects/${projectId}/applications`, { name, description }).then(
      response => {
        segment.track('Application Created');
        return response;
      }
    ),

  updateApplication: ({ projectId, applicationId, data }) =>
    patch(`projects/${projectId}/applications/${applicationId}`, data),

  deleteApplication: ({ projectId, applicationId }) =>
    del(`projects/${projectId}/applications/${applicationId}`),

  roles: ({ projectId }) => get(`projects/${projectId}/roles`),

  role: ({ projectId, roleId }) => get(`projects/${projectId}/roles/${roleId}`),

  createRole: ({ projectId, data: { name, description, config } }) =>
    post(`projects/${projectId}/roles`, { name, description, config }).then(
      response => {
        segment.track('Role Created');

        return response;
      }
    ),

  updateRole: ({ projectId, roleId, data: { name, description, config } }) =>
    put(`projects/${projectId}/roles/${roleId}`, { name, description, config }),

  deleteRole: ({ projectId, roleId }) =>
    del(`projects/${projectId}/roles/${roleId}`),

  memberships: ({ projectId }) => get(`projects/${projectId}/memberships?full`),

  membership: ({ projectId, userId }) =>
    get(`projects/${projectId}/memberships/${userId}?full`),

  addMember: ({ projectId, data: { email } }) =>
    post(`projects/${projectId}/memberships`, { email }).then(response => {
      segment.track('Member Added');
      return response;
    }),

  removeMember: ({ projectId, userId }) =>
    del(`projects/${projectId}/memberships/${userId}`),

  addMembershipRoleBindings: ({ projectId, userId, roleId }) =>
    post(
      `projects/${projectId}/memberships/${userId}/roles/${roleId}/membershiprolebindings`,
      {}
    ),

  removeMembershipRoleBindings: ({ projectId, userId, roleId }) =>
    del(
      `projects/${projectId}/memberships/${userId}/roles/${roleId}/membershiprolebindings`
    ),

  serviceAccounts: ({ projectId }) =>
    get(`projects/${projectId}/serviceaccounts?full`),

  serviceAccount: ({ projectId, serviceId }) =>
    get(`projects/${projectId}/serviceaccounts/${serviceId}?full`),

  createServiceAccount: ({ projectId, data }) =>
    post(`projects/${projectId}/serviceaccounts`, data).then(response => {
      segment.track('Service Account Created');
      return response;
    }),

  updateServiceAccount: ({
    projectId,
    serviceId,
    data: { name, description },
  }) =>
    put(`projects/${projectId}/serviceaccounts/${serviceId}`, {
      name,
      description,
    }),

  deleteServiceAccount: ({ projectId, serviceId }) =>
    del(`projects/${projectId}/serviceaccounts/${serviceId}`),

  addServiceAccountRoleBindings: ({ projectId, serviceId, roleId }) =>
    post(
      `projects/${projectId}/serviceaccounts/${serviceId}/roles/${roleId}/serviceaccountrolebindings`,
      {}
    ),

  removeServiceAccountRoleBindings: ({ projectId, serviceId, roleId }) =>
    del(
      `projects/${projectId}/serviceaccounts/${serviceId}/roles/${roleId}/serviceaccountrolebindings`
    ),

  serviceAccountAccessKeys: ({ projectId, serviceId }) =>
    get(
      `projects/${projectId}/serviceaccounts/${serviceId}/serviceaccountaccesskeys`
    ),

  createServiceAccountAccessKey: ({ projectId, serviceId }) =>
    post(
      `projects/${projectId}/serviceaccounts/${serviceId}/serviceaccountaccesskeys`,
      {}
    ),

  deleteServiceAccountAccessKey: ({ projectId, serviceId, accessKeyId }) =>
    del(
      `projects/${projectId}/serviceaccounts/${serviceId}/serviceaccountaccesskeys/${accessKeyId}`
    ),

  releases: ({ projectId, applicationId }) =>
    get(`projects/${projectId}/applications/${applicationId}/releases?full`),

  release: ({ projectId, applicationId, releaseId }) =>
    get(
      `projects/${projectId}/applications/${applicationId}/releases/${releaseId}?full`
    ),

  createRelease: ({ projectId, applicationId, data: { rawConfig } }) =>
    post(`projects/${projectId}/applications/${applicationId}/releases`, {
      rawConfig,
    }).then(response => {
      segment.track('Release Created');
      return response;
    }),

  latestReleases: ({ projectId, applicationId }) =>
    get(`projects/${projectId}/applications/${applicationId}/releases/latest`),

  userAccessKeys: () => get(`useraccesskeys`),

  createUserAccessKey: () => post(`useraccesskeys`, {}),

  deleteUserAccessKey: ({ id }) => del(`useraccesskeys/${id}`),

  hostMetrics: ({ projectId, deviceId }) =>
    get(`projects/${projectId}/devices/${deviceId}/metrics/host`),

  serviceMetrics: ({ projectId, deviceId, applicationId, serviceId }) =>
    get(
      `projects/${projectId}/devices/${deviceId}/applications/${applicationId}/services/${serviceId}/metrics`
    ),

  projectMetricsConfig: ({ projectId }) =>
    get(`projects/${projectId}/configs/project-metrics-config`),

  updateProjectMetricsConfig: ({ projectId, data }) =>
    put(`projects/${projectId}/configs/project-metrics-config`, {
      exposedMetrics: data,
    }),

  deviceMetricsConfig: ({ projectId }) =>
    get(`projects/${projectId}/configs/device-metrics-config`),

  updateDeviceMetricsConfig: ({ projectId, data }) =>
    put(`projects/${projectId}/configs/device-metrics-config`, {
      exposedMetrics: data,
    }),

  serviceMetricsConfig: ({ projectId }) =>
    get(`projects/${projectId}/configs/service-metrics-config`),

  updateServiceMetricsConfig: ({ projectId, data }) =>
    put(`projects/${projectId}/configs/service-metrics-config`, data),
};

export default api;
