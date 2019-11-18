const development = {
  endpoint: 'http://localhost:8080/api',
  wsEndpoint: 'ws://localhost:8080/api'
};

const endpointBase = window.location.port
  ? `${window.location.hostname}:${window.location.port}/api`
  : `${window.location.hostname}/api`;
const production = {
  endpoint: `${window.location.protocol}//${endpointBase}`,
  wsEndpoint: window.location.protocol === 'http:' ? `ws://${endpointBase}` : `wss://${endpointBase}`,
};

const config = process.env.REACT_APP_ENVIRONMENT === 'development'
  ? development
  : production;

export default {
  agentVersion: '1.7.2',
  cliEndpoint: 'https://cli.deviceplane.com',
  ...config
};
