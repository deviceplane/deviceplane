import React, { useEffect } from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { endpoints, useRequest, useMutation } from '../api';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import { Column, Row, Form, Button } from '../components/core';

const validationSchema = yup.object().shape({
  email: validators.email.required(),
  password: yup.string().required(),
});

const Login = ({
  route: {
    data: { params, context },
  },
}) => {
  const [login, { data: loggedIn, error }] = useMutation(endpoints.login(), {
    errors: {
      403: 'Email confirmation required',
      default: 'Invalid credentials',
    },
  });
  const { data: user } = useRequest(loggedIn ? endpoints.user() : null);
  const navigation = useNavigation();

  useEffect(() => {
    if (user) {
      context.setCurrentUser(user);
      navigation.navigate(
        params.redirectTo ? decodeURIComponent(params.redirectTo) : '/projects'
      );
    }
  }, [user]);

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
        actions={[{ href: '/signup', title: 'Sign up', variant: 'tertiary' }]}
      >
        <Form
          onSubmit={login}
          validationSchema={validationSchema}
          alert={error}
          submitLabel="Login"
        >
          <Field
            required
            autoFocus
            autoComplete="on"
            name="email"
            type="email"
            label="Email address"
          />
          <Field required name="password" type="password" label="Password" />
        </Form>
        <Row marginTop={5}>
          <Button variant="text" href="/forgot" title="Forgot your password?" />
        </Row>
      </Card>
    </Column>
  );
};

export default Login;
