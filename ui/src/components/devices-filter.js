import React, { useRef, useMemo, useState } from 'react';

import utils from '../utils';
import Card from './card';
import Popup from './popup';
import { Column, Row, Group, Button, Input, Text, Select, Icon } from './core';

export const DevicePropertyCondition = 'DevicePropertyCondition';
export const LabelValueCondition = 'LabelValueCondition';
export const LabelExistenceCondition = 'LabelExistenceCondition';

export const OperatorIs = 'is';
export const OperatorIsNot = 'is not';
export const OperatorExists = 'exists';
export const OperatorNotExists = 'does not exist';

const DefaultDevicePropertyConditionParams = () => {
  return {
    property: 'status',
    operator: OperatorIs,
    value: 'online',
  };
};

const DefaultLabelValueConditionParams = () => {
  return {
    key: '',
    operator: OperatorIs,
    value: '',
  };
};

const DefaultLabelExistenceConditionParams = () => {
  return {
    key: '',
    operator: OperatorExists,
  };
};

const OperatorExistenceOptions = [
  { label: OperatorExists, value: OperatorExists },
  {
    label: OperatorNotExists,
    value: OperatorNotExists,
  },
];

export const OperatorIsOptions = [
  { label: OperatorIs, value: OperatorIs },
  { label: OperatorIsNot, value: OperatorIsNot },
];

const DevicePropertyOptions = [{ label: 'Status', value: 'status' }];

const DeviceValueOptions = [
  {
    label: 'Online',
    value: 'online',
  },
  {
    label: 'Offline',
    value: 'offline',
  },
];

