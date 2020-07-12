import React, { useMemo, useState, useEffect, useCallback } from 'react';
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
import ServiceState, {
  ServiceStatePullingImage,
} from '../../components/service-state';

import storage from '../../storage';

const ApplicationServices = ({ projectId, device, applicationStatusInfo }) => {
  const [services, setServices] = useState([]);
  const [showProgress, setShowProgress] = useState({});

  const getImagePullProgress = async ({ applicationId, serviceId }) => {
    try {
      const { data } = await api.imagePullProgress({
        projectId: projectId,
        deviceId: device.name,
        applicationId,
        serviceId,
      });
      return data;
    } catch (e) {
      console.error(e);
    }
    return null;
  };

  const getServices = async () => {
    let appStatusInfo = [];
    try {
      const { data } = await api.device({
        projectId,
        deviceId: device.id,
      });
      appStatusInfo = data.applicationStatusInfo;
    } catch (e) {
      console.error(e);
      return [];
    }

    const newServices = [];

    for (let i = 0; i < appStatusInfo.length; i++) {
      const info = appStatusInfo[i];
      if (info.serviceStates && info.serviceStates.length) {
        for (let j = 0; j < info.serviceStates.length; j++) {
          const s = info.serviceStates[j];

          newServices.push({
            ...s,
            id: `${info.application.name} / ${s.service}`,
            currentRelease: {
              number:
                info.serviceStatuses && info.serviceStatuses.length
                  ? info.serviceStatuses[0].currentRelease.number
                  : null,
            },
            application: info.application,
          });
        }
      } else if (info.serviceStatuses && info.serviceStatuses.length) {
        info.serviceStatuses.forEach(s => {
          newServices.push({
            ...s,
            id: `${info.application.name} / ${s.service}`,
            application: info.application,
          });
        });
      }
    }
    return newServices;
  };

  const servicesPolling = async () => {
    const newServices = await getServices();
    newServices.forEach(newService => {
      setServices(services => {
        const existingService = services.find(({ id }) => id === newService.id);
        if (existingService) {
          return services.map(s =>
            s.id === newService.id
              ? {
                  ...s,
                  ...newService,
                }
              : s
          );
        }
        return [...services, newService];
      });
    });

    for (let i = 0; i < newServices.length; i++) {
      const newService = newServices[i];
      if (newService.state === ServiceStatePullingImage) {
        const imagePullProgress = await getImagePullProgress({
          applicationId: newService.application.id,
          serviceId: newService.service,
        });
        setServices(services =>
          services.map(s =>
            s.id === newService.id ? { ...s, imagePullProgress } : s
          )
        );
      }
    }

    setTimeout(servicesPolling, 5000);
  };

  useEffect(() => {
    servicesPolling();
  }, []);

  const columns = useMemo(() => {
    const cols = [];
    cols.push({
      Header: 'Service',
      accessor: 'id',
      Cell: ({ cell: { value }, row: { original } }) => (
        <Link href={`/${projectId}/applications/${original.application.name}`}>
          {value}
        </Link>
      ),
    });
    if (
      applicationStatusInfo.length &&
      applicationStatusInfo[0].serviceStates &&
      applicationStatusInfo[0].serviceStates.length
    ) {
      cols.push({
        Header: 'State',
        accessor: 'state',
        Cell: ({
          cell: { value },
          row: {
            original: { service, imagePullProgress, errorMessage },
          },
        }) => {
          let label = 'Pulling image';
          let layers = [];
          if (imagePullProgress) {
            layers = Object.values(imagePullProgress);
            const tag = layers.find(({ status }) =>
              status.includes('Pulling from')
            );
            if (tag) {
              label = tag.status;
              layers = layers.filter(
                ({ status }) => !status.includes('Pulling from')
              );
            }
          }
          return (
            <Column flex={1}>
              <ServiceState state={value} />
              {errorMessage && (
                <Text color="red" fontSize={0} marginTop={2}>
                  {errorMessage}
                </Text>
              )}
              {imagePullProgress && (
                <>
                  <Row
                    flex={1}
                    alignItems="center"
                    justifyContent="space-between"
                  >
                    <Text
                      marginTop={2}
                      fontWeight={2}
                      fontSize={0}
                      color="primary"
                    >
                      {label}
                    </Text>
                    <Button
                      title={
                        <Icon
                          icon={
                            showProgress[service.id]
                              ? 'caret-down'
                              : 'caret-right'
                          }
                          size={18}
                          color="primary"
                        />
                      }
                      onClick={() =>
                        setShowProgress(sp => ({
                          ...sp,
                          [service.id]: !sp[service.id],
                        }))
                      }
                      variant="icon"
                    />
                  </Row>
                  <Column height={showProgress[service.id] ? 'auto' : 0}>
                    {layers.map(({ id, status }) => (
                      <Text fontSize={0} marginTop={1}>
                        {id}: {status}
                      </Text>
                    ))}
                  </Column>
                </>
              )}
            </Column>
          );
        },
      });
    }
    cols.push({
      Header: 'Release',
      accessor: 'currentRelease.number',
      Cell: ({ cell: { value }, row: { original } }) =>
        value ? (
          <Link
            href={`/${projectId}/applications/${original.application.name}/releases/${value}`}
          >
            {value}
          </Link>
        ) : (
          '-'
        ),
      maxWidth: '100px',
      minWidth: '100px',
      cellStyle: { justifyContent: 'flex-end' },
    });
    return cols;
  }, [showProgress]);

  const tableData = useMemo(() => services, [services]);

  const tableProps = useTable(
    {
      columns,
      data: tableData,
    },
    useSortBy
  );

  return (
    <Table
      {...tableProps}
      placeholder={
        <Text>
          There are no <strong>Services</strong>.
        </Text>
      }
    />
  );
};

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

      {(storage.get('legacy') || false) && (
        <Card title="Application Services" size="xlarge" marginBottom={5}>
          <ApplicationServices
            projectId={params.project}
            device={device}
            applicationStatusInfo={device.applicationStatusInfo}
          />
        </Card>
      )}

      {(storage.get('legacy') || false) && (
        <EditableLabelTable
          title="Environment Variables"
          dataName="Environment Variable"
          data={device.environmentVariables}
          onAdd={environmentVariable =>
            api.addEnvironmentVariable({
              projectId: params.project,
              deviceId: device.id,
              data: environmentVariable,
            })
          }
          onRemove={key =>
            api.removeEnvironmentVariable({
              projectId: params.project,
              deviceId: device.id,
              key,
            })
          }
        />
      )}
    </>
  );
};

export default DeviceOverview;
