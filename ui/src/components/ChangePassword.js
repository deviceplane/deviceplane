import React, { Component } from 'react';
import './../App.css';
import config from '../config.js';
import utils from '../utils.js';
import { toaster, Pane, majorScale, Heading, TextInputField, Alert, Button } from 'evergreen-ui';
import axios from 'axios';

export default class ChangePassword extends Component {
  constructor(props) {
    super(props);
    this.state = {
      currentPassword: '',
      password: '',
      passwordConfirmation: '',
      currentPasswordValidationMessage: null,
      passwordValidationMessage: null,
      passwordConfirmationValidationMessage: null,
      changePasswordError: null
    };
  }

  handleUpdateCurrentPassword = (event) => {
    this.setState({
      currentPassword: event.target.value
    });
  }

  handleUpdatePassword = (event) => {
    this.setState({
      password: event.target.value
    });
  }

  handleUpdatePasswordConfirmation = (event) => {
    this.setState({
      passwordConfirmation: event.target.value
    });
  }

  handleSubmit = () => {
    var currentPasswordValidationMessage = null;
    var passwordValidationMessage = null;
    var passwordConfirmationValidationMessage = null;

    if (this.state.currentPassword === '') {
      currentPasswordValidationMessage = 'Please enter your current password.'
    }

    if (!utils.passwordRegex.test(this.state.password)) {
      passwordValidationMessage = 'Please enter a valid password.'
    }

    if (!utils.passwordRegex.test(this.state.passwordConfirmation)) {
      passwordConfirmationValidationMessage = 'Please enter a valid password.'
    }

    if (this.state.passwordConfirmation === '' || this.state.password !== this.state.passwordConfirmation) {
      passwordConfirmationValidationMessage = 'Passwords must match.'
    }

    if (this.state.currentPassword === this.state.password) {
      passwordValidationMessage = 'Current and new password are the same. Please pick a new password.'
    }

    this.setState({
      currentPasswordValidationMessage: currentPasswordValidationMessage,
      passwordValidationMessage: passwordValidationMessage,
      passwordConfirmationValidationMessage: passwordConfirmationValidationMessage,
      changePasswordError: null
    });

    if (currentPasswordValidationMessage === null && passwordValidationMessage === null && passwordConfirmationValidationMessage === null) {
      axios.patch(`${config.endpoint}/me`, {
        password: this.state.password,
        currentPassword: this.state.currentPassword
      }, {
        withCredentials: true
      })
        .then((response) => {
          this.setState({
            currentPassword: '',
            password: '',
            passwordConfirmation: ''
          });
          toaster.success('Password updated.');
        })
        .catch((error) => {
          if (utils.is4xx(error.response.status)) {
            this.setState({
              changePasswordError: utils.convertErrorMessage(error.response.data)
            });
          } else {
            toaster.danger('Password was not updated.')
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
            <Heading size={600}>Change Password</Heading>
          </Pane>
        </Pane>
        <Pane
          display="flex"
          flexDirection="column"
          margin={majorScale(4)}
        >
          <TextInputField
            type="password"
            label="Current Password"
            onChange={this.handleUpdateCurrentPassword}
            value={this.state.currentPassword}
            isInvalid={this.state.currentPasswordValidationMessage !== null}
            validationMessage={this.state.currentPasswordValidationMessage}
          />
          <TextInputField
            type="password"
            label="New Password"
            hint="Password must be at least 8 characters, contain at least one lower case letter, one upper case letter, and no spaces."
            onChange={this.handleUpdatePassword}
            value={this.state.password}
            isInvalid={this.state.passwordValidationMessage !== null}
            validationMessage={this.state.passwordValidationMessage}
          />
          <TextInputField
            type="password"
            label="Re-enter New Password"
            onChange={this.handleUpdatePasswordConfirmation}
            value={this.state.passwordConfirmation}
            isInvalid={this.state.passwordConfirmationValidationMessage !== null}
            validationMessage={this.state.passwordConfirmationValidationMessage}
          />
          {this.state.changePasswordError && (
            <Alert marginBottom={majorScale(2)} paddingTop={majorScale(2)} paddingBottom={majorScale(2)} intent="warning" title={this.state.changePasswordError} />
          )}
          <Button appearance="primary" justifyContent="center" onClick={this.handleSubmit}>Submit</Button>
        </Pane>
      </React.Fragment>
    );
  }
}