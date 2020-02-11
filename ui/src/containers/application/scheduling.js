import React, { useState, useEffect, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigation } from 'react-navi';

import api from '../../api';
import utils from '../../utils';
import { renderLabels } from '../../helpers/labels';
import {
  OperatorIs,
  OperatorIsOptions,
  LabelValueCondition,
} from '../../components/devices-filter';
import Card from '../../components/card';
import Alert from '../../components/alert';
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
  Icon,
  toaster,
  Form,
} from '../../components/core';

const ScheduleTypeConditional = 'Conditional';
const ScheduleTypeAll = 'AllDevices';
const ScheduleTypeNone = 'NoDevices';
const InitialCondition = {
  type: LabelValueCondition,
  params: {
    key: '',
    operator: OperatorIs,
    value: '',
  },
};
const InitialFilter = [InitialCondition];

const Scheduling = ({
  route: {
    data: { params, application },
  },
}) => {
  const [conditionalQuery, setConditionalQuery] = useState(
    application.schedulingRule.conditionalQuery &&
      application.schedulingRule.conditionalQuery.length
      ? application.schedulingRule.conditionalQuery
      : [InitialFilter]
  );
  const [scheduledDevices, setScheduledDevices] = useState([]);
  const [backendError, setBackendError] = useState();
  const [searchInput, setSearchInput] = useState('');
  const [searchFocused, setSearchFocused] = useState();
  const [showPreview, setShowPreview] = useState();
  const { register, watch, control, handleSubmit, getValues } = useForm({
    defaultValues: {
      scheduleType: application.schedulingRule.scheduleType || 'NoDevices',
      conditionalQuery,
    },
  });
  const watchScheduleType = watch('scheduleType');

  const columns = useMemo(
    () => [
      {
        Header: 'Status',
        accessor: 'status',
        Cell: ({ cell: { value } }) => <DeviceStatus status={value} />,
      },
      {
        Header: 'Name',
        accessor: 'name',
      },
      {
        Header: 'Labels',
        accessor: 'labels',
        Cell: ({ cell: { value } }) => (value ? renderLabels(value) : null),
      },
    ],
    []
  );
  const tableData = useMemo(() => scheduledDevices, [scheduledDevices]);

  const navigation = useNavigation();

  const addTypeToConditionalQuery = conditionalQuery =>
    conditionalQuery.map(filters =>
      filters.map(condition => ({ type: LabelValueCondition, ...condition }))
    );

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
            scheduleType: data.scheduleType,
            conditionalQuery: data.conditionalQuery
              ? addTypeToConditionalQuery(data.conditionalQuery)
              : application.schedulingRule.conditionalQuery,
          },
        },
      });

      toaster.success('Scheduling applied.');
      navigation.navigate(
        `/${params.project}/applications/${application.name}`
      );
    } catch (error) {
      setBackendError(utils.parseError(error, 'Scheduling failed.'));
      console.error(error);
    }
  };

  const getScheduledDevices = async () => {
    const { scheduleType, conditionalQuery } = getValues({ nest: true });
    try {
      const { data: devices } = await api.scheduledDevices({
        projectId: params.project,
        applicationId: application.name,
        schedulingRule: {
          ...application.schedulingRule,
          scheduleType,
          conditionalQuery: addTypeToConditionalQuery(conditionalQuery),
        },
        search: searchInput,
      });
      setScheduledDevices(devices);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    getScheduledDevices();
  }, [searchInput]);

  const isSubmitDisabled =
    watchScheduleType === ScheduleTypeConditional
      ? false
      : watchScheduleType === application.schedulingRule.scheduleType;

  return (
    <>
      <Card
        title="Scheduling"
        size="xlarge"
        actions={
          watchScheduleType === ScheduleTypeConditional
            ? [
                {
                  title: 'Preview',
                  variant: 'secondary',
                  onClick: () => {
                    getScheduledDevices();
                    setShowPreview(true);
                  },
                },
              ]
            : []
        }
      >
        <Row alignSelf="stretch">
          <Column flex={1}>
            <Alert
              show={backendError}
              variant="error"
              description={backendError}
            />
            <Form
              onSubmit={e => {
                handleSubmit(submit)(e);
              }}
            >
              <Group>
                <Label>This application will run on</Label>
                <Field
                  type="radiogroup"
                  name="scheduleType"
                  control={control}
                  options={[
                    { label: 'No Devices', value: ScheduleTypeNone },
                    { label: 'All Devices', value: ScheduleTypeAll },
                    {
                      label: 'Some Devices',
                      value: ScheduleTypeConditional,
                    },
                  ]}
                />
              </Group>
              {watchScheduleType === ScheduleTypeConditional &&
                conditionalQuery.map((filter, i) => (
                  <Column flex={1}>
                    <Row marginBottom={4}>
                      <Text width="165px" paddingTop={1}>
                        {i === 0 ? 'If devices match' : 'and devices match'}
                      </Text>
                      <Column flex={1}>
                        {filter.map((_, j) => (
                          <Column alignItems="flex-start" flex={1}>
                            {j > 0 && (
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
                                name={`conditionalQuery[${i}][${j}].params.key`}
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
                                name={`conditionalQuery[${i}][${j}].params.operator`}
                                placeholder="Operator"
                                options={OperatorIsOptions}
                                ref={register}
                              />

                              <Field
                                inline
                                required
                                flex={1}
                                variant="small"
                                name={`conditionalQuery[${i}][${j}].params.value`}
                                placeholder="Label Value"
                                ref={register}
                              />

                              {(i > 0 || j > 0) && (
                                <Button
                                  type="button"
                                  marginLeft={2}
                                  variant="icon"
                                  title={
                                    <Icon icon="cross" size={14} color="red" />
                                  }
                                  onClick={() =>
                                    setConditionalQuery(query =>
                                      query
                                        .map((filters, filterIndex) =>
                                          filterIndex === i
                                            ? filters.filter(
                                                (_, conditionIndex) =>
                                                  conditionIndex !== j
                                              )
                                            : filters
                                        )
                                        .filter(
                                          (filters, filterIndex) =>
                                            filterIndex === 0 || filters.length
                                        )
                                    )
                                  }
                                />
                              )}
                            </Row>
                            {j === filter.length - 1 && (
                              <Button
                                marginTop={2}
                                type="button"
                                title="+ OR"
                                color="primary"
                                opacity={1}
                                variant="text"
                                onClick={() =>
                                  setConditionalQuery(query =>
                                    query.map((filter, filterIndex) =>
                                      filterIndex === i
                                        ? [...filter, InitialCondition]
                                        : filter
                                    )
                                  )
                                }
                              />
                            )}
                          </Column>
                        ))}
                      </Column>
                    </Row>

                    {i === conditionalQuery.length - 1 && (
                      <Row marginBottom={2}>
                        <Button
                          type="button"
                          title="+ AND"
                          color="primary"
                          opacity={1}
                          variant="text"
                          onClick={() =>
                            setConditionalQuery(filters => [
                              ...filters,
                              InitialFilter,
                            ])
                          }
                        />
                      </Row>
                    )}
                  </Column>
                ))}
              <Button
                marginTop={6}
                title="Apply Scheduling"
                disabled={isSubmitDisabled}
              />
            </Form>
          </Column>
        </Row>
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
              onFocus={e => setSearchFocused(true)}
              onBlur={e => setSearchFocused(false)}
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

export default Scheduling;
