import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import { useForm } from 'react-hook-form';
import * as yup from 'yup';

import api from '../api';
import utils from '../utils';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import Google from '../components/icons/google';
import Github from '../components/icons/github';
import {
  Box,
  Row,
  Column,
  Button,
  Form,
  Text,
  Link,
  toaster,
} from '../components/core';
import * as auth0 from '../lib/auth0';

const validationSchema = yup.object().shape({
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
        toaster.success('Check your email to confirm your registration.');
      }
    } catch (error) {
      setBackendError(
        utils.parseError(
          error,
          'Registration failed. Please contact us at support@deviceplane.com.'
        )
      );
      console.error(error);
    }
  };

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
        actions={[
          { href: '/login', title: 'Log in', variant: 'tertiaryGreen' },
        ]}
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
        <Text
          textAlign="center"
          marginTop={5}
          marginBottom={3}
          paddingTop={3}
          borderTop={0}
          borderColor="grays.5"
        >
          Sign up with
        </Text>
        <Row justifyContent="space-between">
          <Button
            variant="primaryGray"
            flex={1}
            marginRight={5}
            title={
              <Row
                alignItems="center"
                position="relative"
                flex={1}
                justifyContent="center"
              >
                <Box position="absolute" left={0}>
                  <Google />
                </Box>

                <Text marginLeft={3} color="white" textAlign="center">
                  Google
                </Text>
              </Row>
            }
            onClick={auth0.api.signup.google}
          />
          <Button
            variant="primaryGray"
            flex={1}
            title={
              <Row
                alignItems="center"
                position="relative"
                flex={1}
                justifyContent="center"
              >
                <Box position="absolute" left={0}>
                  <Github />
                </Box>
                <Text marginLeft={3} color="white" textAlign="center">
                  Github
                </Text>
              </Row>
            }
            onClick={auth0.api.signup.github}
          />
        </Row>
        <Text marginTop={5} fontSize={1} fontWeight={0}>
          By signing up you agree to the{' '}
          <Link href="https://deviceplane.com/terms/">Terms of Service</Link>{' '}
          and{' '}
          <Link href="https://deviceplane.com/privacy/">Privacy Policy</Link>
        </Text>
      </Card>
    </Column>
  );
};

export default Signup;