export const DevicesFilter = ({
  title = 'Filter Devices',
  onClose,
  onSubmit,
  buttonTitle,
  whitelistedConditions,
  ...props
}) => {
  const filterListEndRef = useRef();

  const conditionOptions = useMemo(
    () =>
      [
        {
          value: DevicePropertyCondition,
          label: 'Device Property',
        },
        {
          value: LabelValueCondition,
          label: 'Label Value',
        },
        {
          value: LabelExistenceCondition,
          label: 'Label Existence',
        },
      ].filter(c => {
        if (!whitelistedConditions) {
          return true;
        }
        return whitelistedConditions.includes(c.value);
      }),
    [whitelistedConditions]
  );

  const defaultCondition = useMemo(
    () =>
      [
        {
          type: DevicePropertyCondition,
          params: DefaultDevicePropertyConditionParams(),
        },
        {
          type: LabelValueCondition,
          params: DefaultLabelValueConditionParams(),
        },
        {
          type: LabelExistenceCondition,
          params: DefaultLabelExistenceConditionParams(),
        },
      ].filter(c => {
        if (!whitelistedConditions) {
          return true;
        }
        return whitelistedConditions.includes(c.type);
      })[0],
    [whitelistedConditions]
  );

  if (!defaultCondition) {
    throw 'No default condition was whitelisted';
  }

  const [filter, setFilter] = useState(
    props.filter || [utils.deepClone(defaultCondition)]
  );

  const resetFilter = () => setFilter([utils.deepClone(defaultCondition)]);

  const renderCondition = (condition, index) => {
    if (condition.type === LabelValueCondition) {
      return (
        <>
          <Row flex={2}>
            <Input
              placeholder="Label Key"
              padding={2}
              value={condition.params.key}
              onChange={event => {
                const { value: key } = event.target;
                setFilter(
                  filter.map((condition, i) => {
                    if (i === index) {
                      condition.params.key = key;
                    }
                    return condition;
                  })
                );
              }}
            />
          </Row>

          <Select
            multi
            marginX={2}
            flex="0 0 130px"
            placeholder="Operator"
            options={OperatorIsOptions}
            value={condition.params.operator}
            onChange={e => {
              setFilter(
                filter.map((condition, i) => {
                  if (i === index) {
                    condition.params.operator = e.target.value;
                  }
                  return condition;
                })
              );
            }}
          />

          <Row flex={2}>
            <Input
              placeholder="Label Value"
              padding={2}
              value={condition.params.value}
              onChange={event => {
                const { value: value } = event.target;
                setFilter(
                  filter.map((condition, i) => {
                    if (i === index) {
                      condition.params.value = value;
                    }
                    return condition;
                  })
                );
              }}
            />
          </Row>
        </>
      );
    }

    if (condition.type === LabelExistenceCondition) {
      return (
        <>
          <Input
            placeholder="Label Key"
            padding={2}
            marginRight={2}
            value={condition.params.key}
            onChange={event => {
              const { value: key } = event.target;
              setFilter(
                filter.map((condition, i) => {
                  if (i === index) {
                    condition.params.key = key;
                  }
                  return condition;
                })
              );
            }}
          />

          <Select
            flex="0 0 150px"
            value={condition.params.operator}
            placeholder="Operator"
            options={OperatorExistenceOptions}
            onChange={e => {
              setFilter(
                filter.map((condition, i) => {
                  if (i === index) {
                    condition.params.operator = e.target.value;
                  }
                  return condition;
                })
              );
            }}
          />
        </>
      );
    }

    if (condition.type === DevicePropertyCondition) {
      return (
        <>
          <Select
            value={condition.params.property}
            placeholder="Property"
            options={DevicePropertyOptions}
            onChange={e => {
              setFilter(
                filter.map((condition, i) => {
                  if (i === index) {
                    condition.params.property = e.target.value;
                  }
                  return condition;
                })
              );
            }}
          />

          <Select
            marginX={2}
            value={condition.params.operator}
            placeholder="Operator"
            options={OperatorIsOptions}
            onChange={e => {
              setFilter(
                filter.map((condition, i) => {
                  if (i === index) {
                    condition.params.operator = e.target.value;
                  }
                  return condition;
                })
              );
            }}
          />

          <Select
            value={condition.params.value}
            placeholder="Value"
            options={DeviceValueOptions}
            onChange={e => {
              setFilter(
                filter.map((condition, i) => {
                  if (i === index) {
                    condition.params.value = e.target.value;
                  }
                  return condition;
                })
              );
            }}
          />
        </>
      );
    }
  };

  return (
    <Popup
      show={true}
      onClose={() => {
        onClose();
        resetFilter();
      }}
    >
      <Card
        border
        size="xlarge"
        title={title}
        actions={[
          {
            title: 'Add Condition',
            variant: 'secondary',
            onClick: () => {
              setFilter([...filter, utils.deepClone(defaultCondition)]);
              setTimeout(
                () =>
                  filterListEndRef.current.scrollIntoView({
                    behavior: 'smooth',
                  }),
                100
              );
            },
          },
        ]}
      >
        <Column flex={1} overflowY="auto" maxHeight="100%">
          {filter.map((condition, index) => (
            <Group key={index} marginBottom={4}>
              <Row justifyContent="space-between" alignItems="center">
                {conditionOptions.length > 1 && (
                  <Select
                    marginRight={2}
                    flex="0 0 200px"
                    value={condition.type}
                    placeholder="Type"
                    options={conditionOptions}
                    onChange={e => {
                      const type = e.target.value;
                      setFilter(
                        filter.map((condition, i) => {
                          if (i !== index) {
                            return condition;
                          }
                          if (condition.type === type) {
                            return condition;
                          }

                          let params;
                          switch (type) {
                            case LabelValueCondition:
                              params = DefaultLabelValueConditionParams();
                              break;
                            case LabelExistenceCondition:
                              params = DefaultLabelExistenceConditionParams();
                              break;
                            case DevicePropertyCondition:
                            default:
                              params = DefaultDevicePropertyConditionParams();
                              break;
                          }
                          condition = {
                            type,
                            params,
                          };
                          return condition;
                        })
                      );
                    }}
                  />
                )}

                {renderCondition(condition, index)}

                {index > 0 && (
                  <Button
                    title={<Icon icon="cross" size={16} color="red" />}
                    marginLeft={2}
                    variant="iconDanger"
                    onClick={() =>
                      setFilter(filter.filter((_, i) => i !== index))
                    }
                  />
                )}
              </Row>
              {index < filter.length - 1 && (
                <Row marginTop={4}>
                  <Text fontWeight={2} fontSize={2} color="primary">
                    OR
                  </Text>
                </Row>
              )}
            </Group>
          ))}
          <Row ref={filterListEndRef} />
        </Column>

        <Button
          marginTop={3}
          title={buttonTitle || (props.filter ? 'Edit Filter' : 'Apply Filter')}
          onClick={() => {
            const validFilter = filter.filter(({ type, params }) => {
              switch (type) {
                case LabelValueCondition:
                  return (
                    params.key !== '' && params.value !== '' && params.operator
                  );
                case LabelExistenceCondition:
                  return params.key !== '' && params.operator;
                case DevicePropertyCondition:
                default:
                  return (
                    params.property && params.operator && params.value !== ''
                  );
              }
            });

            if (validFilter.length) {
              if (onSubmit) {
                onSubmit(validFilter);
              }
            } else {
              onClose();
            }

            resetFilter();
          }}
        />
      </Card>
    </Popup>
  );
};

export default DevicesFilter;
