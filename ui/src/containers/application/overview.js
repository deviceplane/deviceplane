import React from 'react';

import Editor from '../../components/editor';
import Card from '../../components/card';
import { Row, Group, Link, Label, Value } from '../../components/core';
import { DevicesFilterButtons } from '../../components/devices-filter-buttons';

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
        <Label>Scheduling Rule</Label>
        {schedulingRule.length ? (
          <Row bg="grays.0" borderRadius={1} minHeight={9} padding={2}>
            <DevicesFilterButtons
              query={schedulingRule}
              canRemoveFilter={false}
            />
          </Row>
        ) : (
          <Value>
            No scheduling rule set. You can set one in the{' '}
            <Link href={`/${params.project}/applications/${name}/scheduling`}>
              scheduling
            </Link>{' '}
            page.
          </Value>
        )}
      </Group>

      {latestRelease ? (
        <>
          <Group>
            <Label>Current Release</Label>
            <Link
              href={`/${params.project}/applications/${name}/releases/${latestRelease.id}`}
            >
              {latestRelease.id}
            </Link>
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
