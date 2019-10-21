import React, { Component, Fragment } from 'react';
import axios from 'axios';
import {
  Pane,
  Table,
  Dialog,
  majorScale,
  Button,
  Heading,
  Badge,
  IconButton,
  TextInputField,
  toaster
} from 'evergreen-ui';

import utils from '../../utils';
import config from '../../config';
import InnerCard from '../../components/InnerCard';
import CustomSpinner from '../../components/CustomSpinner';

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
              <DeviceLabels projectName={projectName} device={device} />
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

class DeviceLabels extends Component {
  state = {
    labels: []
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}/labels`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          labels: this.initializeLabels(response.data)
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  initializeLabels = keyValues => {
    var labels = [];
    var labelKeys = Object.keys(keyValues);
    for (var i = 0; i < labelKeys.length; i++) {
      labels.push({
        key: labelKeys[i],
        value: keyValues[labelKeys[i]],
        mode: "default",
        keyValidationMessage: null,
        valueValidationMessage: null,
        showRemoveDialog: false
      });
    }
    return labels.sort(function(a, b) {
      if (a.key < b.key) {
        return -1;
      }
      if (a.key > b.key) {
        return 1;
      }
      return 0;
    });
  };

  createNewDeviceLabel = k => {
    var labels = this.state.labels;
    labels.push({
      key: '',
      value: '',
      mode: 'new',
      keyValidationMessage: null,
      valueValidationMessage: null,
      showRemoveDialog: false
    });
    this.setState({
      labels: labels
    });
  };

  handleUpdate = (i, property) => {
    return event => {
      var labels = this.state.labels;
      labels[i][property] = event.target.value;
      this.setState({
        labels: labels
      });
    };
  };

  setEdit = i => {
    var editLabels = this.state.labels;
    editLabels[i]['mode'] = 'edit';
    this.setState({
      labels: editLabels
    });
  };

  cancelEdit = i => {
    var editLabels = this.state.labels;
    editLabels[i]['mode'] = 'default';
    this.setState({
      labels: editLabels
    });
  };

  setShowRemoveDialog = i => {
    var showRemoveDialogLabels = this.state.labels;
    showRemoveDialogLabels[i]['showRemoveDialog'] = true;
    this.setState({
      labels: showRemoveDialogLabels
    });
  };

  hideShowRemoveDialog = i => {
    var showRemoveDialogLabels = this.state.labels;
    showRemoveDialogLabels[i]['showRemoveDialog'] = false;
    this.setState({
      labels: showRemoveDialogLabels
    });
  };

  setDeviceLabel = (key, value, i) => {
    var updatedLabels = this.state.labels;
    var keyValidationMessage = utils.checkName('key', key);
    var valueValidationMessage = utils.checkName('value', value);

    if (keyValidationMessage === null) {
      for (var j = 0; j < updatedLabels.length; j++) {
        if (i !== j && key === updatedLabels[j]['key']) {
          keyValidationMessage = 'Key already exists.';
          break;
        }
      }
    }

    updatedLabels[i]['keyValidationMessage'] = keyValidationMessage;
    updatedLabels[i]['valueValidationMessage'] = valueValidationMessage;

    this.setState({
      labels: updatedLabels
    });

    if (
      keyValidationMessage === null &&
      valueValidationMessage === null &&
      key !== null &&
      value !== null
    ) {
      axios
        .post(
          `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}/labels`,
          {
            key: key,
            value: value
          },
          {
            withCredentials: true
          }
        )
        .then(response => {
          var updatedLabels = this.state.labels;
          updatedLabels[i]['mode'] = 'default';
          this.setState({
            labels: updatedLabels
          });
        })
        .catch(error => {
          console.log(error);
        });
    }
  };

  deleteDeviceLabel = (key, i) => {
    if (key !== '') {
      axios
        .delete(
          `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}/labels/${key}`,
          {
            withCredentials: true
          }
        )
        .then(response => {
          var removedLabels = this.state.labels;
          removedLabels.splice(i, 1);
          this.setState({
            labels: removedLabels
          });
        })
        .catch(error => {
          var hideRemoveDialogLabels = this.state.labels;
          hideRemoveDialogLabels[i]['showRemoveDialog'] = false;
          this.setState({
            labels: hideRemoveDialogLabels
          });
          toaster.danger('Device label was not removed.');
          console.log(error);
        });
    } else {
      var removedLabels = this.state.labels;
      removedLabels.splice(i, 1);
      this.setState({
        labels: removedLabels
      });
    }
  };

  renderDeviceLabel(deviceLabel, i) {
    switch (deviceLabel.mode) {
      case 'default':
        return (
          <Fragment key={deviceLabel.key}>
            <Table.Row>
              <Table.TextCell>{deviceLabel.key}</Table.TextCell>
              <Table.TextCell>{deviceLabel.value}</Table.TextCell>
              <Table.TextCell flexBasis={75} flexShrink={0} flexGrow={0}>
                <Pane display="flex">
                  <IconButton
                    icon="edit"
                    height={24}
                    appearance="minimal"
                    onClick={() => this.setEdit(i)}
                  />
                  <IconButton
                    icon="trash"
                    height={24}
                    appearance="minimal"
                    onClick={() => this.setShowRemoveDialog(i)}
                  />
                </Pane>
              </Table.TextCell>
            </Table.Row>
            <Pane>
              <Dialog
                isShown={deviceLabel.showRemoveDialog}
                title="Remove Label"
                intent="danger"
                onCloseComplete={() => this.hideShowRemoveDialog(i)}
                onConfirm={() => this.deleteDeviceLabel(deviceLabel.key, i)}
                confirmLabel="Remove Label"
              >
                You are about to remove label <strong>{deviceLabel.key}</strong>
                .
              </Dialog>
            </Pane>
          </Fragment>
        );
      case 'edit':
        return (
          <Table.Row key={deviceLabel.key} height="auto">
            <Table.TextCell>{deviceLabel.key}</Table.TextCell>
            <Table.TextCell>
              <TextInputField
                label=""
                name={`edit-${deviceLabel.key}`}
                value={deviceLabel.value}
                onChange={event => this.handleUpdate(i, 'value')(event)}
                isInvalid={deviceLabel.valueValidationMessage !== null}
                validationMessage={deviceLabel.valueValidationMessage}
                marginTop={8}
                marginBottom={8}
              />
            </Table.TextCell>
            <Table.TextCell flexBasis={75} flexShrink={0} flexGrow={0}>
              <Pane display="flex">
                <IconButton
                  icon="floppy-disk"
                  height={24}
                  appearance="minimal"
                  onClick={() =>
                    this.setDeviceLabel(deviceLabel.key, deviceLabel.value, i)
                  }
                />
                <IconButton
                  icon="cross"
                  height={24}
                  appearance="minimal"
                  onClick={() => this.cancelEdit(i)}
                />
              </Pane>
            </Table.TextCell>
          </Table.Row>
        );
      case 'new':
        return (
          <Table.Row key={`new-${i}`} height="auto">
            <Table.TextCell>
              <TextInputField
                label=""
                name={`new-key-${i}`}
                value={deviceLabel.key}
                onChange={event => this.handleUpdate(i, 'key')(event)}
                isInvalid={deviceLabel.keyValidationMessage !== null}
                validationMessage={deviceLabel.keyValidationMessage}
                marginTop={8}
                marginBottom={8}
              />
            </Table.TextCell>
            <Table.TextCell>
              <TextInputField
                label=""
                name={`new-value-${i}`}
                value={deviceLabel.value}
                onChange={event => this.handleUpdate(i, 'value')(event)}
                isInvalid={deviceLabel.valueValidationMessage !== null}
                validationMessage={deviceLabel.valueValidationMessage}
                marginTop={8}
                marginBottom={8}
              />
            </Table.TextCell>
            <Table.TextCell flexBasis={75} flexShrink={0} flexGrow={0}>
              <Pane display="flex">
                <IconButton
                  icon="floppy-disk"
                  height={24}
                  appearance="minimal"
                  onClick={() =>
                    this.setDeviceLabel(deviceLabel.key, deviceLabel.value, i)
                  }
                />
                <IconButton
                  icon="cross"
                  height={24}
                  appearance="minimal"
                  onClick={() => this.deleteDeviceLabel(deviceLabel.key, i)}
                />
              </Pane>
            </Table.TextCell>
          </Table.Row>
        );
      default:
        return <Fragment />;
    }
  }

  render() {
    return (
      <Table>
        <Table.Head>
          <Table.TextHeaderCell>Key</Table.TextHeaderCell>
          <Table.TextHeaderCell>Value</Table.TextHeaderCell>
          <Table.TextHeaderCell
            flexBasis={75}
            flexShrink={0}
            flexGrow={0}
          ></Table.TextHeaderCell>
        </Table.Head>
        <Table.Body>
          {this.state.labels.map((deviceLabel, i) =>
            this.renderDeviceLabel(deviceLabel, i)
          )}
          <Table.Row key="add">
            <Table.TextCell>
              <IconButton
                icon="plus"
                height={24}
                appearance="minimal"
                onClick={() => this.createNewDeviceLabel()}
              />
            </Table.TextCell>
            <Table.TextCell></Table.TextCell>
            <Table.TextCell
              flexBasis={75}
              flexShrink={0}
              flexGrow={0}
            ></Table.TextCell>
          </Table.Row>
        </Table.Body>
      </Table>
    );
  }
}
