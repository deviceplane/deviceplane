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
import DeviceStatus from '../../components/device-status';
import {
  Group,
  Text,
  Form,
  Button,
  Label,
  Value,
  toaster,
} from '../../components/core';

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
  const [isPopup, togglePopup] = useToggle();

  const submit = async data => {
    try {
      await api.updateDevice({
        projectId: params.project,
        deviceId: device.id,
        data,
      });
      toaster.success('Device updated.');
      navigation.navigate(`/${params.project}/devices/${data.name}`);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Device update failed.'));
      console.error(error);
    }
  };

  const submitDelete = async () => {
    setBackendError(null);
    try {
      await api.deleteDevice({
        projectId: params.project,
        deviceId: device.id,
      });
      toaster.success('Device removed.');
      navigation.navigate(`/${params.project}/devices`);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Device removal failed.'));
      console.error(error);
    }
    togglePopup();
  };

  return (
    <>
      <Card
        title="Device Settings"
        subtitle={<DeviceStatus status={device.status} />}
        actions={[
          {
            title: 'Remove',
            onClick: togglePopup,
            variant: 'danger',
          },
        ]}
        maxWidth="560px"
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
      <Popup show={isPopup} onClose={togglePopup}>
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
