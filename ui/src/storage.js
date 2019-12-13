const id = key => `deviceplane:${key}`;

const storage = {
  get: key => {
    try {
      return JSON.parse(localStorage.getItem(id(key)));
    } catch (e) {
      return null;
    }
  },
  set: (key, value) => {
    try {
      localStorage.setItem(id(key), JSON.stringify(value));
      return true;
    } catch (e) {
      return false;
    }
  },
  remove: key => {
    try {
      localStorage.removeItem(id(key));
      return true;
    } catch (e) {
      return false;
    }
  },
};

export default storage;
