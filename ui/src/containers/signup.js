import React, { useState } from 'react';
import { toaster } from 'evergreen-ui';
import { useNavigation } from 'react-navi';
import useForm from 'react-hook-form';
import * as yup from 'yup';

import api from '../api';
import utils from '../utils';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import { Column, Button, Form, Text, Link } from '../components/core';
import validators from '../validators';

const validationSchema = yup.object().shape({
  firstName: yup
    .string()
    .required()
    .max(64),
  lastName: yup
    .string()
    .required()
    .max(64),
  company: yup.string().max(64),
  email: validators.email.required(),
  password: validators.password.required(),
});

const Signup = () => {
  const { register, handleSubmit, errors } = useForm({
    validationSchema,
  });
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      const response = await api.signup(data);
      navigation.navigate('/login');
      if (!response.data.registrationCompleted) {
        toaster.success(
          'Please check your email to confirm your registration.'
        );
      }
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger(
        'Something went wrong with your registration. Please contact us at support@deviceplane.com.'
      );
      console.log(error);
    }
  };

  return (
    <Column
      alignItems="center"
      flex={1}
      paddingY={[0, 9]}
      height={['initial', '100%']}
      overflow="auto"
      bg={['black', 'pageBackground']}
    >
      <Card
        logo
        size="medium"
        title="Sign up"
        actions={[{ href: '/login', title: 'Log in', variant: 'secondary' }]}
      >
        <Alert show={backendError} variant="error" description={backendError} />
        <Form
          onSubmit={e => {
            setBackendError(null);
            handleSubmit(submit)(e);
          }}
        >
          <Field
            required
            autoFocus
            autoComplete="on"
            autoCapitalize="on"
            label="First Name"
            name="firstName"
            ref={register}
            errors={errors.firstName}
            maxLength={64}
          />
          <Field
            required
            autoComplete="on"
            autoCapitalize="on"
            label="Last Name"
            name="lastName"
            ref={register}
            errors={errors.lastName}
            maxLength={64}
          />
          <Field
            autoComplete="on"
            autoCapitalize="on"
            label="Company"
            name="company"
            ref={register}
            errors={errors.company}
            maxLength={64}
          />
          <Field
            required
            autoComplete="on"
            type="email"
            label="Email"
            name="email"
            ref={register}
            errors={errors.email}
            maxLength={64}
          />
          <Field
            required
            type="password"
            label="Password"
            name="password"
            ref={register}
            errors={errors.password}
            hint="Password must be at least 8 characters."
            maxLength={64}
          />
          <Button title="Sign up" justifyContent="center" />
        </Form>
        <Text marginTop={5} fontSize={1} fontWeight={0}>
          By signing up you agree to the{' '}
          <Link href="https://deviceplane.com/legal/terms">
            Terms of Service
          </Link>{' '}
          and{' '}
          <Link href="https://deviceplane.com/legal/privacy">
            Privacy Policy
          </Link>
        </Text>
      </Card>
    </Column>
  );
};

export default Signup;
