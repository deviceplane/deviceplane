import React, { Fragment, Component } from 'react';
import {
  Pane,
  Dialog,
  Select,
  IconButton,
  Button,
  majorScale,
  Strong,
  TextInput,
  minorScale
  // @ts-ignore
} from 'evergreen-ui';
import utils from '../utils';

export type Query = Filter[]

export type Filter = Condition[]

export type Condition = {
  type:  ConditionType
  params: ConditionParams
}

type ConditionType = string
export const DevicePropertyCondition: ConditionType = "DevicePropertyCondition";
export const	LabelValueCondition: ConditionType = "LabelValueCondition";
export const LabelExistenceCondition: ConditionType = "LabelExistenceCondition";

export type ConditionParams = DevicePropertyConditionParams |
  LabelValueConditionParams |
  LabelExistenceConditionParams;

export type DevicePropertyConditionParams = {
  property: string
  operator: Operator
  value:    string
}

export type LabelValueConditionParams = {
  key: string
  operator: Operator
  value: string
}

export type LabelExistenceConditionParams = {
  key: string
  operator: Operator
}

type Operator = string;
export const OperatorIs:  Operator = "is";
export const OperatorIsNot:  Operator = "is not";
export const OperatorExists:  Operator = "exists";
export const OperatorNotExists:  Operator = "does not exist";

interface Props {
  show: boolean
  onClose: () => void
  onSubmit?: (filter: Filter) => void
}

interface State {
  filter: Filter
}

function initialCondition(): Condition {
  return {
    type: DevicePropertyCondition,
    params: {
      property: 'status',
      operator: OperatorIs,
      value: 'online',
    }
  }
}

function initialState(): State {
  return {
      filter: [
      initialCondition(),
    ]
  }
}

export class DevicesFilter extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = initialState();
  }

  renderCondition = (condition: Condition, index: number) => {
    if (condition.type === LabelValueCondition) {
      let cond = condition.params as LabelValueConditionParams;
      const selectClassName: string = utils.randomClassName();
      return (
        <>
          <Pane
            display="flex"
            flexDirection="column"
            flex="1"
          >
            <TextInput
              width="auto"
              placeholder="Key"
              marginBottom={minorScale(2)}
              onChange={(event: any) => {
                const { value: key } = event.target;
                this.setState({
                  filter: this.state.filter.map((condition: any, i) => {
                    if (i === index) {
                      condition.params.key = key;
                    }
                    return condition;
                  })
                })
              }}
            />

            <Select
              className={selectClassName}
              marginBottom={minorScale(2)}
              value={cond.operator}
              onChange={(event: any) => {
                const { value: operator } = event.target;
                this.setState({
                  filter: this.state.filter.map((condition, i) => {
                    if (i === index) {
                      condition.params.operator = operator;
                    }
                    return condition;
                  })
                });
              }}
            >
              <option value={OperatorIs}>{OperatorIs}</option>
              <option value={OperatorIsNot}>{OperatorIsNot}</option>
            </Select>
            <style>{`
              .${selectClassName} > select {
                padding-top: 7px;
                padding-bottom: 7px;
              }
            `}</style>
            <TextInput
              width="auto"
              placeholder="Value"
              onChange={(event: any) => {
                const { value: value } = event.target;
                this.setState({
                  filter: this.state.filter.map((condition: any, i) => {
                    if (i === index) {
                      condition.params.value = value;
                    }
                    return condition;
                  })
                })
              }}
            />
          </Pane>
        </>
      );
    }

    if (condition.type === LabelExistenceCondition) {
      let cond = condition.params as LabelExistenceConditionParams;
      return (
        <>
          <Pane display="flex" flex="1" marginRight={majorScale(1)}>
            <TextInput
              width="auto"
              placeholder="Key"
              onChange={(event: any) => {
                const { value: key } = event.target;
                this.setState({
                  filter: this.state.filter.map((condition: any, i) => {
                    if (i === index) {
                      condition.params.key = key;
                    }
                    return condition;
                  })
                })
              }}
            />
          </Pane>
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
                })
              });
            }}
            marginRight={majorScale(1)}
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
                })
              });
            }}
            marginRight={majorScale(1)}
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
                })
              });
            }}
            marginRight={majorScale(1)}
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
                })
              })
            }}
            marginRight={majorScale(1)}
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
      <Pane>
        <Dialog
          isShown={show}
          title="Filter Devices"
          onCloseComplete={onClose}
          onConfirm={() => {
            if (onSubmit) {
              onSubmit(filter);
            }
            this.setState(initialState());
          }}
          confirmLabel="Filter"
          hasCancel={false}
          style={{ maxHeight: majorScale(12), overflowY: 'auto' }}
        >
          <Pane display="flex" flexDirection="column">
            {filter.map((condition, index) => (
              <Fragment key={index}>
                <Pane display="flex" justifyContent="space-around">
                  <Select
                    value={condition.type}
                    onChange={(event: Event) => {
                      if (event.target == null) {
                        return;
                      }
                      const { value: property } = event.target as HTMLSelectElement;
                      this.setState({
                        filter: filter.map((condition, i) => {
                          if (i === index) {
                            if (condition.type !== property) {

                              let params: ConditionParams;
                              if (property === DevicePropertyCondition) {
                                let p: DevicePropertyConditionParams = {
                                  property: '',
                                  operator: OperatorIs,
                                  value: '',
                                }
                                params = p;
                              } else if (property === LabelValueCondition) {
                                let p: LabelValueConditionParams = {
                                  key: '',
                                  operator: OperatorIs,
                                  value: ''
                                }
                                params = p;
                              } else if (true || property === LabelExistenceCondition) {
                                let p: LabelExistenceConditionParams = {
                                  key: '',
                                  operator: OperatorExists,
                                }
                                params = p;
                              }
                              condition = {
                                type: property,
                                params,
                              }
                            }
                          }
                          return condition;
                        })
                      });
                    }}
                    className={selectClassName}
                    marginRight={majorScale(1)}
                  >
                    <option value={DevicePropertyCondition}>Device Property</option>
                    <option value={LabelValueCondition}>Label Value</option>
                    <option value={LabelExistenceCondition}>Label Existence</option>
                  </Select>
                  <style>{`
                    .${selectClassName} > select {
                      width: auto;
                    }
                  `}</style>

                  {this.renderCondition(condition, index)}

                  {index > 0 ? (
                    <IconButton
                      icon="cross"
                      intent="danger"
                      appearance="minimal"
                      onClick={() =>
                        this.setState({
                          filter: filter.filter((_, i) => i !== index)
                        })
                      }
                    />
                  ) : (
                    <Pane width={majorScale(4)} />
                  )}
                </Pane>
                {index < filter.length - 1 && (
                  <Pane marginY={majorScale(2)}>
                    <Strong
                      size={300}
                      paddingX={majorScale(2)}
                      paddingY={majorScale(1)}
                      backgroundColor="#E4E7EB"
                      borderRadius={3}
                    >
                      OR
                    </Strong>
                  </Pane>
                )}
              </Fragment>
            ))}
          </Pane>
          <Pane
          display="flex"
          flexDirection="column"
          marginTop={majorScale(4)}>
            <Pane>
              <Button
                intent="none"
                iconBefore="plus"
                onClick={() => {
                  this.setState({
                    filter: [
                      ...filter,
                      initialCondition(),
                    ]
                  });
                }}
              >
                Add Condition
              </Button>
            </Pane>
          </Pane>
        </Dialog>
      </Pane>
    );
  }
}
