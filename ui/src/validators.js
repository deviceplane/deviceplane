import * as yup from 'yup';

const nameRegex = /^[a-zA-Z0-9-]+$/;
//const noSpacesRegex = /^(?!.*\s).{1,100}$/;
//const usernameRegex = /^[a-zA-Z]+$/;

export default {
  name: yup
    .string()
    .max(100)
    .matches(nameRegex, {
      message: 'Can only include letters, numbers, and -.',
    }),
  email: yup
    .string()
    .email()
    .max(64),
  password: yup
    .string()
    .min(8, 'Password should be 8 characters minimum.')
    .max(100),
};
