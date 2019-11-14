import React, { Component, Fragment } from 'react';
import {
  Button,
  Pane,
  Heading,
  majorScale,
  Link,
  Alert,
  minorScale,
  TextInputField,
  toaster
} from 'evergreen-ui';
import axios from 'axios';

import config from '../config';
import utils from '../utils';
import Logo from '../components/logo';

export default class PasswordRecovery extends Component {
  state = {
    invalidToken: false,
    password: '',
    passwordConfirmation: '',
    passwordValidationMessage: null,
    passwordConfirmationValidationMessage: null,
    passwordRecoveryError: null
  };

  componentDidMount() {
    axios
      .get(`${config.endpoint}/passwordrecoverytokens/${this.props.token}`)
      .then(response => {
        this.setState({
          invalidToken: this.isTokenExpired(response.data.expiresAt)
        });
      })
      .catch(error => {
        this.setState({
          invalidToken: true
        });
        console.log(error);
      });
  }

  isTokenExpired(expiration) {
    const expiratonTime = new Date(expiration).getTime();
    const currentTime = new Date().getTime();
    return currentTime > expiratonTime;
  }

  handleUpdatePassword = event => {
    this.setState({
      password: event.target.value
    });
  };

  handleUpdatePasswordConfirmation = event => {
    this.setState({
      passwordConfirmation: event.target.value
    });
  };

  handleSubmit = () => {
    var passwordValidationMessage = null;
    var passwordConfirmationValidationMessage = null;

    if (!utils.passwordRegex.test(this.state.password)) {
      passwordValidationMessage = 'Please enter a valid password.';
    }

    if (
      this.state.passwordConfirmation === '' ||
      this.state.password !== this.state.passwordConfirmation
    ) {
      passwordConfirmationValidationMessage = 'Passwords must match.';
    }

    this.setState({
      passwordValidationMessage: passwordValidationMessage,
      passwordConfirmationValidationMessage: passwordConfirmationValidationMessage,
      passwordRecoveryError: null
    });

    if (
      passwordValidationMessage === null &&
      passwordConfirmationValidationMessage === null
    ) {
      axios
        .post(`${config.endpoint}/changepassword`, {
          passwordRecoveryTokenValue: this.props.token,
          password: this.state.password
        })
        .then(response => {
          toaster.success('Password successfully changed.');
          this.props.history.push(`/login`);
        })
        .catch(error => {
          if (utils.is4xx(error.response.status)) {
            this.setState({
              passwordRecoveryError: utils.convertErrorMessage(
                error.response.data
              )
            });
          } else {
            toaster.danger(
              'Something went wrong with changing your password. Please contact us at support@deviceplane.com.'
            );
            console.log(error);
          }
        });
    }
  };

  render() {
    return (
      <Pane
        display="flex"
        flexDirection="column"
        alignItems="center"
        padding={majorScale(6)}
        width={majorScale(50)}
        marginX="auto"
        marginY={200}
        elevation={2}
      >
        <Pane
          display="flex"
          alignItems="center"
          paddingBottom={majorScale(2)}
          marginX="auto"
        >
          <Pane paddingRight={minorScale(3)}>
            <Logo />
          </Pane>
          <Heading size={600}>Deviceplane</Heading>
        </Pane>
        {this.state.passwordRecoveryError && (
          <Alert
            marginBottom={majorScale(2)}
            paddingTop={majorScale(2)}
            paddingBottom={majorScale(2)}
            intent="warning"
            title={this.state.passwordRecoveryError}
          />
        )}
        {this.state.invalidToken ? (
          <Fragment>
            <Alert
              marginBottom={majorScale(2)}
              paddingTop={majorScale(2)}
              paddingBottom={majorScale(2)}
              intent="warning"
              title="Your recovery token has expired. Please reset your password again."
            />
            <Link href="/forgot">Reset your password</Link>
          </Fragment>
        ) : (
          <Fragment>
            <Pane>
              <TextInputField
                type="password"
                label="New Password"
                hint="Password must be at least 8 characters, contain at least one lower case letter, one upper case letter, and no spaces."
                required
                onChange={this.handleUpdatePassword}
                value={this.state.password}
                isInvalid={this.state.passwordValidationMessage !== null}
                validationMessage={this.state.passwordValidationMessage}
              />
              <TextInputField
                type="password"
                label="Re-enter Password"
                required
                onChange={this.handleUpdatePasswordConfirmation}
                value={this.state.passwordConfirmation}
                isInvalid={
                  this.state.passwordConfirmationValidationMessage !== null
                }
                validationMessage={
                  this.state.passwordConfirmationValidationMessage
                }
              />
            </Pane>
            <Button
              appearance="primary"
              justifyContent="center"
              onClick={this.handleSubmit}
            >
              Submit
            </Button>
          </Fragment>
        )}
      </Pane>
    );
  }
}
