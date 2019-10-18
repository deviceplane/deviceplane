import React, { Component } from 'react';
import {
  Button,
  Pane,
  Heading,
  majorScale,
  Alert,
  toaster,
  Label,
  Textarea,
  TextInputField,
  Dialog
} from 'evergreen-ui';
import axios from 'axios';

import config from '../../config.js';
import utils from '../../utils.js';
import InnerCard from '../../components/InnerCard.js';

export default class DeviceRegistrationTokenSettings extends Component {
  state = {
    name: this.props.deviceRegistrationToken.name,
    description: this.props.deviceRegistrationToken.description,
    maxRegistrations: (typeof this.props.deviceRegistrationToken.maxRegistrations === 'number' ?
      String(this.props.deviceRegistrationToken.maxRegistrations) :
      ''),
    nameValidationMessage: null,
    maxRegistrationsValidationMessage: null,
    unchanged: true,
    showDeleteDialog: false,
    backendError: null
  };

  handleUpdateName = event => {
    this.setState({
      name: event.target.value,
      unchanged: false
    });
  };

  handleUpdateDescription = event => {
    this.setState({
      description: event.target.value,
      unchanged: false
    });
  };

  handleUpdateMaxRegistrations = event => {
    this.setState({
      maxRegistrations: event.target.value,
      unchanged: false
    });
  };

  handleUpdate = () => {
    var nameValidationMessage = utils.checkName('Device Registration Token', this.state.name);

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return;
    }

    // Convert max registrations to int or undefined
    var maxRegistrationsCleaned;
    if (this.state.maxRegistrations === '') {
      maxRegistrationsCleaned = null;
    } else {
      maxRegistrationsCleaned = Number(this.state.maxRegistrations);

      if (isNaN(maxRegistrationsCleaned)) {
        this.setState({
          maxRegistrationsValidationMessage: 'Max Registrations should either be a number or be left empty.'
        });
        return;
      }
    }

    axios
      .put(
        `${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens/${this.props.deviceRegistrationToken.name}`,
        {
          name: this.state.name,
          description: this.state.description,
          maxRegistrations: maxRegistrationsCleaned,
          settings: this.props.deviceRegistrationToken.settings,
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Device Registration Token updated successfully.');
        this.props.history.push(
          // TODO: change this to an overview page once we get one
          `/${this.props.projectName}/provisioning`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Device Registration Token was not updated.');
          console.log(error);
        }
      });
  };

  handleDelete() {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens/${this.props.deviceRegistrationToken.name}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          showDeleteDialog: false
        });
        toaster.success('Successfully deleted Device Registration Token.');
        this.props.history.push(
          `/${this.props.projectName}/provisioning/deviceregistrationtokens`
        );
      })
      .catch(error => {
        this.setState({
          showDeleteDialog: false
        });
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Device Registration Token was not deleted.');
          console.log(error);
        }
      });
  }

  render() {
    return (
      <Pane width="50%">
        <InnerCard>
          <Pane padding={majorScale(4)}>
            {this.state.backendError && (
              <Alert
                marginBottom={majorScale(2)}
                paddingTop={majorScale(2)}
                paddingBottom={majorScale(2)}
                intent="warning"
                title={this.state.backendError}
              />
            )}
            <Heading paddingBottom={majorScale(4)} size={600}>
              Device Registration Token Settings
            </Heading>
            <TextInputField
              label="Name"
              onChange={this.handleUpdateName}
              value={this.state.name}
              validationMessage={this.state.nameValidationMessage}
            />
            <Label
              htmlFor="description-textarea"
              marginBottom="4"
              display="block"
            >
              Description (optional)
            </Label>
            <Textarea
              id="description-textarea"
              height="100px"
              onChange={this.handleUpdateDescription}
              value={this.state.description}
            />
            <TextInputField
              paddingTop={majorScale(4)}
              label="Maximum Device Registrations"
              description="Limit the number of devices that can be registered using this token"
              hint="Leave empty to allow unlimited registrations"
              onChange={this.handleUpdateMaxRegistrations}
              value={this.state.maxRegistrations}
              validationMessage={this.state.maxRegistrationsValidationMessage}
            />
            <Button
              marginTop={majorScale(2)}
              appearance="primary"
              disabled={this.state.unchanged}
              onClick={this.handleUpdate}
            >
              Update Settings
            </Button>
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
              onClick={() => this.setState({ showDeleteDialog: true })}
            >
              Delete Device Registration Token...
            </Button>
          </Pane>
          <Pane>
            <Dialog
              isShown={this.state.showDeleteDialog}
              title="Delete Device Registration Token"
              intent="danger"
              onCloseComplete={() => this.setState({ showDeleteDialog: false })}
              onConfirm={() => this.handleDelete()}
              confirmLabel="Delete Device Registration Token"
            >
              You are about to delete the{" "}
              <strong>{this.props.deviceRegistrationToken.name}</strong> Device
              Registration Token.
            </Dialog>
          </Pane>
        </InnerCard>
      </Pane>
    );
  }
}
