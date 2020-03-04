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
import {
  Column,
  Form,
  Button,
  Text,
  Label,
  toaster,
} from '../../components/core';
import { labelColor } from '../../helpers/labels';

const ProjectMetricOverview = ({
  projectId,
  devices,
  metrics,
  metric,
  close,
}) => {
  const { control, handleSubmit, errors, watch } = useForm({
    defaultValues: {
      enabled: metric.enabled,
      labels: metric.labels.map(label => ({ label, value: label })),
      properties: metric.properties.reduce(
        (obj, property) => ({ ...obj, [property]: true }),
        {}
      ),
    },
  });
  const enabled = watch('enabled');
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
      await api.updateProjectMetricsConfig({
        projectId,
        data: metrics.map(m =>
          m.name === metric.name
            ? {
                ...metric,
                enabled: data.enabled,
                properties: Object.keys(data.properties).filter(
                  property => data.properties[property]
                ),
                labels: data.labels
                  ? data.labels.map(({ value }) => value)
                  : [],
              }
            : m
        ),
      });
      toaster.success('Metric updated.');
      close();
      navigation.refresh();
    } catch (error) {
      setBackendError(utils.parseError(error, 'Updating metric failed.'));
      console.error(error);
    }
  };

  return (
    <Card
      border
      title={getMetricLabel(metric.name)}
      subtitle={`deviceplane.${metric.name}`}
      size="xlarge"
      overflow="scroll"
    >
      <Alert show={backendError} variant="error" description={backendError} />
      <Form
        onSubmit={e => {
          setBackendError(null);
          handleSubmit(submit)(e);
        }}
      >
        <Field
          type="checkbox"
          label="Enabled"
          name="enabled"
          control={control}
        />

        <Column
          style={{ pointerEvents: enabled ? 'auto' : 'none' }}
          opacity={enabled ? 1 : 0.4}
        >
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
        </Column>

        <Button marginTop={3} title="Update" type="submit" />
      </Form>
    </Card>
  );
};

export default ProjectMetricOverview;
