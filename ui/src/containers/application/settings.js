import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import api from '../../api';
import utils from '../../utils';
import validators from '../../validators';
import useToggle from '../../hooks/useToggle';
import Card from '../../components/card';
import Popup from '../../components/popup';
import Field from '../../components/field';
import Alert from '../../components/alert';
import { Button, Text, Form, toaster } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
});

const ApplicationSettings = ({
  route: {
    data: { params, application },
  },
}) => {
  const { register, handleSubmit, errors, formState } = useForm({
    validationSchema,
    defaultValues: {
      name: application.name,
      description: application.description,
    },
  });
  const [backendError, setBackendError] = useState();
  const [isDeletePopup, toggleDeletePopup] = useToggle();
  const navigation = useNavigation();

  const submit = async data => {
    try {
      await api.updateApplication({
        projectId: params.project,
        applicationId: application.id,
        data,
      });
      toaster.success('Application updated.');
      navigation.navigate(`/${params.project}/applications/${data.name}`);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Application update failed.'));
      console.error(error);
    }
  };

  const submitDelete = async () => {
    setBackendError(null);
    try {
      await api.deleteApplication({
        projectId: params.project,
        applicationId: application.name,
      });
      toaster.success('Application deleted.');
      navigation.navigate(`/${params.project}/applications`);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Application deletion failed.'));
      console.error(error);
    }
    toggleDeletePopup();
  };

  return (
    <>
      <Card
        width="100%"
        maxWidth="575px"
        title="Application Settings"
        actions={[
          {
            title: 'Delete',
            onClick: toggleDeletePopup,
            variant: 'danger',
          },
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
            label="Name"
            name="name"
            ref={register}
            errors={errors.name}
          />
          <Field
            type="textarea"
            label="Description"
            name="description"
            ref={register}
            errors={errors.description}
          />
          <Button
            marginTop={3}
            title="Update"
            type="submit"
            disabled={!formState.dirty}
          />
        </Form>
      </Card>
      <Popup show={isDeletePopup} onClose={toggleDeletePopup}>
        <Card title="Delete Application" border size="large">
          <Text>
            You are about to delete the <strong>{application.name}</strong>{' '}
            application.
          </Text>
          <Button
            marginTop={5}
            title="Delete"
            onClick={submitDelete}
            variant="danger"
          />
        </Card>
      </Popup>
    </>
  );
};

export default ApplicationSettings;
