import React, { useState } from 'react';
import moment from 'moment';
import { useNavigation } from 'react-navi';

import api from '../../api';
import utils from '../../utils';
import useToggle from '../../hooks/useToggle';
import Editor from '../../components/editor';
import Card from '../../components/card';
import Popup from '../../components/popup';
import Alert from '../../components/alert';
import {
  Row,
  Group,
  Text,
  Button,
  Link,
  Label,
  Value,
  toaster,
} from '../../components/core';

const ReleasedBy = ({ project, release }) => {
  if (release) {
    if (release.createdByUser) {
      return (
        <Link href={`/${project}/iam/members/${release.createdByUser.id}`}>
          {release.createdByUser.name}
        </Link>
      );
    } else if (release.createdByServiceAccount) {
      return (
        <Link
          href={`/${project}/iam/service-accounts/${release.createdByServiceAccount.name}`}
        >
          {release.createdByServiceAccount.name}
        </Link>
      );
    }
  }
  return '-';
};

const Release = ({
  route: {
    data: { params, release, application },
  },
}) => {
  const [backendError, setBackendError] = useState();
  const [isConfirmPopup, toggleConfirmPopup] = useToggle();
  const navigation = useNavigation();

  const revertRelease = async () => {
    setBackendError(null);
    try {
      await api.createRelease({
        projectId: params.project,
        applicationId: application.id,
        data: { rawConfig: release.rawConfig },
      });
      navigation.navigate(`/${params.project}/applications/${application.id}`);
      toaster.success('Release reverted.');
    } catch (error) {
      setBackendError(utils.parseError(error, 'Release reversion failed.'));
      console.error(error);
    }
    toggleConfirmPopup();
  };

  return (
    <>
      <Card size="xlarge">
        <Text fontWeight={2} fontSize={5} marginBottom={6}>
          Release {release.number}
        </Text>

        <Alert show={backendError} variant="error" description={backendError} />

        <Group>
          <Label>Released By</Label>
          <Row>
            <ReleasedBy project={params.project} release={release} />
          </Row>
        </Group>

        <Group>
          <Label>Started</Label>
          <Value>{moment(release.createdAt).fromNow()}</Value>
        </Group>

        <Label>Config</Label>
        <Editor readOnly width="100%" value={release.rawConfig} />
        <Button
          marginTop={6}
          title="Revert to this Release"
          onClick={toggleConfirmPopup}
        />
      </Card>

      <Popup show={isConfirmPopup} onClose={toggleConfirmPopup}>
        <Card title="Revert Release" border size="large">
          <Text>
            This will create a new release to application{' '}
            <strong>{params.application}</strong> using the config from release{' '}
            <strong>{release.number}</strong>.
          </Text>
          <Button marginTop={5} title="Revert" onClick={revertRelease} />
        </Card>
      </Popup>
    </>
  );
};

export default Release;
