import React, { useState } from 'react';
import moment from 'moment';
import { useNavigation } from 'react-navi';

import api from '../../api';
import utils from '../../utils';
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
} from '../../components/core';
import { toaster } from 'evergreen-ui';

const ReleasedBy = ({ project, release }) => {
  if (release) {
    if (release.createdByUser) {
      return (
        <Link href={`/${project}/iam/members/${release.createdByUser.id}`}>
          {release.createdByUser.firstName} {release.createdByUser.lastName}
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
  const [showConfirmPopup, setShowConfirmPopup] = useState();
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
      toaster.success('Release was reverted successfully.');
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Release was not reverted.');
      console.log(error);
    }
    setShowConfirmPopup(false);
  };

  return (
    <>
      <Card size="xlarge">
        <Text fontWeight={3} fontSize={5} marginBottom={6}>
          {release.id}
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
        <Editor
          readOnly
          width="100%"
          height="200px"
          value={release.rawConfig}
        />
        <Button
          marginTop={6}
          title="Revert to this Release"
          onClick={() => setShowConfirmPopup(true)}
        />
      </Card>

      <Popup show={showConfirmPopup} onClose={() => setShowConfirmPopup(false)}>
        <Card title="Revert Release" border size="large">
          <Text>
            This will create a new release to application{' '}
            <strong>{params.application}</strong> using the config from release{' '}
            <strong>{release.id}</strong>.
          </Text>
          <Button marginTop={5} title="Revert" onClick={revertRelease} />
        </Card>
      </Popup>
    </>
  );
};

export default Release;
