import React from 'react';

import api from '../../api';
import { Column, Label, Group, Value } from '../../components/core';
import Card from '../../components/card';
import EditableLabelTable from '../../components/editable-label-table';

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
        marginBottom={4}
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
      />
    </>
  );
};

export default RegistrationTokenOverview;
