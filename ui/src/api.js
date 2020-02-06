import { useState } from 'react';
import axios from 'axios';
import fetch from 'unfetch';
import useSWR, { trigger } from 'swr';

import segment from './lib/segment';
import config from './config';
import utils from './utils';
import { toaster } from './components/core/toast';

axios.defaults.withCredentials = true;

const url = path => `${config.endpoint}/${path}`;
const get = (path, ...rest) => axios.get(url(path), ...rest);
const post = (path, ...rest) => axios.post(url(path), ...rest);
const del = (path, ...rest) => axios.delete(url(path), ...rest);
const put = (path, ...rest) => axios.put(url(path), ...rest);
const patch = (path, ...rest) => axios.patch(url(path), ...rest);

const responseHandler = async response => {
  let json, text;

  try {
    json = await response.json();
  } catch {}

  if (!json) {
    try {
      text = await response.text();
    } catch {}
  }

  if (response.status >= 400) {
    if (response.status >= 500) {
      toaster.danger('This service is currently down, please try again later.');
    }
    return {
      error: json || utils.capitalize(text) || 'Default error',
    };
  }

  return {
    data: json || text,
    success: true,
    headers: response.headers,
  };
};

export const useRequest = (endpoint, config = {}) => {
  const { data: response, error } = useSWR(
    endpoint,
    async () => {
      const res = await fetch(endpoint);
      const resp = await responseHandler(res);
      return resp;
    },
    config
  );

  return {
    data: response && response.data,
    headers: response && response.headers,
    error,
  };
};

export const useMutation = (endpoint, config = {}) => {
  const [result, setResult] = useState({});

  const mutate = async (body = {}) => {
    const res = await fetch(endpoint, {
      method: 'POST' || config.method,
      headers:
        {
          'Content-Type': 'application/json',
        } || config.headers,
      body: JSON.stringify(body),
    });
    console.log('response', res);
    const result = await responseHandler(res);
    setResult(result);
    if (config.triggers) {
      config.triggers.forEach(key => trigger(key));
    }
  };

  console.log(result);

  return [mutate, result];
};

export const endpoints = {
  login: () => url('login'),

  user: () => url('me'),

  projects: () => url(`memberships?full`),

  project: ({ projectId }) => url(`projects/${projectId}`),

  createProject: () => url(`projects`),

  applications: ({ projectId }) =>
    url(`projects/${projectId}/applications?full`),

  application: ({ projectId, applicationId }) =>
    url(`projects/${projectId}/applications/${applicationId}?full`),

  releases: ({ projectId, applicationId }) =>
    url(`projects/${projectId}/applications/${applicationId}/releases?full`),

  release: ({ projectId, applicationId, releaseId }) =>
    url(
      `projects/${projectId}/applications/${applicationId}/releases/${releaseId}?full`
    ),

  memberships: ({ projectId }) => url(`projects/${projectId}/memberships?full`),

  membership: ({ projectId, userId }) =>
    url(`projects/${projectId}/memberships/${userId}?full`),

  roles: ({ projectId }) => url(`projects/${projectId}/roles`),

  role: ({ projectId, roleId }) => url(`projects/${projectId}/roles/${roleId}`),

  serviceAccounts: ({ projectId }) =>
    url(`projects/${projectId}/serviceaccounts?full`),

  serviceAccount: ({ projectId, serviceId }) =>
    url(`projects/${projectId}/serviceaccounts/${serviceId}?full`),

  devices: ({ projectId, queryString }) =>
    url(`projects/${projectId}/devices${queryString}`),

  device: ({ projectId, deviceId }) =>
    url(`projects/${projectId}/devices/${deviceId}?full`),

  registrationTokens: ({ projectId }) =>
    url(`projects/${projectId}/deviceregistrationtokens?full`),

  registrationToken: ({ projectId, tokenId }) =>
    url(`projects/${projectId}/deviceregistrationtokens/${tokenId}?full`),

  createRegistrationToken: ({ projectId }) =>
    url(`projects/${projectId}/deviceregistrationtokens`),
};

const api = {
  logout: () => post('logout'),

  signup: ({ email, password, firstName, lastName }) =>
    post(`register`, {
      email,
      password,
      firstName,
      lastName,
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

  updateProject: ({ projectId, data }) => put(`projects/${projectId}`, data),

  deleteProject: ({ projectId }) => del(`projects/${projectId}`),

  devices: ({ projectId, queryString = '' }) =>
    get(`projects/${projectId}/devices${queryString}`),

  scheduledDevices: ({ projectId, applicationId, schedulingRule, search }) =>
    get(
      `projects/${projectId}/devices/previewscheduling/${applicationId}?search=${encodeURIComponent(
        search
      )}&schedulingRule=${encodeURIComponent(
        btoa(JSON.stringify(schedulingRule))
      )}`
    ),

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

  addEnvironmentVariable: ({ projectId, deviceId, data }) =>
    put(
      `projects/${projectId}/devices/${deviceId}/environmentvariables`,
      data
    ).then(response => {
      segment.track('Environment Variable Added');
      return response;
    }),

  removeEnvironmentVariable: ({ projectId, deviceId, key }) =>
    del(
      `projects/${projectId}/devices/${deviceId}/environmentvariables/${key}`
    ),

  defaultRegistrationToken: ({ projectId }) =>
    get(`projects/${projectId}/deviceregistrationtokens/default`),

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

  addRegistrationTokenEnvironmentVariable: ({ projectId, tokenId, data }) =>
    put(
      `projects/${projectId}/deviceregistrationtokens/${tokenId}/environmentvariables`,
      data
    ).then(response => {
      segment.track('Registration Token Environment Variable Added');
      return response;
    }),

  removeRegistrationTokenEnvironmentVariable: ({ projectId, tokenId, key }) =>
    del(
      `projects/${projectId}/deviceregistrationtokens/${tokenId}/environmentvariables/${key}`
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

  imagePullProgress: ({ projectId, deviceId, applicationId, serviceId }) =>
    get(
      `projects/${projectId}/devices/${deviceId}/applications/${applicationId}/services/${serviceId}/imagepullprogress`
    ),
};

export default api;
