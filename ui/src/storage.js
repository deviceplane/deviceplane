const id = (key, scope) => `deviceplane:${scope ? `${scope}:${key}` : key}`;

const storage = {
  get: (key, scope) => {
    try {
      return JSON.parse(localStorage.getItem(id(key, scope)));
    } catch (e) {
      return null;
    }
  },
  set: (key, value, scope) => {
    try {
      localStorage.setItem(id(key, scope), JSON.stringify(value));
      return true;
    } catch (e) {
      return false;
    }
  },
  remove: (key, scope) => {
    try {
      localStorage.removeItem(id(key, scope));
      return true;
    } catch (e) {
      return false;
    }
  },
};

export default storage;
