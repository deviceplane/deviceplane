// @ts-nocheck

import React, { Component } from 'react';
import {
  Icon,
  Select,
  // @ts-ignore
} from 'evergreen-ui';

import utils from '../utils';
import { Column, Row, Group, Button, Input, Text } from './core';
import Card from './card';
import Popup from './popup';

export type Query = Filter[];

export type Filter = Condition[];

export type Condition = {
  type: ConditionType;
  params: ConditionParams;
};

type ConditionType = string;
export const DevicePropertyCondition: ConditionType = 'DevicePropertyCondition';
export const LabelValueCondition: ConditionType = 'LabelValueCondition';
export const LabelExistenceCondition: ConditionType = 'LabelExistenceCondition';

export type ConditionParams =
  | DevicePropertyConditionParams
  | LabelValueConditionParams
  | LabelExistenceConditionParams;

export type DevicePropertyConditionParams = {
  property: string;
  operator: Operator;
  value: string;
};

export type LabelValueConditionParams = {
  key: string;
  operator: Operator;
  value: string;
};

export type LabelExistenceConditionParams = {
  key: string;
  operator: Operator;
};

type Operator = string;
export const OperatorIs: Operator = 'is';
export const OperatorIsNot: Operator = 'is not';
export const OperatorExists: Operator = 'exists';
export const OperatorNotExists: Operator = 'does not exist';

const DefaultDevicePropertyConditionParams = (): DevicePropertyConditionParams => {
  return {
    property: 'status',
    operator: OperatorIs,
    value: 'online',
  };
};

const DefaultLabelValueConditionParams = (): LabelValueConditionParams => {
  return {
    key: '',
    operator: OperatorIs,
    value: '',
  };
};

const DefaultLabelExistenceConditionParams = (): LabelExistenceConditionParams => {
  return {
    key: '',
    operator: OperatorExists,
  };
};

interface Props {
  whitelistedConditions?: ConditionType[];
  show: boolean;
  onClose: () => void;
  onSubmit?: (filter: Filter) => void;
}

interface State {
  filter: Filter;
}

export class DevicesFilter extends Component<Props, State> {
  constructor(props: Props) {
    super(props);

    this.conditionOptions = [
      {
        type: DevicePropertyCondition,
        text: 'Device Property',
      },
      {
        type: LabelValueCondition,
        text: 'Label Value',
      },
      {
        type: LabelExistenceCondition,
        text: 'Label Existence',
      },
    ]
      .filter(c => {
        if (!this.props.whitelistedConditions) {
          return true;
        }
        return this.props.whitelistedConditions.includes(c.type);
      })
      .map(c => (
        <option key={c.type} value={c.type}>
          {c.text}
        </option>
      ));

    this.defaultCondition = [
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
      if (!this.props.whitelistedConditions) {
        return true;
      }
      return this.props.whitelistedConditions.includes(c.type);
    })[0];
    if (!this.defaultCondition) {
      throw 'No default condition was whitelisted';
    }

