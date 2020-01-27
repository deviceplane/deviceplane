import React, { useState } from 'react';
import useForm from 'react-hook-form';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';
import { toaster } from 'evergreen-ui';

import api from '../../api';
import utils from '../../utils';
import validators from '../../validators';
import Card from '../../components/card';
import Field from '../../components/field';
import Popup from '../../components/popup';
import Alert from '../../components/alert';
import DeviceStatus from '../../components/device-status';
import { Group, Text, Form, Button, Label, Value } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
});

const DeviceSettings = ({
  route: {
    data: { params, device },
  },
}) => {
  const { register, handleSubmit, formState, errors } = useForm({
    validationSchema,
    defaultValues: {
      name: device.name,
    },
  });
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();
  const [showPopup, setShowPopup] = useState();

  const submit = async data => {
    try {
      await api.updateDevice({
        projectId: params.project,
        deviceId: device.id,
        data,
      });
      toaster.success('Device updated successfully.');
      navigation.navigate(`/${params.project}/devices/${data.name}`);
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Device was not updated.');
      console.log(error);
    }
  };

  const submitDelete = async () => {
    setBackendError(null);
    try {
      await api.deleteDevice({
        projectId: params.project,
        deviceId: device.id,
      });
      toaster.success('Successfully deleted device.');
      navigation.navigate(`/${params.project}/devices`);
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Device was not removed.');
      console.log(error);
    }
    setShowPopup(false);
  };

  return (
    <>
      <Card
        title="Device Settings"
        subtitle={<DeviceStatus status={device.status} />}
        size="large"
        actions={[
          {
            title: 'Remove',
            onClick: () => setShowPopup(true),
            variant: 'danger',
          },
        ]}
      >
        <Alert show={backendError} variant="error" description={backendError} />

        <Group>
          <Label>ID</Label>
          <Value>{device.id}</Value>
        </Group>

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
          <Button
            marginTop={3}
            type="submit"
            title="Update"
            disabled={!formState.dirty}
          />
        </Form>
      </Card>
      <Popup show={showPopup} onClose={() => setShowPopup(false)}>
        <Card title="Remove Device" border size="large">
          <Text>
            You are about to remove the <strong>{device.name}</strong> device.
          </Text>

          <Button
            marginTop={5}
            title="Remove"
            onClick={submitDelete}
            variant="danger"
          />
        </Card>
      </Popup>
    </>
  );
};

export default DeviceSettings;
