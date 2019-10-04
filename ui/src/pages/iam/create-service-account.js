import React, { Component } from 'react';
import axios from 'axios';
import {
  Pane,
  majorScale,
  Button,
  Heading,
  Alert,
  toaster,
  TextInputField,
  Textarea,
  Label
} from 'evergreen-ui';

import utils from '../../utils';
import config from '../../config';
import segment from '../../segment';
import InnerCard from '../../components/InnerCard';

export default class CreateServiceAccount extends Component {
  state = {
    name: '',
    nameValidationMessage: null,
    description: '',
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

  handleSubmit = event => {
    event.preventDefault();

    var nameValidationMessage = utils.checkName(
      'service account',
      this.state.name
    );

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return;
    }

    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts`,
        {
          name: this.state.name,
          description: this.state.description
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        segment.track('Service Account Created');
        this.props.history.push(
          `/${this.props.projectName}/iam/serviceaccounts/`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Service Account was not created.');
          console.log(error);
        }
      });
  };

  handleCancel() {
    this.props.history.push(`/${this.props.projectName}/iam/serviceaccounts/`);
  }

  render() {
    return (
      <Pane width="50%">
        <InnerCard>
          <Pane padding={majorScale(4)} is="form">
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
              Create Service Account
            </Heading>
            <TextInputField
              label="Name"
              onChange={this.handleUpdateName}
              isInvalid={this.state.nameValidationMessage !== null}
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
            <Pane display="flex" flex="row">
              <Button
                marginTop={majorScale(2)}
                appearance="primary"
                onClick={this.handleSubmit}
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
