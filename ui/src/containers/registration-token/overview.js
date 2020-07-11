import React from 'react';

import api from '../../api';
import { Column, Label, Group, Value } from '../../components/core';
import Card from '../../components/card';
import EditableLabelTable from '../../components/editable-label-table';

import storage from '../../storage';

const RegistrationTokenOverview = ({
  route: {
    data: { params, registrationToken },
  },
}) => {
  return (
    <>
      <Card
        size="xlarge"
        title={registrationToken.name}
        subtitle={registrationToken.description}
        marginBottom={5}
      >
        <Group>
          <Label>ID</Label>
          <Value>{registrationToken.id}</Value>
        </Group>
        <Group>
          <Label>Devices Registered</Label>
          <Value>{registrationToken.deviceCounts.allCount}</Value>
        </Group>
        <Group>
          <Label>Maximum Device Registerations</Label>
          <Value>{registrationToken.maxRegistrations || 'Unlimited'}</Value>
        </Group>
      </Card>
      <EditableLabelTable
        data={registrationToken.labels}
        onAdd={label =>
          api.addRegistrationTokenLabel({
            projectId: params.project,
            tokenId: registrationToken.id,
            data: label,
          })
        }
        onRemove={labelId =>
          api.removeRegistrationTokenLabel({
            projectId: params.project,
            tokenId: registrationToken.id,
            labelId,
          })
        }
        marginBottom={5}
      />
      {(storage.get('legacy') || false) && (
        <EditableLabelTable
          title="Environment Variables"
          dataName="Environment Variable"
          data={registrationToken.environmentVariables}
          onAdd={environmentVariable =>
            api.addRegistrationTokenEnvironmentVariable({
              projectId: params.project,
              tokenId: registrationToken.id,
              data: environmentVariable,
            })
          }
          onRemove={key =>
            api.removeRegistrationTokenEnvironmentVariable({
              projectId: params.project,
              tokenId: registrationToken.id,
              key,
            })
          }
        />
      )}
    </>
  );
};

export default RegistrationTokenOverview;
