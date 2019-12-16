import React from 'react';

import api from '../../api';
import { Column, Label, Value } from '../../components/core';
import Card from '../../components/card';
import EditableLabelTable from '../../components/EditableLabelTable';

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
      >
        <Column marginBottom={6}>
          <Label>ID</Label>
          <Value>{registrationToken.id}</Value>
        </Column>
        <Column marginBottom={6}>
          <Label>Devices Registered</Label>
          <Value>{registrationToken.deviceCounts.allCount}</Value>
        </Column>
        <Column>
          <Label>Maximum Device Registerations</Label>
          <Value>{registrationToken.maxRegistrations || 'Unlimited'}</Value>
        </Column>
      </Card>
      <Column marginTop={4}>
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
      </Column>
    </>
  );
};

export default RegistrationTokenOverview;
