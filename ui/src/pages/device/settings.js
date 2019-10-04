import React, { Component, Fragment } from 'react';
import axios from 'axios';
import {
  Pane,
  Table,
  Dialog,
  majorScale,
  Button,
  Heading,
  Alert,
  toaster,
  Badge,
  Text,
  TextInputField
} from 'evergreen-ui';

import utils from '../../utils';
import config from '../../config';
import InnerCard from '../../components/InnerCard';

export default class DeviceSettings extends Component {
  state = {
    name: this.props.device.name,
    nameValidationMessage: null,
    unchanged: true,
    showRemoveDialog: false,
    backendError: null
  };

  handleUpdateName = event => {
    this.setState({
      name: event.target.value,
      unchanged: false
    });
  };

  handleUpdate = () => {
    var nameValidationMessage = utils.checkName('device', this.state.name);

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return;
    }

    axios
      .patch(
        `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}`,
        {
          name: this.state.name
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Device updated successfully.');
        this.props.history.push(
          `/${this.props.projectName}/devices/${this.state.name}`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Device was not updated.');
          console.log(error);
        }
      });
  };

  handleRemove = () => {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Successfully deleted device.');
        this.props.history.push(`/${this.props.projectName}/devices`);
      })
      .catch(error => {
        this.setState({
          showRemoveDialog: false
        });
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Device was not removed.');
          console.log(error);
        }
      });
  };

  render() {
    const device = this.props.device;
    return (
      <Pane width="50%">
        <InnerCard>
          <Pane display="flex" flexDirection="column" padding={majorScale(4)}>
            {this.state.backendError && (
              <Alert
                marginBottom={majorScale(2)}
                paddingTop={majorScale(2)}
                paddingBottom={majorScale(2)}
                intent="warning"
                title={this.state.backendError}
              />
            )}
            <Table paddingBottom={majorScale(2)}>
              <Table.Row>
                <Table.Cell paddingLeft="0">
                  <Heading size={600}>Device Settings</Heading>
                </Table.Cell>
                <Table.Cell flexBasis={90} flexShrink={0} flexGrow={0}>
                  {device.status === 'offline' ? (
                    <Badge color="red">offline</Badge>
                  ) : (
                    <Badge color="green">online</Badge>
                  )}
                </Table.Cell>
              </Table.Row>
            </Table>
            <Text>
              <strong>ID: </strong>
              {device.id}
            </Text>
            <TextInputField
              label="Name"
              onChange={this.handleUpdateName}
              value={this.state.name}
              validationMessage={this.state.nameValidationMessage}
              paddingTop={majorScale(2)}
            />
            <Pane marginTop={majorScale(2)}>
              <Button
                appearance="primary"
                disabled={this.state.unchanged}
                onClick={this.handleUpdate}
              >
                Update Settings
              </Button>
            </Pane>
          </Pane>
          <Pane
            borderTop="default"
            marginRight={majorScale(4)}
            marginLeft={majorScale(4)}
            marginBottom={majorScale(4)}
          >
            <Button
              marginTop={majorScale(4)}
              iconBefore="trash"
              intent="danger"
              onClick={() => this.setState({ showRemoveDialog: true })}
            >
              Remove Device...
            </Button>
          </Pane>
          <Pane>
            <Dialog
              isShown={this.state.showRemoveDialog}
              title="Remove Device"
              intent="danger"
              onCloseComplete={() => this.setState({ showRemoveDialog: false })}
              onConfirm={() => this.handleRemove()}
              confirmLabel="Remove Device"
            >
              You are about to remove the <strong>{device.name}</strong> device.
            </Dialog>
          </Pane>
        </InnerCard>
      </Pane>
    );
  }
}
