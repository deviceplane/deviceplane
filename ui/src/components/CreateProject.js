import React, { Component } from 'react';
import config from '../config.js';
import segment from '../segment.js';
import utils from '../utils.js';
import TopHeader from './TopHeader.js';
import InnerCard from './InnerCard.js';
import { toaster, Pane, majorScale, Alert, TextInputField, Button } from 'evergreen-ui'
import axios from 'axios';

export default class CreateProject extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: '',
      nameValidationMessage: null,
      backendError: null
    };
  }

  handleUpdateName = (event) => {
    this.setState({
      name: event.target.value
    });
  }

  handleSubmit = (event) => {
    event.preventDefault();

    var nameValidationMessage = utils.checkName("project", this.state.name)

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return
    }

    axios.post(`${config.endpoint}/projects`, {
      name: this.state.name
    }, {
      withCredentials: true
    })
      .then((response) => {
        segment.track('Project Created');
        this.props.history.push(`/${response.data.name}`);
      })
      .catch((error) => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Project was not created.')
          console.log(error);
        }
      });
  }

  handleCancel() {
    this.props.history.push('/projects');
  }

  render() {
    const heading = "Create Project"
    return (
      <Pane
        display="flex"
        flexDirection="column"
        alignItems="center"
        background="tint1"
        flex={1}
        justifyContent="stretch"
      >
        <TopHeader user={this.props.user} heading={heading} showLogo={true} hideSwitchProjects={true} history={this.props.history} />
        <Pane width={majorScale(50)}>
          <InnerCard>
            <Pane
              padding={majorScale(4)}
              is="form"
            >
              {this.state.backendError && (
                <Alert marginBottom={majorScale(2)} paddingTop={majorScale(2)} paddingBottom={majorScale(2)} intent="warning" title={this.state.backendError} />
              )}
              <TextInputField
                label="Name"
                onChange={this.handleUpdateName}
                value={this.state.name}
                isInvalid={this.state.nameValidationMessage !== null}
                validationMessage={this.state.nameValidationMessage}
              />
              <Pane display="flex" flex="row">
                <Button appearance="primary" onClick={this.handleSubmit}>Submit</Button>
                <Button marginLeft={majorScale(2)} onClick={() => this.handleCancel()}>Cancel</Button>
              </Pane>
            </Pane>
          </InnerCard>
        </Pane>
      </Pane>
    );
  }
}
