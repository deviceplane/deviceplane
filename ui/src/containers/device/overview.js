import React, { useMemo, useState, useEffect } from 'react';
import { useTable, useSortBy } from 'react-table';
import parsePrometheusTextFormat from 'parse-prometheus-text-format';

import api from '../../api';
import {
  Column,
  Row,
  Group,
  Button,
  Value,
  Link,
  Label,
  Text,
  Icon,
  toaster,
} from '../../components/core';
import Card from '../../components/card';
import Table from '../../components/table';
import Popup from '../../components/popup';
import Editor from '../../components/editor';
import EditableLabelTable from '../../components/editable-label-table';
import DeviceStatus from '../../components/device-status';

const DeviceOverview = ({
  route: {
    data: { params, device },
  },
}) => {
  return (
    <>
      <Card
        size="xlarge"
        title={device.name}
        subtitle={
          <DeviceStatus
            inline
            status={device.status}
            lastSeenAt={device.lastSeenAt}
          />
        }
        actions={[
          {
            title: 'Reboot',
            variant: 'secondary',
            disabled: device.status === 'offline',
            onClick: async () => {
              try {
                await api.reboot({
                  projectId: params.project,
                  deviceId: device.id,
                });
                toaster.success('Reboot was initiated.');
              } catch (error) {
                toaster.danger('Reboot failed.');
                console.error(error);
              }
            },
          },
          {
            title: 'SSH',
            variant: 'secondary',
            disabled: device.status === 'offline',
            newTab: true,
            href: `/${params.project}/ssh?devices=${device.name}`,
          },
        ]}
        marginBottom={5}
      >
        <Group>
          <Label>Agent Version</Label>
          <Value>
            {device.info.hasOwnProperty('agentVersion')
              ? device.info.agentVersion
              : ''}
          </Value>
        </Group>

        <Group>
          <Label>IP Address</Label>
          <Value>
            {device.info.hasOwnProperty('ipAddress')
              ? device.info.ipAddress
              : ''}
          </Value>
        </Group>

        <Group>
          <Label>Operating System</Label>
          <Value>
            {device.info.hasOwnProperty('osRelease') &&
            device.info.osRelease.hasOwnProperty('prettyName')
              ? device.info.osRelease.prettyName
              : '-'}
          </Value>
        </Group>
      </Card>

      <EditableLabelTable
        data={device.labels}
        onAdd={label =>
          api.addDeviceLabel({
            projectId: params.project,
            deviceId: device.id,
            data: label,
          })
        }
        onRemove={labelId =>
          api.removeDeviceLabel({
            projectId: params.project,
            deviceId: device.id,
            labelId,
          })
        }
        marginBottom={5}
      />
    </>
  );
};

export default DeviceOverview;
