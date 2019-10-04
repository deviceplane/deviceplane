import React, { Component, Fragment } from 'react';
import axios from 'axios';
import {
  Pane,
  Dialog,
  majorScale,
  Button,
  Heading,
  Alert,
  toaster,
  TextInputField,
  Label,
  Textarea
} from 'evergreen-ui';

import utils from '../../utils';
import config from '../../config';
import CustomSpinner from '../../components/CustomSpinner';
import Editor from '../../components/Editor';
import InnerCard from '../../components/InnerCard';

export default class Role extends Component {
  state = {
    role: null,
    name: '',
    nameValidationMessage: null,
    description: '',
    config: '',
    unchanged: true,
    showDeleteDialog: false,
    backendError: null
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/roles/${this.props.roleName}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          role: response.data,
          name: response.data.name,
          description: response.data.description,
          config: response.data.config
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

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

  handleUpdate() {
    var nameValidationMessage = utils.checkName('role', this.state.name);

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
        `${config.endpoint}/projects/${this.props.projectName}/roles/${this.state.role.id}`,
        {
          name: this.state.name,
          description: this.state.description,
          config: this.state.config
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Successfully updated role.');
        this.props.history.push(`/${this.props.projectName}/iam/roles`);
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Role was not updated.');
          console.log(error);
        }
      });
  }

  handleDelete() {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/roles/${this.state.role.id}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          showDeleteDialog: false
        });
        toaster.success('Successfully deleted role.');
        this.props.history.push(`/${this.props.projectName}/iam/roles`);
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
          toaster.danger('Role was not deleted.');
          console.log(error);
        }
      });
  }

  render() {
    const role = this.state.role;
    return (
      <Pane width="70%">
        {role ? (
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
                Role / {role.name}
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
              <Heading paddingTop={majorScale(2)} marginBottom={majorScale(2)}>
                Config
              </Heading>
              <Editor
                width="100%"
                height="300px"
                value={this.state.config}
                onChange={value =>
                  this.setState({ config: value, unchanged: false })
                }
              />
              <Button
                marginTop={majorScale(2)}
                appearance="primary"
                disabled={this.state.unchanged}
                onClick={() => this.handleUpdate()}
              >
                Update Role
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
                Delete Role...
              </Button>
            </Pane>
            <Pane>
              <Dialog
                isShown={this.state.showDeleteDialog}
                title="Delete Role"
                intent="danger"
                onCloseComplete={() =>
                  this.setState({ showDeleteDialog: false })
                }
                onConfirm={() => this.handleDelete()}
                confirmLabel="Delete Role"
              >
                You are about to delete the{' '}
                <strong>{this.props.roleName}</strong> role.
              </Dialog>
            </Pane>
          </InnerCard>
        ) : (
          <CustomSpinner />
        )}
      </Pane>
    );
  }
}
