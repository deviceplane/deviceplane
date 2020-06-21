const id = (key, scope) => `deviceplane:${scope ? `${scope}:${key}` : key}`;

const storage = {
  get: (key, scope) => {
    if (!scope) {
      scope = 'global';
    }
    try {
      return JSON.parse(localStorage.getItem(id(key, scope)));
    } catch (e) {
      return null;
    }
  },
  set: (key, value, scope) => {
    if (!scope) {
      scope = 'global';
    }
    try {
      localStorage.setItem(id(key, scope), JSON.stringify(value));
      return true;
    } catch (e) {
      return false;
    }
  },
  remove: (key, scope) => {
    if (!scope) {
      scope = 'global';
    }
    try {
      localStorage.removeItem(id(key, scope));
      return true;
    } catch (e) {
      return false;
    }
  },
};

export default storage;
