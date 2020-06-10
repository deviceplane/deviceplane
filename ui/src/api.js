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

  loginSSO: data => post('loginsso', data),

  logout: () =>
    post('logout').then(() => {
      segment.reset();
    }),

  signup: ({ email, password }) =>
    post(`register`, {
      email,
      password,
    }),

  signupSSO: data => post('registersso', data),

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
    put(`projects/${projectId}/devices/${deviceId}/labels`, data).then(
      response => {
        segment.track('Device Label Added');
        return response;
      }
    ),

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
    }).then(response => {
      segment.track('Registration Token Created');

      return response;
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

  addMember: ({ projectId, data: { email, userId } }) =>
    post(`projects/${projectId}/memberships`, { email, userId }).then(
      response => {
        segment.track('Member Added');
        return response;
      }
    ),

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

  userAccessKeys: () => get(`useraccesskeys`),

  createUserAccessKey: () => post(`useraccesskeys`, {}),

  deleteUserAccessKey: ({ id }) => del(`useraccesskeys/${id}`),
};

export default api;
