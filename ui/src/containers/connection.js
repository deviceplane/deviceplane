import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';
import { useForm } from 'react-hook-form';

import api from '../api';
import utils from '../utils';
import validators from '../validators';
import Layout from '../components/layout';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import Popup from '../components/popup';
import { Text, Button, Row, Form, toaster } from '../components/core';

const ProtocolTCP = 'tcp';
const ProtocolHTTP = 'http';

const ProtocolOptions = [
  {
    value: ProtocolHTTP,
    label: 'HTTP',
  },
  {
    value: ProtocolTCP,
    label: 'TCP',
  },
];

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  port: yup
    .number()
    .integer()
    .min(1)
    .max(65535),
  protocol: yup.mixed().oneOf([ProtocolHTTP, ProtocolTCP]),
});

const Connection = ({
  route: {
    data: { params, connection },
  },
}) => {
  const { register, handleSubmit, errors, formState } = useForm({
    validationSchema,
    defaultValues: connection,
  });
  const navigation = useNavigation();
  const [showDeletePopup, setShowDeletePopup] = useState();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.updateConnection({
        projectId: params.project,
        connectionId: connection.name,
        data,
      });
      navigation.navigate(`/${params.project}/connections`);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Connection update failed.'));
      console.error(error);
    }
  };

  const submitDelete = async () => {
    try {
      await api.deleteConnection({
        projectId: params.project,
        connectionId: connection.name,
      });
      navigation.navigate(`/${params.project}/connections`);
    } catch (error) {
      console.error(error);
      toaster.danger('Connection deletion failed.');
    }
  };

  return (
    <Layout alignItems="center">
      <Card
        title={connection.name}
        size="large"
        actions={[
          {
            title: 'Delete',
            variant: 'danger',
            onClick: () => setShowDeletePopup(true),
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
            autoFocus
            label="Name"
            name="name"
            ref={register}
            errors={errors.name}
          />
          <Field
            required
            label="Port"
            name="port"
            ref={register}
            errors={errors.port}
          />
          <Field
            required
            label="Protocol"
            name="protocol"
            type="select"
            options={ProtocolOptions}
            ref={register}
            errors={errors.protocol}
          />

          <Button
            marginTop={3}
            type="submit"
            title="Update Connection"
            disabled={!formState.dirty}
          />
          <Row marginTop={4}>
            <Button
              title="Cancel"
              variant="text"
              href={`/${params.project}/connections`}
            />
          </Row>
        </Form>
      </Card>
      <Popup show={showDeletePopup} onClose={() => setShowDeletePopup(false)}>
        <Card title="Delete Connection" border size="large">
          <Text>
            You are about to delete the <strong>{connection.name}</strong>{' '}
            connection.
          </Text>
          <Button
            marginTop={5}
            title="Delete"
            onClick={submitDelete}
            variant="danger"
          />
        </Card>
      </Popup>
    </Layout>
  );
};

export default Connection;
