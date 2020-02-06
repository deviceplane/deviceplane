import React, { useState, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigation } from 'react-navi';

import utils from '../../utils';
import api from '../../api';
import Field from '../../components/field';
import Card from '../../components/card';
import Alert from '../../components/alert';
import { Form, Button, toaster } from '../../components/core';

const ServiceMetricsSettings = ({
  projectId,
  applicationId,
  metricEndpointConfigs,
  service,
  close,
}) => {
  const endpointConfig = metricEndpointConfigs[service] || {};
  const { register, handleSubmit, errors } = useForm({
    defaultValues: {
      port: endpointConfig.port || 2112,
      path: endpointConfig.path || '/metrics',
    },
  });
  const [backendError, setBackendError] = useState();
  const navigation = useNavigation();

  const submit = async data => {
    try {
      await api.updateApplication({
        projectId,
        applicationId,
        data: {
          metricEndpointConfigs: {
            ...metricEndpointConfigs,
            [service]: {
              port: +data.port,
              path: data.path,
            },
          },
        },
      });
      toaster.success('Service metrics settings updated.');
      close();
      navigation.refresh();
    } catch (error) {
      setBackendError(
        utils.parseError(error, 'Service metrics settings update failed.')
      );
      console.error(error);
    }
  };

  return (
    <Card
      border
      title="Service Metrics Settings"
      size="xlarge"
      overflow="visible"
    >
      <Alert show={backendError} variant="error" description={backendError} />
      <Form
        onSubmit={e => {
          setBackendError(null);
          handleSubmit(submit)(e);
        }}
      >
        <Field
          label="Port"
          name="port"
          description="The port that your service exposes metrics from."
          ref={register}
          errors={errors.port}
        />
        <Field
          label="Path"
          name="path"
          description="The path that your service exposes metrics from."
          ref={register}
          errors={errors.path}
        />
        <Button marginTop={3} title="Update" type="submit" />
      </Form>
    </Card>
  );
};

export default ServiceMetricsSettings;
