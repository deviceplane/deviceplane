import React, { useState, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigation } from 'react-navi';

import utils from '../../utils';
import api from '../../api';
import config from '../../config';
import Field from '../../components/field';
import Card from '../../components/card';
import Alert from '../../components/alert';
import { DeviceLabelMulti } from '../../components/device-label';
import { getMetricLabel } from '../../helpers/metrics';
import { Form, Button, Text, Label, toaster } from '../../components/core';
import { labelColor } from '../../helpers/labels';

const metricsOptions = config.supportedDeviceMetrics.map(value => ({
  label: getMetricLabel(value),
  value,
}));

const DeviceMetricsForm = ({ params, devices, metrics, close }) => {
  const { control, handleSubmit, errors } = useForm({});
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();

  const labelsOptions = useMemo(
    () =>
      [
        ...new Set(
          devices.reduce(
            (options, device) => [...options, ...Object.keys(device.labels)],
            []
          )
        ),
      ].map(
        label => ({
          label,
          value: label,
          props: {
            color: labelColor(label),
          },
        }),
        []
      ),
    [devices]
  );

  const submit = async data => {
    try {
      await api.updateDeviceMetricsConfig({
        projectId: params.project,
        data: [
          ...data.metrics.map(({ value }) => ({
            name: value,
            properties: Object.keys(data.properties).filter(
              property => data.properties[property]
            ),
            labels: data.labels ? data.labels.map(({ value }) => value) : [],
          })),
          ...metrics.filter(
            ({ name }) => !data.metrics.find(({ value }) => value === name)
          ),
        ],
      });
      toaster.success('Metrics added successfully.');
      close();
      navigation.refresh();
    } catch (error) {
      console.log(error);

      const errorMessage = utils.parseError(error);

      toaster.danger('Metrics were not added.');

      if (errorMessage) {
        setBackendError(errorMessage);
      } else {
        close();
      }
    }
  };

  return (
    <Card border title="Add Device Metrics" size="xlarge">
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
          type="multiselect"
          label="Metrics"
          name="metrics"
          options={metricsOptions}
          placeholder="Select metrics"
          control={control}
          errors={errors.metrics}
        />
        <Field
          type="multiselect"
          label="Labels"
          name="labels"
          control={control}
          options={labelsOptions}
          multiComponent={DeviceLabelMulti}
          placeholder="Select labels"
          noOptionsMessage={() => (
            <Text>
              There are no <strong>Labels</strong>.
            </Text>
          )}
          errors={errors.description}
        />

        <Label>Properties</Label>
        {config.supportedDeviceMetricProperties.map(property => (
          <Field
            multi
            type="checkbox"
            key={property.id}
            label={property.label}
            name={`properties[${property.id}]`}
            control={control}
            hint={property.description}
          />
        ))}

        <Button marginTop={3} title="Add" type="submit" />
      </Form>
    </Card>
  );
};

export default DeviceMetricsForm;
