import React, { Component } from 'react';
import './../App.css';
import config from '../config.js';
import utils from '../utils.js';
import { toaster, Pane, majorScale, Heading, TextInputField, Alert, Button } from 'evergreen-ui';
import axios from 'axios';

export default class EditProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      firstName: this.props.user.firstName,
      lastName: this.props.user.lastName,
      company: this.props.user.company,
      firstNameValidationMessage: null,
      lastNameValidationMessage: null,
      invalidFirstName: false,
      invalidLastName: false,
      backendError: null
    };
  }

  handleUpdateFirstName = (event) => {
    this.setState({
      firstName: event.target.value
    });
  }

  handleUpdateLastName = (event) => {
    this.setState({
      lastName: event.target.value
    });
  }

  handleUpdateCompany = (event) => {
    this.setState({
      company: event.target.value
    });
  }

  handleSubmit = () => {
    var firstNameValidationMessage = null;
    var lastNameValidationMessage = null;
    var invalidFirstName = false;
    var invalidLastName = false;

    if (!utils.usernameRegex.test(this.state.firstName)) {
      firstNameValidationMessage = 'Please enter a valid first name. Only letters are allowed.'
      invalidFirstName = true
    }

    if (!utils.usernameRegex.test(this.state.lastName)) {
      lastNameValidationMessage = 'Please enter a valid last name. Only letters are allowed.'
      invalidLastName = true
    }

    this.setState({
      firstNameValidationMessage: firstNameValidationMessage,
      lastNameValidationMessage: lastNameValidationMessage,
      invalidFirstName: invalidFirstName,
      invalidLastName: invalidLastName,
      backendError: null
    });

    if (!invalidFirstName && !invalidLastName) {
      axios.patch(`${config.endpoint}/me`, {
        firstName: this.state.firstName,
        lastName: this.state.lastName,
        company: this.state.company
      }, {
        withCredentials: true
      })
        .then((response) => {
          toaster.success('Profile updated.');
        })
        .catch((error) => {
          if (utils.is4xx(error.response.status)) {
            this.setState({
              backendError: utils.convertErrorMessage(error.response.data)
            });
          } else {
            toaster.danger('Profile was not updated.')
            console.log(error);
          }
        });
    }
  }

  render() {
    return (
      <React.Fragment>
        <Pane zIndex={1} flexShrink={0} elevation={0} backgroundColor="white">
          <Pane padding={majorScale(2)}>
            <Heading size={600}>Edit Profile</Heading>
          </Pane>
        </Pane>
        <Pane
          display="flex"
          flexDirection="column"
          margin={majorScale(4)}
        >
          <TextInputField
            label="First Name"
            required
            onChange={this.handleUpdateFirstName}
            value={this.state.firstName}
            isInvalid={this.state.invalidFirstName}
            validationMessage={this.state.firstNameValidationMessage}
          />
          <TextInputField
            label="Last Name"
            required
            onChange={this.handleUpdateLastName}
            value={this.state.lastName}
            isInvalid={this.state.invalidLastName}
            validationMessage={this.state.lastNameValidationMessage}
          />
          <TextInputField
            label="Company"
            onChange={this.handleUpdateCompany}
            value={this.state.company}
          />
          {this.state.backendError && (
            <Alert marginBottom={majorScale(2)} paddingTop={majorScale(2)} paddingBottom={majorScale(2)} intent="warning" title={this.state.backendError} />
          )}
          <Button appearance="primary" justifyContent="center" onClick={this.handleSubmit}>Submit</Button>
        </Pane>
      </React.Fragment>
    );
  }
} 