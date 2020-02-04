import React, { useState, useEffect, useMemo } from 'react';
import { useForm } from 'react-hook-form';

import api from '../../api';
import utils from '../../utils';
import { renderLabels } from '../../helpers/labels';
import {
  OperatorIs,
  OperatorIsNot,
  LabelValueCondition,
} from '../../components/devices-filter';
import Card from '../../components/card';
import Table from '../../components/table';
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
  Form,
  Icon,
  toaster,
} from '../../components/core';

const InitialFilter = [
  {
    params: {
      key: '',
      value: '',
      operator: { label: OperatorIs, value: OperatorIs },
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
    data: { params, application, releases, devices },
  },
}) => {
  const [releaseSelectors, setReleaseSelectors] = useState(
    application.schedulingRule.releaseSelectors.length > 0
      ? application.schedulingRule.releaseSelectors
      : []
  );
  const { register, handleSubmit, setValue, getValues } = useForm({
    defaultValues: {
      defaultReleaseId:
        application.schedulingRule.defaultReleaseId || LatestReleaseId,
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
        Cell: ({ cell: { value } }) => (value ? renderLabels(value) : null),
        style: {
          flex: 2,
          overflow: 'hidden',
        },
        cellStyle: {
          marginBottom: '-8px',
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
      ...releases.map(({ number }) => ({ value: number, label: number })),
    ],
    [releases]
  );

  const getScheduledDevices = async () => {
    const formValues = getValues({ nest: true }); // use form submit for this instead?
    console.log(formValues);
    // try {
    //   const { data: devices } = await api.scheduledDevices({
    //     projectId: params.project,
    //     applicationId: application.name,
    //     schedulingRule: {
    //       ...application.schedulingRule,
    //       releaseSelectors: releaseSelectors.map((_, index) => ({
    //         releaseId: `${
    //           formValues[`releaseSelectors[${index}].releaseId`].value
    //         }`,
    //         releaseQuery: [
    //           [
    //             {
    //               type: LabelValueCondition,
    //               params: {
    //                 key:
    //                   formValues[
    //                     `releaseSelectors[${index}].releaseQuery[0].params.key`
    //                   ],
    //                 value:
    //                   formValues[
    //                     `releaseSelectors[${index}].releaseQuery[0].params.value`
    //                   ],
    //                 operator:
    //                   formValues[
    //                     `releaseSelectors[${index}].releaseQuery[0].params.operator`
    //                   ].value,
    //               },
    //             },
    //           ],
    //         ],
    //       })),
    //       defaultReleaseId: `${formValues.defaultReleaseId.value}`,
    //     },
    //     search: searchInput,
    //   });
    //   setScheduledDevices(devices);
    // } catch (error) {
    //   console.log(error);
    //   toaster.danger('Fetching device preview was unsuccessful.');
    // }
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
              <Label>Release Pins</Label>
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
                  title="Add Release Pin"
                  type="button"
                  variant="secondary"
                  onClick={() =>
                    setReleaseSelectors(rs => [InitialReleaseSelector, ...rs])
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
              <Text fontWeight={1}>No releases pinned</Text>
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
                  <Row marginBottom={5} flex={1}>
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
                          <Row alignItems="center" flex={1} alignSelf="stretch">
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
                              width="80px"
                              marginX={4}
                              type="select"
                              variant="small"
                              name={`releaseSelectors[${i}].releaseQuery[${j}][${k}].params.operator`}
                              placeholder="Operator"
                              options={[
                                { label: OperatorIs, value: OperatorIs },
                                {
                                  label: OperatorIsNot,
                                  value: OperatorIsNot,
                                },
                              ]}
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
                              marginLeft={4}
                              variant="icon"
                              title={
                                <Icon icon="cross" size={14} color="red" />
                              }
                            />
                          </Row>
                          {k === filter.length - 1 && (
                            <Button
                              marginTop={2}
                              type="button"
                              title="+ OR"
                              color="primary"
                              opacity={1}
                              variant="text"
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
                        title="+ AND"
                        color="primary"
                        opacity={1}
                        variant="text"
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
              width="125px"
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

      <Popup show={showPreview} onClose={() => setShowPreview(false)}>
        <Card border size="xxlarge" title="Preview">
          <Row position="relative" alignItems="center" marginBottom={4}>
            <Icon
              icon="search"
              size={16}
              color={searchFocused ? 'primary' : 'white'}
              style={{ position: 'absolute', left: 16 }}
            />
            <Input
              bg="black"
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
