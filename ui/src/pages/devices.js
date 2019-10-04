import React, { Fragment, Component } from 'react';
import axios from 'axios';
import { Button, Pane, Table, Heading, Badge, majorScale } from 'evergreen-ui';

import config from '../config.js';
import InnerCard from '../components/InnerCard.js';
import TopHeader from '../components/TopHeader.js';

export default class Devices extends Component {
  state = {
    devices: []
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/devices?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          devices: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

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
    return (
      <Fragment>
        <TopHeader
          user={this.props.user}
          heading="Devices"
          history={this.props.history}
        />
        <Pane width="70%">
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading paddingLeft={majorScale(2)}>Devices</Heading>
              <Button
                margin={majorScale(2)}
                appearance="primary"
                onClick={() =>
                  this.props.history.push(
                    `/${this.props.projectName}/devices/add`
                  )
                }
              >
                Add Device
              </Button>
            </Pane>
            {this.state.devices && this.state.devices.length > 0 && (
              <Table>
                <Table.Head background="tint2">
                  <Table.TextHeaderCell
                    flexBasis={90}
                    flexShrink={0}
                    flexGrow={0}
                  >
                    Status
                  </Table.TextHeaderCell>
                  <Table.TextHeaderCell>Name</Table.TextHeaderCell>
                  <Table.TextHeaderCell>IP Address</Table.TextHeaderCell>
                  <Table.TextHeaderCell>OS</Table.TextHeaderCell>
                </Table.Head>
                <Table.Body>
                  {this.state.devices.map(device => (
                    <Table.Row
                      key={device.id}
                      isSelectable
                      onSelect={() =>
                        this.props.history.push(
                          `/${this.props.projectName}/devices/${device.name}`
                        )
                      }
                    >
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
                      <Table.TextCell>{device.name}</Table.TextCell>
                      <Table.TextCell>
                        {device.info.hasOwnProperty('ipAddress')
                          ? device.info.ipAddress
                          : ''}
                      </Table.TextCell>
                      <Table.TextCell>
                        {this.renderDeviceOs(device)}
                      </Table.TextCell>
                    </Table.Row>
                  ))}
                </Table.Body>
              </Table>
            )}
          </InnerCard>
        </Pane>
      </Fragment>
    );
  }
}
