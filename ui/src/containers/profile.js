import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import * as yup from 'yup';
import { useCurrentRoute } from 'react-navi';
import api from '../api';
import utils from '../utils';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import { Form, Button, toaster } from '../components/core';

const validationSchema = yup.object().shape({
  name: yup
    .string()
    .required()
    .max(100),
});

const Profile = ({ close }) => {
  const {
    data: {
      context: { currentUser, setCurrentUser },
    },
  } = useCurrentRoute();
  const { register, handleSubmit, formState, errors } = useForm({
    validationSchema,
    defaultValues: {
      name: currentUser.name,
    },
  });
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.updateUser(data);
      setCurrentUser({ ...currentUser, ...data });
      toaster.success('Profile updated.');
      close();
    } catch (error) {
      setBackendError(utils.parseError(error, 'Profile update failed.'));
      console.error(error);
    }
  };

  return (
    <Card title="Profile" border size="large">
      <Alert show={backendError} variant="error" description={backendError} />
      <Form
        onSubmit={e => {
          setBackendError(null);
          handleSubmit(submit)(e);
        }}
      >
        <Field
          required
          autoCapitalize="words"
          label="Name"
          name="name"
          ref={register}
          errors={errors.name}
        />
        <Button
          marginTop={3}
          title="Update"
          type="submit"
          disabled={!formState.dirty}
        />
      </Form>
    </Card>
  );
};

export default Profile;
