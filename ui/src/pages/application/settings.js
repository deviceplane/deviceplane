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

export default class ApplicationSettings extends Component {
  state = {
    name: this.props.application.name,
    nameValidationMessage: null,
    description: this.props.application.description,
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

  handleUpdate = () => {
    var nameValidationMessage = utils.checkName('application', this.state.name);

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return;
    }

    axios
      .put(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.application.name}`,
        {
          name: this.state.name,
          description: this.state.description,
          settings: this.props.application.settings
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Application updated successfully.');
        this.props.history.push(
          `/${this.props.projectName}/applications/${this.state.name}`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Application was not updated.');
          console.log(error);
        }
      });
  };

  handleDelete() {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.application.name}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          showDeleteDialog: false
        });
        toaster.success('Successfully deleted application.');
        this.props.history.push(`/${this.props.projectName}/applications`);
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
          toaster.danger('Application was not deleted.');
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
              Application Settings
            </Heading>
            <TextInputField
              label="Application Name"
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
              Delete Application...
            </Button>
          </Pane>
          <Pane>
            <Dialog
              isShown={this.state.showDeleteDialog}
              title="Delete Application"
              intent="danger"
              onCloseComplete={() => this.setState({ showDeleteDialog: false })}
              onConfirm={() => this.handleDelete()}
              confirmLabel="Delete Application"
            >
              You are about to delete the{' '}
              <strong>{this.props.application.name}</strong> application.
            </Dialog>
          </Pane>
        </InnerCard>
      </Pane>
    );
  }
}
