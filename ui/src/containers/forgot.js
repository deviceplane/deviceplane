import React, { useState } from 'react';
import useForm from 'react-hook-form';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';
import { toaster } from 'evergreen-ui';

import api from '../api';
import utils from '../utils';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import { Text, Row, Column, Button, Form } from '../components/core';

const validationSchema = yup.object().shape({
  email: validators.email.required(),
});

const PasswordReset = () => {
  const { register, handleSubmit, errors } = useForm({ validationSchema });
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    setBackendError(null);
    try {
      await api.resetPassword(data);
      navigation.navigate(`/login`);
      toaster.success(
        'Password recovery email sent. Please check your email to reset your password.'
      );
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger(
        'There was an error with your e-mail. Please contact us at info@deviceplane.com.'
      );
      console.log(error);
    }
  };

  return (
    <Column
      flex={1}
      alignItems="center"
      paddingY={[0, 9]}
      height="100%"
      overflow="auto"
      bg={['black', 'pageBackground']}
    >
      <Card logo size="medium" title="Reset Password">
        <Alert
          show={backendError}
          variant="error"
          description="There is no user with that email address."
        />
        <Text marginBottom={6} fontWeight={1}>
          You will receive an email with a link to reset your password.
        </Text>
        <Form
          onSubmit={e => {
            setBackendError(null);
            handleSubmit(submit)(e);
          }}
        >
          <Field
            autoFocus
            autoComplete="on"
            type="email"
            label="Email address"
            name="email"
            errors={errors.email}
            ref={register}
          />
          <Button title="Reset Password" />
        </Form>

        <Row marginTop={4}>
          <Button href="/login" variant="text" title="Cancel" />
        </Row>
      </Card>
    </Column>
  );
};

export default PasswordReset;
