import React, { useState, useEffect, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { useTable, useSortBy } from 'react-table';

import api from '../../api';
import utils from '../../utils';
import useToggle from '../../hooks/useToggle';
import { renderLabels } from '../../helpers/labels';
import {
  OperatorIs,
  OperatorIsNot,
  OperatorIsOptions,
  LabelValueCondition,
} from '../../components/devices-filter';
import Card from '../../components/card';
import Table from '../../components/table';
import DeviceStatus from '../../components/device-status';
import Popup from '../../components/popup';
import Field from '../../components/field';
import Alert from '../../components/alert';
import {
  Column,
  Row,
  Group,
  Label,
  Button,
  Text,
  Input,
  Form,
  Icon,
  toaster,
} from '../../components/core';

const InitialFilter = [
  {
    params: {
      key: '',
      value: '',
      operator: OperatorIs,
    },
  },
];

const InitialReleaseSelector = {
  releaseQuery: [InitialFilter],
  releaseId: '',
};

const LatestReleaseId = 'latest';

const ReleasePinning = ({
  route: {
    data: { params, application, releases },
  },
}) => {
  const [releaseSelectors, setReleaseSelectors] = useState(
    application.schedulingRule.releaseSelectors.length > 0
      ? application.schedulingRule.releaseSelectors
      : []
  );
  const { register, handleSubmit, getValues } = useForm({
    defaultValues: {
      defaultReleaseId:
        application.schedulingRule.defaultReleaseId || LatestReleaseId,
      releaseSelectors,
    },
  });
  const [isPreview, togglePreview] = useToggle();
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
        maxWidth: '72px',
      },
      {
        Header: 'Release',
        accessor: ({ releaseId }) => (isNaN(releaseId) ? 'Latest' : releaseId),
        maxWidth: '120px',
      },
      {
        Header: 'Name',
        accessor: 'name',
        minWidth: '200px',
      },
      {
        Header: 'Labels',
        accessor: 'labels',
        Cell: ({ cell: { value } }) => (value ? renderLabels(value) : null),
        minWidth: '300px',
      },
    ],
    []
  );
  const tableData = useMemo(() => scheduledDevices, [scheduledDevices]);

  const tableProps = useTable(
    {
      columns,
      data: tableData,
    },
    useSortBy
  );

  const releaseOptions = useMemo(
    () => [
      {
        value: 'latest',
        label: 'Latest',
      },
      ...releases.map(({ number }) => ({ value: number, label: number })),
    ],
    [releases]
  );

  const getSchedulingRuleFromFormData = data => ({
    ...application.schedulingRule,
    releaseSelectors: data.releaseSelectors
      ? data.releaseSelectors.map(({ releaseId, releaseQuery }) => ({
          releaseId: `${releaseId}`,
          releaseQuery: releaseQuery.map(filter =>
            filter.map(({ params }) => ({
              type: LabelValueCondition,
              params: {
                key: params.key,
                operator: params.operator,
                value: params.value,
              },
            }))
          ),
        }))
      : [],
    defaultReleaseId: `${data.defaultReleaseId}`,
  });

  const getScheduledDevices = async () => {
    const formData = getValues({ nest: true }); // use form submit for this instead?
    try {
      const { data } = await api.scheduledDevices({
        projectId: params.project,
        applicationId: application.name,
        schedulingRule: getSchedulingRuleFromFormData(formData),
        search: searchInput,
      });
      setScheduledDevices(data);
    } catch (error) {
      toaster.danger('Preview failed to load.');
      console.error(error);
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
          schedulingRule: getSchedulingRuleFromFormData(data),
        },
      });

      toaster.success('Releases pinned.');
    } catch (error) {
      setBackendError(utils.parseError(error, 'Release pinning failed.'));
      console.error(error);
    }
  };

  useEffect(() => {
    if (isPreview) {
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
              togglePreview();
            },
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
          <Group>
            <Row justifyContent="space-between" alignItems="center">
              <Label>Pinned Releases</Label>
              <Row>
                {releaseSelectors.length > 1 && (
                  <Button
                    title="Clear"
                    type="button"
                    variant="text"
                    onClick={() => setReleaseSelectors(InitialReleaseSelectors)}
                    marginRight={4}
                  />
                )}
                <Button
                  title="Add Release"
                  type="button"
                  variant="secondary"
                  onClick={() =>
                    setReleaseSelectors(rs => [...rs, InitialReleaseSelector])
                  }
                />
              </Row>
            </Row>
          </Group>
          {releaseSelectors.length === 0 && (
            <Row
              bg="grays.0"
              alignItems="center"
              justifyContent="center"
              padding={4}
              borderRadius={1}
              marginBottom={5}
            >
              <Text fontWeight={1}>No pinned releases</Text>
            </Row>
          )}
          {releaseSelectors.map(({ releaseQuery }, i) => (
            <Column
              marginBottom={6}
              paddingBottom={6}
              borderBottom={0}
              borderColor="grays.5"
            >
              {releaseQuery.map((filter, j) => (
                <Column flex={1}>
                  <Row marginBottom={4}>
                    <Text width="165px" paddingTop={1}>
                      {j === 0 ? 'If devices match' : 'and devices match'}
                    </Text>
                    <Column flex={1}>
                      {filter.map((_, k) => (
                        <Column alignItems="flex-start" flex={1}>
                          {k > 0 && (
                            <Text marginY={2} fontSize={0} fontWeight={2}>
                              OR
                            </Text>
                          )}
                          <Row alignItems="center" alignSelf="stretch">
                            <Field
                              inline
                              required
                              flex={1}
                              variant="small"
                              name={`releaseSelectors[${i}].releaseQuery[${j}][${k}].params.key`}
                              placeholder="Label Key"
                              ref={register}
                            />

                            <Field
                              inline
                              required
                              width="100px"
                              marginX={4}
                              type="select"
                              variant="small"
                              name={`releaseSelectors[${i}].releaseQuery[${j}][${k}].params.operator`}
                              placeholder="Operator"
                              options={OperatorIsOptions}
                              ref={register}
                            />

                            <Field
                              inline
                              required
                              flex={1}
                              variant="small"
                              name={`releaseSelectors[${i}].releaseQuery[${j}][${k}].params.value`}
                              placeholder="Label Value"
                              ref={register}
                            />

                            <Button
                              type="button"
                              marginLeft={2}
                              variant="iconDanger"
                              title={
                                <Icon icon="cross" size={16} color="red" />
                              }
                              onClick={() =>
                                setReleaseSelectors(releaseSelectors =>
                                  releaseSelectors
                                    .map((selector, selectorIndex) =>
                                      selectorIndex === i
                                        ? {
                                            ...selector,
                                            releaseQuery: selector.releaseQuery.map(
                                              (query, queryIndex) =>
                                                queryIndex === j
                                                  ? query.filter(
                                                      (_, filterIndex) =>
                                                        filterIndex !== k
                                                    )
                                                  : query
                                            ),
                                          }
                                        : selector
                                    )
                                    .map(selector => ({
                                      ...selector,
                                      releaseQuery: selector.releaseQuery.filter(
                                        arr => arr.length
                                      ),
                                    }))
                                    .filter(
                                      selector => selector.releaseQuery.length
                                    )
                                )
                              }
                            />
                          </Row>
                          {k === filter.length - 1 && (
                            <Button
                              marginTop={2}
                              type="button"
                              title={
                                <>
                                  <Text
                                    color="primary"
                                    marginRight={1}
                                    fontSize={1}
                                    fontWeight={3}
                                  >
                                    +
                                  </Text>
                                  <Text color="primary">OR</Text>
                                </>
                              }
                              variant="tertiary"
                              onClick={() =>
                                setReleaseSelectors(releaseSelectors =>
                                  releaseSelectors.map(
                                    (selector, selectorIndex) =>
                                      selectorIndex === i
                                        ? {
                                            ...selector,
                                            releaseQuery: selector.releaseQuery.map(
                                              (query, queryIndex) =>
                                                queryIndex === j
                                                  ? [...query, ...InitialFilter]
                                                  : query
                                            ),
                                          }
                                        : selector
                                  )
                                )
                              }
                            />
                          )}
                        </Column>
                      ))}
                    </Column>
                  </Row>

                  {j === releaseQuery.length - 1 && (
                    <Row marginBottom={2}>
                      <Button
                        type="button"
                        title={
                          <>
                            <Text
                              color="primary"
                              marginRight={1}
                              fontSize={1}
                              fontWeight={3}
                            >
                              +
                            </Text>
                            <Text color="primary">AND</Text>
                          </>
                        }
                        variant="tertiary"
                        onClick={() =>
                          setReleaseSelectors(releaseSelectors =>
                            releaseSelectors.map((selector, index) =>
                              index === i
                                ? {
                                    ...selector,
                                    releaseQuery: [
                                      ...selector.releaseQuery,
                                      InitialFilter,
                                    ],
                                  }
                                : selector
                            )
                          )
                        }
                      />
                    </Row>
                  )}
                </Column>
              ))}
              <Row alignItems="center" alignSelf="flex-start">
                <Text marginRight={4}>Pin devices to</Text>
                <Field
                  inline
                  required
                  width="120px"
                  type="select"
                  variant="small"
                  name={`releaseSelectors[${i}].releaseId`}
                  options={releaseOptions}
                  placeholder="Release"
                  ref={register}
                />
              </Row>
            </Column>
          ))}
          <Row alignItems="center" alignSelf="flex-start">
            <Text marginRight={4}>
              Pin {releaseSelectors.length > 0 ? 'remaining' : 'all'} devices to
            </Text>
            <Field
              inline
              required
              width="120px"
              type="select"
              variant="small"
              name="defaultReleaseId"
              options={releaseOptions}
              placeholder="Release"
              ref={register}
            />
          </Row>
          <Button
            type="submit"
            marginTop={6}
            title="Pin Releases"
            disabled={isSubmitDisabled}
          />
        </Form>
      </Card>

      <Popup show={isPreview} onClose={togglePreview}>
        <Card border size="xxlarge" title="Preview">
          <Row position="relative" alignItems="center" marginBottom={4}>
            <Icon
              icon="search"
              size={16}
              color={searchFocused ? 'primary' : 'white'}
              style={{ position: 'absolute', left: 16 }}
            />
            <Input
              placeholder="Search devices by name or labels"
              paddingLeft={7}
              value={searchInput}
              width="325px"
              onChange={e => setSearchInput(e.target.value)}
              onFocus={() => setSearchFocused(true)}
              onBlur={() => setSearchFocused(false)}
            />
          </Row>

          <Table
            {...tableProps}
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
