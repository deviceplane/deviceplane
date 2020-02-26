const uuidv4 = require('uuid/v4');

const id = () =>
  Math.random()
    .toString(36)
    .substr(2, 9);

const capitalize = message => message.replace(/^\w/, c => c.toUpperCase());

const randomClassName = () =>
  'rcn_' +
  uuidv4()
    .replace(/-/g, '')
    .substring(0, 10);

const deepClone = object => JSON.parse(JSON.stringify(object));

const deepEqual = (a, b) => JSON.stringify(a) === JSON.stringify(b);

const parseError = (
  error,
  placeholder = 'Something went wrong. Please contact us at support@deviceplane.com.'
) => {
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
  return placeholder;
};

export default {
  randomClassName,
  deepClone,
  deepEqual,
  capitalize,
  parseError,
  id,
};
