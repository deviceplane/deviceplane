import React, { useEffect, useState, useMemo } from 'react';

import config from '../config';
import Layout from '../components/layout';
import Card from '../components/card';
import { Label, Group, Text, Link, Code, Select } from '../components/core';

const getLocalCommand = ({ id, projectId }) =>
  [
    'go run cmd/agent/main.go',
    '--controller=http://localhost:8080/api',
    '--conf-dir=./cmd/agent/conf',
    '--state-dir=./cmd/agent/state',
    '--log-level=debug',
    `--project=${projectId}`,
    `--registration-token=${id}`,
    '# note, this is the local version',
  ].join(' ');

const getCommand = ({ id, projectId }) =>
  [
    'curl https://downloads.deviceplane.com/install.sh',
    '|',
    `VERSION=${config.agentVersion}`,
    `PROJECT=${projectId}`,
    `REGISTRATION_TOKEN=${id}`,
    'bash',
  ].join(' ');

const AddDevice = ({
  route: {
    data: { params, registrationTokens },
  },
}) => {
  const selectOptions = useMemo(
    () =>
      registrationTokens.map(token => ({
        label: token.name,
        value: JSON.stringify(token),
      })),
    []
  );
  const [selectValue, setSelectValue] = useState(() => {
    if (registrationTokens.length) {
      const defaultToken = registrationTokens.find(
        ({ name }) => name === 'default'
      );
      return defaultToken
        ? JSON.stringify(defaultToken)
        : JSON.stringify(registrationTokens[0]);
    }
    return null;
  });
  const selection = selectValue && JSON.parse(selectValue);

  useEffect(() => {
    if (selection) {
      console.log(getLocalCommand(selection));
    }
  }, [selection]);

  return (
    <Layout alignItems="center">
      <Card title="Register Device" size="large">
        {registrationTokens && registrationTokens.length > 0 ? (
          <>
            <Group>
              <Label>Registration Token</Label>
              <Select
                options={selectOptions}
                value={selectValue}
                onChange={e => setSelectValue(e.target.value)}
              />
            </Group>
            {selection && (
              <>
                <Text marginBottom={2} fontWeight={1}>
                  Run the following command on the device you want to register:
                </Text>
                <Code>{getCommand(selection)}</Code>
              </>
            )}
          </>
        ) : (
          <>
            <Text>
              Create a <strong>Registration Token</strong> from the{' '}
              <Link href={`/${params.project}/provisioning`}>Provisioning</Link>{' '}
              page to enable device registration.{' '}
            </Text>
          </>
        )}
      </Card>
    </Layout>
  );
};

export default AddDevice;
