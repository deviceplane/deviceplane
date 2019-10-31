import React, { Fragment, Component } from 'react';
import {
  Pane,
  Text,
  Dialog,
  Select,
  IconButton,
  Button,
  majorScale,
  Strong,
  TextInput,
  Icon,
  minorScale
  // @ts-ignore
} from 'evergreen-ui';
import utils from '../utils';
import { Query, Condition, LabelValueCondition, LabelValueConditionParams, LabelExistenceCondition, LabelExistenceConditionParams, DevicePropertyCondition, DevicePropertyConditionParams } from './DevicesFilter';

interface Props {
  query: Query
  canRemoveFilter: boolean,
  removeFilter?: (index: number) => void;
}

interface State {

}

export class DevicesFilterButtons extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
  }

  renderCondition = (condition: Condition) => {
    if (condition.type === LabelValueCondition) {
      let cond = condition.params as LabelValueConditionParams;
      return (
        <Fragment>
          <Text
            fontWeight={700}
            marginRight={minorScale(1)}
          >
            {cond.key}
          </Text>

          <Text fontWeight={500} marginRight={minorScale(1)}>
            {cond.operator}
          </Text>

          <Text fontWeight={700}>{cond.value}</Text>
        </Fragment>
      );
    }

    if (condition.type === LabelExistenceCondition) {
      let cond = condition.params as LabelExistenceConditionParams;
      return (
        <Fragment>
          <Text
            fontWeight={700}
            marginRight={minorScale(1)}
          >
            {cond.key}
          </Text>

          <Text fontWeight={500} marginRight={minorScale(1)}>
            {cond.operator}
          </Text>
        </Fragment>
      );
    }

    if (condition.type === DevicePropertyCondition) {
      let cond = condition.params as DevicePropertyConditionParams;
      return (
        <Fragment>
          <Text
            fontWeight={700}
            marginRight={minorScale(1)}
            style={{ textTransform: 'capitalize' }}
          >
            {cond.property}
          </Text>

          <Text fontWeight={500} marginRight={minorScale(1)}>
            {cond.operator}
          </Text>

          <Text style={{ textTransform: 'capitalize' }} fontWeight={700}>
            {cond.value}
          </Text>
        </Fragment>
      );
    }

    return (
      <Fragment>
        <Text
          fontWeight={500}
          marginRight={minorScale(1)}
          style={{ textTransform: 'capitalize' }}
        >
          Error rendering label.
        </Text>
      </Fragment>
    );
  };

  render() {
    return (
      <Pane
        paddingX={majorScale(2)}
        paddingBottom={majorScale(2)}
        display="flex"
        flexWrap="wrap"
        padding={5}
      >
        {this.props.query.map((filter, index) => (
          <Pane
            display="flex"
            alignItems="center"
            marginRight={minorScale(3)}
            key={index}
            margin={3}
          >
            <Pane
              backgroundColor="#B7D4EF"
              borderRadius={3}
              paddingX={minorScale(2)}
              paddingY={minorScale(1)}
              display="flex"
              alignItems="center"
            >
              {filter.map((condition, i) => (
                <Fragment key={i}>
                  {this.renderCondition(condition)}
                  {i < filter.length - 1 && (
                    <Text
                      fontSize={10}
                      fontWeight={700}
                      marginX={minorScale(3)}
                      color="white"
                    >
                      OR
                    </Text>
                  )
                }
                </Fragment>
              ))}
              {this.props.canRemoveFilter && (<Icon
                marginLeft={minorScale(3)}
                icon="cross"
                appearance="minimal"
                cursor="pointer"
                color="white"
                size={14}
                onClick={() => this.props.removeFilter ? this.props.removeFilter(index) : null}
              />)}
            </Pane>
          </Pane>
        ))}
      </Pane>
    )
  }
}