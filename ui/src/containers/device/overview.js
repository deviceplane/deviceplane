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
import ServiceState, {
  ServiceStatePullingImage,
} from '../../components/service-state';
import { getMetricLabel } from '../../helpers/metrics';

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

    for (let i = 0; i < appStatusInfo.length; i++) {
      const info = appStatusInfo[i];
      if (info.serviceStates && info.serviceStates.length) {
        for (let j = 0; j < info.serviceStates.length; j++) {
          const s = info.serviceStates[j];

          const newService = {
            ...s,
            id: `${info.application.name} / ${s.service}`,
            currentRelease: {
              number:
                info.serviceStatuses && info.serviceStatuses.length
                  ? info.serviceStatuses[0].currentRelease.number
                  : null,
            },
            application: info.application,
          };

          setServices(services => {
            const existingService = services.find(
              ({ id }) => id === newService.id
            );
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
        }
      } else if (info.serviceStatuses && info.serviceStatuses.length) {
        info.serviceStatuses.forEach(s => {
          const newService = {
            ...s,
            id: `${info.application.name} / ${s.service}`,
            application: info.application,
          };

          setServices(services => {
            const existingService = services.find(
              ({ id }) => id === newService.id
            );
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
      }
    }
  };

  const updateServices = async () => {
    await getServices();
    for (let i = 0; i < services.length; i++) {
      const service = services[i];
      if (service.state === ServiceStatePullingImage) {
        const imagePullProgress = await getImagePullProgress({
          applicationId: service.application.id,
          serviceId: service.service,
        });

        setServices(
          services.map(s =>
            s.id === service.id ? { ...s, imagePullProgress } : s
          )
        );
      }
    }

    setTimeout(updateServices, 3000);
  };

  useEffect(updateServices, []);

  const [serviceMetrics, setServiceMetrics] = useState({});

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
    cols.push({
      Header: ' ',
      Cell: ({ row: { original } }) => (
        <Button
          disabled={device.status === 'offline'}
          title={<Icon icon="pulse" size={16} color="primary" />}
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
              toaster.danger('Service Metrics are currently unavailable.');
              console.error(error);
            }
          }}
        />
      ),
      maxWidth: '50px',
      minWidth: '50px',
      cellStyle: {
        justifyContent: 'flex-end',
      },
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
    <>
      <Table
        {...tableProps}
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
            maxLines={30}
            value={serviceMetrics.metrics}
            readOnly
          />
        </Card>
      </Popup>
    </>
  );
};

const parseMetrics = data =>
  JSON.stringify(
    parsePrometheusTextFormat(data).reduce(
      (obj, { name, help, metrics }) => ({
        ...obj,
        [getMetricLabel(name)]: {
          description: help,
          metrics,
        },
      }),
      {}
    ),
    null,
    '\t'
  );

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
        subtitle={
          <DeviceStatus
            inline
            status={device.status}
            lastSeenAt={device.lastSeenAt}
          />
        }
        actions={[
          {
            title: <Icon icon="pulse" size={18} color="primary" />,
            variant: 'icon',
            onClick: async () => {
              try {
                const { data } = await api.hostMetrics({
                  projectId: params.project,
                  deviceId: device.id,
                });
                setHostMetrics(parseMetrics(data));
              } catch (error) {
                toaster.danger('Current device metrics are unavailable.');
                console.error(error);
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

      <Card title="Application Services" size="xlarge" marginBottom={5}>
        <ApplicationServices
          projectId={params.project}
          device={device}
          applicationStatusInfo={device.applicationStatusInfo}
        />
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

      <Popup show={!!hostMetrics} onClose={() => setHostMetrics(null)}>
        <Card
          border
          title="Current Device Metrics"
          subtitle={device.name}
          size="xxlarge"
          overflow="scroll"
        >
          <Editor
            width="100%"
            value={hostMetrics}
            fontSize={12}
            mode="json"
            readOnly
            maxLines={30}
          />
        </Card>
      </Popup>
    </>
  );
};

export default DeviceOverview;
