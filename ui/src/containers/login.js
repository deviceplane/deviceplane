import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import { useForm } from 'react-hook-form';
import * as yup from 'yup';

import api from '../api';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import Google from '../components/icons/google';
import Github from '../components/icons/github';
import { Box, Column, Row, Form, Button, Text } from '../components/core';
import * as auth0 from '../lib/auth0';

const validationSchema = yup.object().shape({
  email: validators.email.required(),
  password: yup.string().required(),
});

const Login = ({
  route: {
    data: { params, context },
  },
}) => {
  const { register, handleSubmit, errors } = useForm({
    validationSchema,
  });
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.login(data);
      const response = await api.user();
      context.setCurrentUser(response.data);
      navigation.navigate(
        params.redirectTo ? decodeURIComponent(params.redirectTo) : '/projects'
      );
    } catch (error) {
      if (error.response && error.response.status === 403) {
        setBackendError('Email confirmation required');
      } else {
        setBackendError('Invalid credentials');
      }
    }
  };

  return (
    <Column
      flex={1}
      alignItems="center"
      justifyContent="center"
      paddingY={[0, 0, 6]}
      height="100%"
      overflow="auto"
      bg={['black', 'black', 'pageBackground']}
    >
      <Card
        logo
        title="Log in"
        size="medium"
        actions={[
          { href: '/signup', title: 'Sign up', variant: 'tertiaryGreen' },
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
            autoFocus
            autoComplete="on"
            ref={register}
            name="email"
            type="email"
            label="Email address"
            errors={errors.email}
          />

          <Field
            required
            ref={register}
            name="password"
            type="password"
            label="Password"
            errors={errors.password}
            marginBottom={2}
          />
          <Row marginBottom={5}>
            <Button
              variant="text"
              href="/forgot"
              title="Forgot your password?"
            />
          </Row>

          <Button justifyContent="center" title="Log in" />
        </Form>
        <Text
          textAlign="center"
          marginTop={5}
          marginBottom={3}
          paddingTop={3}
          borderTop={0}
          borderColor="grays.5"
        >
          Log in with
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
            onClick={auth0.api.login.google}
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
            onClick={auth0.api.login.github}
          />
        </Row>
      </Card>
    </Column>
  );
};

export default Login;
