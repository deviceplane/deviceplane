import React from 'react';

import theme from '../../theme';
import Editor from '../../components/editor';
import Card from '../../components/card';
import { Row, Group, Link, Label, Value, Text } from '../../components/core';
import { DevicesFilterButtons } from '../../components/devices-filter-buttons';

const getSchedulingRule = schedulingRule => {
  if (
    schedulingRule.scheduleType === 'Conditional' &&
    schedulingRule.conditionalQuery
  ) {
    return (
      <Row bg="grays.0" borderRadius={1} minHeight={7} padding={2}>
        <DevicesFilterButtons
          query={schedulingRule.conditionalQuery}
          canRemoveFilter={false}
        />
      </Row>
    );
  } else {
    return (
      <Value>
        This application is running on{' '}
        <strong style={{ color: theme.colors.white }}>
          {schedulingRule.scheduleType === 'AllDevices'
            ? 'all devices'
            : 'no devices'}
        </strong>
        .
      </Value>
    );
  }
};

const ApplicationOverview = ({
  route: {
    data: {
      params,
      application: { description, latestRelease, name, schedulingRule },
    },
  },
}) => {
  return (
    <Card title={name} subtitle={description} size="xlarge">
      <Group>
        <Label>Scheduling</Label>
        {getSchedulingRule(schedulingRule)}
      </Group>

      {latestRelease ? (
        <>
          <Group>
            <Label>Current Release</Label>
            <Row>
              <Link
                href={`/${params.project}/applications/${name}/releases/${latestRelease.number}`}
              >
                {latestRelease.number}
              </Link>
            </Row>
          </Group>

          <Group>
            <Label>Current Config</Label>
            <Editor width="100%" value={latestRelease.rawConfig} readOnly />
          </Group>
        </>
      ) : (
        <Group>
          <Label>Current Release</Label>
          <Value>
            Create your first release on the{' '}
            <Link href={`/${params.project}/applications/${name}/releases`}>
              releases
            </Link>{' '}
            page.
          </Value>
        </Group>
      )}
    </Card>
  );
};

export default ApplicationOverview;
