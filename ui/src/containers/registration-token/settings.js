import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import api from '../../api';
import utils from '../../utils';
import validators from '../../validators';
import useToggle from '../../hooks/useToggle';
import Card from '../../components/card';
import Field from '../../components/field';
import Popup from '../../components/popup';
import Alert from '../../components/alert';
import { Text, Button, Form, toaster } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
  maxRegistrations: yup
    .mixed()
    .notRequired()
    .nullable(),
});

const RegistrationTokenSettings = ({
  route: {
    data: { params, registrationToken },
  },
}) => {
  const navigation = useNavigation();
  const { register, handleSubmit, formState, errors } = useForm({
    validationSchema,
    defaultValues: {
      name: registrationToken.name,
      description: registrationToken.description,
      maxRegistrations: registrationToken.maxRegistrations,
    },
  });
  const [isDeletePopup, toggleDeletePopup] = useToggle();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.updateRegistrationToken({
        projectId: params.project,
        tokenId: registrationToken.id,
        data: {
          ...data,
          settings: registrationToken.settings,
        },
      });
      navigation.navigate(`/${params.project}/provisioning`);
      toaster.success('Registration token updated.');
    } catch (error) {
      setBackendError(
        utils.parseError(error, 'Registration token update failed.')
      );
      console.error(error);
    }
  };

  const submitDelete = async () => {
    setBackendError(null);
    try {
      await api.deleteRegistrationToken({
        projectId: params.project,
        tokenId: registrationToken.id,
      });
      toaster.success('Registration token deleted.');
      navigation.navigate(`/${params.project}/provisioning`);
    } catch (error) {
      setBackendError(
        utils.parseError(error, 'Registration token deletion failed.')
      );
      console.error(error);
    }
    toggleDeletePopup();
  };

  return (
    <Card
      title="Registration Token Settings"
      maxWidth="560px"
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
        <Field label="Name" name="name" ref={register} errors={errors.name} />
        <Field
          type="textarea"
          label="Description"
          name="description"
          ref={register}
          errors={errors.description}
        />
        <Field
          type="number"
          label="Maximum Device Registrations"
          name="maxRegistrations"
          description="Limits the number of devices that can be registered using this token."
          hint="Leave empty to allow unlimited registrations."
          errors={errors.maxRegistrations}
          ref={register}
        />
        <Button title="Update" disabled={!formState.dirty} />
      </Form>
      <Popup show={isDeletePopup} onClose={toggleDeletePopup}>
        <Card title="Delete Registration Token" border size="large">
          <Text>
            You are about to delete the{' '}
            <strong>{registrationToken.name}</strong> Registration Token.
          </Text>
          <Button
            title="Delete"
            marginTop={5}
            onClick={submitDelete}
            variant="danger"
          />
        </Card>
      </Popup>
    </Card>
  );
};

export default RegistrationTokenSettings;
