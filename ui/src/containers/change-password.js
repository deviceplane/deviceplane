import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import * as yup from 'yup';

import api from '../api';
import utils from '../utils';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import { Form, Button, toaster } from '../components/core';

const validationSchema = yup.object().shape({
  currentPassword: yup.string().required(),
  password: validators.password.required(),
});

const ChangePassword = ({ close }) => {
  const { register, handleSubmit, errors } = useForm({ validationSchema });
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.updateUser(data);
      toaster.success('Password updated.');
      close();
    } catch (error) {
      setBackendError(utils.parseError(error, 'Password update failed.'));
      console.error(error);
    }
  };

  return (
    <Card title="Change Password" border size="large">
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
          type="password"
          label="Current Password"
          name="currentPassword"
          ref={register}
          errors={errors.currentPassword}
        />
        <Field
          required
          type="password"
          label="New Password"
          name="password"
          ref={register}
          errors={errors.password}
          hint="Password must be at least 8 characters."
        />
        <Button marginTop={3} title="Change Password" type="submit" />
      </Form>
    </Card>
  );
};

export default ChangePassword;
