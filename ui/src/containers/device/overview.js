import React, { useMemo, useState } from 'react';
import { Icon, toaster } from 'evergreen-ui';

import api from '../../api';
import theme from '../../theme';
import {
  Group,
  Column,
  Button,
  Value,
  Link,
  Label,
  Text,
} from '../../components/core';
import Card from '../../components/card';
import Table from '../../components/table';
import Popup from '../../components/popup';
import Editor from '../../components/editor';
import EditableLabelTable from '../../components/editable-label-table';
import DeviceStatus from '../../components/device-status';

const DeviceServices = ({ projectId, device, applicationStatusInfo }) => {
  const [serviceMetrics, setServiceMetrics] = useState({});
  const columns = useMemo(
    () => [
      {
        Header: 'Service',
        accessor: ({ application, service }) =>
          `${application.name} / ${service}`,
        Cell: ({ cell: { value }, row: { original } }) => (
          <Link
            href={`/${projectId}/applications/${original.application.name}`}
          >
            {value}
          </Link>
        ),
      },
      {
        Header: 'Current Release',
        accessor: 'currentReleaseId',
        Cell: ({ cell: { value }, row: { original } }) => (
          <Link
            href={`/${projectId}/applications/${original.application.name}/releases/${original.currentReleaseId}`}
          >
            {value}
          </Link>
        ),
      },
      {
        Header: ' ',
        Cell: ({ row: { original } }) => (
          <Button
            disabled={device.status === 'offline'}
            title={<Icon icon="pulse" size={18} color={theme.colors.white} />}
            variant="icon"
            onClick={async () => {
              try {
                const response = await api.serviceMetrics({
                  projectId,
                  deviceId: device.id,
                  applicationId: original.application.name,
                  serviceId: original.service,
                });
                setServiceMetrics({
                  service: original.service,
                  metrics: response.data,
                });
              } catch (error) {
                console.log(error);
                toaster.danger('Service Metrics are currently unavailable.');
              }
            }}
          />
        ),
        style: {
          flex: '0 0 50px',
          justifyContent: 'flex-end',
        },
      },
    ],
    []
  );
  const tableData = useMemo(
    () =>
      applicationStatusInfo.reduce(
        (data, curr) => [
          ...data,
          ...curr.serviceStatuses.map(status => ({
            ...status,
            application: curr.application,
          })),
        ],
        []
      ),
    [applicationStatusInfo]
  );

  return (
    <>
      <Table
        columns={columns}
        data={tableData}
        placeholder={
          <Text>
            There are no <strong>Services</strong>.
          </Text>
        }
      />
      <Popup
        show={!!serviceMetrics.service}
        onClose={() => setServiceMetrics({})}
      >
        <Card
          border
          title="Service Metrics"
          subtitle={serviceMetrics.service}
          size="xxlarge"
        >
          <Editor
            width="100%"
            height="70vh"
            value={serviceMetrics.metrics}
            readOnly
          />
        </Card>
      </Popup>
    </>
  );
};

const DeviceOverview = ({
  route: {
    data: { params, device },
  },
}) => {
  const [hostMetrics, setHostMetrics] = useState();
  return (
    <>
      <Card
        size="xlarge"
        title={device.name}
        subtitle={<DeviceStatus status={device.status} />}
        marginBottom={4}
        actions={[
          {
            title: <Icon icon="pulse" size={18} color={theme.colors.white} />,
            variant: 'icon',
            onClick: async () => {
              try {
                const response = await api.hostMetrics({
                  projectId: params.project,
                  deviceId: device.id,
                });
                setHostMetrics(response.data);
              } catch (error) {
                toaster.danger('Current device metrics are unavailable.');
                console.log(error);
              }
            },
            disabled: device.status === 'offline',
          },
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
                toaster.success('Reboot was initiated successfully.');
              } catch (error) {
                toaster.danger('Reboot was not successful.');
                console.log(error);
              }
            },
          },
        ]}
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
      <Column marginBottom={4}>
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
        />
      </Column>
      <Card title="Services" size="xlarge">
        <DeviceServices
          projectId={params.project}
          device={device}
          applicationStatusInfo={device.applicationStatusInfo}
        />
      </Card>
      <Popup show={!!hostMetrics} onClose={() => setHostMetrics(null)}>
        <Card
          border
          title="Current Device Metrics"
          subtitle={device.name}
          size="xxlarge"
        >
          <Editor width="100%" height="70vh" value={hostMetrics} readOnly />
        </Card>
      </Popup>
    </>
  );
};

export default DeviceOverview;
