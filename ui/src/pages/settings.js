import React, { Component, Fragment } from 'react';
import {
  Button,
  Pane,
  Heading,
  majorScale,
  TextInput,
  TextInputField,
  Alert,
  Dialog,
  toaster
} from 'evergreen-ui';
import axios from 'axios';

import config from '../config.js';
import utils from '../utils.js';
import InnerCard from '../components/InnerCard.js';
import TopHeader from '../components/TopHeader.js';

export default class ProjectSettings extends Component {
  state = {
    name: this.props.projectName,
    nameValidationMessage: null,
    unchanged: true,
    showDeleteDialog: false,
    disableDeleteConfirm: true,
    backendError: null
  };

  handleUpdateName = event => {
    this.setState({
      name: event.target.value,
      unchanged: false
    });
  };

  handleUpdate = () => {
    var nameValidationMessage = utils.checkName('project', this.state.name);

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
        `${config.endpoint}/projects/${this.props.projectName}`,
        {
          name: this.state.name
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Project updated successfully.');
        this.props.history.push(`/${this.state.name}/settings`);
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Project was not updated.');
          console.log(error);
        }
      });
  };

  handleDelete() {
    this.setState({
      backendError: null
    });

    axios
      .delete(`${config.endpoint}/projects/${this.props.projectName}`, {
        withCredentials: true
      })
      .then(response => {
        this.setState({
          showDeleteDialog: false
        });
        toaster.success('Successfully deleted project.');
        this.props.history.push(`/`);
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
          toaster.danger('Project was not deleted.');
          console.log(error);
        }
      });
  }

  handleDeleteName = event => {
    this.setState({
      disableDeleteConfirm: event.target.value !== this.props.projectName
    });
  };

  render() {
    const { projectName, history, user } = this.props;
    const {
      backendError,
      name,
      nameValidationMessage,
      unchanged,
      disableDeleteConfirm,
      showDeleteDialog
    } = this.state;

    return (
      <Fragment>
        <TopHeader
          user={user}
          heading={`Project / ${projectName}`}
          history={history}
        />
        <Pane width="50%">
          <InnerCard>
            <Pane padding={majorScale(4)}>
              {backendError && (
                <Alert
                  marginBottom={majorScale(2)}
                  paddingTop={majorScale(2)}
                  paddingBottom={majorScale(2)}
                  intent="warning"
                  title={backendError}
                />
              )}
              <Heading paddingBottom={majorScale(4)} size={600}>
                Project Settings
              </Heading>
              <TextInputField
                label="Project Name"
                onChange={this.handleUpdateName}
                value={name}
                validationMessage={nameValidationMessage}
              />
              <Button
                marginTop={majorScale(2)}
                appearance="primary"
                disabled={unchanged}
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
                Delete Project...
              </Button>
            </Pane>
            <Pane>
              <Dialog
                isShown={showDeleteDialog}
                title="Delete Project"
                intent="danger"
                onCloseComplete={() =>
                  this.setState({ showDeleteDialog: false })
                }
                onConfirm={() => this.handleDelete()}
                confirmLabel="Delete Project"
                isConfirmDisabled={disableDeleteConfirm}
              >
                This action <strong>cannot</strong> be undone. This will
                permanently delete the <strong>{projectName}</strong> project.
                <p></p>Please type in the name of the project to confirm.
                <TextInput
                  marginTop={majorScale(1)}
                  onChange={this.handleDeleteName}
                />
              </Dialog>
            </Pane>
          </InnerCard>
        </Pane>
      </Fragment>
    );
  }
}
