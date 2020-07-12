const development = {
  endpoint: 'http://localhost:8080/api',
  wsEndpoint: 'ws://localhost:8080/api',
};

const baseURL = window.location.port
  ? `${window.location.hostname}:${window.location.port}`
  : `${window.location.hostname}`;

const endpointBase = baseURL + '/api';

const production = {
  endpoint: `${window.location.protocol}//${endpointBase}`,
  wsEndpoint:
    window.location.protocol === 'http:'
      ? `ws://${endpointBase}`
      : `wss://${endpointBase}`,
};

const frontendURL = `${window.location.protocol}//${baseURL}`;

const config =
  process.env.NODE_ENV === 'development' ? development : production;

const auth0_domain = process.env.AUTH0_DOMAIN
  ? new URL(process.env.AUTH0_DOMAIN).host
  : '';
const auth0_client_id = process.env.AUTH0_AUDIENCE || '';

const auth0_login_callback_url = frontendURL + '/login/sso-callback';
const auth0_signup_callback_url = frontendURL + '/signup/sso-callback';

export default {
  agentVersion: '1.16.0',
  cliEndpoint: 'https://downloads.deviceplane.com/cli',
  auth0_domain,
  auth0_client_id,
  auth0_login_callback_url,
  auth0_signup_callback_url,
  ...config,
};