    this.state = {
      filter: [utils.deepClone(this.defaultCondition)],
    };
  }

  resetFilter() {
    this.setState({
      filter: [utils.deepClone(this.defaultCondition)],
    });
  }

  defaultCondition: Condition;
  conditionOptions: JSX.Element[];

  renderCondition = (condition: Condition, index: number) => {
    if (condition.type === LabelValueCondition) {
      let cond = condition.params as LabelValueConditionParams;
      const selectClassName: string = utils.randomClassName();
      return (
        <>
          <Column>
            <Input
              width="auto"
              placeholder="Key"
              padding={2}
              onChange={(event: any) => {
                const { value: key } = event.target;
                this.setState({
                  filter: this.state.filter.map((condition: any, i) => {
                    if (i === index) {
                      condition.params.key = key;
                    }
                    return condition;
                  }),
                });
              }}
            />

            <Select
              className={selectClassName}
              marginY={12}
              value={cond.operator}
              onChange={(event: any) => {
                const { value: operator } = event.target;
                this.setState({
                  filter: this.state.filter.map((condition, i) => {
                    if (i === index) {
                      condition.params.operator = operator;
                    }
                    return condition;
                  }),
                });
              }}
            >
              <option value={OperatorIs}>{OperatorIs}</option>
              <option value={OperatorIsNot}>{OperatorIsNot}</option>
            </Select>
            <Input
              width="auto"
              placeholder="Value"
              padding={2}
              onChange={(event: any) => {
                const { value: value } = event.target;
                this.setState({
                  filter: this.state.filter.map((condition: any, i) => {
                    if (i === index) {
                      condition.params.value = value;
                    }
                    return condition;
                  }),
                });
              }}
            />
          </Column>
        </>
      );
    }

    if (condition.type === LabelExistenceCondition) {
      let cond = condition.params as LabelExistenceConditionParams;
      return (
        <>
          <Input
            width="auto"
            placeholder="Key"
            padding={2}
            marginRight={3}
            onChange={(event: any) => {
              const { value: key } = event.target;
              this.setState({
                filter: this.state.filter.map((condition: any, i) => {
                  if (i === index) {
                    condition.params.key = key;
                  }
                  return condition;
                }),
              });
            }}
          />
          <Select
            value={cond.operator}
            onChange={(event: any) => {
              const { value: operator } = event.target;
              this.setState({
                filter: this.state.filter.map((condition, i) => {
                  if (i === index) {
                    condition.params.operator = operator;
                  }
                  return condition;
                }),
              });
            }}
          >
            <option value={OperatorExists}>{OperatorExists}</option>
            <option value={OperatorNotExists}>{OperatorNotExists}</option>
          </Select>
        </>
      );
    }

    if (condition.type === DevicePropertyCondition) {
      let cond = condition.params as DevicePropertyConditionParams;
      return (
        <>
          <Select
            value={cond.property}
            onChange={(event: any) => {
              const { value: property } = event.target;
              this.setState({
                filter: this.state.filter.map((condition: any, i) => {
                  if (i === index) {
                    condition.params.property = property;
                  }
                  return condition;
                }),
              });
            }}
            marginRight={16}
          >
            <option value={'status'}>Status</option>
          </Select>
          <Select
            value={cond.operator}
            onChange={(event: any) => {
              const { value: operator } = event.target;
              this.setState({
                filter: this.state.filter.map((condition, i) => {
                  if (i === index) {
                    condition.params.operator = operator;
                  }
                  return condition;
                }),
              });
            }}
            marginRight={16}
          >
            <option value={OperatorIs}>{OperatorIs}</option>
            <option value={OperatorIsNot}>{OperatorIsNot}</option>
          </Select>
          <Select
            value={cond.value}
            onChange={(event: any) => {
              const { value: value } = event.target;
              this.setState({
                filter: this.state.filter.map((condition: any, i) => {
                  if (i === index) {
                    condition.params.value = value;
                  }
                  return condition;
                }),
              });
            }}
          >
            <option value="online">Online</option>
            <option value="offline">Offline</option>
          </Select>
        </>
      );
    }
  };

  render() {
    const { show, onClose, onSubmit } = this.props;
    const { filter } = this.state;
    const selectClassName: string = utils.randomClassName();

    return (
      <Popup
        show={show}
        onClose={() => {
          onClose();
          this.resetFilter();
        }}
      >
        <Card
          border
          size="xlarge"
          title="Filter Devices"
          actions={[
            {
              title: 'Add Condition',
              variant: 'secondary',
              onClick: () => {
                this.setState({
                  filter: [...filter, utils.deepClone(this.defaultCondition)],
                });
              },
            },
          ]}
        >
          {filter.map((condition, index) => (
            <Group key={index}>
              <Row justifyContent="space-between" alignItems="center">
                <Select
                  value={condition.type}
                  onChange={event => {
                    if (event.target == null) {
                      return;
                    }
                    var { value: property } = event.target as HTMLSelectElement;
                    this.setState({
                      filter: filter.map((condition, i) => {
                        if (i !== index) {
                          return condition;
                        }
                        if (condition.type === property) {
                          return condition;
                        }

                        let params: ConditionParams;
                        switch (property) {
                          case DevicePropertyCondition:
                            params = DefaultDevicePropertyConditionParams();
                            break;
                          case LabelValueCondition:
                            params = DefaultLabelValueConditionParams();
                            break;
                          case LabelExistenceCondition:
                            params = DefaultLabelExistenceConditionParams();
                            break;
                          default:
                            property = DevicePropertyCondition;
                            params = DefaultDevicePropertyConditionParams();
                        }
                        condition = {
                          type: property,
                          params,
                        };
                        return condition;
                      }),
                    });
                  }}
                  className={selectClassName}
                  marginRight={16}
                >
                  {this.conditionOptions}
                </Select>
                <style>{`
                    .${selectClassName} > select {
                      width: auto;
                    }
                  `}</style>

                {this.renderCondition(condition, index)}

                {index > 0 && (
                  <Button
                    title={<Icon icon="cross" size={18} color="white" />}
                    marginLeft={2}
                    variant="icon"
                    onClick={() =>
                      this.setState({
                        filter: filter.filter((_, i) => i !== index),
                      })
                    }
                  />
                )}
              </Row>
              {index < filter.length - 1 && (
                <Row marginTop={6}>
                  <Text fontWeight={4} fontSize={3}>
                    OR
                  </Text>
                </Row>
              )}
            </Group>
          ))}
          <Button
            title="Apply Filter"
            onClick={() => {
              if (onSubmit) {
                onSubmit(filter);
              }
              this.resetFilter();
            }}
          />
        </Card>
      </Popup>
    );
  }
}
