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
import { Button, Row, Form, toaster } from '../components/core';

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

const CreateConnection = ({
  route: {
    data: { params },
  },
}) => {
  const { register, handleSubmit, errors } = useForm({
    validationSchema,
    defaultValues: {
      name: '',
      protocol: ProtocolHTTP,
    },
  });
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.createConnection({ projectId: params.project, data });
      navigation.navigate(`/${params.project}/connections`);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Connection creation failed.'));
      console.error(error);
    }
  };

  return (
    <Layout alignItems="center">
      <Card title="Create Connection" size="large">
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

          <Button marginTop={3} type="submit" title="Create Connection" />
          <Row marginTop={4}>
            <Button
              title="Cancel"
              variant="text"
              href={`/${params.project}/connections`}
            />
          </Row>
        </Form>
      </Card>
    </Layout>
  );
};

export default CreateConnection;
