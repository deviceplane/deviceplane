// @ts-nocheck

import React from 'react';
import {
  Icon,
  // @ts-ignore
} from 'evergreen-ui';
import {
  Query,
  Condition,
  LabelValueCondition,
  LabelValueConditionParams,
  LabelExistenceCondition,
  LabelExistenceConditionParams,
  DevicePropertyCondition,
  DevicePropertyConditionParams,
} from './DevicesFilter';

import { Row, Text, Button, Badge } from './core';

interface Props {
  query: Query;
  canRemoveFilter: boolean;
  removeFilter?: (index: number) => void;
}

interface State {}

const ConditionComp = ({ type, params }) => {
  switch (type) {
    case LabelValueCondition:
      return (
        <>
          <Text fontWeight={3} marginRight={2}>
            {params.key}
          </Text>

          <Text fontWeight={2} marginRight={2}>
            {params.operator}
          </Text>

          <Text fontWeight={3}>{params.value}</Text>
        </>
      );
    case LabelExistenceCondition:
      return (
        <>
          <Text fontWeight={3} marginRight={2}>
            {params.key}
          </Text>

          <Text fontWeight={2} marginRight={2}>
            {params.operator}
          </Text>
        </>
      );
    case DevicePropertyCondition:
      return (
        <>
          <Text fontWeight={3} marginRight={2}>
            {params.property}
          </Text>

          <Text fontWeight={2} marginRight={2}>
            {params.operator}
          </Text>

          <Text fontWeight={3}>{params.value}</Text>
        </>
      );
    default:
      return (
        <Text fontWeight={2} marginRight={2}>
          Error rendering label.
        </Text>
      );
  }
};

export const DevicesFilterButtons = ({
  query,
  removeFilter,
  canRemoveFilter,
}) => {
  return (
    <Row flexWrap="wrap">
      {query.map((filter, index) => (
        <Row alignItems="center" key={index} marginRight={4}>
          <Row
            border={0}
            borderRadius={1}
            borderColor="white"
            padding={2}
            alignItems="center"
          >
            {filter.map((condition, i) => (
              <React.Fragment key={i}>
                <Badge bg="grays.3">
                  <ConditionComp {...condition} />
                </Badge>

                {i < filter.length - 1 && (
                  <Text fontSize={0} fontWeight={3} marginX={4}>
                    OR
                  </Text>
                )}
              </React.Fragment>
            ))}
            {canRemoveFilter && (
              <Button
                marginLeft={2}
                variant="icon"
                title={<Icon icon="cross" color="white" size={14} />}
                onClick={() => (removeFilter ? removeFilter(index) : null)}
              />
            )}
          </Row>
        </Row>
      ))}
    </Row>
  );
};
