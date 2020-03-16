import auth0 from 'auth0-js';

const AUTH0_DOMAIN = 'deviceplane-dev.auth0.com'; // TODO: prod vs dev vs local
const AUTH0_CLIENT_ID = 'uvYKum4oRaWM4gDxgcGHZ73PDC1ZRcJf';
const AUTH0_CALLBACK_URL = 'http://localhost:3000' + '/login/sso-callback';

var webAuth = new auth0.WebAuth({
  domain: AUTH0_DOMAIN,
  clientID: AUTH0_CLIENT_ID,
  redirectUri: AUTH0_CALLBACK_URL, // TODO: current url + "/login/sso-callback";
  responseType: 'token id_token',
  scope: 'openid profile email',
  leeway: 60,
  connection: 'google',
});

function login() {
  webAuth.authorize({
    connection: 'google-oauth2',
  });
}

export { login };
