// Modified from https://auth0.com/docs/quickstart/spa/react?framed=1&sq=1#configure-auth0

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
});

const PREFIX = 'auth0.';
function get(key) {
  return localStorage.getItem(PREFIX + key);
}
function set(key, value) {
  return localStorage.setItem(PREFIX + key, value);
}
function remove(key) {
  return localStorage.removeItem(PREFIX + key);
}

// Set the time that the access token will expire at
function setSession(authResult) {
  var expiresAt = JSON.stringify(
    authResult.expiresIn * 1000 + new Date().getTime()
  );
  set('raw_session', JSON.stringify(authResult));
  set('access_token', authResult.accessToken);
  set('id_token', authResult.idToken);
  set('expires_at', expiresAt);
}

function login() {
  webAuth.authorize();
}

// Remove tokens and expiry time from localStorage
function logout() {
  remove('raw_session');
  remove('access_token');
  remove('id_token');
  remove('expires_at');
}

// Check whether the current time is past the
// access token's expiry time
function isAuthenticated() {
  try {
    var expiresAt = JSON.parse(get('expires_at'));
    return new Date().getTime() < expiresAt;
  } catch (e) {
    remove('expires_at');
    return false;
  }
}

function rawSession() {
  if (!isAuthenticated()) {
    return false;
  }

  try {
    return JSON.parse(get('raw_session'));
  } catch (e) {
    remove('raw_session');
    return undefined;
  }
}

function handleAuthentication(callback) {
  webAuth.parseHash(function(err, authResult) {
    if (authResult && authResult.accessToken && authResult.idToken) {
      setSession(authResult);
      callback(authResult);
    } else if (err) {
      callback(undefined, err.error);
    }
  });
}

export { login, logout, isAuthenticated, handleAuthentication, rawSession };
