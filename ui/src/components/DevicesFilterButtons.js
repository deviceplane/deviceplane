import React from 'react';
import { Icon } from 'evergreen-ui';
import {
  LabelValueCondition,
  LabelExistenceCondition,
  DevicePropertyCondition,
} from './DevicesFilter';

import theme from '../theme';
import { Row, Text, Button, Badge } from './core';

const ConditionComp = ({ type, params }) => {
  switch (type) {
    case LabelValueCondition:
      return (
        <>
          <Text
            fontWeight={3}
            marginRight={2}
            color="primary"
            style={{ textTransform: 'none' }}
          >
            {params.key}
          </Text>

          <Text fontWeight={2} marginRight={2}>
            {params.operator}
          </Text>

          <Text fontWeight={3} style={{ textTransform: 'none' }}>
            {params.value}
          </Text>
        </>
      );
    case LabelExistenceCondition:
      return (
        <>
          <Text
            fontWeight={3}
            marginRight={2}
            color="primary"
            style={{ textTransform: 'none' }}
          >
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
          <Text fontWeight={3} marginRight={2} color="primary">
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
  onEdit,
}) => {
  return (
    <Row flexWrap="wrap">
      {query.map((filter, index) => (
        <Row alignItems="center" key={index} margin={2} marginLeft={0}>
          <Row
            border={0}
            borderRadius={1}
            borderColor="white"
            padding={2}
            alignItems="center"
            style={{ cursor: canRemoveFilter ? 'pointer' : 'default' }}
            onClick={canRemoveFilter ? () => onEdit(index) : () => {}}
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
          </Row>
          {canRemoveFilter && (
            <Button
              marginLeft={2}
              variant="icon"
              title={<Icon icon="cross" color={theme.colors.red} size={14} />}
              onClick={() => (removeFilter ? removeFilter(index) : null)}
            />
          )}
        </Row>
      ))}
    </Row>
  );
};
