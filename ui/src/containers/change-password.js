import React from 'react';
import * as yup from 'yup';

import { endpoints } from '../api';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import { Form, toaster } from '../components/core';

const endpoint = endpoints.updateUser();
const validationSchema = yup.object().shape({
  currentPassword: yup.string().required(),
  password: validators.password.required(),
});

const ChangePassword = ({ close }) => {
  return (
    <Card title="Change Password" border size="large">
      <Form
        endpoint={endpoint}
        onSuccess={() => {
          toaster.success('Password updated.');
          close();
        }}
        validationSchema={validationSchema}
        errorMessages={{ default: 'Password update failed.' }}
        submitLabel="Change Password"
      >
        <Field
          required
          autoFocus
          type="password"
          label="Current Password"
          name="currentPassword"
        />
        <Field
          required
          type="password"
          label="New Password"
          name="password"
          hint="Password must be at least 8 characters."
        />
      </Form>
    </Card>
  );
};

export default ChangePassword;
