import React, { useState, useMemo } from 'react';
import useForm from 'react-hook-form';
import { useNavigation } from 'react-navi';
import { toaster } from 'evergreen-ui';

import utils from '../../utils';
import api from '../../api';
import Field from '../../components/field';
import Card from '../../components/card';
import Alert from '../../components/alert';
import { Form, Button } from '../../components/core';

const ServiceMetricsSettings = ({
  projectId,
  applicationId,
  metricEndpointConfigs,
  service,
  close,
}) => {
  const endpointConfig = metricEndpointConfigs[service] || {};
  const { register, handleSubmit, errors, setValue } = useForm({
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
      toaster.success('Service metrics settings updated successfully.');
      close();
      navigation.refresh();
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Service metrics settings were not updated.');
      console.log(error);
    }
  };

  return (
    <Card
      title="Service Metrics Settings"
      size="xlarge"
      border
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
        <Button title="Update" type="submit" />
      </Form>
    </Card>
  );
};

export default ServiceMetricsSettings;
