import React, { Component, Fragment } from 'react';
import axios from 'axios';
import {
  Pane,
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
import segment from '../../segment';
import InnerCard from '../../components/InnerCard';
import Editor from '../../components/Editor';

export default class CreateRole extends Component {
  state = {
    name: '',
    nameValidationMessage: null,
    description: '',
    config: '',
    backendError: null
  };

  handleUpdateName = event => {
    this.setState({
      name: event.target.value
    });
  };

  handleUpdateDescription = event => {
    this.setState({
      description: event.target.value
    });
  };

  handleSubmit() {
    var nameValidationMessage = utils.checkName('role', this.state.name);
    var configError = utils.checkConfig('role config', this.state.config);

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: configError
    });

    if (nameValidationMessage !== null || configError !== null) {
      return;
    }

    if (this.state.config === '') {
      this.setState({
        backendError: 'Please define the role config.'
      });
      return;
    }

    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/roles`,
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
        segment.track('Role Created');
        this.props.history.push(`/${this.props.projectName}/iam/roles/`);
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Role was not created.');
          console.log(error);
        }
      });
  }

  handleCancel() {
    this.props.history.push(`/${this.props.projectName}/iam/roles/`);
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
              Create Role
            </Heading>
            <TextInputField
              label="Name"
              onChange={this.handleUpdateName}
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
            />
            <Heading paddingTop={majorScale(2)} marginBottom={majorScale(2)}>
              Config
            </Heading>
            <Editor
              width="100%"
              height="300px"
              value={this.state.config}
              onChange={value => this.setState({ config: value })}
            />
            <Pane display="flex" flex="row">
              <Button
                marginTop={majorScale(2)}
                appearance="primary"
                onClick={() => this.handleSubmit()}
              >
                Submit
              </Button>
              <Button
                marginTop={majorScale(2)}
                marginLeft={majorScale(2)}
                onClick={() => this.handleCancel()}
              >
                Cancel
              </Button>
            </Pane>
          </Pane>
        </InnerCard>
      </Pane>
    );
  }
}
