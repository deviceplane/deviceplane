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
  Row,
  Group,
  Form,
  Button,
  Text,
  Label,
  toaster,
} from '../../components/core';
import { labelColor } from '../../helpers/labels';

const MetricOverview = ({
  projectId,
  devices,
  metrics,
  metric,
  close,
  service,
  application,
}) => {
  const { control, handleSubmit, errors, watch } = useForm({
    defaultValues: {
      labels: metric.labels.map(label => ({ label, value: label })),
      properties: metric.properties.reduce(
        (obj, property) => ({ ...obj, [property]: true }),
        {}
      ),
      whitelistedTags: metric.whitelistedTags
        ? metric.whitelistedTags.map(({ key, values }) => ({
            label: key,
            value: key,
            values: values.map(value => ({ label: value, value })),
          }))
        : [],
    },
  });
  console.log(metric.properties);
  const whitelistedTags = watch('whitelistedTags');
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
    const update = {
      properties: Object.keys(data.properties).filter(
        property => data.properties[property]
      ),
      labels: data.labels ? data.labels.map(({ value }) => value) : [],
      whitelistedTags: data.whitelistedTags.map(({ value, values = [] }) => ({
        key: value,
        values: values.map(({ value }) => value),
      })),
    };
    const updatedMetrics = metrics.map(m =>
      m.name === metric.name
        ? {
            ...metric,
            ...update,
          }
        : m
    );
    try {
      if (service) {
        await api.updateServiceMetricsConfig({
          projectId,
          data: [
            ...metrics.filter(m =>
              m.applicationId === application.id ? m.service !== service : true
            ),
            {
              applicationId: application.id,
              service,
              exposedMetrics: updatedMetrics,
            },
          ],
        });
      } else {
        await api.updateDeviceMetricsConfig({
          projectId,
          data: updatedMetrics,
        });
      }
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
      subtitle={service ? null : metric.name}
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

        <Field
          required
          creatable
          type="multiselect"
          label="Whitelisted Tags"
          name="whitelistedTags"
          options={[]}
          placeholder="Whitelist tags"
          description="When empty, all tags for this metric are sent to Datadog. Whitelisting causes only the selected tags to be sent."
          noOptionsMessage={() => (
            <Text>
              Start typing to add a <strong>Tag</strong>.
            </Text>
          )}
          formatCreateLabel={value => (
            <Text>
              Add <strong>{value}</strong> Tag
            </Text>
          )}
          control={control}
          errors={errors.whitelistedTags}
        />
        {whitelistedTags && whitelistedTags.length > 0 && (
          <Group>
            <Row flex={1}>
              <Label width={11} marginRight={3}>
                Tag
              </Label>
              <Label>Whitelisted Values</Label>
            </Row>
            {whitelistedTags.map(({ value }, i) => {
              return (
                <Row alignItems="center" marginTop={i > 0 ? 4 : 0}>
                  <Row width={11} marginRight={3}>
                    <Text
                      padding="4px 6px"
                      fontSize={1}
                      fontWeight={1}
                      bg="white"
                      color="black"
                      borderRadius={1}
                    >
                      {value}
                    </Text>
                  </Row>

                  <Field
                    inline
                    flex={1}
                    required
                    creatable
                    type="multiselect"
                    name={`whitelistedTags[${i}].values`}
                    options={[]}
                    placeholder="Add values to whitelist"
                    noOptionsMessage={() => (
                      <Text>
                        Start typing to add a <strong>Value</strong>.
                      </Text>
                    )}
                    formatCreateLabel={value => (
                      <Text>
                        Add <strong>{value}</strong> Value
                      </Text>
                    )}
                    control={control}
                  />
                </Row>
              );
            })}
          </Group>
        )}

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

        <Button marginTop={3} title="Update" type="submit" />
      </Form>
    </Card>
  );
};

export default MetricOverview;
