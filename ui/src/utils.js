const uuidv4 = require("uuid/v4");

const emailRegex = /^([a-zA-Z0-9_.-])+@(([a-zA-Z0-9-])+\.)+([a-zA-Z0-9]{2,4})+$/;
const usernameRegex = /^[a-zA-Z]+$/;
const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?!.*\s).{8,100}$/;
const nameRegex = /^[a-z0-9-]+$/;
const minEightCharRegex = /^.{8,100}$/;
const lowercaseRegex = /^(?=.*[a-z]).{1,100}$/;
const uppercaseRegex = /^(?=.*[A-Z]).{1,100}$/;
const noSpacesRegex = /^(?!.*\s).{1,100}$/;

export default {
  emailRegex: emailRegex,
  usernameRegex: usernameRegex,
  passwordRegex: passwordRegex,
  nameRegex: nameRegex,
  minEightCharRegex: minEightCharRegex,
  lowercaseRegex: lowercaseRegex,
  uppercaseRegex: uppercaseRegex,
  noSpacesRegex: noSpacesRegex,

  checkName: (objectName, name) => {
    if (name === '') {
      return 'Please enter a name.'
    }

    if (!nameRegex.test(name)) {
      return 'Invalid ' + objectName + ' name. Name can only contain lowercase letters, numbers and -.'
    }

    if (name.length > 100) {
      return 'The ' + objectName + ' name must be less than 100 characters.'
    }

    return null
  },

  customValidate: (name, validationRegex, containsMsg, maxLength, value) => {
    if (value === '') {
      return `Please enter a ${name}.`
    }

    if (!validationRegex.test(value)) {
      return `Invalid  ${name}. ${name} can only contain ${containsMsg}.`
    }

    if (maxLength && value.length > maxLength) {
      return `The ${name} must be less than ${maxLength} characters.`
    }

    return null
  },

  checkConfig: (objectName, config) => {
    if (config === '') {
      return 'Please enter a ' + objectName + '.'
    }

    return null;
  },

  convertErrorMessage: (errorMessage) => {
    if (errorMessage !== '') {
      return errorMessage.charAt(0).toUpperCase() + errorMessage.slice(1).trim() + '.'
    }
    return ''
  },

  is4xx: (status) => {
    return status >= 400 && status < 500
  },

  randomClassName: () => {
    return 'rcn_' + uuidv4().replace(/-/g, '').substring(0, 10);
  },
};