import React, { useState, useEffect, useMemo } from 'react';
import { toaster, Icon } from 'evergreen-ui';
import useForm from 'react-hook-form';

import api from '../../api';
import utils from '../../utils';
import theme, { labelColors } from '../../theme';
import { buildLabelColorMap, renderLabels } from '../../helpers/labels';
import { DevicesFilterButtons } from '../../components/devices-filter-buttons';
import {
  OperatorIs,
  OperatorIsNot,
  DevicesFilter,
  LabelValueCondition,
} from '../../components/devices-filter';
import Card from '../../components/card';
import Alert from '../../components/alert';
import Table from '../../components/table';
import RadioGroup from '../../components/radio-group';
import DeviceStatus from '../../components/device-status';
import Popup from '../../components/popup';
import Field from '../../components/field';
import {
  Column,
  Row,
  Group,
  Label,
  Button,
  Text,
  Input,
  Select,
  Link,
  Form,
} from '../../components/core';

const initialReleaseSelectors = [
  {
    releaseQuery: [
      [
        {
          params: {
            key: '',
            value: '',
            operator: { label: OperatorIs, value: OperatorIs },
          },
        },
      ],
    ],
    releaseId: '',
  },
];

const ReleasePinning = ({
  route: {
    data: { params, application, releases, devices },
  },
}) => {
  const [labelColorMap, setLabelColorMap] = useState(
    buildLabelColorMap({}, labelColors, devices)
  );
  const [releaseSelectors, setReleaseSelectors] = useState(
    application.schedulingRule.releaseSelectors.length > 0
      ? application.schedulingRule.releaseSelectors
      : initialReleaseSelectors
  );
  const {
    register,
    handleSubmit,
    formState,
    errors,
    setValue,
    getValues,
  } = useForm({
    defaultValues: {
      defaultReleaseId: application.schedulingRule.defaultReleaseId
        ? {
            value: application.schedulingRule.defaultReleaseId,
            label: isNaN(application.schedulingRule.defaultReleaseId)
              ? 'Latest'
              : application.schedulingRule.defaultReleaseId,
          }
        : {
            value: 'latest',
            label: 'Latest',
          },
      releaseSelectors,
    },
  });
  const [showPreview, setShowPreview] = useState();
  const [scheduledDevices, setScheduledDevices] = useState([]);
  const [searchInput, setSearchInput] = useState('');
  const [searchFocused, setSearchFocused] = useState();
  const [backendError, setBackendError] = useState();

  const columns = useMemo(
    () => [
      {
        Header: 'Status',
        accessor: 'status',
        Cell: ({ cell: { value } }) => <DeviceStatus status={value} />,
        style: {
          flex: '0 0 72px',
        },
      },
      {
        Header: 'Release',
        accessor: ({ releaseId }) => (isNaN(releaseId) ? 'Latest' : releaseId),
        style: {
          flex: '0 0 60px',
        },
      },
      {
        Header: 'Name',
        accessor: 'name',
      },
      {
        Header: 'Labels',
        accessor: 'labels',
        Cell: ({ cell: { value } }) =>
          value ? renderLabels(value, labelColorMap) : null,
        style: {
          flex: 2,
          overflow: 'hidden',
        },
      },
    ],
    []
  );
  const tableData = useMemo(() => scheduledDevices, [scheduledDevices]);

  const releaseOptions = useMemo(
    () => [
      {
        value: 'latest',
        label: 'Latest',
      },
      ...releases
        .filter(({ number }) => number !== application.latestRelease.number)
        .map(({ number }) => ({ value: number, label: number })),
    ],
    [releases]
  );

  const getScheduledDevices = async () => {
    const formValues = getValues(); // use form submit for this instead?
    try {
      const { data: devices } = await api.scheduledDevices({
        projectId: params.project,
        applicationId: application.name,
        schedulingRule: {
          ...application.schedulingRule,
          releaseSelectors: releaseSelectors.map((_, index) => ({
            releaseId: `${
              formValues[`releaseSelectors[${index}].releaseId`].value
            }`,
            releaseQuery: [
              [
                {
                  type: LabelValueCondition,
                  params: {
                    key:
                      formValues[
                        `releaseSelectors[${index}].releaseQuery[0].params.key`
                      ],
                    value:
                      formValues[
                        `releaseSelectors[${index}].releaseQuery[0].params.value`
                      ],
                    operator:
                      formValues[
                        `releaseSelectors[${index}].releaseQuery[0].params.operator`
                      ].value,
                  },
                },
              ],
            ],
          })),
          defaultReleaseId: `${formValues.defaultReleaseId.value}`,
        },
        search: searchInput,
      });
      setScheduledDevices(devices);
    } catch (error) {
      console.log(error);
      toaster.danger('Fetching device preview was unsuccessful.');
    }
  };

  const submit = async data => {
    try {
      await api.updateApplication({
        projectId: params.project,
        applicationId: application.name,
        data: {
          name: application.name,
          description: application.description,
          schedulingRule: {
            ...application.schedulingRule,
            releaseSelectors: data.releaseSelectors.map(
              ({ releaseId, releaseQuery }) => ({
                releaseId: `${releaseId.value}`,
                releaseQuery: [
                  releaseQuery.map(({ params }) => ({
                    type: LabelValueCondition,
                    params: {
                      key: params.key,
                      operator: params.operator.value,
                      value: params.value,
                    },
                  })),
                ],
              })
            ),
            defaultReleaseId: `${data.defaultReleaseId.value}`,
          },
        },
      });

      toaster.success('Release pinning successful.');
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Release pinning was not successful.');
      console.log(error);
    }
  };

  useEffect(() => {
    if (showPreview) {
      getScheduledDevices();
    }
  }, [searchInput]);

  let isSubmitDisabled = false;

  return (
    <>
      <Card
        size="xlarge"
        title="Release Pinning"
        subtitle="Schedule releases to run on devices based on conditions"
        actions={[
          {
            title: 'Preview',
            variant: 'secondary',
            disabled: isSubmitDisabled,
            onClick: () => {
              getScheduledDevices();
              setShowPreview(true);
            },
          },
        ]}
      >
        <Form
          onSubmit={e => {
            setBackendError(null);
            handleSubmit(submit)(e);
          }}
        >
          <Group>
            <Row justifyContent="space-between" alignItems="center">
              <Label>Conditions</Label>
              <Row>
                {releaseSelectors.length > 1 && (
                  <Button
                    title="Clear"
                    type="button"
                    variant="text"
                    onClick={() => setReleaseSelectors(initialReleaseSelectors)}
                    marginRight={4}
                  />
                )}
                <Button
                  title="Add Condition"
                  type="button"
                  variant="secondary"
                  onClick={() =>
                    setReleaseSelectors(rs => [
                      { query: [], releaseId: '' },
                      ...rs,
                    ])
                  }
                />
              </Row>
            </Row>
          </Group>
          {releaseSelectors.map((_, index) => (
            <Column
              marginBottom={6}
              paddingBottom={6}
              borderBottom={0}
              borderColor="white"
            >
              <Row alignItems="center" marginBottom={4}>
                <Text>If devices match</Text>
                <Row marginLeft={4} alignItems="center" flex={1}>
                  <Row flex={1}>
                    <Field
                      inline
                      name={`releaseSelectors[${index}].releaseQuery[0].params.key`}
                      placeholder="Label Key"
                      ref={register}
                    />
                  </Row>

                  <Row width="125px" marginX={4}>
                    <Field
                      inline
                      name={`releaseSelectors[${index}].releaseQuery[0].params.operator`}
                      as={
                        <Select
                          placeholder="Operator"
                          options={[
                            { label: OperatorIs, value: OperatorIs },
                            { label: OperatorIsNot, value: OperatorIsNot },
                          ]}
                        />
                      }
                      setValue={setValue}
                      register={register}
                    />
                  </Row>

                  <Row flex={1}>
                    <Field
                      inline
                      name={`releaseSelectors[${index}].releaseQuery[0].params.value`}
                      placeholder="Label Value"
                      ref={register}
                    />
                  </Row>
                </Row>
              </Row>

              <Row alignItems="center" alignSelf="flex-start">
                <Text marginRight="34px">Pin devices to</Text>
                <Row width="120px">
                  <Field
                    inline
                    name={`releaseSelectors[${index}].releaseId`}
                    as={
                      <Select options={releaseOptions} placeholder="Release" />
                    }
                    setValue={setValue}
                    register={register}
                  />
                </Row>
              </Row>
            </Column>
          ))}
          {releaseSelectors.length > 0 && (
            <Row alignItems="center" alignSelf="flex-start">
              <Text marginRight={4}>Pin remaining devices to</Text>
              <Row width="125px">
                <Field
                  inline
                  name="defaultReleaseId"
                  as={<Select options={releaseOptions} placeholder="Release" />}
                  setValue={setValue}
                  register={register}
                />
              </Row>
            </Row>
          )}
          <Button
            type="submit"
            marginTop={6}
            title="Pin Releases"
            disabled={isSubmitDisabled}
          />
        </Form>
      </Card>

      <Popup show={showPreview} onClose={() => setShowPreview(false)}>
        <Card
          border
          size="xxlarge"
          title="Preview"
          subtitle="This application will run on the devices listed below."
        >
          <Row position="relative" alignItems="center" marginBottom={4}>
            <Icon
              icon="search"
              size={16}
              color={searchFocused ? theme.colors.primary : theme.colors.white}
              style={{ position: 'absolute', left: 16 }}
            />
            <Input
              bg="black"
              placeholder="Search devices by name or labels"
              paddingLeft={8}
              value={searchInput}
              width="350px"
              onChange={e => setSearchInput(e.target.value)}
              onFocus={() => setSearchFocused(true)}
              onBlur={() => setSearchFocused(false)}
            />
          </Row>

          <Table
            columns={columns}
            data={tableData}
            placeholder={
              <Text>
                There are no <strong>Devices</strong>.
              </Text>
            }
          />
        </Card>
      </Popup>
    </>
  );
};

export default ReleasePinning;
