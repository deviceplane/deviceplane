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
import { Label, Form, Button, Text, toaster } from '../../components/core';
import { labelColor } from '../../helpers/labels';

const ServiceMetricsForm = ({
  params,
  allMetrics,
  metrics,
  devices,
  application,
  service,
  close,
}) => {
  const { control, handleSubmit, errors } = useForm({});
  const [backendError, setBackendError] = useState();
  const navigation = useNavigation();

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
      await api.updateServiceMetricsConfig({
        projectId: params.project,
        data: [
          ...allMetrics.filter(m =>
            m.applicationId === application.id ? m.service !== service : true
          ),
          {
            applicationId: application.id,
            service,
            exposedMetrics: [
              ...data.metrics.map(({ value }) => ({
                name: value,
                properties: Object.keys(data.properties).filter(
                  property => data.properties[property]
                ),
                labels: data.labels
                  ? data.labels.map(({ value }) => value)
                  : [],
              })),
              ...metrics.filter(
                ({ name }) => !data.metrics.find(({ value }) => value === name)
              ),
            ],
          },
        ],
      });
      toaster.success('Metrics added.');
      close();
      navigation.refresh();
    } catch (error) {
      setBackendError(utils.parseError(error, 'Adding metric failed.'));
      console.error(error);
    }
  };

  return (
    <Card title="Add Service Metrics" size="xlarge" border overflow="visible">
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
          creatable
          type="multiselect"
          label="Metrics"
          name="metrics"
          options={[]}
          placeholder="Add metrics"
          noOptionsMessage={() => (
            <Text>
              Start typing to add a <strong>Metric</strong>.
            </Text>
          )}
          formatCreateLabel={value => (
            <Text>
              Add <strong>{value}</strong> Metric
            </Text>
          )}
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
        {config.supportedServiceMetricProperties.map(property => (
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

export default ServiceMetricsForm;
