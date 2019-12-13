import React, { useState } from 'react';
import useForm from 'react-hook-form';
import { useNavigation } from 'react-navi';
import { toaster } from 'evergreen-ui';
import * as yup from 'yup';

import api from '../../api';
import utils from '../../utils';
import validators from '../../validators';
import Card from '../../components/card';
import Field from '../../components/field';
import Popup from '../../components/popup';
import Alert from '../../components/alert';
import { Text, Button, Form } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
  maxRegistrations: yup
    .mixed()
    .notRequired()
    .nullable(),
});

const DeviceRegistrationTokenSettings = ({
  route: {
    data: { params, deviceRegistrationToken },
  },
}) => {
  const navigation = useNavigation();
  const { register, handleSubmit, formState, errors } = useForm({
    validationSchema,
    defaultValues: {
      name: deviceRegistrationToken.name,
      description: deviceRegistrationToken.description,
      maxRegistrations: deviceRegistrationToken.maxRegistrations,
    },
  });
  const [showDeletePopup, setShowDeletePopup] = useState();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.updateDeviceRegistrationToken({
        projectId: params.project,
        tokenId: deviceRegistrationToken.id,
        data: {
          ...data,
          settings: deviceRegistrationToken.settings,
        },
      });
      navigation.navigate(`/${params.project}/provisioning`);
      toaster.success('Device Registration Token updated successfully.');
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Device Registration Token was not updated.');
      console.log(error);
    }
  };

  const submitDelete = async () => {
    setBackendError(null);
    try {
      await api.deleteDeviceRegistrationToken({
        projectId: params.project,
        tokenId: deviceRegistrationToken.id,
      });
      toaster.success('Successfully deleted Device Registration Token.');
      navigation.navigate(`/${params.project}/provisioning`);
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Device Registration Token was not deleted.');
      console.log(error);
    }
    setShowDeletePopup(false);
  };

  return (
    <Card
      title="Device Registration Token Settings"
      actions={[
        {
          title: 'Delete',
          onClick: () => setShowDeletePopup(true),
          variant: 'secondary',
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
      <Popup show={showDeletePopup} onClose={() => setShowDeletePopup(false)}>
        <Card title="Delete Device Registration Token" border>
          <Text>
            You are about to delete the{' '}
            <strong>{deviceRegistrationToken.name}</strong> Device Registration
            Token.
          </Text>
          <Button title="Delete" marginTop={5} onClick={submitDelete} />
        </Card>
      </Popup>
    </Card>
  );
};

export default DeviceRegistrationTokenSettings;
