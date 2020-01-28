import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import useForm from 'react-hook-form';
import * as yup from 'yup';

import api from '../api';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
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
      setBackendError('Invalid credentials');
      console.log(error);
    }
  };

  return (
    <Column
      flex={1}
      alignItems="center"
      justifyContent="center"
      paddingY={[0, 6]}
      height="100%"
      overflow="auto"
      bg={['black', 'pageBackground']}
    >
      <Card
        logo
        title="Log in"
        size="medium"
        actions={[{ href: '/signup', title: 'Sign up', variant: 'secondary' }]}
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
          />

          <Button justifyContent="center" title="Log in" />
        </Form>
        <Row marginTop={5}>
          <Button variant="text" href="/forgot" title="Forgot your password?" />
        </Row>
      </Card>
    </Column>
  );
};

export default Login;
