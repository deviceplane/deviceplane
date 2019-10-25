import React, { Component, Fragment } from 'react';
import {
  Pane,
  Table,
  majorScale,
  Button,
  Heading,
  Badge,
} from 'evergreen-ui';

import config from '../../config';
import InnerCard from '../../components/InnerCard';
import CustomSpinner from '../../components/CustomSpinner';
import { EditableLabelTable } from "../../components/EditableLabelTable";

export default class DeviceOverview extends Component {
  renderDeviceOs(device) {
    var innerText = '-';
    if (
      device.info.hasOwnProperty('osRelease') &&
      device.info.osRelease.hasOwnProperty('prettyName')
    ) {
      innerText = device.info.osRelease.prettyName;
    }
    return <Pane>{innerText}</Pane>;
  }

  render() {
    const { device, projectName, history } = this.props;

    return (
      <Pane width="70%" display="flex" flexDirection="column">
        {device ? (
          <Fragment>
            <InnerCard>
              <Heading padding={majorScale(2)}>Device Info</Heading>
              <Table>
                <Table.Head>
                  <Table.TextHeaderCell
                    flexBasis={90}
                    flexShrink={0}
                    flexGrow={0}
                  >
                    Status
                  </Table.TextHeaderCell>
                  <Table.TextHeaderCell>IP Address</Table.TextHeaderCell>
                  <Table.TextHeaderCell>OS</Table.TextHeaderCell>
                </Table.Head>
                <Table.Body>
                  <Table.Row>
                    <Table.TextCell
                      flexBasis={90}
                      flexShrink={0}
                      flexGrow={0}
                      alignItems="center"
                      paddingRight="0"
                    >
                      {device.status === 'offline' ? (
                        <Badge color="red">offline</Badge>
                      ) : (
                        <Badge color="green">online</Badge>
                      )}
                    </Table.TextCell>
                    <Table.TextCell>
                      {device.info.hasOwnProperty('ipAddress')
                        ? device.info.ipAddress
                        : ''}
                    </Table.TextCell>
                    <Table.TextCell>
                      {this.renderDeviceOs(device)}
                    </Table.TextCell>
                  </Table.Row>
                </Table.Body>
              </Table>
            </InnerCard>
            <InnerCard>
              <Heading padding={majorScale(2)}>Labels</Heading>
              <EditableLabelTable
                getEndpoint={`${config.endpoint}/projects/${projectName}/devices/${device.id}`}
                setEndpoint={`${config.endpoint}/projects/${projectName}/devices/${device.id}/labels`}
                deleteEndpoint={`${config.endpoint}/projects/${projectName}/devices/${device.id}/labels`}
              />
            </InnerCard>
            <InnerCard>
              <Heading padding={majorScale(2)}>Services</Heading>
              <DeviceServices
                projectName={projectName}
                device={device}
                history={history}
              />
            </InnerCard>
          </Fragment>
        ) : (
          <CustomSpinner />
        )}
      </Pane>
    );
  }
}

class DeviceServices extends Component {
  checkServices(applicationStatusInfo) {
    for (var i = 0; i < applicationStatusInfo.length; i++) {
      if (
        applicationStatusInfo[i].serviceStatuses &&
        applicationStatusInfo[i].serviceStatuses.length > 0
      ) {
        return true;
      }
    }
    return false;
  }

  render() {
    const applicationStatusInfo = this.props.device.applicationStatusInfo;
    return (
      <Pane>
        {this.checkServices(applicationStatusInfo) && (
          <Table>
            <Table.Head>
              <Table.TextHeaderCell>Service</Table.TextHeaderCell>
              <Table.TextHeaderCell>Current Release</Table.TextHeaderCell>
            </Table.Head>
            {applicationStatusInfo.map(applicationInfo => (
              <Pane key={applicationInfo.application.id}>
                <Table.Body>
                  {applicationInfo.serviceStatuses.map(
                    (serviceStatus, index) => (
                      <Table.Row key={index}>
                        <Table.TextCell>
                          <Button
                            color="#425A70"
                            appearance="minimal"
                            onClick={() =>
                              this.props.history.push(
                                `/${this.props.projectName}/applications/${applicationInfo.application.name}`
                              )
                            }
                          >
                            {applicationInfo.application.name} /{' '}
                            {serviceStatus.service}
                          </Button>
                        </Table.TextCell>
                        <Table.TextCell>
                          {serviceStatus.currentReleaseId}
                        </Table.TextCell>
                      </Table.Row>
                    )
                  )}
                </Table.Body>
              </Pane>
            ))}
          </Table>
        )}
      </Pane>
    );
  }
}