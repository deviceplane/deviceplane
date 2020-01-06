import React, { useState, useMemo } from 'react';
import useForm from 'react-hook-form';
import { useNavigation } from 'react-navi';
import { toaster } from 'evergreen-ui';

import utils from '../../utils';
import api from '../../api';
import config from '../../config';
import Field from '../../components/field';
import Card from '../../components/card';
import Alert from '../../components/alert';
import { DeviceLabelMulti } from '../../components/device-label';
import {
  Label,
  Form,
  Button,
  Select,
  Text,
  Checkbox,
} from '../../components/core';

const ServiceMetricsForm = ({
  params,
  allMetrics,
  metrics,
  devices,
  application,
  service,
  close,
  labelColorMap,
}) => {
  const { register, handleSubmit, errors, setValue } = useForm({});
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
            color: labelColorMap[label],
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
          label="Metrics"
          name="metrics"
          as={
            <Select
              multi
              creatable
              options={[]}
              placeholder={'Add metrics'}
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
            />
          }
          setValue={setValue}
          register={register}
          errors={errors.metrics}
        />
        <Field
          label="Labels"
          name="labels"
          setValue={setValue}
          register={register}
          as={
            <Select
              multi
              options={labelsOptions}
              multiComponent={DeviceLabelMulti}
              placeholder="Select labels"
              noOptionsMessage={() => (
                <Text>
                  There are no <strong>Labels</strong>.
                </Text>
              )}
            />
          }
          errors={errors.description}
        />

        <Label>Properties</Label>
        {config.supportedServiceMetricProperties.map(property => (
          <Field
            multi
            key={property.id}
            name={`properties[${property.id}]`}
            as={<Checkbox label={property.label} />}
            register={register}
            setValue={setValue}
            hint={property.description}
          />
        ))}

        <Button title="Add" type="submit" />
      </Form>
    </Card>
  );
};

export default ServiceMetricsForm;
