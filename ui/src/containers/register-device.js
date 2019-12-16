import React, { useEffect } from 'react';

import config from '../config';
import Layout from '../components/layout';
import Card from '../components/card';
import { Row, Text, Link, Code } from '../components/core';

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

const getDockerCommand = ({ id, projectId }) =>
  [
    'curl https://install.deviceplane.com',
    '|',
    `VERSION=${config.agentVersion}`,
    `PROJECT=${projectId}`,
    `REGISTRATION_TOKEN=${id}`,
    'bash',
  ].join(' ');

const AddDevice = ({
  route: {
    data: { params, registrationToken },
  },
}) => {
  useEffect(() => {
    console.log(getLocalCommand(registrationToken));
  }, []);

  return (
    <Layout alignItems="center">
      <Card title="Register Device">
        {registrationToken ? (
          <>
            <Row marginBottom={4}>
              <Text>
                Default registration token with ID{' '}
                <Code>{registrationToken.id}</Code> is being used.
              </Text>
            </Row>
            <Text marginBottom={2}>
              Run the following command on the device you want to register:
            </Text>
            <Code>{getDockerCommand(registrationToken)}</Code>
          </>
        ) : (
          <>
            <Text>
              Create a <strong>default</strong> Registration Token from the{' '}
              <Link href={`/${params.project}/provisioning`}>Provisioning</Link>{' '}
              page to enable device registration from the UI.{' '}
            </Text>
          </>
        )}
      </Card>
    </Layout>
  );
};

export default AddDevice;
