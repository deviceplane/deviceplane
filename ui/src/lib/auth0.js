import auth0 from 'auth0-js';
import config from '../config';

console.log(config.auth0_login_callback_url);

var login = new auth0.WebAuth({
  domain: config.auth0_domain,
  clientID: config.auth0_client_id,
  redirectUri: config.auth0_login_callback_url,
  responseType: 'token id_token',
  scope: 'openid profile email',
  leeway: 60,
});

var signup = new auth0.WebAuth({
  domain: config.auth0_domain,
  clientID: config.auth0_client_id,
  redirectUri: config.auth0_signup_callback_url,
  responseType: 'token id_token',
  scope: 'openid profile email',
  leeway: 60,
});

const api = {
  login: {
    google: () =>
      login.authorize({
        connection: 'google-oauth2',
      }),
    github: () =>
      login.authorize({
        connection: 'github',
      }),
  },
  signup: {
    google: () =>
      signup.authorize({
        connection: 'google-oauth2',
      }),
    github: () =>
      signup.authorize({
        connection: 'github',
      }),
  },
};

export { api };
