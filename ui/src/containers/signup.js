import React from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { endpoints } from '../api';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import { Column, Form, Text, Link, toaster } from '../components/core';

const endpoint = endpoints.signup();
const validationSchema = yup.object().shape({
  fullName: yup
    .string()
    .required()
    .max(128),
  email: validators.email.required(),
  password: validators.password.required(),
});
const errorMessages = {
  default: 'Registration failed. Please contact us at support@deviceplane.com.',
};

const Signup = () => {
  const navigation = useNavigation();
  return (
    <Column
      alignItems="center"
      justifyContent="center"
      flex={1}
      paddingY={[0, 0, 5]}
      height={['initial', 'initial', '100%']}
      overflow="auto"
      bg={['black', 'black', 'pageBackground']}
    >
      <Card
        logo
        size="medium"
        title="Sign up"
        actions={[{ href: '/login', title: 'Log in', variant: 'tertiary' }]}
      >
        <Form
          endpoint={endpoint}
          onSuccess={data => {
            navigation.navigate('/login');
            if (!data.registrationCompleted) {
              toaster.success('Check your email to confirm your registration.');
            }
          }}
          onData={({ fullName }) => {
            let firstName, lastName;
            const firstSpace = fullName.indexOf(' ');
            if (firstSpace === -1) {
              firstName = fullName;
              lastName = ' ';
            } else {
              firstName = fullName.substr(0, firstSpace);
              lastName = fullName.substr(firstSpace + 1);
            }
            return {
              firstName,
              lastName,
            };
          }}
          errorMessages={errorMessages}
          validationSchema={validationSchema}
          submitLabel="Sign up"
        >
          <Field
            required
            autoFocus
            autoComplete="on"
            autoCapitalize="on"
            label="Full Name"
            name="fullName"
            maxLength={128}
          />
          <Field
            required
            autoComplete="on"
            type="email"
            label="Email"
            name="email"
            maxLength={64}
          />
          <Field
            required
            type="password"
            label="Password"
            name="password"
            hint="Password must be at least 8 characters."
            maxLength={64}
          />
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
