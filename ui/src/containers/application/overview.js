import React from 'react';

import Editor from '../../components/editor';
import Card from '../../components/card';
import { Group, Link, Label, Value } from '../../components/core';
import { DevicesFilterButtons } from '../../components/DevicesFilterButtons';

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
          <DevicesFilterButtons
            query={schedulingRule}
            canRemoveFilter={false}
          />
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
            <Label>Current Release ID</Label>
            <Link
              href={`/${params.project}/applications/${name}/releases/${latestRelease.id}`}
            >
              {latestRelease.id}
            </Link>
          </Group>

          <Group>
            <Label>Current Config</Label>
            <Editor
              width="100%"
              height="150px"
              value={latestRelease.rawConfig}
              readOnly
            />
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
