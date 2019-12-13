const uuidv4 = require('uuid/v4');

const capitalize = message => message.replace(/^\w/, c => c.toUpperCase());

const randomClassName = () => {
  return (
    'rcn_' +
    uuidv4()
      .replace(/-/g, '')
      .substring(0, 10)
  );
};

const deepClone = object => {
  return JSON.parse(JSON.stringify(object));
};

const deepEqual = (a, b) => JSON.stringify(a) === JSON.stringify(b);

const parseError = error => {
  if (
    error &&
    typeof error === 'object' &&
    error.response &&
    typeof error.response === 'object' &&
    error.response.data &&
    typeof error.response.data === 'string'
  ) {
    return capitalize(error.response.data);
  }

  return null;
};

export default {
  randomClassName,
  deepClone,
  deepEqual,
  capitalize,
  parseError,
};
