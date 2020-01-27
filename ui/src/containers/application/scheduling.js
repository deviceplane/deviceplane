import React, { useState, useEffect, useMemo } from 'react';
import { toaster, Icon } from 'evergreen-ui';
import { useNavigation } from 'react-navi';

import api from '../../api';
import utils from '../../utils';
import theme from '../../theme';
import { renderLabels } from '../../helpers/labels';
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
} from '../../components/core';

const ScheduleTypeConditional = 'Conditional';
const ScheduleTypeAll = 'AllDevices';
const ScheduleTypeNone = 'NoDevices';

const Scheduling = ({
  route: {
    data: { params, application, devices, releases },
  },
}) => {
  const [conditionalQuery, setConditionalQuery] = useState(
    application.schedulingRule.conditionalQuery || []
  );
  const [scheduleType, setScheduleType] = useState(
    application.schedulingRule.scheduleType || 'NoDevices'
  );
  const [scheduledDevices, setScheduledDevices] = useState([]);
  const [backendError, setBackendError] = useState();
  const [showConditionPopup, setShowConditionPopup] = useState();
  const [conditionToEdit, setConditionToEdit] = useState(null);
  const [searchInput, setSearchInput] = useState('');
  const [searchFocused, setSearchFocused] = useState();
  const [showPreview, setShowPreview] = useState();

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
      },
    ],
    []
  );
  const tableData = useMemo(() => scheduledDevices, [scheduledDevices]);

  const navigation = useNavigation();

  let isSubmitDisabled = false;

  if (scheduleType === ScheduleTypeConditional) {
    if (
      conditionalQuery.length === 0 ||
      utils.deepEqual(
        conditionalQuery,
        application.schedulingRule.conditionalQuery
      )
    ) {
      isSubmitDisabled = true;
    }
  } else if (scheduleType === application.scheduleType) {
    isSubmitDisabled = true;
  }

  const submit = async () => {
    setBackendError(null);

    try {
      await api.updateApplication({
        projectId: params.project,
        applicationId: application.name,
        data: {
          name: application.name,
          description: application.description,
          schedulingRule: {
            ...application.schedulingRule,
            scheduleType,
            conditionalQuery,
          },
        },
      });

      toaster.success('Scheduling was successful.');
      navigation.navigate(
        `/${params.project}/applications/${application.name}`
      );
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Scheduling was not successful.');
      console.log(error);
    }
  };

  const getScheduledDevices = async () => {
    if (
      scheduleType === ScheduleTypeConditional &&
      conditionalQuery.length === 0
    ) {
      setScheduledDevices([]);
      return;
    }

    try {
      const { data: devices } = await api.scheduledDevices({
        projectId: params.project,
        applicationId: application.name,
        schedulingRule: {
          ...application.schedulingRule,
          scheduleType,
          conditionalQuery: conditionalQuery,
        },
        search: searchInput,
      });
      setScheduledDevices(devices);
    } catch (error) {
      console.log(error);
    }
  };

  const removeCondition = index => {
    setConditionalQuery(conditionalQuery.filter((_, i) => i !== index));
  };

  const addCondition = condition => {
    setShowConditionPopup(false);
    if (conditionToEdit !== null) {
      setConditionalQuery(conditionalQuery =>
        conditionalQuery.map((rule, index) =>
          index === conditionToEdit ? condition : rule
        )
      );
    } else {
      setConditionalQuery(conditionalQuery => [...conditionalQuery, condition]);
    }
    setConditionToEdit(null);
  };

  const clearConditions = () => {
    setConditionalQuery([]);
  };

  useEffect(() => {
    getScheduledDevices();
  }, [searchInput]);

  return (
    <>
      <Card
        title="Scheduling"
        size="xlarge"
        actions={
          scheduleType === ScheduleTypeConditional
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
            <Group>
              <Label>This application will run on</Label>
              <RadioGroup
                onChange={setScheduleType}
                value={scheduleType}
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
            {scheduleType === ScheduleTypeConditional && (
              <Column>
                <Group>
                  <Row justifyContent="space-between" alignItems="center">
                    <Label marginBottom={0}>Conditions</Label>
                    <Row>
                      {conditionalQuery.length > 0 && (
                        <Button
                          title="Clear"
                          variant="text"
                          onClick={clearConditions}
                          marginRight={4}
                        />
                      )}
                      <Button
                        title="Add Conditions"
                        variant="secondary"
                        onClick={() => setShowConditionPopup(true)}
                      />
                    </Row>
                  </Row>
                  <Row
                    bg="grays.0"
                    borderRadius={1}
                    minHeight={7}
                    padding={2}
                    marginTop={4}
                  >
                    {conditionalQuery.length ? (
                      <DevicesFilterButtons
                        canRemoveFilter
                        query={conditionalQuery}
                        removeFilter={removeCondition}
                        onEdit={index => {
                          setShowConditionPopup(true);
                          setConditionToEdit(index);
                        }}
                      />
                    ) : (
                      <Row flex={1} justifyContent="center" alignItems="center">
                        <Text>
                          Add <strong>Conditions</strong>.
                        </Text>
                      </Row>
                    )}
                  </Row>
                </Group>
              </Column>
            )}
            <Button
              marginTop={6}
              title="Apply Scheduling"
              onClick={submit}
              disabled={isSubmitDisabled}
            />
          </Column>
          {showConditionPopup && (
            <DevicesFilter
              title="Conditions"
              buttonTitle="Set Conditions"
              filter={
                conditionToEdit !== null && conditionalQuery[conditionToEdit]
              }
              whitelistedConditions={[LabelValueCondition]}
              onClose={() => {
                setShowConditionPopup(false);
                setConditionToEdit(null);
              }}
              onSubmit={addCondition}
            />
          )}
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
              color={searchFocused ? theme.colors.primary : theme.colors.white}
              style={{ position: 'absolute', left: 16 }}
            />
            <Input
              bg="black"
              placeholder="Search devices by name or labels"
              paddingLeft={7}
              value={searchInput}
              width="350px"
              onChange={e => setSearchInput(e.target.value)}
              onFocus={e => setSearchFocused(true)}
              onBlur={e => setSearchFocused(false)}
            />
          </Row>

          <Table
            columns={columns}
            data={tableData}
            placeholder={
              conditionalQuery.length > 0 ? (
                <Text>
                  There are no <strong>Devices</strong>.
                </Text>
              ) : (
                <Text>
                  Add <strong>Conditions</strong> to preview{' '}
                  <strong>Devices</strong>.
                </Text>
              )
            }
          />
        </Card>
      </Popup>
    </>
  );
};

export default Scheduling;
