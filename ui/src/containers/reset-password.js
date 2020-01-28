import React, { useState, useEffect } from 'react';
import { useNavigation } from 'react-navi';
import useForm from 'react-hook-form';
import * as yup from 'yup';
import { toaster } from 'evergreen-ui';

import api from '../api';
import utils from '../utils';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import { Column, Button, Form } from '../components/core';

const validationSchema = yup.object().shape({
  password: validators.password.required(),
});

const isTokenExpired = expiration => {
  const expiratonTime = new Date(expiration).getTime();
  const currentTime = new Date().getTime();
  return currentTime > expiratonTime;
};

const PasswordRecovery = ({
  route: {
    data: {
      params: { token },
    },
  },
}) => {
  const navigation = useNavigation();
  const { register, handleSubmit, errors } = useForm({ validationSchema });
  const [invalidToken, setInvalidToken] = useState();
  const [backendError, setBackendError] = useState();

  useEffect(() => {
    api
      .verifyPasswordResetToken({ token })
      .then(response => {
        if (isTokenExpired(response.data.expiresAt)) {
          setInvalidToken(true);
        }
      })
      .catch(error => {
        setInvalidToken(true);
        console.log(error);
      });
  }, []);

  const submit = async ({ password }) => {
    try {
      await api.updatePassword({ password, token });
      toaster.success('Password changed successfully.');
      navigation.navigate(`/login`);
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger(
        'Something went wrong with changing your password. Please contact us at support@deviceplane.com.'
      );
      console.log(error);
    }
  };

  if (invalidToken) {
    return (
      <Column flex={1} alignItems="center" paddingTop={9}>
        <Card title="Recover Password" logo size="large">
          <Alert
            show
            variant="error"
            description="Your recovery token has expired. Please reset your password again."
          />
          <Button marginTop={3} href="/forgot" title="Reset your password" />
        </Card>
      </Column>
    );
  }

  return (
    <Column
      flex={1}
      alignItems="center"
      justifyContent="center"
      paddingY={[0, 6]}
      height="100%"
      overflow="auto"
    >
      <Card title="Recover Password" logo size="large">
        <Alert show={backendError} variant="error" description={backendError} />
        <Form
          onSubmit={e => {
            setBackendError(null);
            handleSubmit(submit)(e);
          }}
        >
          <Field
            required
            type="password"
            label="New Password"
            name="password"
            ref={register}
            errors={errors.password}
            hint="Password must be at least 8 characters."
          />
          <Button marginTop={3} title="Submit" type="submit" />
        </Form>
      </Card>
    </Column>
  );
};

export default PasswordRecovery;
