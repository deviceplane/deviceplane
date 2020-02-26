import React from 'react';

import {
  LabelValueCondition,
  LabelExistenceCondition,
  DevicePropertyCondition,
} from './devices-filter';
import { labelColor } from '../helpers/labels';
import { Row, Text, Button, Badge, Icon } from './core';

const ConditionComp = ({ type, params }) => {
  switch (type) {
    case LabelValueCondition:
      return (
        <>
          <Text
            fontSize={0}
            fontWeight={2}
            marginRight={2}
            color={labelColor(params.key)}
          >
            {params.key}
          </Text>

          <Text fontSize={0} fontWeight={1} marginRight={2} color="grays.12">
            {params.operator}
          </Text>

          <Text fontWeight={2} fontSize={0}>
            {params.value}
          </Text>
        </>
      );
    case LabelExistenceCondition:
      return (
        <>
          <Text
            fontWeight={2}
            marginRight={2}
            color={labelColor(params.key)}
            fontSize={0}
          >
            {params.key}
          </Text>

          <Text fontSize={0} fontWeight={1} color="grays.12">
            {params.operator}
          </Text>
        </>
      );
    case DevicePropertyCondition:
      return (
        <>
          <Text
            fontWeight={2}
            marginRight={2}
            fontSize={0}
            style={{ textTransform: 'uppercase' }}
          >
            {params.property}
          </Text>

          <Text fontSize={0} fontWeight={1} marginRight={2} color="grays.12">
            {params.operator}
          </Text>

          <Text
            fontWeight={2}
            fontSize={0}
            style={{ textTransform: 'capitalize' }}
          >
            {params.value}
          </Text>
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
  onEdit,
}) => {
  return (
    <Row flexWrap="wrap">
      {query.map((filter, index) => (
        <Row
          key={index}
          marginY={2}
          marginRight={2}
          border={0}
          borderRadius={5}
          borderColor="primary"
          paddingY={2}
          paddingX={3}
          alignItems="center"
          style={{ cursor: canRemoveFilter ? 'pointer' : 'default' }}
          onClick={canRemoveFilter ? () => onEdit(index) : () => {}}
        >
          {filter.map((condition, i) => (
            <React.Fragment key={i}>
              <ConditionComp {...condition} />

              {i < filter.length - 1 && (
                <Text fontSize={0} fontWeight={2} marginX={4} color="primary">
                  OR
                </Text>
              )}
            </React.Fragment>
          ))}
          {canRemoveFilter && (
            <Button
              marginLeft={2}
              variant="text"
              title={<Icon icon="cross" color="red" size={14} />}
              onClick={e => {
                e.stopPropagation();
                if (removeFilter) {
                  removeFilter(index);
                }
              }}
            />
          )}
        </Row>
      ))}
    </Row>
  );
};
